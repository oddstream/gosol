package sol

type SubtypeAPI interface {
	CanMoveTail([]*Card) (bool, error)
	CanAcceptCard(*Card) (bool, error)
	CanAcceptTail([]*Card) (bool, error)
	TailTapped([]*Card)
	Collect()
	Conformant() bool
	Complete() bool
	UnsortedPairs() int
	Reset()
}
