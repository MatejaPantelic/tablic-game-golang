package models

type Images struct {
	Svg string `json:"svg"`
	Png string `json:"png"`
}
type CardList struct {
	Image  string `json:"image"`
	Value  string `json:"value"`
	Suit   string `json:"suit"`
	Code   string `json:"code"`
	Images Images `json:"images"`
}
type PileList struct {
	Cards     []CardList `json:"cards"`
	Remaining int        `json:"remaining"`
}

type Piles struct {
	Hand1  PileList `json:"hand1"`
	Hand2  PileList `json:"hand2"`
	Taken1 PileList `json:"taken1"`
	Taken2 PileList `json:"taken2"`
	Table  PileList `json:"table"`
}
type DrawResponse struct {
	Success   string `json:"success"`
	DeckId    string `json:"deck_id"`
	Remaining int    `json:"remaining"`
	Piles     Piles  `json:"piles"`
}
