package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/luizhreis/clue-solver/solver"
)

// Run is the entry point for the CLI. It sets up the game and runs
// the main round loop until the mystery is solved or the user quits.
func Run() {
	reader := bufio.NewReader(os.Stdin)

	game, allCards, players := Setup(reader)

	fmt.Println()
	fmt.Println("Game configured. Let's solve the mystery!")
	DisplayState(game, allCards)

	engine    := solver.NewEngine()
	hypothesis := solver.NewHypothesisEngine()

	for {
		if game.IsSolved() {
			DisplaySolution(game)
			break
		}

		suggestion := Round(reader, game, allCards, players)

		snapshot := SnapshotStates(game, allCards)

		suggestion.Process(game)
		engine.Run(game)

		if !anyNewDeductions(snapshot, game, allCards) {
			hypothesis.Run(game)
			engine.Run(game)
		}

		fmt.Println()
		DisplayRoundResult(snapshot, currentStates(game, allCards))
		DisplayState(game, allCards)
		DisplayUnknowns(game, allCards)

		if game.IsSolved() {
			DisplaySolution(game)
			break
		}

		if !continueGame(reader) {
			fmt.Println("Goodbye. The mystery remains unsolved.")
			break
		}
	}
}

// anyNewDeductions returns true if any card changed state after a round.
func anyNewDeductions(before map[solver.Card]solver.CardState, game *solver.Game, allCards []solver.Card) bool {
	for _, card := range allCards {
		if before[card] != game.StateOf(card) {
			return true
		}
	}
	return false
}

// currentStates returns the current state of all cards as a map.
func currentStates(game *solver.Game, allCards []solver.Card) map[solver.Card]solver.CardState {
	return SnapshotStates(game, allCards)
}

// continueGame asks the user whether to proceed to the next round.
func continueGame(reader *bufio.Reader) bool {
	fmt.Println()
	fmt.Println("What would you like to do?")
	fmt.Println("  1. Enter next round")
	fmt.Println("  2. Quit")

	fmt.Print("Choice: ")
	return readInt(reader, 1, 2) == 1
}