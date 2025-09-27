package game

import (
	"fmt"
	"math/rand"
	"time"
)

// Suits defines the four suits in a deck of cards.
var Suits = []string{"Hearts", "Diamonds", "Clubs", "Spades"}

// Ranks defines the 13 ranks in a deck of cards.
var Ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

// NewDeck creates and returns a new, unshuffled deck of 52 cards.
// It returns a pointer to the Deck to encourage modification only through methods.
func NewDeck() *Deck {
	cards := make([]Card, 0, 52)
	for _, suit := range Suits {
		for _, rank := range Ranks {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	return &Deck{cards: cards}
}

// Shuffle randomizes the order of cards in the deck.
// It uses the Fisher-Yates shuffle algorithm.
func (d *Deck) Shuffle() {
	// Seed the random number generator to ensure different shuffles each time.
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Deal removes and returns the top card from the deck.
// It returns an error if the deck is empty.
func (d *Deck) Deal() (Card, error) {
	if len(d.cards) == 0 {
		return Card{}, fmt.Errorf("cannot deal from an empty deck")
	}
	card := d.cards[0]
	d.cards = d.cards[1:]
	return card, nil
}

// Len returns the number of cards remaining in the deck.
func (d *Deck) Len() int {
	return len(d.cards)
}
