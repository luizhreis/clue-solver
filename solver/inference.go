package solver

// InferenceRule is the interface that every deduction rule must implement.
// Apply attempts to derive new information from the current game state.
// It returns true if any change was made, signaling the engine to run again.
type InferenceRule interface {
	Apply(g *Game) bool
}

// Engine runs all inference rules in a loop until no new deductions
// can be made (fixed point). Rules are applied in order on every pass.
type Engine struct {
	Rules []InferenceRule
}

// NewEngine creates an Engine with the default set of inference rules
// in the recommended application order.
func NewEngine() *Engine {
	return &Engine{
		Rules: []InferenceRule{
			&RuleDirectElimination{},
			&RuleSetCollapse{},
			&RuleCategoryDeduction{},
			&RuleCrossDeduction{},
			&RuleFullHand{},
		},
	}
}

// Run applies all rules repeatedly until a full pass produces no changes.
func (e *Engine) Run(g *Game) {
	for {
		changed := false
		for _, rule := range e.Rules {
			if rule.Apply(g) {
				changed = true
			}
		}
		if !changed {
			break
		}
	}
}

// RuleDirectElimination removes Innocent cards from all pending
// constraint sets across all players.
//
// Example: card X is proven Innocent elsewhere → remove X from every
// ConstraintSet that still lists it as a candidate.
type RuleDirectElimination struct{}

func (r *RuleDirectElimination) Apply(g *Game) bool {
	changed := false
	for card, state := range g.Cards {
		if state == Innocent {
			for _, p := range g.Players {
				for _, cs := range p.Constraints {
					if cs.Cards[card] {
						cs.Remove(card)
						changed = true
					}
				}
			}
		}
	}
	return changed
}

// RuleSetCollapse promotes a card to Innocent when its ConstraintSet
// has been reduced to a single candidate. That player must hold it.
//
// Example: constraint {A, B, C} → B and C proven Innocent → {A} →
// player must have A → A is Innocent.
type RuleSetCollapse struct{}

func (r *RuleSetCollapse) Apply(g *Game) bool {
	changed := false
	for _, p := range g.Players {
		for _, card := range p.ResolvedConstraints() {
			if g.StateOf(card) == Unknown {
				g.SetInnocent(card)
				p.AddToHand(card)
				changed = true
			}
		}
	}
	return changed
}

// RuleCategoryDeduction marks the last Unknown card in a category as
// Guilty when all other cards in that category are Innocent.
//
// Example: suspects = {A: Innocent, B: Innocent, C: Unknown} → C is Guilty.
type RuleCategoryDeduction struct{}

func (r *RuleCategoryDeduction) Apply(g *Game) bool {
	changed := false
	for _, category := range []Category{Suspect, Location, Weapon} {
		if _, solved := g.Solution[category]; solved {
			continue
		}
		cards := g.CardsByCategory(category)
		unknowns := make([]Card, 0)
		for _, card := range cards {
			if g.StateOf(card) == Unknown {
				unknowns = append(unknowns, card)
			}
		}
		if len(unknowns) == 1 {
			g.SetGuilty(unknowns[0])
			changed = true
		}
	}
	return changed
}

// RuleCrossDeduction ensures that once a Guilty card is confirmed,
// all remaining Unknown cards in the same category become Innocent.
//
// This is largely handled by SetGuilty already, but this rule catches
// any residual Unknown cards that may appear due to ordering of operations.
type RuleCrossDeduction struct{}

func (r *RuleCrossDeduction) Apply(g *Game) bool {
	changed := false
	for card, state := range g.Cards {
		if state == Guilty {
			for _, other := range g.CardsByCategory(card.Category) {
				if other != card && g.StateOf(other) == Unknown {
					g.SetInnocent(other)
					changed = true
				}
			}
		}
	}
	return changed
}

