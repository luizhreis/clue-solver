package cmd

import (
	"fmt"
	"strings"

	"github.com/luizhreis/clue-solver/solver"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// DisplayState prints the current known state of all cards,
// grouped by category, with color-coded states.
func DisplayState(game *solver.Game, allCards []solver.Card) {
	fmt.Println()
	fmt.Println(bold("=== Current Game State ==="))

	for _, category := range []solver.Category{solver.Suspect, solver.Location, solver.Weapon} {
		fmt.Println()
		fmt.Printf("%s[%s]%s\n", colorCyan, category, colorReset)

		for _, card := range allCards {
			if card.Category != category {
				continue
			}
			state := game.StateOf(card)
			fmt.Printf("  %-22s %s\n", card.Name, coloredState(state))
		}
	}
	fmt.Println()
}

// DisplaySolution prints the confirmed solution in a highlighted block.
func DisplaySolution(game *solver.Game) {
	if !game.IsSolved() {
		return
	}

	fmt.Println()
	fmt.Println(bold("╔══════════════════════════════╗"))
	fmt.Println(bold("║        MYSTERY SOLVED!       ║"))
	fmt.Println(bold("╚══════════════════════════════╝"))
	fmt.Println()

	solution := game.Solution
	for _, category := range []solver.Category{solver.Suspect, solver.Location, solver.Weapon} {
		card := solution[category]
		fmt.Printf("  %s%-10s%s %s%s%s\n",
			colorYellow, category.String()+":", colorReset,
			colorGreen, card.Name, colorReset,
		)
	}
	fmt.Println()
}

// DisplayUnknowns prints all cards whose state is still Unknown.
// Useful at the end of a round to show what remains to be deduced.
func DisplayUnknowns(game *solver.Game, allCards []solver.Card) {
	unknowns := make([]string, 0)
	for _, card := range allCards {
		if game.StateOf(card) == solver.Unknown {
			unknowns = append(unknowns, card.Name)
		}
	}

	if len(unknowns) == 0 {
		return
	}

	fmt.Printf("%sStill unknown:%s %s\n",
		colorYellow, colorReset,
		strings.Join(unknowns, ", "),
	)
}

// DisplayRoundResult prints a brief summary of what was deduced
// after processing a suggestion.
func DisplayRoundResult(before, after map[solver.Card]solver.CardState) {
	newInnocent := make([]string, 0)
	newGuilty := make([]string, 0)

	for card, stateBefore := range before {
		stateAfter := after[card]
		if stateBefore == solver.Unknown && stateAfter == solver.Innocent {
			newInnocent = append(newInnocent, card.Name)
		}
		if stateBefore == solver.Unknown && stateAfter == solver.Guilty {
			newGuilty = append(newGuilty, card.Name)
		}
	}

	if len(newInnocent) == 0 && len(newGuilty) == 0 {
		fmt.Printf("%sNo new deductions this round.%s\n", colorYellow, colorReset)
		return
	}

	if len(newInnocent) > 0 {
		fmt.Printf("%s✓ Proven innocent:%s %s\n",
			colorGreen, colorReset,
			strings.Join(newInnocent, ", "),
		)
	}
	if len(newGuilty) > 0 {
		fmt.Printf("%s★ Confirmed guilty:%s %s\n",
			colorRed, colorReset,
			strings.Join(newGuilty, ", "),
		)
	}
}

// SnapshotStates captures the current CardState of all cards.
// Used to compare before/after a round for DisplayRoundResult.
func SnapshotStates(game *solver.Game, allCards []solver.Card) map[solver.Card]solver.CardState {
	snapshot := make(map[solver.Card]solver.CardState, len(allCards))
	for _, card := range allCards {
		snapshot[card] = game.StateOf(card)
	}
	return snapshot
}

// coloredState returns a color-coded string for a CardState.
func coloredState(state solver.CardState) string {
	switch state {
	case solver.Innocent:
		return colorGreen + "Innocent" + colorReset
	case solver.Guilty:
		return colorRed + colorBold + "GUILTY" + colorReset
	default:
		return colorYellow + "Unknown" + colorReset
	}
}

// bold wraps a string in the bold ANSI code.
func bold(s string) string {
	return colorBold + s + colorReset
}