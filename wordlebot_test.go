package wordlebot

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestGetWordList(t *testing.T) {
	f, err := os.Create("testFile")
	if err != nil {
		log.Fatalf("Failed to create testFile: %e\n", err)
	}
	_, err = f.Write([]byte("apple\ngrape\nfruit"))
	if err != nil {
		log.Fatalf("Failed to write to testFile: %e\n", err)
	}
	got := GetWordList("testFile")
	want := []Word{{letters: [5]rune{'a', 'p', 'p', 'l', 'e'}}, {letters: [5]rune{'g', 'r', 'a', 'p', 'e'}}, {letters: [5]rune{'f', 'r', 'u', 'i', 't'}}}
	for i, gotWord := range got {
		if gotWord != want[i] {
			t.Fatalf("Got % v, want % v\n", gotWord, want[i])
		}
	}
}

func TestCreatePatterns(t *testing.T) {
	allPatterns := CreatePatterns()
	fmt.Println("Number of patterns: ", len(allPatterns))
	m := make(map[Pattern]int)
	for i, p := range allPatterns {
		if _, patternExists := m[p]; patternExists {
			t.Fatalf("Duplicate Pattern!: % v\n", p)
		} else {
			m[p] = i
		}
	}
}

func TestCompareWord(t *testing.T) {
	examples := []struct {
		name   string
		guess  Word
		answer Word
		want   Pattern
	}{
		{
			name:   "yelps/yells",
			guess:  Word{Name: "yelps", letters: [5]rune{'y', 'e', 'l', 'p', 's'}},
			answer: Word{Name: "yells", letters: [5]rune{'y', 'e', 'l', 'l', 's'}},
			want:   Pattern{2, 2, 2, 0, 2},
		},
		{
			name:   "wants/wants",
			guess:  Word{Name: "wants", letters: [5]rune{'w', 'a', 'n', 't', 's'}},
			answer: Word{Name: "wants", letters: [5]rune{'w', 'a', 'n', 't', 's'}},
			want:   Pattern{2, 2, 2, 2, 2},
		},
		{
			name:   "sheet/grain",
			guess:  Word{Name: "sheet", letters: [5]rune{'s', 'h', 'e', 'e', 't'}},
			answer: Word{Name: "grain", letters: [5]rune{'g', 'r', 'a', 'i', 'n'}},
			want:   Pattern{0, 0, 0, 0, 0},
		},
		{
			name:   "babes/abbey",
			guess:  Word{Name: "babes", letters: [5]rune{'b', 'a', 'b', 'e', 's'}},
			answer: Word{Name: "abbey", letters: [5]rune{'a', 'b', 'b', 'e', 'y'}},
			want:   Pattern{1, 1, 2, 2, 0},
		},
		{
			name:   "kebab/abbey",
			guess:  Word{Name: "kebab", letters: [5]rune{'k', 'e', 'b', 'a', 'b'}},
			answer: Word{Name: "abbey", letters: [5]rune{'a', 'b', 'b', 'e', 'y'}},
			want:   Pattern{0, 1, 2, 1, 1},
		},
		{
			name:   "speed/abide",
			guess:  Word{Name: "speed", letters: [5]rune{'s', 'p', 'e', 'e', 'd'}},
			answer: Word{Name: "abide", letters: [5]rune{'a', 'b', 'i', 'd', 'e'}},
			want:   Pattern{0, 0, 1, 0, 1},
		},
		{
			name:   "speed/erase",
			guess:  Word{Name: "speed", letters: [5]rune{'s', 'p', 'e', 'e', 'd'}},
			answer: Word{Name: "erase", letters: [5]rune{'e', 'r', 'a', 's', 'e'}},
			want:   Pattern{1, 0, 1, 1, 0},
		},
		{
			name:   "speed/steal",
			guess:  Word{Name: "speed", letters: [5]rune{'s', 'p', 'e', 'e', 'd'}},
			answer: Word{Name: "steal", letters: [5]rune{'s', 't', 'e', 'a', 'l'}},
			want:   Pattern{2, 0, 2, 0, 0},
		},
	}
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			if got := compareWord(ex.guess, ex.answer); got != ex.want {
				log.Fatalf("%s: Patterns don't match: got %+v, want %+v", ex.name, got, ex.want)
			}
		})
	}

}