// RuleFullHand resolves all remaining constraints for a player whose
// entire hand is already known. If a player has no more unknown cards,
// any unresolved ConstraintSet they hold can only be satisfied by cards
// already in their confirmed hand.
//
// Example: player has HandSize=3, Hand={A, B, C} → any pending
// constraint {X, Y} where none of X or Y is in their hand is a
// contradiction (should not happen in a valid game), and constraints
// that include a hand card are resolved to that card.
type RuleFullHand struct{}

func (r *RuleFullHand) Apply(g *Game) bool {
	changed := false
	for _, p := range g.Players {
		if !p.IsHandComplete() {
			continue
		}
		for _, cs := range p.Constraints {
			if cs.IsResolved() {
				continue
			}
			// Keep only cards confirmed in the player's hand.
			for card := range cs.Cards {
				if !p.HasCard(card) {
					cs.Remove(card)
					changed = true
				}
			}
		}
	}
	return changed
}

// HypothesisEngine attempts to find contradictions by assuming a card
// is Guilty and checking if that leads to an inconsistent state.
// If a contradiction is found, the card must be Innocent.
//
// This is used when the main Engine reaches a fixed point without
// fully solving the game.
type HypothesisEngine struct {
	engine *Engine
}

// NewHypothesisEngine creates a HypothesisEngine backed by a standard Engine.
func NewHypothesisEngine() *HypothesisEngine {
	return &HypothesisEngine{engine: NewEngine()}
}

// Run iterates over all Unknown cards, hypothesizes each as Guilty,
// and checks for contradictions on a game snapshot.
// If a contradiction is found, the card is marked Innocent in the real game.
// Returns true if any new information was derived.
func (h *HypothesisEngine) Run(g *Game) bool {
	changed := false
	for _, card := range g.UnknownCards() {
		snapshot := cloneGame(g)
		snapshot.SetGuilty(card)
		h.engine.Run(snapshot)
		if hasContradiction(snapshot) {
			g.SetInnocent(card)
			changed = true
		}
	}
	return changed
}

// hasContradiction returns true if the game state is logically inconsistent.
// A contradiction occurs when a category has more than one Guilty card,
// or when a category has no remaining candidates (all Innocent, none Guilty).
func hasContradiction(g *Game) bool {
	for _, category := range []Category{Suspect, Location, Weapon} {
		guiltyCount := 0
		unknownOrGuilty := 0
		for card, state := range g.Cards {
			if card.Category != category {
				continue
			}
			if state == Guilty {
				guiltyCount++
			}
			if state != Innocent {
				unknownOrGuilty++
			}
		}
		if guiltyCount > 1 {
			return true
		}
		if unknownOrGuilty == 0 {
			return true
		}
	}
	return false
}

// cloneGame produces a deep copy of the game state for use in hypothesis
// testing. The clone shares no mutable references with the original.
func cloneGame(g *Game) *Game {
	cards := make(map[Card]CardState, len(g.Cards))
	for k, v := range g.Cards {
		cards[k] = v
	}

	solution := make(map[Category]Card, len(g.Solution))
	for k, v := range g.Solution {
		solution[k] = v
	}

	players := make([]*Player, len(g.Players))
	for i, p := range g.Players {
		players[i] = clonePlayer(p)
	}

	return &Game{
		Cards:    cards,
		Players:  players,
		MyHand:   g.MyHand,
		Solution: solution,
	}
}

// clonePlayer produces a deep copy of a Player for hypothesis testing.
func clonePlayer(p *Player) *Player {
	hand := make(map[Card]bool, len(p.Hand))
	for k, v := range p.Hand {
		hand[k] = v
	}

	constraints := make([]*ConstraintSet, len(p.Constraints))
	for i, cs := range p.Constraints {
		cloned := make(map[Card]bool, len(cs.Cards))
		for k, v := range cs.Cards {
			cloned[k] = v
		}
		constraints[i] = &ConstraintSet{Cards: cloned}
	}

	return &Player{
		Name:        p.Name,
		Hand:        hand,
		Constraints: constraints,
		HandSize:    p.HandSize,
	}
}