package main

import (
	"fmt"

	"github.com/deuta/goTexas/src/game"
)

// a helper to map HandRank to a string for printing
var rankToString = map[game.HandRank]string{
	game.HighCard:      "高牌",
	game.OnePair:       "一对",
	game.TwoPair:       "两对",
	game.ThreeOfAKind:  "三条",
	game.Straight:      "顺子",
	game.Flush:         "同花",
	game.FullHouse:     "葫芦",
	game.FourOfAKind:   "四条",
	game.StraightFlush: "同花顺",
}

func main() {
	// --- 开始一局新的德州扑克 ---
	fmt.Printf("--- 开始一局新的德州扑克 ---\n")
	// 1. 创建并洗牌
	deck := game.NewDeck()
	fmt.Printf("--- 洗牌 ---\n")
	deck.Shuffle()
	// 2. 创建玩家
	Alice := &game.Player{
		ID:    "player1",
		Name:  "Alice",
		Chips: 1000,
	}
	Bob := &game.Player{
		ID:    "player2",
		Name:  "Bob",
		Chips: 1000,
	}
	fmt.Printf("--- 创建玩家 ---\n")
	fmt.Printf("Alice: %+v\n", Alice)
	fmt.Printf("Bob: %+v\n", Bob)
	// 3. 发底牌 (Hole Cards)
	fmt.Printf("--- 发底牌 ---\n")
	card1, _ := deck.Deal()
	card2, _ := deck.Deal()
	Alice.Hand = append(Alice.Hand, card1, card2)
	fmt.Printf("Alice's hand: %+v\n", Alice.Hand)
	card3, _ := deck.Deal()
	card4, _ := deck.Deal()
	Bob.Hand = append(Bob.Hand, card3, card4)
	fmt.Printf("Bob's hand: %+v\n", Bob.Hand)
	// 4. 发公共牌 (Community Cards)
	fmt.Printf("--- 发公共牌 ---\n")
	communityCards := make([]game.Card, 0, 5)
	// 翻牌 (Flop)
	fmt.Printf("--- 翻牌 ---\n")
	for i := 0; i < 3; i++ {
		card, _ := deck.Deal()
		communityCards = append(communityCards, card)
	}
	// 转牌 (Turn)
	fmt.Printf("--- 转牌 ---\n")
	turnCard, _ := deck.Deal()
	communityCards = append(communityCards, turnCard)
	// 河牌 (River)
	fmt.Printf("--- 河牌 ---\n")
	riverCard, _ := deck.Deal()
	communityCards = append(communityCards, riverCard)
	// 5. 摊牌与裁决 (Showdown)\
	fmt.Printf("--- 摊牌与裁决 ---\n")
	// 评估玩家A的最佳手牌
	cardsForAlice := append(Alice.Hand, communityCards...)
	fmt.Printf("--- 评估玩家A的最佳手牌 ---%+v\n", cardsForAlice)
	bestValueHandForAlice := game.EvaluateBestHand(cardsForAlice)

	fmt.Printf("Alice's best hand: %+v\n", bestValueHandForAlice)
	// 评估玩家B的最佳手牌
	cardsForBob := append(Bob.Hand, communityCards...)
	bestValueHandForBob := game.EvaluateBestHand(cardsForBob)
	fmt.Printf("Bob's best hand: %+v\n", bestValueHandForBob)
	// 6. 宣布胜者
	bestValueHandForAlice.IsBetterThan(bestValueHandForBob)
	fmt.Printf("--- 宣布胜者 ---\n")
	if bestValueHandForAlice.IsBetterThan(bestValueHandForBob) {
		fmt.Printf("Alice wins! cards are %+v, 关键牌:%+v \n", rankToString[bestValueHandForAlice.HandRank], bestValueHandForAlice.HighCards)
	} else if bestValueHandForBob.IsBetterThan(bestValueHandForAlice) {
		fmt.Printf("Bob wins! cards are %+v, 关键牌:%+v \n", rankToString[bestValueHandForBob.HandRank], bestValueHandForBob.HighCards)
	} else {
		fmt.Printf("平局！\n")
	}
}
