package model

type SearchFundRequest struct {
	Term string `json:"term"`
	Fund string `json:"fund"`
}
