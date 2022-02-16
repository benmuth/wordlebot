package wordlebot

import (
	"log"
	"math"
	"os"
	"strings"
)

//GetWordList reads the named file and returns each newline separated string as a slice of words
func GetWordList(fileName string) []Word {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Failed to read %s: %e\n", fileName, err)
	}
	strings := strings.Split(string(data), "\n")
	words := make([]Word, len(strings))
	for i, s := range strings {
		words[i].Name = s
		for j, c := range s {
			words[i].letters[j] = c
		}
	}
	return words
}

type Word struct {
	Name    string
	letters [5]rune
}

type Pattern [5]int8

type patternCount struct {
	p     Pattern
	count int
}

func CreatePatterns() []Pattern {
	//https://rosettacode.org/wiki/Permutations_with_repetitions#Go
	allPatterns := make([]Pattern, 0, 243)
	n := 5
	values := []int8{0, 1, 2}
	k := len(values)

	pn := make([]int, n)
	var p Pattern
	//var j int
	for {
		// generate permutaton
		for i, x := range pn {
			p[i] = values[x]
		}
		// add permutation to collection of patterns
		allPatterns = append(allPatterns, p)

		// increment permutation number
		for i := 0; ; {
			pn[i]++
			if pn[i] < k {
				break
			}
			pn[i] = 0
			i++
			if i == n {
				return allPatterns // all permutations generated
			}
		}
	}
}

//Given two words, return a Pattern of Colors according to Wordle rules
func compareWord(guess Word, answer Word) Pattern {
	letterCounts := make(map[rune]int)
	for _, c := range answer.letters {
		letterCounts[c]++
	}
	var p Pattern
	for i, c := range guess.letters {
		if c == answer.letters[i] {
			p[i] = 2
			letterCounts[c]--
		} else if count, answerContains := letterCounts[c]; answerContains {
			if count > 0 {
				p[i] = 1
			}
			letterCounts[c]--
		}
	}
	return p
}

//For each Pattern of Colors, count how many words in a given wordlist match that Pattern
func CountMatches(wordList []Word, allPatterns []Pattern) [][]patternCount {
	matchCounts := make(map[Pattern]int, 243)
	allWordsPatternCounts := make([][]patternCount, len(wordList))

	for i, w := range wordList {
		//init matchCounts
		for _, p := range allPatterns {
			matchCounts[p] = 0
		}
		//go through every possible answer word and count how many times there is a match for each Pattern
		for _, x := range wordList {
			p := compareWord(w, x)
			matchCounts[p]++
		}
		//store the counts of each Pattern in a slice
		patternCounts := make([]patternCount, 243)
		var j int
		for p, count := range matchCounts {
			patternCounts[j].p = p
			patternCounts[j].count = count
			j++
		}
		//store each slice of Pattern counts in a slice where the index corresponds to the index of the wordList
		allWordsPatternCounts[i] = patternCounts
	}
	return allWordsPatternCounts
}

func CalcEntropies(allWordsPatternCounts [][]patternCount) []float64 {
	entropies := make([]float64, len(allWordsPatternCounts))
	for i, wordPatternCounts := range allWordsPatternCounts {
		var entropy float64
		for _, patternCount := range wordPatternCounts {
			probability := float64(patternCount.count) / float64(len(allWordsPatternCounts))
			if probability != 0 {
				entropy += probability * math.Log2(float64(1/probability))
			} else {
				entropy += 0
			}
		}
		entropies[i] = entropy
	}
	return entropies
}

func MakeGuess(entropies []float64) int {
	var guessIndex int
	var greatestEntropy float64
	for i, v := range entropies {
		if v > greatestEntropy {
			greatestEntropy = v
			guessIndex = i
		}
	}
	return guessIndex
}

func PruneWordList(responsePattern Pattern, wordList []Word, guessIndex int, allWordsPatternCounts [][]patternCount) []Word {
	guess := wordList[guessIndex]
	guessWordPatternCounts := allWordsPatternCounts[guessIndex]
	var guessWordPatternMatches int
	for _, pc := range guessWordPatternCounts {
		if pc.p == responsePattern {
			guessWordPatternMatches = pc.count
		}
	}
	newWordList := make([]Word, 0, guessWordPatternMatches)
	for _, possibleAnswer := range wordList {
		p := compareWord(guess, possibleAnswer)
		if p == responsePattern {
			newWordList = append(newWordList, possibleAnswer)
		}
	}
	return newWordList
}
