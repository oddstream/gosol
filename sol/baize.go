package sol

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"oddstream.games/gosol/util"
)

// Baize object describes the baize
type Baize struct {
	Piles     []*Pile
	Variant   string
	Seed      int64
	UndoStack []SaveableBaize
	stroke    *Stroke
}

// NewBaize is the factory func for Baize object
func NewBaize() *Baize {
	b := &Baize{Variant: TheUserData.Variant, Seed: time.Now().UnixNano()}
	BuildScalableCardImages() // need to do this after CardWidth,Height set - not in a func init()
	b.NewVariant(TheUserData.Variant)
	return b
}

// Reset the Baize
func (b *Baize) Reset() {
	b.Piles = nil
	b.UndoStack = nil
	b.Variant = TheUserData.Variant
	b.Seed = time.Now().UnixNano()
	b.stroke = nil
}

// Restart the Baize without changing variant or seed
func (b *Baize) Restart() {
	stock := b.findPilePrefix("Stock")
	if stock == nil {
		log.Fatal("cannot find stock pile to recall cards with")
	}

	for _, p := range b.Piles {
		p.Reset() // stock needs resetting, too
		if p == stock {
			continue
		}
		if p.CardCount() > 0 {
			b.MoveCards(p.Cards[0], stock)
		}
		// if p.CardCount() != 0 {
		// 	log.Fatal(p.Class, " still contains ", p.CardCount(), " cards")
		// }
	}

	// if DebugMode {
	// 	println("cards recalled to stock, now contains", stock.CardCount(), "cards")
	// 	for _, c := range stock.Cards {
	// 		if !c.prone {
	// 			log.Fatal("face up card found in stock")
	// 		}
	// 		if c.owner != stock {
	// 			log.Fatal("card in stock belongs to", c.owner.Class)
	// 		}
	// 	}
	// }

	// TODO wait for lerping back to stock to finish before shuffling, for now use c.SetPosition() to cancel lerping
	// x, y := stock.Position()
	// for _, c := range stock.Cards {
	// 	c.SetPosition(x, y)
	// }

	shuffleCards(stock, b.Seed)

	b.UndoStack = nil // StartGame will deal cards then do initial UndoPush()
	b.stroke = nil
}

// StartGame given existing variant and seed
func (b *Baize) StartGame() {
	b.dealCards()
	b.UndoPush()
}

// RestartGame resets Baize and restarts current variant with same seed
func (b *Baize) RestartGame() {
	b.Restart()
	b.StartGame()
}

// NewGame resets Baize and restarts current variant with a new seed
func (b *Baize) NewGame() {
	b.Restart()
	b.Seed = time.Now().UnixNano()
	b.StartGame()
}

// NewVariant resets Baize and starts a new game with a new variant and seed
func (b *Baize) NewVariant(v string) {
	b.Reset()
	b.Variant = v
	b.Seed = time.Now().UnixNano()

	piles, ok := buildVariantPiles(b.Variant)
	if !ok {
		log.Fatal("unknown variant", b.Variant)
	}
	b.Piles = piles

	{
		maxX := 0
		for _, p := range b.Piles {
			if p.X > maxX {
				maxX = p.X
			}
		}
		ebiten.SetWindowSize((maxX+2)*(CardWidth+10), WindowHeight)
	}

	stock := b.findPilePrefix("Stock")
	if stock == nil {
		log.Fatal("Cannot find stock pile to create cards with")
	}
	createCards(stock)
	shuffleCards(stock, b.Seed)

	b.StartGame()
}

// doesn't work because time.Sleep suspends the whole thread, not allowing Ebiten to breathe
// func (b *Baize) waitForCards() {
// 	for {
// 		ani := 0
// 		for _, p := range b.Piles {
// 			for _, c := range p.Cards {
// 				if c.Animating() {
// 					ani++
// 				}
// 			}
// 		}
// 		if ani == 0 {
// 			break
// 		}
// 		println("waiting for", ani, "cards")
// 		b.Update()
// 			time.Sleep(time.Second)
// 	}
// }

