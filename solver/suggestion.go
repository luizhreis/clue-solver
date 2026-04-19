package solver

// Suggestion represents a single guess made during the game,
// along with its outcome.
type Suggestion struct {
	Suspect  Card
	Location Card
	Weapon   Card
	Guesser  *Player
	Refuter  *Player
	ShownCard *Card
}

// NewSuggestion creates a Suggestion. Both refuter and shownCard are
// optional — pass nil when no one refuted, or when the card was not
// shown to you directly.
func NewSuggestion(
	suspect Card,
	location Card,
	weapon Card,
	guesser *Player,
	refuter *Player,
	shownCard *Card,
) *Suggestion {
	return &Suggestion{
		Suspect:   suspect,
		Location:  location,
		Weapon:    weapon,
		Guesser:   guesser,
		Refuter:   refuter,
		ShownCard: shownCard,
	}
}

// cards returns the three cards involved in this suggestion.
// Internal helper to avoid repeating the same slice literal.
func (s *Suggestion) cards() []Card {
	return []Card{s.Suspect, s.Location, s.Weapon}
}

// Process applies the suggestion's outcome to the game state.
// It covers three scenarios:
//
//  1. No one refuted → all three cards are Guilty (solution found).
//  2. A card was shown to us → that card is Innocent.
//  3. Someone refuted but didn't show us → add a ConstraintSet to the refuter.
func (s *Suggestion) Process(g *Game) {
	switch {
	case s.Refuter == nil:
		s.processNoRefutation(g)
	case s.ShownCard != nil:
		s.processCardShown(g)
	default:
		s.processHiddenRefutation(g)
	}
}

// processNoRefutation marks all three suggested cards as Guilty.
// This only happens when no player — including ourselves — can refute.
func (s *Suggestion) processNoRefutation(g *Game) {
	for _, card := range s.cards() {
		g.SetGuilty(card)
	}
}

// processCardShown marks the revealed card as Innocent and also
// adds it to the refuter's confirmed hand.
func (s *Suggestion) processCardShown(g *Game) {
	g.SetInnocent(*s.ShownCard)
	if s.Refuter != nil {
		s.Refuter.AddToHand(*s.ShownCard)
	}
}

// processHiddenRefutation registers a ConstraintSet on the refuter,
// meaning they hold at least one of the three suggested cards.
// Cards already known to be Innocent are excluded from the set upfront.
func (s *Suggestion) processHiddenRefutation(g *Game) {
	candidates := make([]Card, 0, 3)
	for _, card := range s.cards() {
		if g.StateOf(card) != Innocent {
			candidates = append(candidates, card)
		}
	}

	if len(candidates) == 0 {
		return
	}

	cs := NewConstraintSet(candidates)
	s.Refuter.AddConstraint(cs)
}