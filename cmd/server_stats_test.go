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

func TestParseRemoteStatsUnsupported(t *testing.T) {
	s := parseRemoteStatsOutput("STATERR=unsupported\n")
	if s.Err != "non-linux" {
		t.Fatalf("got %q", s.Err)
	}
	if s.CPUString() != "—" || s.RAMString() != "—" {
		t.Fatal("expected em dash for failed stats")
	}
}
