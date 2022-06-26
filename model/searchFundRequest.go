package model

type SearchFundRequest struct {
	Term   string `json:"term"`
	Fund   string `json:"fund"`
	OnlyEU bool   `json:"only_eu"`
	OnlyUS bool   `json:"only_us"`
}
