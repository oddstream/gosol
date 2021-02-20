package sol

// CardTransitionQueue Card Transition Queue
type CardTransitionQueue struct {
	q []item
}

type item struct {
	c    *Card
	x, y int
}

// Add puts a Card transition request into the queue
func (ctq *CardTransitionQueue) Add(c *Card, x, y int) {
	// TODO check if this card is already in queue
	found := false
	for _, i := range ctq.q {
		if i.c == c && !c.lerping {
			found = true
			i.x, i.y = x, y
			break
		}
	}
	if !found {
		ctq.q = append(ctq.q, item{c, x, y})
	}
}

// Update triggers a Card transition once per tick
func (ctq *CardTransitionQueue) Update() {
	if len(ctq.q) > 0 {
		i := ctq.q[0]
		ctq.q = ctq.q[1:]
		i.c.TransitionTo(i.x, i.y)
	}
}
