package solver

import "fmt"

// Category represents the type of a Clue card.
type Category int

const (
	Suspect Category = iota
	Location Category = iota
	Weapon Category = iota
)

func (c Category) String() string {
	switch c {
		case Suspect:
			return "Suspect"
		case Location:
			return "Location"
		case Weapon:
			return "Weapon"
		default:
			return fmt.Sprintf("Unknown Category(%d)", int(c))
	}
}

// CardState represents the state of a card in the game.
type CardState int

const (
	Unknown CardState = iota
	Innocent CardState = iota
	Guilty CardState = iota
)

func (cs CardState) String() string {
	switch cs {
		case Unknown:
			return "Unknown"
		case Innocent:
			return "Innocent"
		case Guilty:
			return "Guilty"
		default:
			return fmt.Sprintf("Unknown CardState(%d)", int(cs))
	}
}

// Card is a struct with a Name, Category, and State. It represents a card in the Clue game, which can be a suspect, location, or weapon. The State indicates whether the card is still unknown, has been determined to be innocent, or has been determined to be guilty.
type Card struct {
	Name     string
	Category Category
}

// String implements fmt.Stringer for Card.
func (c Card) String() string {
	return fmt.Sprintf("%s(%s)", c.Name, c.Category)
}