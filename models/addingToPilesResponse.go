package models

type Discard struct {
	Remaining int `json:"remaining"`
}

type Pile struct {
	Discard Discard `json:"discard"`
}

type AddingToPilesResponse struct {
	Success   bool   `json:"success"`
	DeckId    string `json:"deck_id"`
	Remaining int    `json:"remaining"`
	Piles     Pile   `json:"piles"`
}
