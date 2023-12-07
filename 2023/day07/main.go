package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type Card = int8

const (
	CardJ Card = 1
	Card2 Card = 2
	Card3 Card = 3
	Card4 Card = 4
	Card5 Card = 5
	Card6 Card = 6
	Card7 Card = 7
	Card8 Card = 8
	Card9 Card = 9
	CardT Card = 10
	CardQ Card = 11
	CardK Card = 12
	CardA Card = 13
)

func CardString(card Card) string {
	switch card {
	case CardJ:
		return "J"
	case Card2:
		return "2"
	case Card3:
		return "3"
	case Card4:
		return "4"
	case Card5:
		return "5"
	case Card6:
		return "6"
	case Card7:
		return "7"
	case Card8:
		return "8"
	case Card9:
		return "9"
	case CardT:
		return "T"
	case CardQ:
		return "Q"
	case CardK:
		return "K"
	case CardA:
		return "A"
	default:
		log.Fatalf("Invalid card: %d", card)
		return "Invalid"
	}
}

type Combination = int8

const (
	CombinationFiveOfAKind  Combination = 7 // Five of a kind, where all five cards have the same label: AAAAA
	CombinationFourOfAKind  Combination = 6 // Four of a kind, where four cards have the same label and one card has a different label: AA8AA
	CombinationFullHouse    Combination = 5 // Full house, where three cards have the same label, and the remaining two cards share a different label: 23332
	CombinationThreeOfAKind Combination = 4 // Three of a kind, where three cards have the same label, and the remaining two cards are each different from any other card in the hand: TTT98
	CombinationTwoPair      Combination = 3 // Two pair, where two cards share one label, two other cards share a second label, and the remaining card has a third label: 23432
	CombinationOnePair      Combination = 2 // One pair, where two cards share one label, and the other three cards have a different label from the pair and each other: A23A4
	CombinationHighCard     Combination = 1 // High card, where all cards' labels are distinct: 23456
)

func CombinationString(combo Combination) string {
	switch combo {
	case CombinationFiveOfAKind:
		return "FiveOfAKind"
	case CombinationFourOfAKind:
		return "FourOfAKind"
	case CombinationFullHouse:
		return "FullHouse"
	case CombinationThreeOfAKind:
		return "ThreeOfAKind"
	case CombinationTwoPair:
		return "TwoPair"
	case CombinationOnePair:
		return "OnePair"
	case CombinationHighCard:
		return "HighCard"
	default:
		log.Fatalf("Invalid combination: %d", combo)
		return "Invalid"
	}
}

type Hand struct {
	cards []Card
	combo Combination
	Bid   int64
}

func (h Hand) String() string {
	cardsStr := "Cards("
	for _, card := range h.cards {
		cardsStr += CardString(card)
	}
	return cardsStr + ") Bid(" + strconv.FormatInt(h.Bid, 10) + ") Combo(" + CombinationString(h.combo) + ")"
}

func (h Hand) Score() int64 {
	return int64(h.combo)*10000000000 + int64(h.cards[0])*100000000 + int64(h.cards[1])*1000000 + int64(h.cards[2])*10000 + int64(h.cards[3])*100 + int64(h.cards[4])
}

