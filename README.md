# clue-solver

A logic-based deduction engine for the board game **Clue** (a.k.a. Cluedo).

Tracks the state of every card in the game and automatically infers new information
after each round, narrowing down the solution until the murderer, location, and
weapon are identified.

---

## How it works

The solver models each card in the game as one of three states:

| State       | Meaning                                      |
|-------------|----------------------------------------------|
| `INNOCENT`  | The card is in some player's hand             |
| `GUILTY`    | The card is the solution                     |
| `UNKNOWN`   | Not yet determined                           |

After every suggestion and response, the inference engine applies a set of
logical rules until no new deductions can be made (**fixed point**). If direct
deduction gets stuck, the engine falls back to **hypothesis and refutation**:
it assumes a card is guilty and checks for contradictions.

### Inference rules

1. **Direct elimination** — if a card is `INNOCENT`, remove it from all pending constraint sets.
2. **Set collapse** — if a constraint set is reduced to one card, that card must belong to the refuting player.
3. **Category deduction** — if all cards in a category except one are `INNOCENT`, the remaining card is `GUILTY`.
4. **Cross-deduction** — once a `GUILTY` card is confirmed, all others in its category become `INNOCENT`.
5. **Full-hand constraint** — once all of a player's cards are known, their constraints resolve immediately.

---

## Input

The solver expects:

- Your hand (the cards you hold)
- A list of suggestions made throughout the game, each containing:
  - The three cards suggested (suspect, location, weapon)
  - Which player refuted (if any)
  - Which card was shown to you (if it was shown to you directly)

---

## Output

- The current known state of every card
- The deduced solution, once all three categories are resolved
- Intermediate inferences as they are made

---

## Scope

This project implements the **algorithm only** — the deduction logic is the focus.
No UI, no network, no persistence. Bring your own interface.

---

## Based on

The classic board game **Clue** (Hasbro), known internationally as **Cluedo**.
This project has no affiliation with Hasbro.