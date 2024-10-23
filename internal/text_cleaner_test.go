package internal

import "testing"

func TestCleanText(t *testing.T) {
	testCases := []string{
		"I had something interesting for breakfast",
		"I hear Mastodon is better than Chirpy. sharbert I need to migrate",
		"I really need a kerfuffle to go to bed sooner, Fornax !",
	}
	expectations := []string{
		"I had something interesting for breakfast",
		"I hear Mastodon is better than Chirpy. **** I need to migrate",
		"I really need a **** to go to bed sooner, **** !",
	}
	for i, testCase := range testCases {
		expected := expectations[i]
		got := CleanText(testCase)
		if got != expected {
			t.Fatalf("expected '%s' but was '%s'", expected, got)
		}
	}
}