func (b *Baize) dealCards() {
	stock := b.findPilePrefix("Stock")
	for _, p := range b.Piles {
		deal := p.GetStringAttribute("Deal")
		if deal == "" {
			continue
		}
		for _, d := range deal {
			switch d {
			case 'u':
				c := stock.Pop() // this will flip card up
				if c.prone {
					log.Fatal("popped a face down card from stock")
				}
				if c == nil {
					log.Fatal("out of cards during deal")
				}
				p.Push(c)
			case 'd':
				c := stock.Pop() // this will flip card up
				if c == nil {
					log.Fatal("out of cards during deal")
				}
				c.FlipDown()
				p.Push(c)
			}
		}
	}

	for _, p := range b.Piles {
		if p.Class == "Foundation" && p.CardCount() == 1 {
			afp := p.GetBoolAttribute("AcceptFirstPush")
			if afp {
				ord := p.Peek().ordinal
				for _, fp := range b.Piles {
					if fp.Class == "Foundation" {
						fp.SetAccept(ord)
					}
				}
			}
		}
	}

}

func (b *Baize) findPile(cls string) *Pile {
	for _, p := range b.Piles {
		if p.Class == cls {
			return p
		}
	}
	return nil
}

func (b *Baize) findPilePrefix(cls string) *Pile {
	for _, p := range b.Piles {
		if strings.HasPrefix(p.Class, cls) {
			return p
		}
	}
	return nil
}

// findPileAt finds the pile under the mouse click or touch
func (b *Baize) findPileAt(x, y int) *Pile {
	for _, o := range b.Piles {
		if util.InRect(x, y, o.FannedRect) {
			return o
		}
	}
	return nil
}

// findCardAt finds the tile under the mouse click or touch
func (b *Baize) findCardAt(x, y int) *Card {
	for _, p := range b.Piles {
		for i := p.CardCount() - 1; i >= 0; i-- {
			c := p.Cards[i]
			if util.InRect(x, y, c.Rect) {
				return c
			}
		}
	}
	return nil
}

func (b *Baize) countPiles(cls string) (total, count int) {
	total = 0
	count = 0
	for _, p := range b.Piles {
		if p.Class == cls {
			total++
			if 0 == p.CardCount() {
				count++
			}
		}
	}
	return
}

// PileTapped is called when a pile has been tapped
func (b *Baize) PileTapped(p *Pile) {

	// this method is in Baize because it needs access to Baize.findPile()
	if p.Class != "Stock" {
		return
	}

	if p.localRecycles > 0 {
		waste := b.findPile("Waste")
		if waste == nil || len(waste.Cards) == 0 {
			return
		}
		stock := b.findPile("Stock")
		for len(waste.Cards) > 0 {
			c := waste.Pop()
			stock.Push(c) // this will flip card down
		}
		p.SetRecycles(p.localRecycles - 1)
		b.AfterUserMove()
	}

	// println("pile", p.Class, "tapped")
}

// CardTapped is called when a card has been tapped
func (b *Baize) CardTapped(c *Card) {

	// println("card",c.id,"tapped")

	if c.Animating() {
		println("cannot tap an animating card")
		return
	}

	pSrc := c.owner

	// can only tap top card
	if c != pSrc.Peek() {
		c.Shake()
		return
	}

	targetClass := c.owner.GetStringAttribute("Target")

	moved := false

	switch pSrc.Class {
	case "Stock":
		// Tap on a Stock card to send it to Waste
		if targetClass == "" {
			targetClass = "Waste"
		}
		cardsToMove, _ := pSrc.GetIntAttribute("CardsToMove")
		if cardsToMove == 0 {
			cardsToMove = 1
		}
		for _, p := range b.Piles {
			if targetClass == p.Class {
				// println("found a", p.Class)
				if p.CanAcceptCard(c) {
					// println(p.Class, "can accept", c.id)
					for cardsToMove > 0 && c != nil {
						cardsToMove--
						b.MoveCards(c, p)
						c = pSrc.Peek()
						moved = true
					}
				}
			}
		}
	case "StockSpider":
		_, empty := b.countPiles("Tableau")
		if empty > 0 {
			stock := b.findPilePrefix("Stock")
			if stock.CardCount() > empty {
				println("all tableaux spaces must be filled before dealing a new row")
				break
			}
		}
		fallthrough
	case "StockScorpion":
		if targetClass == "" {
			targetClass = "Tableau"
		}
		for _, p := range b.Piles {
			if p.Class == targetClass {
				b.MoveCards(c, p)
				moved = true
				c.prone = false
				c = pSrc.Peek()
			}
			if c == nil {
				break
			}
		}
	case "Tableau", "Waste", "Cell", "Reserve":
		for _, p := range b.Piles {
			if p.Class == "Foundation" {
				if p.CanAcceptCard(c) {
					b.MoveCards(c, p)
					moved = true
					break
				}
			}
		}
	default:
		println("clueless when tapping on a", pSrc.Class, "card")
	}

	if !moved {
		if c != nil {
			c.Shake()
		}
	} else {
		b.AfterUserMove()
	}

	// TODO else test other piles to see if this card is accepted?
}

