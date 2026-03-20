package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ServerStats holds CPU/RAM metrics gathered from a remote Linux host over SSH.
// Values are best-effort; Err is set when the host is unreachable or metrics cannot be read.
// Ready is false until the first fetch completes (used for progressive --stats display).
type ServerStats struct {
	Ready      bool
	CPUCores   int
	Load1      float64
	MemTotalKB uint64
	MemAvailKB uint64
	Err        string
}

// remoteStatsScript runs on the SSH target. It assumes a Linux system with /proc.
// Output lines: CORES=, LOAD=, MTOTAL=, MAVAIL=, optional STATERR=unsupported
const remoteStatsScript = `bash -c 'if ! test -r /proc/meminfo; then echo STATERR=unsupported; exit 0; fi; echo CORES=$(nproc); echo LOAD=$(cut -d" " -f1 /proc/loadavg); echo MTOTAL=$(grep MemTotal /proc/meminfo | awk "{print \$2}"); MA=$(grep MemAvailable /proc/meminfo | awk "{print \$2}"); if [ -z "$MA" ]; then MA=$(grep MemFree /proc/meminfo | awk "{print \$2}"); fi; echo MAVAIL=$MA'`

const statsSessionTimeout = 15 * time.Second

// errStatsSkipGit is set when CPU/RAM fetch is skipped because alias or host matches the git filter.
const errStatsSkipGit = "git-skip"

// skipStatsDueToGitHost returns true if we should not SSH for metrics (alias or hostname contains "git", case-insensitive).
func skipStatsDueToGitHost(c Config) bool {
	a := strings.ToLower(strings.TrimSpace(c.Alias))
	h := strings.ToLower(strings.TrimSpace(c.Hostname))
	return strings.Contains(a, "git") || strings.Contains(h, "git")
}

// fetchServerStats connects once, runs the remote probe, and parses the result.
func fetchServerStats(config Config) ServerStats {
	if skipStatsDueToGitHost(config) {
		return ServerStats{Ready: true, Err: errStatsSkipGit}
	}

	client, err := dialSSHClient(config)
	if err != nil {
		return ServerStats{Err: shortErr(err), Ready: true}
	}
	defer func() { _ = client.Close() }()

	session, err := client.NewSession()
	if err != nil {
		return ServerStats{Err: shortErr(err), Ready: true}
	}
	defer func() { _ = session.Close() }()

	outCh := make(chan []byte, 1)
	errCh := make(chan error, 1)
	go func() {
		b, e := session.CombinedOutput(remoteStatsScript)
		if e != nil {
			errCh <- e
			return
		}
		outCh <- b
	}()

	var out []byte
	select {
	case out = <-outCh:
	case err := <-errCh:
		return ServerStats{Err: shortErr(err), Ready: true}
	case <-time.After(statsSessionTimeout):
		_ = session.Close()
		return ServerStats{Err: "timeout", Ready: true}
	}

	return parseRemoteStatsOutput(string(out))
}

func shortErr(err error) string {
	s := err.Error()
	if len(s) > 48 {
		return s[:45] + "..."
	}
	return s
}

func parseRemoteStatsOutput(raw string) ServerStats {
	var s ServerStats
	for _, line := range strings.Split(strings.TrimSpace(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "STATERR="):
			s.Err = strings.TrimPrefix(line, "STATERR=")
			if s.Err == "unsupported" {
				s.Err = "non-linux"
			}
		case strings.HasPrefix(line, "CORES="):
			n, _ := strconv.Atoi(strings.TrimPrefix(line, "CORES="))
			s.CPUCores = n
		case strings.HasPrefix(line, "LOAD="):
			f, _ := strconv.ParseFloat(strings.TrimPrefix(line, "LOAD="), 64)
			s.Load1 = f
		case strings.HasPrefix(line, "MTOTAL="):
			s.MemTotalKB = parseUint(strings.TrimPrefix(line, "MTOTAL="))
		case strings.HasPrefix(line, "MAVAIL="):
			s.MemAvailKB = parseUint(strings.TrimPrefix(line, "MAVAIL="))
		}
	}
	if s.Err != "" {
		s.Ready = true
		return s
	}
	if s.MemTotalKB == 0 && s.CPUCores == 0 {
		s.Err = "empty"
	}
	s.Ready = true
	return s
}

func parseUint(s string) uint64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

// CPUString formats load average and core count for table display.
func (s ServerStats) CPUString() string {
	if !s.Ready {
		return SymbolLoading
	}
	if s.Err == "empty" {
		return SymbolEmpty
	}
	if s.Err == errStatsSkipGit {
		return SymbolSkipped
	}
	if s.Err != "" {
		return SymbolIssue
	}
	if s.CPUCores <= 0 {
		return SymbolEmpty
	}
	return fmt.Sprintf("%.2f/%d", s.Load1, s.CPUCores)
}

// RAMString formats used/total and usage percent (KiB from /proc/meminfo).
func (s ServerStats) RAMString() string {
	if !s.Ready {
		return SymbolLoading
	}
	if s.Err == "empty" {
		return SymbolEmpty
	}
	if s.Err == errStatsSkipGit {
		return SymbolSkipped
	}
	if s.Err != "" {
		return SymbolIssue
	}
	if s.MemTotalKB == 0 {
		return SymbolEmpty
	}
	used := s.MemTotalKB - s.MemAvailKB
	if s.MemAvailKB > s.MemTotalKB {
		used = 0
	}
	pct := float64(used) * 100 / float64(s.MemTotalKB)
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	return fmt.Sprintf("%.0f%% %s/%s", pct, formatKiBShort(used), formatKiBShort(s.MemTotalKB))
}

func formatKiBShort(kib uint64) string {
	const kibPerGib = 1024 * 1024
	if kib >= kibPerGib {
		return fmt.Sprintf("%.1fG", float64(kib)/float64(kibPerGib))
	}
	if kib >= 1024 {
		return fmt.Sprintf("%.0fM", float64(kib)/1024)
	}
	return fmt.Sprintf("%dK", kib)
}

// fetchAllServerStats queries every config in parallel (bounded concurrency).
func fetchAllServerStats(configs []Config) []ServerStats {
	out := make([]ServerStats, len(configs))
	if len(configs) == 0 {
		return out
	}

	const maxConcurrent = 8
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	for i := range configs {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			out[i] = fetchServerStats(configs[i])
		}()
	}
	wg.Wait()
	return out
}
