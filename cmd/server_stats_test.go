package cmd

import "testing"

func TestParseRemoteStatsOutput(t *testing.T) {
	raw := `CORES=4
LOAD=0.52
MTOTAL=16384236
MAVAIL=8192000`
	s := parseRemoteStatsOutput(raw)
	if s.Err != "" {
		t.Fatalf("unexpected err: %q", s.Err)
	}
	if s.CPUCores != 4 || s.Load1 != 0.52 || s.MemTotalKB != 16384236 || s.MemAvailKB != 8192000 {
		t.Fatalf("parse mismatch: %+v", s)
	}
	if s.CPUString() != "0.52/4" {
		t.Fatalf("CPUString: got %q", s.CPUString())
	}
}

func TestSkipStatsGitHost(t *testing.T) {
	s := fetchServerStats(Config{Alias: "prod-git", Hostname: "10.0.0.1", User: "u", IdentityFile: "~/.ssh/id_rsa"})
	if s.Err != errStatsSkipGit || !s.Ready {
		t.Fatalf("expected git-skip: %+v", s)
	}
	if s.CPUString() != SymbolSkipped || s.RAMString() != SymbolSkipped {
		t.Fatalf("expected skipped symbol, got CPU=%q RAM=%q", s.CPUString(), s.RAMString())
	}
	s2 := fetchServerStats(Config{Alias: "app", Hostname: "git.example.com", User: "u", IdentityFile: "~/.ssh/id_rsa"})
	if s2.Err != errStatsSkipGit {
		t.Fatalf("hostname git: %+v", s2)
	}
}

func TestPendingStatsPlaceholder(t *testing.T) {
	s := ServerStats{} // Ready false
	if s.CPUString() != SymbolLoading || s.RAMString() != SymbolLoading {
		t.Fatalf("pending should show loading symbol, got CPU=%q RAM=%q", s.CPUString(), s.RAMString())
	}
}

func TestParseRemoteStatsUnsupported(t *testing.T) {
	s := parseRemoteStatsOutput("STATERR=unsupported\n")
	if s.Err != "non-linux" {
		t.Fatalf("got %q", s.Err)
	}
	if s.CPUString() != SymbolIssue || s.RAMString() != SymbolIssue {
		t.Fatalf("expected issue symbol for failed stats, got CPU=%q RAM=%q", s.CPUString(), s.RAMString())
	}
}
