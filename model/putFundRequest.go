package model

type PutFundRequest struct {
	ExternalId      string   `json:"external_id"`
	Link            string   `json:"link"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	Fund            string   `json:"fund"`
	Program         string   `json:"program"`
	Call            string   `json:"call"`
	TypeOfEffort    string   `json:"type_of_effort"`
	Description     string   `json:"description"`
	TotalBudget     float64  `json:"total_budget"`
	GrantBudgetLow  float64  `json:"grant_budget_low"`
	GrantBudgetHigh float64  `json:"grant_budget_high"`
	Currency        string   `json:"currency"`
	DueDate         int64    `json:"due_date"`
	SubmissionType  string   `json:"submission_type"`
	Keywords        []string `json:"keywords"`
	Tags            []string `json:"tags"`
}
