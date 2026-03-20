package cmd

const (
	SymbolEmpty   = "∅" // Missing or no data
	SymbolIssue   = "⚠" // Fetch / connection / unsupported host issue
	SymbolLoading = "…" // Metrics still loading (--stats)
	SymbolSkipped = "—" // Stats not fetched (e.g. git alias/host filter)
)
