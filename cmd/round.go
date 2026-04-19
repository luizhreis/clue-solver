package cmd

import (
	"bufio"
	"fmt"

	"github.com/luizhreis/clue-solver/solver"
)

// Round collects all information about a single suggestion from the user
// and returns a processed Suggestion ready to be applied to the game.
func Round(reader *bufio.Reader, game *solver.Game, allCards []solver.Card, players []*solver.Player) *solver.Suggestion {
	fmt.Println()
	fmt.Println("--- New Round ---")

	suspect := collectCard(reader, allCards, solver.Suspect, "Which suspect was suggested?")
	location := collectCard(reader, allCards, solver.Location, "Which location was suggested?")
	weapon := collectCard(reader, allCards, solver.Weapon, "Which weapon was suggested?")

	fmt.Println()
	fmt.Printf("Suggestion: %s, %s, %s\n", suspect.Name, location.Name, weapon.Name)

	guesser := collectOptionalPlayer(reader, players, "Who made this suggestion? (0 = you)")

	refuter, shownCard := collectRefutation(reader, allCards, players, suspect, location, weapon)

	suggestion := solver.NewSuggestion(suspect, location, weapon, guesser, refuter, shownCard)

	printRoundSummary(suggestion)
	if !confirmRound(reader) {
		fmt.Println("Round discarded. Please re-enter.")
		return Round(reader, game, allCards, players)
	}

	return suggestion
}

// collectCard presents all cards of a given category and returns
// the one chosen by the user.
func collectCard(reader *bufio.Reader, allCards []solver.Card, category solver.Category, prompt string) solver.Card {
	filtered := cardsByCategory(allCards, category)

	fmt.Println()
	fmt.Println(prompt)
	for i, card := range filtered {
		fmt.Printf("  %2d. %s\n", i+1, card.Name)
	}

	fmt.Print("Choice: ")
	n := readInt(reader, 1, len(filtered))
	return filtered[n-1]
}

// collectOptionalPlayer presents the player list and returns the chosen
// player, or nil if the user selects 0 (meaning themselves).
func collectOptionalPlayer(reader *bufio.Reader, players []*solver.Player, prompt string) *solver.Player {
	fmt.Println()
	fmt.Println(prompt)
	fmt.Println("  0. Me")
	PrintPlayers(players)

	fmt.Print("Choice: ")
	n := readInt(reader, 0, len(players))
	if n == 0 {
		return nil
	}
	return PlayerByIndex(players, n)
}

// collectRefutation asks whether anyone refuted the suggestion,
// and if so, whether the card was shown to us directly.
func collectRefutation(
	reader *bufio.Reader,
	allCards []solver.Card,
	players []*solver.Player,
	suspect, location, weapon solver.Card,
) (*solver.Player, *solver.Card) {
	fmt.Println()
	fmt.Println("Did anyone refute the suggestion?")
	fmt.Println("  0. No — nobody refuted")
	PrintPlayers(players)

	fmt.Print("Choice: ")
	n := readInt(reader, 0, len(players))
	if n == 0 {
		return nil, nil
	}

	refuter := PlayerByIndex(players, n)
	shownCard := collectShownCard(reader, allCards, suspect, location, weapon)
	return refuter, shownCard
}

// collectShownCard asks whether the refuter showed us a card directly.
// If yes, the user selects which of the three suggested cards was shown.
func collectShownCard(
	reader *bufio.Reader,
	allCards []solver.Card,
	suspect, location, weapon solver.Card,
) *solver.Card {
	fmt.Println()
	fmt.Println("Did the refuter show YOU the card?")
	fmt.Println("  1. Yes")
	fmt.Println("  2. No")

	fmt.Print("Choice: ")
	n := readInt(reader, 1, 2)
	if n == 2 {
		return nil
	}

	suggested := []solver.Card{suspect, location, weapon}
	fmt.Println()
	fmt.Println("Which card was shown to you?")
	for i, card := range suggested {
		fmt.Printf("  %d. %s\n", i+1, card.Name)
	}

	fmt.Print("Choice: ")
	m := readInt(reader, 1, len(suggested))
	chosen := suggested[m-1]
	return &chosen
}

// printRoundSummary displays a summary of the round before confirmation.
func printRoundSummary(s *solver.Suggestion) {
	fmt.Println()
	fmt.Println("=== Round Summary ===")
	fmt.Printf("  Suggestion : %s, %s, %s\n", s.Suspect.Name, s.Location.Name, s.Weapon.Name)

	if s.Guesser != nil {
		fmt.Printf("  Guesser    : %s\n", s.Guesser.Name)
	} else {
		fmt.Println("  Guesser    : Me")
	}

	if s.Refuter == nil {
		fmt.Println("  Refuter    : Nobody")
	} else {
		fmt.Printf("  Refuter    : %s\n", s.Refuter.Name)
	}

	if s.ShownCard != nil {
		fmt.Printf("  Shown card : %s\n", s.ShownCard.Name)
	} else if s.Refuter != nil {
		fmt.Println("  Shown card : Unknown (not shown to us)")
	}
}

// confirmRound asks the user to confirm the round data before processing.
func confirmRound(reader *bufio.Reader) bool {
	fmt.Println()
	fmt.Println("Confirm this round?")
	fmt.Println("  1. Yes")
	fmt.Println("  2. No — re-enter")

	fmt.Print("Choice: ")
	return readInt(reader, 1, 2) == 1
}

// cardsByCategory filters a card list by category.
// Kept here to avoid exposing it unnecessarily in setup.go.
func cardsByCategory(cards []solver.Card, category solver.Category) []solver.Card {
	result := make([]solver.Card, 0)
	for _, card := range cards {
		if card.Category == category {
			result = append(result, card)
		}
	}
	return result
}
