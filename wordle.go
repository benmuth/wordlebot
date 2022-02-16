package main

import (
	"Users/ben/Documents/WordleBot/wordlebot"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Playing Wordle!")
	wordList := wordlebot.GetWordList("/Users/ben/Documents/WordleBot/all_words.txt")
	allPatterns := wordlebot.CreatePatterns()

	for len(wordList) > 2 {
		allWordsPatternCounts := wordlebot.CountMatches(wordList, allPatterns)
		entropies := wordlebot.CalcEntropies(allWordsPatternCounts)
		guessIndex := wordlebot.MakeGuess(entropies)
		fmt.Printf("YOU SHOULD GUESS: %s\n", wordList[guessIndex].Name)
		fmt.Println("Input the feedback pattern as a continuous string of 5 numbers; 0s for greys, 1s for yellows, and 2s for greens:")
		var feedback string
		fmt.Scanln(&feedback)

		chars := strings.Split(feedback, "")
		var responsePattern wordlebot.Pattern
		for i, c := range chars {
			n, err := strconv.ParseInt(c, 10, 8)
			if err != nil {
				log.Fatalf("Couldn't convert user input %s to int: %e", feedback, err)
			}
			responsePattern[i] = int8(n)
		}

		wordList = wordlebot.PruneWordList(responsePattern, wordList, guessIndex, allWordsPatternCounts)
	}

	if len(wordList) == 0 {
		fmt.Printf("Whoops! No matches! D:\n")
	} else if len(wordList) == 1 {
		fmt.Printf("The word is probably: \n%s\n", wordList[0].Name)
	} else {
		fmt.Printf("The word is probably one of: \n")
		for _, w := range wordList {
			fmt.Printf("%s\n", w.Name)
		}
	}
}
