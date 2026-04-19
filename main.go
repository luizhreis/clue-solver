package main

import (
	"fmt"

	"github.com/luizhreis/clue-solver/solver"
)

func main() {
	// --- Card definitions ---

	// Suspects
	scarlett := solver.Card{Name: "Miss Scarlett", Category: solver.Suspect}
	mustard := solver.Card{Name: "Col. Mustard", Category: solver.Suspect}
	white := solver.Card{Name: "Mrs. White", Category: solver.Suspect}
	green := solver.Card{Name: "Mr. Green", Category: solver.Suspect}
	peacock := solver.Card{Name: "Mrs. Peacock", Category: solver.Suspect}
	plum := solver.Card{Name: "Prof. Plum", Category: solver.Suspect}

	// Locations
	kitchen := solver.Card{Name: "Kitchen", Category: solver.Location}
	ballroom := solver.Card{Name: "Ballroom", Category: solver.Location}
	library := solver.Card{Name: "Library", Category: solver.Location}
	study := solver.Card{Name: "Study", Category: solver.Location}
	hall := solver.Card{Name: "Hall", Category: solver.Location}
	lounge := solver.Card{Name: "Lounge", Category: solver.Location}

	// Weapons
	candlestick := solver.Card{Name: "Candlestick", Category: solver.Weapon}
	knife := solver.Card{Name: "Knife", Category: solver.Weapon}
	pipe := solver.Card{Name: "Lead Pipe", Category: solver.Weapon}
	revolver := solver.Card{Name: "Revolver", Category: solver.Weapon}
	rope := solver.Card{Name: "Rope", Category: solver.Weapon}
	wrench := solver.Card{Name: "Wrench", Category: solver.Weapon}

	allCards := []solver.Card{
		scarlett, mustard, white, green, peacock, plum,
		kitchen, ballroom, library, study, hall, lounge,
		candlestick, knife, pipe, revolver, rope, wrench,
	}

	// --- Players ---
	// 4 players total, 18 cards distributed:
	// - us: 5 cards
	// - each other player: ~4 cards (simplified)
	alice := solver.NewPlayer("Alice", 4)
	bob := solver.NewPlayer("Bob", 4)
	carol := solver.NewPlayer("Carol", 4)

	players := []*solver.Player{alice, bob, carol}

	// --- Our hand ---
	myHand := []solver.Card{scarlett, kitchen, candlestick, mustard, ballroom}

	// --- Game setup ---
	game := solver.NewGame(allCards, myHand, players)

	fmt.Println("Initial state after dealing our hand:")
	fmt.Println(game.Summary())

	// --- Inference engine ---
	engine := solver.NewEngine()
	hypothesis := solver.NewHypothesisEngine()

	// --- Round 1 ---
	// Alice suggests (White, Library, Knife).
	// Bob refutes, but doesn't show us the card.
	s1 := solver.NewSuggestion(white, library, knife, alice, bob, nil)
	s1.Process(game)
	engine.Run(game)
	fmt.Println("After round 1 (Bob refuted Alice, card unknown):")
	fmt.Println(game.Summary())

	// --- Round 2 ---
	// Bob suggests (Green, Study, Rope).
	// Carol refutes and shows us: Rope.
	shownRope := rope
	s2 := solver.NewSuggestion(green, study, rope, bob, carol, &shownRope)
	s2.Process(game)
	engine.Run(game)
	fmt.Println("After round 2 (Carol showed us Rope):")
	fmt.Println(game.Summary())

	// --- Round 3 ---
	// We suggest (Peacock, Hall, Pipe).
	// Nobody refutes → solution found for these three cards.
	s3 := solver.NewSuggestion(peacock, hall, pipe, nil, nil, nil)
	s3.Process(game)
	engine.Run(game)
	fmt.Println("After round 3 (nobody refuted our suggestion):")
	fmt.Println(game.Summary())

	// --- Round 4 ---
	// Alice suggests (Plum, Lounge, Revolver).
	// Nobody refutes.
	s4 := solver.NewSuggestion(plum, lounge, revolver, alice, nil, nil)
	s4.Process(game)
	engine.Run(game)

	// Run hypothesis engine in case direct deduction got stuck.
	hypothesis.Run(game)
	engine.Run(game)

	fmt.Println("Final state:")
	fmt.Println(game.Summary())

	if game.IsSolved() {
		fmt.Println("=== The mystery is solved! ===")
	} else {
		fmt.Println("Not enough information yet. Keep playing!")
		fmt.Printf("Remaining unknown cards: %v\n", game.UnknownCards())
	}
}
