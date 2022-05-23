package model

type Program struct {
	ProgramName string                 `json:"program_name"`
	Funds       map[string]interface{} `json:"funds"`
}
