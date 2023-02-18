package sol

type DarkBaize struct {
	piles    []*Pile
	recycles int
	bookmark int
	script   Scripter
}

type DarkPile struct {
	category string
	vtable   PileVtabler
	label    string
	moveType MoveType
	// needs card []CardID, but at the same time []*Card
}

// DarkCard is CardID
