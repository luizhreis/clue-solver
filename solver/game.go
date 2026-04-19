package solver

import "fmt"

// Game holds the complete state of a Clue solver session.
// It is the single source of truth for card states, players,
// and the emerging solution.
type Game struct {
	Cards   map[Card]CardState
	Players []*Player
	MyHand  []Card
	Solution map[Category]Card
}

// NewGame creates a new Game, initializing all provided cards as Unknown
// and marking the local player's hand as Innocent.
func NewGame(allCards []Card, myHand []Card, players []*Player) *Game {
	g := &Game{
		Cards:    make(map[Card]CardState, len(allCards)),
		Players:  players,
		MyHand:   myHand,
		Solution: make(map[Category]Card),
	}

	for _, card := range allCards {
		g.Cards[card] = Unknown
	}

	for _, card := range myHand {
		g.SetInnocent(card)
	}

	return g
}

// SetInnocent marks a card as Innocent and propagates the information
// to all players' constraint sets.
// It is a no-op if the card is already Innocent.
func (g *Game) SetInnocent(card Card) {
	if g.Cards[card] == Innocent {
		return
	}
	g.Cards[card] = Innocent
	for _, p := range g.Players {
		p.RemoveFromConstraints(card)
	}
}

// SetGuilty marks a card as the solution for its category.
// It also marks all other cards in the same category as Innocent.
func (g *Game) SetGuilty(card Card) {
	if g.Cards[card] == Guilty {
		return
	}
	g.Cards[card] = Guilty
	g.Solution[card.Category] = card

	for c := range g.Cards {
		if c.Category == card.Category && c != card {
			g.SetInnocent(c)
		}
	}
}

// IsSolved returns true when all three categories have a confirmed solution.
func (g *Game) IsSolved() bool {
	return len(g.Solution) == 3
}

// StateOf returns the current CardState of a given card.
func (g *Game) StateOf(card Card) CardState {
	return g.Cards[card]
}

// CardsByCategory returns all cards belonging to a given category.
func (g *Game) CardsByCategory(category Category) []Card {
	result := make([]Card, 0)
	for card := range g.Cards {
		if card.Category == category {
			result = append(result, card)
		}
	}
	return result
}

// UnknownCards returns all cards whose state is still Unknown.
func (g *Game) UnknownCards() []Card {
	result := make([]Card, 0)
	for card, state := range g.Cards {
		if state == Unknown {
			result = append(result, card)
		}
	}
	return result
}

// PlayerByName returns a player by name, or an error if not found.
func (g *Game) PlayerByName(name string) (*Player, error) {
	for _, p := range g.Players {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("player %q not found", name)
}

// Summary returns a human-readable string of the current game state.
func (g *Game) Summary() string {
	summary := "=== Game State ===\n"
	for _, category := range []Category{Suspect, Location, Weapon} {
		summary += fmt.Sprintf("\n[%s]\n", category)
		for card, state := range g.Cards {
			if card.Category == category {
				summary += fmt.Sprintf("  %-20s %s\n", card.Name, state)
			}
		}
	}
	if g.IsSolved() {
		summary += "\n=== SOLUTION ===\n"
		for _, category := range []Category{Suspect, Location, Weapon} {
			summary += fmt.Sprintf("  %s: %s\n", category, g.Solution[category].Name)
		}
	}
	return summary
}