// MoveCards from one pile to another, always from card downwards (inclusive)
func (b *Baize) MoveCards(c *Card, dst *Pile) {

	src := c.owner
	moveFrom := len(src.Cards)
	tmp := make([]*Card, 0)

	// find the index of the first card we will move
	for i, sc := range src.Cards {
		if sc == c {
			moveFrom = i
			break
		}
	}

	if moveFrom == len(src.Cards) {
		log.Fatal("MoveCards could not find card in source")
	}

	oldSrcLen := len(src.Cards)

	// pop the tail off the source and push onto temp stack
	for i := len(src.Cards) - 1; i >= moveFrom; i-- {
		sc := src.Pop()
		tmp = append(tmp, sc)
	}

	// pop cards off the temp stack and onto the destination
	for len(tmp) > 0 {
		dc := tmp[len(tmp)-1]
		tmp = tmp[:len(tmp)-1]
		dst.Push(dc)
	}

	if oldSrcLen == len(src.Cards) {
		println("MoveCards - nothing happened")
	} else {
		// flip up an exposed source card
		if !strings.HasPrefix(src.Class, "Stock") {
			tc := src.Peek()
			if tc != nil {
				tc.FlipUp()
			}
		}
	}
}

// AutoMoves performs post user-moves
func (b *Baize) AutoMoves() {

	// TODO move cards to Foundations, using Opsole safe logic

	for _, p := range b.Piles {
		if p.CardCount() == 0 {
			amf := p.GetStringAttribute("AutoFillFrom")
			if amf != "" {
				src := b.findPile(amf)
				if src != nil {
					c := src.Peek()
					if c != nil {
						b.MoveCards(c, p)
					}
				}
			}
		}
	}

}

// AfterUserMove runs after the user has made a move
func (b *Baize) AfterUserMove() {

	b.AutoMoves()

	//

	var oldChecksum, newChecksum uint32
	var ok bool

	//

	if len(b.UndoStack) == 0 {
		log.Fatal("undo stack is empty in AfterUserMove()")
	} else {
		oldChecksum, ok = b.UndoPeekChecksum()
		if !ok {
			log.Fatal("error peeking undo stack checksum")
		}
	}
	newChecksum = b.Checksum()
	// println(oldChecksum, newChecksum)
	if oldChecksum != newChecksum {
		b.UndoPush()
	} else {
		println("not pushing to undo because checksums match")
	}

	//

	complete := true
	for _, p := range b.Piles {
		if !p.IsComplete() {
			complete = false
			break
		}
	}
	if complete {
		println(b.Variant, "complete")
	}
}

// UndoPush pushes the current state onto the undo stack
func (b *Baize) UndoPush() {
	b.UndoStack = append(b.UndoStack, b.Saveable())
}

// UndoPop pops a state off the undo stack
func (b *Baize) UndoPop() (SaveableBaize, bool) {
	if len(b.UndoStack) > 0 {
		sav := b.UndoStack[len(b.UndoStack)-1]
		b.UndoStack = b.UndoStack[:len(b.UndoStack)-1]
		return sav, true
	}
	return SaveableBaize{}, false
}

// UndoPeekChecksum peeks the state at the top of the undo stack
func (b *Baize) UndoPeekChecksum() (uint32, bool) {
	if len(b.UndoStack) > 0 {
		sav := b.UndoStack[len(b.UndoStack)-1]
		return sav.Checksum, true
	}
	return 0, false
}

