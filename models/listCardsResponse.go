package models

type CardList struct {
	Image string `json:"image"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

type Player1 struct {
	Cards     []CardList `json:"cards"`
	Remaining string     `json:"remaining"`
}

type Player2 struct {
	Cards     []CardList `json:"cards"`
	Remaining string     `json:"remaining"`
}

type Piles struct {
	Player1 Player1 `json:"player1"`
	Player2 Player2 `json:"player2"`
}

type ListCardResponse struct {
	Success   bool   `json:"success"`
	DeckId    string `json:"deck_id"`
	Remaining string `json:"remaining"`
	Piles     Pile   `json:"piles"`
}
