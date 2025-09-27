package game

// Card represents a single playing card with a suit and rank.
type Card struct {
	Suit string // e.g., "Hearts", "Diamonds", "Clubs", "Spades"
	Rank string // e.g., "2", "3", "10", "J", "Q", "K", "A"
}

// Deck represents a deck of cards.
// It encapsulates the card slice to provide safer, more controlled access.
type Deck struct {
	cards []Card
}

// Player represents a player in the game.
type Player struct {
	ID    string
	Name  string
	Chips int
	Hand  []Card // The player's two private hole cards
}

// GameState represents the overall state of a single poker game.
// This will be expanded significantly later.
type GameState struct {
	Players       []*Player
	Deck          Deck
	CommunityCard []Card
	Pot           int
}
