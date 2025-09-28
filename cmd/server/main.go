package main

import (
	"fmt"

	"github.com/deuta/goTexas/src/game"
)

// rankToString 是一个辅助工具，用于将 HandRank 枚举值转换成可读的中文牌型字符串
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
	fmt.Println("--- 开始一局新的德州扑克 ---")

	// 1. 创建并洗牌
	deck := game.NewDeck()
	deck.Shuffle()
	fmt.Println("一副新牌已创建并洗好。")

	// 2. 创建玩家
	player1 := &game.Player{ID: "player-1", Name: "玩家A"}
	player2 := &game.Player{ID: "player-2", Name: "玩家B"}
	fmt.Printf("欢迎 %s 和 %s 加入牌局。\n\n", player1.Name, player2.Name)

	// 3. 发底牌 (Hole Cards)
	card1, _ := deck.Deal()
	card2, _ := deck.Deal()
	player1.Hand = []game.Card{card1, card2}

	card3, _ := deck.Deal()
	card4, _ := deck.Deal()
	player2.Hand = []game.Card{card3, card4}

	fmt.Println("--- 底牌已发出 ---")
	fmt.Printf("  -> %s 的手牌: %v %v\n", player1.Name, player1.Hand[0], player1.Hand[1])
	fmt.Printf("  -> %s 的手牌: %v %v\n\n", player2.Name, player2.Hand[0], player2.Hand[1])

	// 4. 发公共牌 (Community Cards)
	communityCards := make([]game.Card, 0, 5)
	deck.Deal() // 烧牌
	flop1, _ := deck.Deal()
	flop2, _ := deck.Deal()
	flop3, _ := deck.Deal()
	communityCards = append(communityCards, flop1, flop2, flop3)
	fmt.Println("--- 翻牌圈 (Flop) ---")
	fmt.Printf("  公共牌: %v %v %v\n\n", communityCards[0], communityCards[1], communityCards[2])

	deck.Deal() // 烧牌
	turn, _ := deck.Deal()
	communityCards = append(communityCards, turn)
	fmt.Println("--- 转牌圈 (Turn) ---")
	fmt.Printf("  公共牌: %v %v %v %v\n\n", communityCards[0], communityCards[1], communityCards[2], communityCards[3])

	deck.Deal() // 烧牌
	river, _ := deck.Deal()
	communityCards = append(communityCards, river)
	fmt.Println("--- 河牌圈 (River) ---")
	fmt.Printf("  公共牌: %v %v %v %v %v\n\n", communityCards[0], communityCards[1], communityCards[2], communityCards[3], communityCards[4])

	// 5. 摊牌与裁决 (Showdown)
	fmt.Println("--- 摊牌与裁决 ---")

	cardsForPlayer1 := append(player1.Hand, communityCards...)
	bestValuePlayer1 := game.EvaluateBestHand(cardsForPlayer1)
	fmt.Printf("  -> %s 的最佳成手: %s, 关键牌: %v\n", player1.Name, rankToString[bestValuePlayer1.HandRank], bestValuePlayer1.HighCards)

	cardsForPlayer2 := append(player2.Hand, communityCards...)
	bestValuePlayer2 := game.EvaluateBestHand(cardsForPlayer2)
	fmt.Printf("  -> %s 的最佳成手: %s, 关键牌: %v\n\n", player2.Name, rankToString[bestValuePlayer2.HandRank], bestValuePlayer2.HighCards)

	// 6. 宣布胜者
	if bestValuePlayer1.IsBetterThan(bestValuePlayer2) {
		fmt.Printf("🎉 赢家是 %s!\n", player1.Name)
	} else if bestValuePlayer2.IsBetterThan(bestValuePlayer1) {
		fmt.Printf("🎉 赢家是 %s!\n", player2.Name)
	} else {
		fmt.Println("平局! 双方平分彩池。")
	}
}
