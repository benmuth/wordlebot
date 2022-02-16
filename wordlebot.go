package wordlebot

import (
	"log"
	"math"
	"os"
	"strings"
)

//GetWordList reads the named file and returns each newline separated string as a slice of Words
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

//Creates all possible patterns of colors that can be given as feedback for a 5 letter word
func CreatePatterns() []Pattern {
	//https://rosettacode.org/wiki/Permutations_with_repetitions#Go
	//Generate all permutations of colors with repetition (3^5 possibilities)
	allPatterns := make([]Pattern, 0, 243)
	n := 5
	values := []int8{0, 1, 2}
	k := len(values)

	pn := make([]int, n)
	var p Pattern
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
			if i == n { // all permutations generated
				return allPatterns
			}
		}
	}
}

//Given two words, return a pattern of colors according to Wordle rules
func compareWord(guess Word, answer Word) Pattern {
	letterCounts := make(map[rune]int)
	for _, c := range answer.letters { //count letters to deal with multi-letter cases
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

//MakeGuess returns the index of the word that gives the most information on average
func MakeGuess(wordList []Word, allPatterns []Pattern) int {
	//Calculate entropies
	matchCounts := make(map[Pattern]int, 243)

	entropies := make([]float64, len(wordList))
	for i, w := range wordList {
		var entropy float64
		//Go through every possible answer word and count how many times a pattern is matched
		for _, x := range wordList {
			p := compareWord(w, x)
			matchCounts[p]++
		}

		for _, count := range matchCounts {
			//H(X) = SUM(P(X) * log2(1/(P(X))))
			probability := float64(count) / float64(len(wordList))
			entropy += probability * math.Log2(float64(1/probability))
		}
		entropies[i] = entropy
	}

	//Make guess by choosing index of word with highest entropy
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

func PruneWordList(responsePattern Pattern, wordList []Word, guessIndex int) []Word {
	guess := wordList[guessIndex]
	newWordList := make([]Word, 0) //pass the pattern match count to be the capacity for the []Word to reduce allocations
	for _, possibleAnswer := range wordList {
		p := compareWord(guess, possibleAnswer)
		if p == responsePattern {
			newWordList = append(newWordList, possibleAnswer)
		}
	}
	return newWordList
}
