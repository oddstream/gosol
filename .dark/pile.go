package dark

type MoveType int

const (
	MOVE_NONE MoveType = iota
	MOVE_ANY
	MOVE_ONE
	MOVE_ONE_PLUS
	MOVE_ONE_OR_ALL
)

type MovableTail struct {
	dst  *Pile
	tail []*Card
}

// PileVtabler interface for each subpile type, implements the behaviours
// specific to each subtype
type PileVtabler interface {
	CanAcceptTail([]*Card) (bool, error)
	TailTapped([]*Card)
	Conformant() bool
	UnsortedPairs() int
	MovableTails() []*MovableTail
}
