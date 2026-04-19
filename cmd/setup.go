package cmd

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/luizhreis/clue-solver/solver"
)

// defaultCards returns the standard Clue/Cluedo card set.
func defaultCards() []solver.Card {
	suspects := []string{
		"Miss Scarlett", "Col. Mustard", "Mrs. White",
		"Mr. Green", "Mrs. Peacock", "Prof. Plum",
	}
	locations := []string{
		"Kitchen", "Ballroom", "Conservatory",
		"Billiard Room", "Library", "Study",
		"Hall", "Lounge", "Dining Room",
	}
	weapons := []string{
		"Candlestick", "Knife", "Lead Pipe",
		"Revolver", "Rope", "Wrench",
	}

	cards := make([]solver.Card, 0, len(suspects)+len(locations)+len(weapons))
	for _, name := range suspects {
		cards = append(cards, solver.Card{Name: name, Category: solver.Suspect})
	}
	for _, name := range locations {
		cards = append(cards, solver.Card{Name: name, Category: solver.Location})
	}
	for _, name := range weapons {
		cards = append(cards, solver.Card{Name: name, Category: solver.Weapon})
	}
	return cards
}

// Setup guides the user through the initial game configuration and
// returns a ready-to-use Game and the full card list.
func Setup(reader *bufio.Reader) (*solver.Game, []solver.Card, []*solver.Player) {
	fmt.Println("=== Clue Solver ===")
	fmt.Println()

	allCards := defaultCards()

	// --- Players ---
	players := setupPlayers(reader)

	// --- Our hand ---
	myHand := setupHand(reader, allCards)

	game := solver.NewGame(allCards, myHand, players)
	return game, allCards, players
}

// setupPlayers collects player names and hand sizes.
func setupPlayers(reader *bufio.Reader) []*solver.Player {
	fmt.Print("How many other players are in the game? ")
	count := readInt(reader, 1, 5)

	players := make([]*solver.Player, 0, count)
	for i := 0; i < count; i++ {
		fmt.Printf("Name of player %d: ", i+1)
		name := readLine(reader)

		fmt.Printf("How many cards does %s hold? ", name)
		handSize := readInt(reader, 1, 10)

		players = append(players, solver.NewPlayer(name, handSize))
	}
	return players
}

// setupHand presents the full card list and collects which cards
// are in the user's hand.
func setupHand(reader *bufio.Reader, allCards []solver.Card) []solver.Card {
	fmt.Println()
	fmt.Println("Here are all the cards in the game:")
	printCardList(allCards)

	fmt.Println()
	fmt.Println("Enter the numbers of the cards in YOUR hand, one at a time.")
	fmt.Println("Enter 0 when done.")

	chosen := make(map[int]bool)
	hand := make([]solver.Card, 0)

	for {
		fmt.Print("Card number (0 to finish): ")
		n := readInt(reader, 0, len(allCards))
		if n == 0 {
			if len(hand) == 0 {
				fmt.Println("You must enter at least one card.")
				continue
			}
			break
		}
		if chosen[n] {
			fmt.Println("Card already added.")
			continue
		}
		chosen[n] = true
		hand = append(hand, allCards[n-1])
		fmt.Printf("  Added: %s\n", allCards[n-1].Name)
	}

	return hand
}

// printCardList prints all cards grouped by category with their index.
func printCardList(cards []solver.Card) {
	categories := []solver.Category{solver.Suspect, solver.Location, solver.Weapon}
	idx := 1
	for _, category := range categories {
		fmt.Printf("\n[%s]\n", category)
		for _, card := range cards {
			if card.Category == category {
				fmt.Printf("  %2d. %s\n", idx, card.Name)
				idx++
			}
		}
	}
}

// readLine reads a trimmed line of input from the reader.
func readLine(reader *bufio.Reader) string {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Please try again.")
			continue
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			fmt.Print("Input cannot be empty. Try again: ")
			continue
		}
		return trimmed
	}
}

// readInt reads an integer from the reader within [min, max].
func readInt(reader *bufio.Reader, min, max int) int {
	for {
		line := strings.TrimSpace(readLineRaw(reader))
		n, err := strconv.Atoi(line)
		if err != nil || n < min || n > max {
			fmt.Printf("Please enter a number between %d and %d: ", min, max)
			continue
		}
		return n
	}
}

// readLineRaw reads a raw line without empty validation,
// used internally by readInt.
func readLineRaw(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// PrintPlayers prints a numbered list of players.
func PrintPlayers(players []*solver.Player) {
	for i, p := range players {
		fmt.Printf("  %d. %s\n", i+1, p.Name)
	}
}

// CardByIndex returns the card at 1-based index from the list.
func CardByIndex(cards []solver.Card, n int) solver.Card {
	return cards[n-1]
}

// PlayerByIndex returns the player at 1-based index from the list.
func PlayerByIndex(players []*solver.Player, n int) *solver.Player {
	return players[n-1]
}

// FormatCards formats a slice of cards as a readable string.
func FormatCards(cards []solver.Card) string {
	names := make([]string, len(cards))
	for i, c := range cards {
		names[i] = c.Name
	}
	return strings.Join(names, ", ")
}
