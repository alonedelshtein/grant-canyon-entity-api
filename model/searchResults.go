package model

type SearchResults struct {
	Programs []ProgramSearchResult `json:"programs"`
}

type ProgramSearchResult struct {
	ProgramName string        `json:"program_name"`
	Funds       []interface{} `json:"funds"`
}