func getCardsCombination(cards []Card) Combination {
	cardsCount := make(map[Card]int8)
	for _, card := range cards {
		cardsCount[card]++
	}
	threeIdenticalCards := false
	twoIdenticalCards := 0
	nbJokers := cardsCount[CardJ]
	for card, count := range cardsCount {
		if card == CardJ {
			continue
		}
		if count == 5 {
			return CombinationFiveOfAKind
		}
		if count == 4 {
			if nbJokers > 0 {
				return CombinationFiveOfAKind
			}
			return CombinationFourOfAKind
		}
		if count == 3 {
			threeIdenticalCards = true
		}
		if count == 2 {
			twoIdenticalCards += 1
		}
	}
	if nbJokers == 0 {
		if threeIdenticalCards && twoIdenticalCards > 0 {
			return CombinationFullHouse
		} else if threeIdenticalCards {
			return CombinationThreeOfAKind
		} else if twoIdenticalCards == 2 {
			return CombinationTwoPair
		} else if twoIdenticalCards == 1 {
			return CombinationOnePair
		} else if twoIdenticalCards == 0 {
			return CombinationHighCard
		} else {
			log.Fatalf("Invalid combination with 0 jokers: %v", cardsCount)
		}
	} else if nbJokers == 1 {
		if threeIdenticalCards {
			return CombinationFourOfAKind
		} else if twoIdenticalCards == 2 {
			return CombinationFullHouse
		} else if twoIdenticalCards == 1 {
			return CombinationThreeOfAKind
		} else if twoIdenticalCards == 0 {
			return CombinationOnePair
		} else {
			log.Fatalf("Invalid combination with 1 jokers: %v", cardsCount)
		}
	} else if nbJokers == 2 {
		if threeIdenticalCards {
			return CombinationFiveOfAKind
		} else if twoIdenticalCards == 2 {
			log.Fatalf("Invalid combination with 2 jokers: %v", cardsCount)
		} else if twoIdenticalCards == 1 {
			return CombinationFourOfAKind
		} else if twoIdenticalCards == 0 {
			return CombinationThreeOfAKind
		} else {
			log.Fatalf("Invalid combination with 2 jokers: %v", cardsCount)
		}
	} else if nbJokers == 3 {
		if twoIdenticalCards == 1 {
			return CombinationFiveOfAKind
		} else {
			return CombinationFourOfAKind
		}
	} else if nbJokers == 4 || nbJokers == 5 {
		return CombinationFiveOfAKind
	}

	log.Fatalf("Invalid combination: %v", cardsCount)
	return -1
}

func parseNumber(numberStr string) int64 {
	number, errParsing := strconv.ParseInt(numberStr, 10, 64)
	if errParsing != nil {
		log.Fatalf("Unable to parse number: %s", numberStr)
	}
	return number
}

func parseCard(cardStr rune) Card {
	switch cardStr {
	case '2':
		return Card2
	case '3':
		return Card3
	case '4':
		return Card4
	case '5':
		return Card5
	case '6':
		return Card6
	case '7':
		return Card7
	case '8':
		return Card8
	case '9':
		return Card9
	case 'T':
		return CardT
	case 'J':
		return CardJ
	case 'Q':
		return CardQ
	case 'K':
		return CardK
	case 'A':
		return CardA
	default:
		log.Fatalf("Unable to parse card: %v", cardStr)
	}
	return -1
}

var handRegex = regexp.MustCompile(`([2-9TJQKA]{5}) ([0-9]+)`)

func parseHand(handStr string) Hand {
	matches := handRegex.FindStringSubmatch(handStr)
	if len(matches) != 3 {
		log.Fatalf("Unable to parse hand: %s", handStr)
	}
	hand := Hand{
		cards: make([]Card, 5),
	}
	cards := matches[1]
	for i, cardStr := range cards {
		hand.cards[i] = parseCard(cardStr)
	}
	hand.Bid = parseNumber(matches[2])
	hand.combo = getCardsCombination(hand.cards)
	return hand
}

func parseInput(input io.Reader) []Hand {
	scanner := bufio.NewScanner(input)
	var hands []Hand
	for scanner.Scan() {
		line := scanner.Text()
		hands = append(hands, parseHand(line))
	}
	return hands
}

func getResult(input io.Reader) int64 {
	hands := parseInput(input)
	sort.SliceStable(hands, func(i, j int) bool {
		if hands[i].Score() == hands[j].Score() {
			log.Fatalf("Two hands with the same score: %s / %s", hands[i].String(), hands[j].String())
		}
		return hands[i].Score() < hands[j].Score()
	})
	var results int64
	for i, hand := range hands {
		//log.Printf("Hand %d: %d / %s", i+1, hand.Score(), hand.String())
		results += hand.Bid * int64(i+1)
	}
	return results
}

func main() {
	start := time.Now()
	inputFile, errOpeningFile := os.Open("./input.txt")
	if errOpeningFile != nil {
		log.Fatalf("Unable to open input file: %v", errOpeningFile)
	}
	defer inputFile.Close()

	result := getResult(inputFile)

	log.Printf("Final result: %d", result)
	log.Printf("Execution took %s", time.Since(start))
}
