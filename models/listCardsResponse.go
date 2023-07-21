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

type Piles struct { //Ili ovde Pile?
	Players []Players `json:"players"`
}

type ListCardResponse struct {
	Success   bool   `json:"success"`
	DeckId    string `json:"deck_id"`
	Remaining string `json:"remaining"`
	Piles     Piles   `json:"piles"` //Da li ovde treba Piles?
}
//Moze li umesto piles direktno Players[]