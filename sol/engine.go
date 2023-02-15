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
}

type DarkCard struct {
	pack    int
	suit    int
	ordinal int
	prone   bool
	// ID    CardID
	owner *DarkPile
}
