package sol

type DarkBaize struct {
	piles    []*Pile
	recycles int
	bookmark int
	script   Scripter
}

type DarkPile struct {
	// cards    []*DarkCard
	category string
	vtable   PileVtabler
	label    string
	moveType MoveType
}

// DarkCard is CardID
