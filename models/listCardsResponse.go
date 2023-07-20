package models

type CardList struct {
	Image string `json:"image"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

type Players struct {
	Cards     []CardList `json:"cards"`
	Remaining string     `json:"remaining"`
}

type Piles struct {
	Players []Players `json:"players"`
}

type ListCardResponse struct {
	Success   bool   `json:"success"`
	DeckId    string `json:"deck_id"`
	Remaining string `json:"remaining"`
	Piles     Pile   `json:"piles"`
}
