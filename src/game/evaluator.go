package game

import (
	"fmt"
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
var rankToInt = map[string]int{"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "10": 10, "J": 11, "Q": 12, "K": 13, "A": 14}

// HandValue 代表一手5张牌的、可比较的完整牌力。
// 它可以用来决定多个玩家之间的胜负。
type HandValue struct {
	HandRank  HandRank
	HighCards []int
}

// IsBetterThan 比较两个 HandValue，如果 v1 的牌力比 v2 好，则返回 true。
// 这个方法是公开的（首字母大写），所以可以在其他包（如 main）中被调用。
func (V1 HandValue) IsBetterThan(V2 HandValue) bool {
	// 直接比较牌型等级
	if V1.HandRank != V2.HandRank {
		return V1.HandRank > V2.HandRank
	}
	// 相同牌型，逐个比较大小
	for i := 0; i < len(V1.HighCards); i++ {
		fmt.Printf("v1:%+v, v2:%+v\n", V1.HighCards[i], V2.HighCards[i])
		if V1.HighCards[i] != V2.HighCards[i] {

			return V1.HighCards[i] > V2.HighCards[i]
		}
	}
	// 所有都比较过了，平局
	return false
}

// EvaluateBestHand 会接收7张牌，并返回其中能组成的最佳5张牌的 HandValue。
// 这是我们最终要暴露给外部的函数。
func EvaluateBestHand(cards []Card) HandValue {
	// 初始化一个“最差”的牌力作为当前的最佳牌力
	bestHandValue := HandValue{HandRank: HighCard, HighCards: []int{0, 0, 0, 0, 0}}
	// 通过5层嵌套循环，生成 C(7,5) = 21 种所有可能的5张牌组合
	// 这是最直接的组合生成算法，对于7张牌这个固定的小数目是完全足够且清晰的。

	for i := 0; i < len(cards)-4; i++ {
		for j := i + 1; j < len(cards)-3; j++ {
			for k := j + 1; k < len(cards)-2; k++ {
				for l := k + 1; l < len(cards)-1; l++ {
					for m := l + 1; m < len(cards); m++ {
						// 得到一个5张牌的组合
						currentHand := []Card{cards[i], cards[j], cards[k], cards[l], cards[m]}
						// 评估这个组合的牌力
						currentHandValue := evaluateFiveCardHand(currentHand)
						fmt.Printf("currentHand: %+v\n", currentHandValue)
						if currentHandValue.IsBetterThan(bestHandValue) {
							// 如果当前组合比我们记录的最佳牌力还要好，就更新它
							bestHandValue = currentHandValue
						}
					}
				}
			}
		}
	}
	return bestHandValue
}

// evaluateFiveCardHand 判断任意5张牌的 HandValue。
// 这是单个5张牌组合的核心评估逻辑。
func evaluateFiveCardHand(cards []Card) HandValue {
	// 这个函数只应该被传入5张牌。
	if len(cards) != 5 {
		return HandValue{HandRank: HighCard, HighCards: []int{0}}
	}
	// --- 1. 数据准备 ---
	// 创建一个整数切片，用来存放5张牌的点数。
	ranks := make([]int, 5)
	for index, card := range cards {
		ranks[index] = rankToInt[card.Rank]
	}
	// 将点数按从大到小的顺序排序 (例如, [14, 13, 10, 5, 2])
	sort.Sort(sort.Reverse(sort.IntSlice(ranks)))
	// --- 2. 检查是否为同花或顺子 ---
	// 检查是否为同花：所有牌的花色都与第一张牌相同。
	isFlush := cards[0].Suit == cards[1].Suit && cards[0].Suit == cards[2].Suit && cards[0].Suit == cards[3].Suit && cards[0].Suit == cards[4].Suit
	// 检查是否为普通顺子 (例如, 9,8,7,6,5)。
	isStraight := ranks[0]-ranks[4] == 4 && ranks[1]-ranks[2] == 1 && ranks[2]-ranks[3] == 1 && ranks[3]-ranks[4] == 1
	// 检查是否为特殊的 A-5 顺子 (例如, A,5,4,3,2)。此时ranks是 [14,5,4,3,2]。
	if !isStraight && ranks[0] == 14 {
		// 对于 A-5 顺子，为了比较大小，我们把A当作点数1来处理。
		ranks[0] = 1
		sort.Sort(sort.Reverse(sort.IntSlice(ranks)))
		isStraight = true
	}

	// --- 3. 检查高级牌型 (同花顺等) ---
	// 如果既是同花也是顺子，那就是同花顺。
	if isFlush && isStraight {
		// 同花顺的大小由最大的那张牌决定。
		return HandValue{HandRank: StraightFlush, HighCards: []int{ranks[0]}}
	}
	// --- 4. 检查对子、三条、四条、葫芦 ---
	// 创建一个映射表，统计每个点数出现的次数。
	rankCount := make(map[int]int)
	for _, rank := range ranks {
		rankCount[rank]++
	}
	// 创建三个切片，分别存放找到的对子、三条、四条的点数。
	var pairs, threes, fours []int
	for rank, count := range rankCount {
		switch count {
		case 2:
			pairs = append(pairs, rank)
		case 3:
			threes = append(threes, rank)
		case 4:
			fours = append(fours, rank)
			fmt.Printf("rank: %d, four: %d\n", rank, fours[0])
		}
	}
	// 如果找到了多个对子，按大小排序（比如两对，A对和2对，A在前）。
	sort.Sort(sort.Reverse(sort.IntSlice(pairs)))
	// 如果找到了一个四条
	if len(fours) == 1 {
		// 找到那张唯一的“踢脚牌”
		var kicker int
		for _, rank := range ranks {
			if rank != fours[0] {
				kicker = rank
				break
			}
		}
		return HandValue{HandRank: FourOfAKind, HighCards: []int{fours[0], kicker}}
	}

	// 如果找到了一个三条和一个对子，那就是葫芦。
	if len(threes) == 1 && len(pairs) == 1 {
		return HandValue{HandRank: FullHouse, HighCards: []int{threes[0], pairs[0]}}
	}
	// 如果是同花 (但不是同花顺)
	if isFlush {
		// 同花的大小由5张牌从大到小依次比较决定。
		return HandValue{HandRank: Flush, HighCards: ranks}
	}
	// 如果是顺子 (但不是同花顺)
	if isStraight {
		// 顺子的大小由最大的那张牌决定。
		return HandValue{HandRank: Straight, HighCards: []int{ranks[0]}}
	}
	// 如果只找到了一个三条
	if len(threes) == 1 {
		// 找到另外两张“踢脚牌”
		var kickers []int
		for _, rank := range ranks {
			if rank != threes[0] {
				kickers = append(kickers, rank)
			}
		}
		return HandValue{HandRank: ThreeOfAKind, HighCards: append([]int{threes[0]}, kickers...)}
	}

	// 如果找到了两个对子
	if len(pairs) == 2 {
		// 找到最后那张“踢脚牌”
		var kicker int
		for _, rank := range ranks {
			if rank != pairs[0] && rank != pairs[1] {
				kicker = rank
				break
			}
		}
		return HandValue{HandRank: TwoPair, HighCards: append([]int{pairs[0]}, kicker, pairs[1])}
	}

	// 如果只找到了一个对子
	if len(pairs) == 1 {
		// 找到另外三张“踢脚牌”
		var kicker []int
		for _, rank := range ranks {
			if rank != pairs[0] {
				kicker = append(kicker, rank)
				break
			}
		}
		return HandValue{HandRank: OnePair, HighCards: append([]int{pairs[0]}, kicker...)}
	}

	// --- 5. 如果什么都不是，那就是高牌 ---
	// 高牌的大小由5张牌从大到小依次比较决定。
	return HandValue{HandRank: HighCard, HighCards: ranks}
}
