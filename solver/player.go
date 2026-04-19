package solver

// ConstraintSet represents a set of cards where at least one
// must be in a player's hand. It is created when a player refutes
// a suggestion without revealing the card to us.
type ConstraintSet struct {
	Cards map[Card]bool
}

// NewConstraintSet creates a ConstraintSet from a slice of cards.
func NewConstraintSet(cards []Card) *ConstraintSet {
	m := make(map[Card]bool, len(cards))
	for _, c := range cards {
		m[c] = true
	}

	return &ConstraintSet{Cards: m}
}

// Remove eliminates a card from the set, since it was proven innocent
// by another means and therefore is not the card held by this player.
func (cs *ConstraintSet) Remove(card Card) {
	delete(cs.Cards, card)
}

// IsResolved returns true when only one candidate remains in the set.
func (cs *ConstraintSet) IsResolved() bool {
	return len(cs.Cards) == 1
}

// Resolved returns the last remaining card when the set is resolved.
// Returns the zero value of Card and false if not yet resolved.
func (cs *ConstraintSet) Resolved() (Card, bool) {
	if !cs.IsResolved() {
		return Card{}, false
	}
	for card := range cs.Cards {
		return card, true
	}
	return Card{}, false
}

// Player represents another player in the game.
type Player struct {
	Name        string
	Hand        map[Card]bool
	Constraints []*ConstraintSet
	HandSize    int
}

// NewPlayer creates a new player with an empty hand and no constraints.
// handSize is the number of cards this player was dealt.
func NewPlayer(name string, handSize int) *Player {
	return &Player{
		Name:        name,
		Hand:        make(map[Card]bool),
		Constraints: make([]*ConstraintSet, 0),
		HandSize:    handSize,
	}
}

// AddToHand marks a card as confirmed in this player's hand.
func (p *Player) AddToHand(card Card) {
	p.Hand[card] = true
}

// HasCard returns true if the card is confirmed in this player's hand.
func (p *Player) HasCard(card Card) bool {
	return p.Hand[card]
}

// AddConstraint registers a new constraint set for this player.
func (p *Player) AddConstraint(cs *ConstraintSet) {
	p.Constraints = append(p.Constraints, cs)
}

// RemoveFromConstraints eliminates a card from all pending constraint
// sets. Should be called whenever a card is proven innocent.
func (p *Player) RemoveFromConstraints(card Card) {
	for _, cs := range p.Constraints {
		cs.Remove(card)
	}
}

// IsHandComplete returns true when all cards in the player's hand
// are known, meaning HandSize confirmed cards have been identified.
func (p *Player) IsHandComplete() bool {
	return len(p.Hand) >= p.HandSize
}

// ResolvedConstraints returns all constraint sets that have been
// reduced to a single card and are therefore resolved.
func (p *Player) ResolvedConstraints() []Card {
	resolved := make([]Card, 0)
	for _, cs := range p.Constraints {
		if card, ok := cs.Resolved(); ok {
			resolved = append(resolved, card)
		}
	}
	return resolved
}