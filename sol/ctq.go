package sol

const (
	ctqTransition = 0
	ctqFlipUp     = 1
	ctqFlipDown   = 2
)

// CardTransitionQueue Card Transition Queue
type CardTransitionQueue struct {
	q []qItem
}

type qItem struct {
	c      *Card
	x, y   int
	action int
}

// Add puts a Card transition request into the queue
func (ctq *CardTransitionQueue) Add(c *Card, x, y int) {
	// found := false
	// for _, i := range ctq.q {
	// 	if i.c == c && !c.lerping {
	// 		found = true
	// 		i.x, i.y = x, y
	// 		break
	// 	}
	// }
	// if !found {
	// ctq.q = append(ctq.q, qItem{c: c, x: x, y: y, action: ctqTransition})
	// }
	c.TransitionTo(x, y)
}

// AddFlipUp puts a Card transition request into the queue
// func (ctq *CardTransitionQueue) AddFlipUp(c *Card) {
// 	ctq.q = append(ctq.q, qItem{c: c, action: ctqFlipUp})
// }

// AddFlipDown puts a Card transition request into the queue
// func (ctq *CardTransitionQueue) AddFlipDown(c *Card) {
// 	ctq.q = append(ctq.q, qItem{c: c, action: ctqFlipDown})
// }

// Update triggers a Card transition once per tick
func (ctq *CardTransitionQueue) Update() {
	if len(ctq.q) > 0 {
		qi := ctq.q[0]
		ctq.q = ctq.q[1:]
		switch qi.action {
		case ctqTransition:
			qi.c.TransitionTo(qi.x, qi.y)
		case ctqFlipUp:
			qi.c.FlipUp()
		case ctqFlipDown:
			qi.c.FlipDown()
		}
	}
}
