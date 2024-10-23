package internal

import "strings"

func CleanText(dirty string) string {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	inWords := strings.Split(dirty, " ")
	outWords := make([]string, 0, len(inWords))
	for _, word := range inWords {
		isReplaced := false
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				isReplaced = true
				break
			}
		}
		if isReplaced {
			outWords = append(outWords, "****")
		} else {
			outWords = append(outWords, word)
		}
	}
	return strings.Join(outWords, " ")
}
