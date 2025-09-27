package game

import (
	"sort"
)

// HandRank 代表扑克牌型（例如，同花顺，葫芦）。
// 其整数值允许我们直接比较牌型的大小。
type HandRank int

// 定义所有可能的牌型，按牌力从低到高排序。
// iota 是一个Go语言的关键字，它会在每个const声明中自动加一，非常适合定义枚举值。
const (
	HighCard      HandRank = iota + 1 // 高牌
	OnePair                           // 一对
	TwoPair                           // 两对
	ThreeOfAKind                      // 三条
	Straight                          // 顺子
	Flush                             // 同花
	FullHouse                         // 葫芦 (满堂红)
	FourOfAKind                       // 四条 (铁支)
	StraightFlush                     // 同花顺
	// 备注: 皇家同花顺是同花顺的一种特例
)

// rankToInt 是一个映射表，用于将牌的点数（字符串）转换为整数，以便于比较大小。
var rankToInt = map[string]int{
	"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "10": 10,
	"J": 11, "Q": 12, "K": 13, "A": 14, // Ace在这里被当作最大的点数14
}

// HandValue 代表一手5张牌的、可比较的完整牌力。
// 它可以用来决定多个玩家之间的胜负。
type HandValue struct {
	Rank      HandRank // 牌型的主要等级 (例如, FullHouse)。
	HighCards []int    // 决定牌力大小的关键牌和踢脚牌的点数，按重要性降序排列。
}

// isBetterThan 比较两个 HandValue，如果 v1 的牌力比 v2 好，则返回 true。
func (v1 HandValue) isBetterThan(v2 HandValue) bool {
	// 首先比较牌型等级
	if v1.Rank != v2.Rank {
		return v1.Rank > v2.Rank
	}

	// 如果牌型等级相同，则依次比较关键牌和踢脚牌
	for i := 0; i < len(v1.HighCards); i++ {
		if v1.HighCards[i] != v2.HighCards[i] {
			return v1.HighCards[i] > v2.HighCards[i]
		}
	}

	// 如果所有都相同，则为平局，v1不比v2好
	return false
}

// EvaluateBestHand 会接收7张牌，并返回其中能组成的最佳5张牌的 HandValue。
// 这是我们最终要暴露给外部的函数。
func EvaluateBestHand(cards []Card) HandValue {
	if len(cards) < 5 {
		return HandValue{}
	}

	// 初始化一个“最差”的牌力作为当前的最佳牌力
	bestValue := HandValue{Rank: HighCard, HighCards: []int{0}}

	// 通过5层嵌套循环，生成 C(7,5) = 21 种所有可能的5张牌组合
	// 这是最直接的组合生成算法，对于7张牌这个固定的小数目是完全足够且清晰的。
	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {
			for k := j + 1; k < len(cards); k++ {
				for l := k + 1; l < len(cards); l++ {
					for m := l + 1; m < len(cards); m++ {
						// 得到一个5张牌的组合
						currentHand := []Card{cards[i], cards[j], cards[k], cards[l], cards[m]}

						// 评估这个组合的牌力
						currentValue := evaluateFiveCardHand(currentHand)

						// 如果当前组合比我们记录的最佳牌力还要好，就更新它
						if currentValue.isBetterThan(bestValue) {
							bestValue = currentValue
						}
					}
				}
			}
		}
	}

	return bestValue
}

// evaluateFiveCardHand 判断任意5张牌的 HandValue。
// 这是单个5张牌组合的核心评估逻辑。
func evaluateFiveCardHand(cards []Card) HandValue {
	if len(cards) != 5 {
		return HandValue{Rank: HighCard}
	}

	// --- 1. 数据准备 ---
	ranks := make([]int, 5)
	for i, c := range cards {
		ranks[i] = rankToInt[c.Rank]
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ranks)))

	// --- 2. 检查是否为同花或顺子 ---
	isFlush := cards[0].Suit == cards[1].Suit && cards[0].Suit == cards[2].Suit && cards[0].Suit == cards[3].Suit && cards[0].Suit == cards[4].Suit
	isStraight := ranks[0]-ranks[4] == 4 && ranks[0]-ranks[1] == 1 && ranks[1]-ranks[2] == 1 && ranks[2]-ranks[3] == 1 && ranks[3]-ranks[4] == 1
	if !isStraight && ranks[0] == 14 && ranks[1] == 5 && ranks[2] == 4 && ranks[3] == 3 && ranks[4] == 2 {
		isStraight = true
		ranks = []int{5, 4, 3, 2, 1}
	}

	// --- 3. 检查高级牌型 (同花顺等) ---
	if isStraight && isFlush {
		return HandValue{Rank: StraightFlush, HighCards: []int{ranks[0]}}
	}

	// --- 4. 检查对子、三条、四条、葫芦 ---
	rankCounts := make(map[int]int)
	for _, r := range ranks {
		rankCounts[r]++
	}

	var pairs, threes, fours []int
	for rank, count := range rankCounts {
		switch count {
		case 2:
			pairs = append(pairs, rank)
		case 3:
			threes = append(threes, rank)
		case 4:
			fours = append(fours, rank)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(pairs)))

	if len(fours) == 1 {
		kicker := 0
		for _, r := range ranks {
			if r != fours[0] {
				kicker = r
				break
			}
		}
		return HandValue{Rank: FourOfAKind, HighCards: []int{fours[0], kicker}}
	}

	if len(threes) == 1 && len(pairs) == 1 {
		return HandValue{Rank: FullHouse, HighCards: []int{threes[0], pairs[0]}}
	}

	if isFlush {
		return HandValue{Rank: Flush, HighCards: ranks}
	}

	if isStraight {
		return HandValue{Rank: Straight, HighCards: []int{ranks[0]}}
	}

	if len(threes) == 1 {
		kickers := make([]int, 0, 2)
		for _, r := range ranks {
			if r != threes[0] {
				kickers = append(kickers, r)
			}
		}
		return HandValue{Rank: ThreeOfAKind, HighCards: append([]int{threes[0]}, kickers...)}
	}

	if len(pairs) == 2 {
		kicker := 0
		for _, r := range ranks {
			if r != pairs[0] && r != pairs[1] {
				kicker = r
				break
			}
		}
		return HandValue{Rank: TwoPair, HighCards: []int{pairs[0], pairs[1], kicker}}
	}

	if len(pairs) == 1 {
		kickers := make([]int, 0, 3)
		for _, r := range ranks {
			if r != pairs[0] {
				kickers = append(kickers, r)
			}
		}
		return HandValue{Rank: OnePair, HighCards: append([]int{pairs[0]}, kickers...)}
	}

	// --- 5. 如果什么都不是，那就是高牌 ---
	return HandValue{Rank: HighCard, HighCards: ranks}
}