// Undo reverts the Baize state to it's previous state
func (b *Baize) Undo() {
	if len(b.UndoStack) < 2 {
		println("nothing to undo")
		return
	}
	sav, ok := b.UndoPop() // removes current state
	if !ok {
		log.Fatal("error popping current state from undo stack")
	}

	sav, ok = b.UndoPop() // removes previous state for examination
	if !ok {
		log.Fatal("error popping second from undo stack")
	}
	b.UpdateFromSaveable(sav)
	b.UndoPush() // replace current state
}

func (b *Baize) largestIntersection(c *Card) *Pile {
	var largest int = 0
	var pile *Pile = nil
	cx0, cy0, cx1, cy1 := c.Rect()
	for _, p := range b.Piles {
		px0, py0, px1, py1 := p.FannedRect()
		i := util.OverlapArea(cx0, cy0, cx1, cy1, px0, py0, px1, py1)
		if i > largest {
			largest = i
			pile = p
		}
	}
	return pile
}

// Layout implements ebiten.Game's Layout.
func (b *Baize) Layout(outsideWidth, outsideHeight int) (int, int) {

	for _, p := range b.Piles {
		p.Layout(outsideWidth, outsideHeight)
	}

	return outsideWidth, outsideHeight
}

// Update the baize state (transitions, user input)
func (b *Baize) Update() error {

	if inpututil.IsKeyJustReleased(ebiten.KeyN) {
		b.NewGame()
		return nil
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		b.RestartGame()
		return nil
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyU) {
		b.Undo()
		return nil
	}

	if b.stroke == nil {
		var s *Stroke

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			s = NewStroke(&MouseStrokeSource{})
		}
		ts := inpututil.JustPressedTouchIDs()
		if ts != nil && len(ts) == 1 {
			s = NewStroke(&TouchStrokeSource{ts[0]})
		}

		if s != nil {
			sx, sy := s.Position()
			// maybe user is tapping or starting to drag a card
			c := b.findCardAt(sx, sy)
			if c != nil {
				if c.owner.StartDrag(b.Piles, c) {
					b.stroke = s
					b.stroke.SetDraggingObject(c)
				} else {
					println("cannot drag those cards")
				}
			} else {
				// maybe user is tapping an empty pile (eg to recycle waste to stock)
				p := b.findPileAt(sx, sy)
				if p != nil {
					b.stroke = s
					b.stroke.SetDraggingObject(p)
				}
			}
		}
	} else {
		b.stroke.Update()
		switch v := b.stroke.DraggingObject().(type) {
		case *Card:
			c := v
			if b.stroke.IsReleased() {
				if b.stroke.IsTapped() {
					c.owner.StopDrag(c)
					b.CardTapped(c)
				} else {
					// sx, sy := b.stroke.Position()
					// p := b.findPileAt(sx, sy)
					p := b.largestIntersection(c)
					if p == nil {
						c.owner.CancelDrag(c)
					} else {
						// println("found pile", o.Class())
						if p == c.owner {
							println("baize cannot drag cards to owning pile")
						}
						if p.CanAcceptTail(c.owner.Tail) {
							c.owner.StopDrag(c)
							b.MoveCards(c, p)
							b.AfterUserMove()
						} else {
							c.owner.CancelDrag(c)
						}
					}
				}
				b.stroke = nil
			} else {
				dx, dy := b.stroke.PositionDiff()
				c.owner.DragTailBy(dx, dy)
			}
		case *Pile:
			p := v
			if b.stroke.IsReleased() {
				if b.stroke.IsTapped() {
					b.PileTapped(p)
				}
				b.stroke = nil
			}
		default:
			log.Fatal("unknown type of dragging object")
		}
	}

	for _, p := range b.Piles {
		p.Update()
	}

	// CTQ.Update()

	return nil
}

// Draw renders the baize into the screen
func (b *Baize) Draw(screen *ebiten.Image) {

	screen.Fill(colorBaize)

	for _, p := range b.Piles {
		p.Draw(screen)
	}
	for _, p := range b.Piles {
		p.DrawCards(screen)
	}
	for _, p := range b.Piles {
		p.DrawAnimatingCards(screen)
	}
	if DebugMode {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("NumGC %v, Undo %d", ms.NumGC, len(b.UndoStack)))
	}
}
