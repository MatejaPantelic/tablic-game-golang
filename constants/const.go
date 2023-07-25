package constants

const (
	NEW_DECK_URL             = "https://www.deckofcardsapi.com/api/deck/new/"
	LIST_PILE_CARDS_URL      = "https://www.deckofcardsapi.com/api/deck/%s/pile/%s/list/"
	ADD_TO_PILE_URL          = "https://www.deckofcardsapi.com/api/deck/%s/pile/%s/add/?cards=%s"
	DRAW_CARDS_FROM_PILE_URL = "https://www.deckofcardsapi.com/api/deck/%s/pile/%s/draw/?cards=%s"
	AUTH_HEADER_MISSING      = "Auth header missing"
	INVALID_TOKEN            = "Invalid token"
	TOKEN_EXPIRED            = "Token expired"
	FORBIDDEN_ACCESS         = "Forbidden"
)
