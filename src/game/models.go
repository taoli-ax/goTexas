package game

type card struct {
	Suit string
	Rank string
}

type deck struct {
	Cards []card
}

type Player struct {
	ID    string
	Name  string
	Chips int
	Hand  []card
}

type GameState struct {
	Players       []*Player
	Deck          deck
	CommunityCard []card
	Pot           int
}
