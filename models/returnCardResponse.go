package models

type DiscardReturn struct {
	Remaining int `json:"remaining"`
}

type PileReturn struct {
	Discard DiscardReturn `json:"discard"`
}

type returnCardResponse struct {
	Success   bool       `json:"success"`
	DeckId    string     `json:"deck_id"`
	Remaining int        `json:"remaining"`
	Piles     PileReturn `json:"piles"`
}