func TestPruneWordList(t *testing.T) {
	responsePattern := Pattern{0, 1, 0, 1, 0}
	wordList := []Word{
		{Name: "weird", letters: [5]rune{'w', 'e', 'i', 'r', 'd'}},
		{Name: "sheet", letters: [5]rune{'s', 'h', 'e', 'e', 't'}},
		{Name: "grape", letters: [5]rune{'g', 'r', 'a', 'p', 'e'}},
		{Name: "sheer", letters: [5]rune{'s', 'h', 'e', 'e', 'r'}},
		{Name: "abbey", letters: [5]rune{'a', 'b', 'b', 'e', 'y'}},
		{Name: "kebab", letters: [5]rune{'k', 'e', 'b', 'a', 'b'}},
	}
	guessIndex := 0
	newWordList := PruneWordList(responsePattern, wordList, guessIndex)
	fmt.Printf("PruneWordList result: %+v\n", newWordList)
	want := []Word{
		{Name: "grape", letters: [5]rune{'g', 'r', 'a', 'p', 'e'}},
		{Name: "sheer", letters: [5]rune{'s', 'h', 'e', 'e', 'r'}},
	}
	for i, got := range newWordList {
		if got != want[i] {
			log.Fatalf("Failed to prune word list: got %+v, want %+v\n", got, want[i])
		}
	}
}

/*
func TestCountMatches(t *testing.T) {
	examples := []struct {
		name     string
		wordList []Word
	}{
		{
			name: "small list",
			wordList: []Word{{Name: "weird", letters: [5]rune{'w', 'e', 'i', 'r', 'd'}},
				{Name: "sheet", letters: [5]rune{'s', 'h', 'e', 'e', 't'}},
				{Name: "grape", letters: [5]rune{'g', 'r', 'a', 'p', 'e'}},
				{Name: "sleet", letters: [5]rune{'s', 'l', 'e', 'e', 't'}},
				{Name: "abbey", letters: [5]rune{'a', 'b', 'b', 'e', 'y'}},
				{Name: "kebab", letters: [5]rune{'k', 'e', 'b', 'a', 'b'}},
			},
		},
	}
	allPatterns := CreatePatterns()
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			allWordsPatternCounts := CountMatches(ex.wordList, allPatterns)
			for i, patternCounts := range allWordsPatternCounts {
				fmt.Printf("word: %s\n", ex.wordList[i].Name)
				fmt.Printf("	patterns: ")
				for _, patternCount := range patternCounts {
					if patternCount.count != 0 {
						fmt.Printf("%+v\n", patternCount)
					}
				}
			}
		})
	}
}
*/

/*
func TestCalcEntropy(t *testing.T) {
	examples := []struct {
		name     string
		wordList []Word
	}{
		{
			name: "small list",
			wordList: []Word{{Name: "weird", letters: [5]rune{'w', 'e', 'i', 'r', 'd'}},
				{Name: "sheet", letters: [5]rune{'s', 'h', 'e', 'e', 't'}},
				{Name: "grape", letters: [5]rune{'g', 'r', 'a', 'p', 'e'}},
				{Name: "sleet", letters: [5]rune{'s', 'l', 'e', 'e', 't'}},
				{Name: "abbey", letters: [5]rune{'a', 'b', 'b', 'e', 'y'}},
				{Name: "kebab", letters: [5]rune{'k', 'e', 'b', 'a', 'b'}},
			},
		},
	}
	allPatterns := CreatePatterns()
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			allWordsPatternCounts := CountMatches(ex.wordList, allPatterns)
			entropies := CalcEntropies(allWordsPatternCounts)
			for i := 0; i < len(entropies); i++ {
				fmt.Printf("word: %s ; entropy: %+v\n", ex.wordList[i].Name, entropies[i])
			}
		})
	}
}
*/

/*
func TestMakeGuess(t *testing.T) {
	examples := []struct {
		name      string
		entropies []float64
		want      int
	}{
		{
			name:      "ascending",
			entropies: []float64{1.34, 3.25, 4.19, 6.57},
			want:      3,
		},
		{
			name:      "descending",
			entropies: []float64{6.57, 4.19, 3.86, 2.41},
			want:      0,
		},
		{
			name:      "random",
			entropies: []float64{4.19, 3.86, 5.96, 1.42, 1.23, 2.25},
			want:      2,
		},
	}
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			if got := MakeGuess(ex.entropies); got != ex.want {
				fmt.Printf("NAME: %s\n", ex.name)
				log.Fatalf("Made wrong guess: got %v, want %v\n", got, ex.want)
			}

		})
	}
}
*/
