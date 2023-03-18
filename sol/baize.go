package sol

import (
	"fmt"
	"hash/crc32"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/ui"
	"oddstream.games/gosol/util"
)

const (
	dirtyWindowSize = 1 << iota
	dirtyPilePositions
	dirtyCardSizes
	dirtyCardImages
	dirtyPileBackgrounds
	dirtyCardPositions
)

// Baize object describes the baize
type Baize struct {
	variant      string
	piles        []*Pile
	cardCount    int
	recycles     int
	bookmark     int
	script       Scripter
	undoStack    []*SavableBaize
	dirtyFlags   uint32 // what needs doing when we Update
	moves        int    // number of possible (not useless) moves
	fmoves       int    // number of possible moves to a Foundation (for enabling Collect button)
	stroke       *input.Stroke
	dragStart    image.Point
	dragOffset   image.Point
	WindowWidth  int // the most recent window width given to Layout
	WindowHeight int // the most recent window height given to Layout
	// hotCard      *Card
}

//--+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8

// NewBaize is the factory func for the single Baize object
func NewBaize(variant string) *Baize {
	// let WindowWidth, WindowHeight be zero, so that the first Layout will
	// trigger card scaling and pile placement
	var script Scripter
	var ok bool
	if script, ok = Variants[variant]; !ok {
		log.Printf("do not know how to play " + variant)
		return nil
	}
	return &Baize{variant: variant, script: script, dirtyFlags: 0xFFFF}
}

func (b *Baize) flagSet(flag uint32) bool {
	return b.dirtyFlags&flag == flag
}

func (b *Baize) setFlag(flag uint32) {
	b.dirtyFlags |= flag
}

// func (b *Baize) clearFlag(flag uint32) {
// 	b.dirtyFlags &= ^flag
// }

// func (b *Baize) Valid() bool {
// 	return b != nil
// }

func (b *Baize) CRC() uint32 {
	/*
		var crc uint = 0xFFFFFFFF
		var mask uint
		for _, p := range b.piles {
			crc = crc ^ uint(p.Len())
			for j := 7; j >= 0; j-- {
				mask = -(crc & 1)
				crc = (crc >> 1) ^ (0xEDB88320 & mask)
			}
		}
		return ^crc // bitwise NOT
	*/
	var lens []byte
	for _, p := range b.piles {
		lens = append(lens, byte(p.Len()))
	}
	return crc32.ChecksumIEEE(lens)
}

func (b *Baize) AddPile(pile *Pile) {
	b.piles = append(b.piles, pile)
}

func (b *Baize) Refan() {
	b.setFlag(dirtyCardPositions)
}

// NewDeal restarts current variant (ie no pile building) with a new seed
func (b *Baize) NewDeal() {

	// a virgin game has one state on the undo stack
	if len(b.undoStack) > 1 && !b.Complete() {
		percent := b.PercentComplete()
		toastStr := TheGame.Statistics.RecordLostGame(b.variant, percent)
		TheGame.UI.Toast("Fail", toastStr)
	}

	// for {
	b.Reset()

	for _, p := range b.piles {
		p.Reset()
	}

	// Stock.Fill() needs parameters
	packs := b.script.Packs()
	suits := b.script.Suits()
	b.cardCount = b.script.Stock().Fill(packs, suits)
	b.script.Stock().Shuffle()
	b.script.StartGame()
	b.UndoPush()
	b.FindDestinations()

	// 	if b.moves > 0 {
	// 		break
	// 	}
	// 	TheGame.UI.Toast("Glass", "Found a deal with no moves")
	// }

	sound.Play("Fan")

	b.setFlag(dirtyCardPositions)
}

func (b *Baize) MirrorSlots() {
	/*
		0 1 2 3 4 5
		5 4 3 2 1 0

		0 1 2 3 4
		4 3 2 1 0
	*/
	var minX int = 32767
	var maxX int = 0
	for _, p := range b.piles {
		if p.Slot().X < 0 {
			continue // ignore hidden pile
		}
		if p.Slot().X < minX {
			minX = p.Slot().X
		}
		if p.Slot().X > maxX {
			maxX = p.Slot().X
		}
	}
	for _, p := range b.piles {
		slot := p.Slot()
		if slot.X < 0 {
			continue // ignore hidden pile
		}
		p.SetSlot(image.Point{X: maxX - slot.X + minX, Y: slot.Y})
		switch p.FanType() {
		case FAN_RIGHT:
			p.SetFanType(FAN_LEFT)
		case FAN_LEFT:
			p.SetFanType(FAN_RIGHT)
		case FAN_RIGHT3:
			p.SetFanType(FAN_LEFT3)
		case FAN_LEFT3:
			p.SetFanType(FAN_RIGHT3)
		}
	}
}

func (b *Baize) Reset() {
	b.StopSpinning()
	b.undoStack = []*SavableBaize{}
	b.bookmark = 0
	b.recycles = 0
	// leave script intact
}

// StartFreshGame resets Baize and starts a new game with a new seed
func (b *Baize) StartFreshGame() {
	b.Reset()
	b.piles = []*Pile{}
	b.script.BuildPiles()
	if TheGame.Settings.MirrorBaize {
		b.MirrorSlots()
	}
	// b.FindBuddyPiles()

	TheGame.UI.SetTitle(b.variant)
	sound.Play("Fan")
	b.dirtyFlags = 0xFFFF

	b.script.StartGame()
	b.UndoPush()
	b.FindDestinations()
}

func (b *Baize) ChangeVariant(newVariant string) {
	// no longer record a lost game here because variants saved in separate .json files
	var newScript Scripter
	var ok bool
	if newScript, ok = Variants[newVariant]; !ok {
		TheGame.UI.Toast("Error", "Do not know how to play "+newVariant)
		return
	}
	b.Save()
	b.variant = newVariant
	TheGame.Settings.Variant = b.variant
	TheGame.Settings.Save()
	b.script = newScript
	b.StartFreshGame()
	if !NoGameLoad {
		TheGame.Baize.Load()
	}
}

func (b *Baize) SetUndoStack(undoStack []*SavableBaize) {
	b.undoStack = undoStack
	TheGame.UI.Toast("Glass", "Loaded a saved game of "+b.variant)
	sav := b.UndoPeek()
	b.updateFromSavable(sav)
	b.FindDestinations()
	TheGame.UI.HideFAB()
	if b.Complete() {
		TheGame.UI.Toast("Complete", "Complete")
		TheGame.UI.AddButtonToFAB("star", ebiten.KeyN)
		b.StartSpinning()
	} else if b.Conformant() {
		TheGame.UI.AddButtonToFAB("done_all", ebiten.KeyC)
	} else if b.moves == 0 {
		TheGame.UI.Toast("Error", "No movable cards")
		TheGame.UI.AddButtonToFAB("star", ebiten.KeyN)
		TheGame.UI.AddButtonToFAB("restore", ebiten.KeyR)
		if b.bookmark > 0 {
			TheGame.UI.AddButtonToFAB("bookmark", ebiten.KeyL)
		}
	}
}

// findPileAt finds the Pile under the mouse position
func (b *Baize) FindPileAt(pt image.Point) *Pile {
	for _, p := range b.piles {
		if pt.In(p.ScreenRect()) {
			return p
		}
	}
	return nil
}

// FindLowestCardAt finds the bottom-most Card under the mouse position
func (b *Baize) FindLowestCardAt(pt image.Point) *Card {
	for _, p := range b.piles {
		for i := p.Len() - 1; i >= 0; i-- {
			c := p.Get(i)
			if pt.In(c.ScreenRect()) {
				return c
			}
		}
	}
	return nil
}

// FindHighestCardAt finds the top-most Card under the mouse position
func (b *Baize) FindHighestCardAt(pt image.Point) *Card {
	for _, p := range b.piles {
		for _, c := range p.cards {
			if pt.In(c.ScreenRect()) {
				return c
			}
		}
	}
	return nil
}

func (b *Baize) LargestIntersection(c *Card) *Pile {
	var largestArea int = 0
	var pile *Pile = nil
	cardRect := c.BaizeRect()
	for _, p := range b.piles {
		if p == c.Owner() {
			continue
		}
		pileRect := p.FannedBaizeRect()
		intersectRect := pileRect.Intersect(cardRect)
		area := intersectRect.Dx() * intersectRect.Dy()
		if area > largestArea {
			largestArea = area
			pile = p
		}
	}
	return pile
}

// StartDrag return true if the Baize can be dragged
func (b *Baize) StartDrag() bool {
	b.dragStart = b.dragOffset
	return true
}

// DragBy move ('scroll') the Baize by dragging it
// dx, dy is the difference between where the drag started and where the cursor is now
func (b *Baize) DragBy(dx, dy int) {
	b.dragOffset.X = b.dragStart.X + dx
	if b.dragOffset.X > 0 {
		b.dragOffset.X = 0 // DragOffsetX should only ever be 0 or -ve
	}
	b.dragOffset.Y = b.dragStart.Y + dy
	if b.dragOffset.Y > 0 {
		b.dragOffset.Y = 0 // DragOffsetY should only ever be 0 or -ve
	}
}

// StopDrag stop dragging the Baize
func (b *Baize) StopDrag() {
	b.setFlag(dirtyCardPositions)
}

// StartSpinning tells all the cards to start spinning
func (b *Baize) StartSpinning() {
	for _, p := range b.piles {
		// use a method expression, which yields a function value with a regular first parameter taking the place of the receiver
		p.ApplyToCards((*Card).StartSpinning)
	}
}

// StopSpinning tells all the cards to stop spinning and return to their upright position
// debug only
func (b *Baize) StopSpinning() {
	for _, p := range b.piles {
		// use a method expression, which yields a function value with a regular first parameter taking the place of the receiver
		p.ApplyToCards((*Card).StopSpinning)
	}
	b.setFlag(dirtyCardPositions)
}

func (b *Baize) AfterUserMove() {
	b.script.AfterMove()
	b.UndoPush()
	b.FindDestinations()
	TheGame.UI.HideFAB()
	if b.Complete() {
		TheGame.UI.AddButtonToFAB("star", ebiten.KeyN)
		b.StartSpinning()
		{
			var toastStr = TheGame.Statistics.RecordWonGame(b.variant, len(b.undoStack)-1)
			TheGame.UI.Toast("Complete", toastStr)
		}
		ShowStatisticsDrawer()
	} else if b.Conformant() {
		TheGame.UI.AddButtonToFAB("done_all", ebiten.KeyC)
	} else if b.moves == 0 {
		TheGame.UI.ToastError("No movable cards")
		TheGame.UI.AddButtonToFAB("star", ebiten.KeyN)
		TheGame.UI.AddButtonToFAB("restore", ebiten.KeyR)
		if b.bookmark > 0 {
			TheGame.UI.AddButtonToFAB("bookmark", ebiten.KeyL)
		}
	}
}

// AfterAfterMove checks for and executes an automatic collect.
// Kept as separated-out function at the moment, in case this
// creates a horrible recursive loop
func (b *Baize) AfterAfterUserMove() {
	if b.fmoves > 0 && TheGame.Settings.AutoCollect {
		b.Collect2()
	}
}

/*
InputStart finds out what object the user input is starting on
(UI Container > Card > Pile > Baize, in that order)
then tells that object.

If the Input starts on a Card, then a tail of cards is formed.
*/
func (b *Baize) InputStart(v input.StrokeEvent) {
	b.stroke = v.Stroke

	if con := TheGame.UI.FindContainerAt(v.X, v.Y); con != nil {
		if w := con.FindWidgetAt(v.X, v.Y); w != nil {
			b.stroke.SetDraggedObject(w)
		} else {
			con.StartDrag()
			b.stroke.SetDraggedObject(con)
		}
	} else {
		pt := image.Pt(v.X, v.Y)
		if card := b.FindLowestCardAt(pt); card != nil {
			if card.Lerping() {
				TheGame.UI.Toast("Glass", "Confusing to move a moving card")
				v.Stroke.Cancel()
			} else {
				tail := card.Owner().MakeTail(card)
				b.StartTailDrag(tail)
				b.stroke.SetDraggedObject(tail)
			}
		} else {
			if p := b.FindPileAt(pt); p != nil {
				b.stroke.SetDraggedObject(p)
			} else {
				if b.StartDrag() {
					b.stroke.SetDraggedObject(b)
				} else {
					v.Stroke.Cancel()
				}
			}
		}
	}
}

func (b *Baize) InputMove(v input.StrokeEvent) {
	if v.Stroke.DraggedObject() == nil {
		return
		// log.Panic("*** move stroke with nil dragged object ***")
	}
	// for _, p := range b.piles {
	// 	p.target = false
	// }
	switch obj := v.Stroke.DraggedObject().(type) {
	case ui.Containery:
		obj.DragBy(v.Stroke.PositionDiff())
	case ui.Widgety:
		obj.Parent().DragBy(v.Stroke.PositionDiff())
	case []*Card:
		dx, dy := v.Stroke.PositionDiff()
		b.DragTailBy(obj, dx, dy)
		// if c, ok := v.Stroke.DraggedObject().(*Card); ok {
		// 	if p := b.LargestIntersection(c); p != nil {
		// 		p.target = true
		// 	}
		// }
	case *Pile:
		// do nothing
	case *Baize:
		b.DragBy(v.Stroke.PositionDiff())
	default:
		log.Panic("*** unknown move dragging object ***")
	}
}

func (b *Baize) InputStop(v input.StrokeEvent) {
	if v.Stroke.DraggedObject() == nil {
		return
		// log.Panic("*** stop stroke with nil dragged object ***")
	}
	// for _, p := range b.piles {
	// 	p.SetTarget(false)
	// }
	switch obj := v.Stroke.DraggedObject().(type) {
	case ui.Containery:
		obj.StopDrag()
	case ui.Widgety:
		obj.Parent().StopDrag()
	case []*Card:
		tail := obj     // alias for readability
		card := tail[0] // for readability
		if card.WasDragged() {
			src := card.Owner()
			// tap handled elsewhere
			// tap is time-limited
			if dst := b.LargestIntersection(card); dst == nil {
				// println("no intersection for", c.String())
				b.CancelTailDrag(tail)
			} else {
				var ok bool
				var err error
				// generically speaking, can this tail be moved?
				if ok, err = src.CanMoveTail(tail); !ok {
					TheGame.UI.ToastError(err.Error())
					b.CancelTailDrag(tail)
				} else {
					if ok, err = dst.vtable.CanAcceptTail(tail); !ok {
						TheGame.UI.ToastError(err.Error())
						b.CancelTailDrag(tail)
					} else {
						// it's ok to move this tail
						if src == dst {
							b.CancelTailDrag(tail)
						} else if ok, err = b.script.TailMoveError(tail); !ok {
							TheGame.UI.ToastError(err.Error())
							b.CancelTailDrag(tail)
						} else {
							crc := b.CRC()
							if len(tail) == 1 {
								MoveCard(src, dst)
							} else {
								MoveTail(card, dst)
							}
							b.StopTailDrag(tail) // do this before AfterUserMove
							if crc != b.CRC() {
								b.AfterUserMove()
								b.AfterAfterUserMove()
							}
						}
					}
				}
			}
		}
		if DebugMode && card.Dragging() {
			log.Printf("Card %s is still dragging", card.String())
		}
	case *Pile:
		// do nothing
	case *Baize:
		// println("stop dragging baize")
		b.StopDrag()
	default:
		log.Panic("*** stop dragging unknown object ***")
	}
}

func (b *Baize) InputCancel(v input.StrokeEvent) {
	if v.Stroke.DraggedObject() == nil {
		log.Print("*** cancel stroke with nil dragged object ***")
		return
	}
	switch obj := v.Stroke.DraggedObject().(type) { // type switch
	case ui.Containery:
		obj.CancelDrag()
	case ui.Widgety:
		obj.Parent().CancelDrag()
	case []*Card:
		b.CancelTailDrag(obj)
	case *Pile:
		// p := v.Stroke.DraggedObject().(*Pile)
		// println("stop dragging pile", p.Class)
		// do nothing
	case *Baize:
		// println("stop dragging baize")
		b.StopDrag()
	default:
		log.Panic("*** cancel dragging unknown object ***")
	}
}

func (b *Baize) InputTap(v input.StrokeEvent) {
	// stroke sends a tap event, and later sends a cancel event
	// println("Baize.NotifyCallback() tap", v.X, v.Y)
	switch obj := v.Stroke.DraggedObject().(type) {
	case ui.Containery:
		obj.Tapped()
	case ui.Widgety:
		obj.Tapped()
	case []*Card:
		// offer TailTapped to the script first
		// to implement things like Stock.TailTapped
		// if the script doesn't want to do anything, it can call pile.vtable.TailTapped
		// which will either ignore it (eg Foundation, Discard)
		// or use Pile.DefaultTailTapped
		crc := b.CRC()
		b.script.TailTapped(obj)
		if crc != b.CRC() {
			sound.Play("Slide")
			b.AfterUserMove()
			b.AfterAfterUserMove()
		} else {
			TheGame.UI.Toast("Error", "Attention!")
		}
	case *Pile:
		crc := b.CRC()
		b.script.PileTapped(obj)
		if crc != b.CRC() {
			sound.Play("Shove")
			b.AfterUserMove()
			b.AfterAfterUserMove()
		}
	case *Baize:
		pt := image.Pt(v.X, v.Y)
		// a tap outside any open ui drawer (ie on the baize) closes the drawer
		if con := TheGame.UI.VisibleDrawer(); con != nil && !pt.In(image.Rect(con.Rect())) {
			con.Hide()
		}
	default:
		log.Panic("*** tap unknown object ***")
	}
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (b *Baize) NotifyCallback(v input.StrokeEvent) {
	switch v.Event {
	case input.Start:
		b.InputStart(v)
	case input.Move:
		b.InputMove(v)
	case input.Stop:
		b.InputStop(v)
	case input.Cancel:
		b.InputCancel(v)
	case input.Tap:
		b.InputTap(v)
	default:
		log.Panic("*** unknown stroke event ***", v.Event)
	}
}

// ForeachCard applys a function to each card
func (b *Baize) ForeachCard(fn func(*Card)) {
	for _, p := range b.piles {
		for _, c := range p.cards {
			fn(c)
		}
	}
}

// ApplyToTail applies a method func to this card and all the others after it in the tail
func (b *Baize) ApplyToTail(tail []*Card, fn func(*Card)) {
	// https://golang.org/ref/spec#Method_expressions
	// (*Card).CancelDrag yields a function with the signature func(*Card)
	// fn passed as a method expression so add the receiver explicitly
	for _, c := range tail {
		fn(c)
	}
}

// DragTailBy repositions all the cards in the tail: dx, dy is the position difference from the start of the drag
func (b *Baize) DragTailBy(tail []*Card, dx, dy int) {
	// println("Baize.DragTailBy(", dx, dy, ")")
	for _, c := range tail {
		c.DragBy(dx, dy)
	}
}

func (b *Baize) StartTailDrag(tail []*Card) {
	// hiding the mouse cursor creates flickering when tapping
	// ebiten.SetCursorMode(ebiten.CursorModeHidden)
	b.ApplyToTail(tail, (*Card).StartDrag)
}

func (b *Baize) StopTailDrag(tail []*Card) {
	// ebiten.SetCursorMode(ebiten.CursorModeVisible)
	b.ApplyToTail(tail, (*Card).StopDrag)
}

func (b *Baize) CancelTailDrag(tail []*Card) {
	// ebiten.SetCursorMode(ebiten.CursorModeVisible)
	b.ApplyToTail(tail, (*Card).CancelDrag)
}

func (b *Baize) powerMoves(pDraggingTo *Pile) int {
	// (1 + number of empty freecells) * 2 ^ (number of empty columns)
	// see http://ezinearticles.com/?Freecell-PowerMoves-Explained&id=104608
	// and http://www.solitairecentral.com/articles/FreecellPowerMovesExplained.html
	var emptyCells, emptyCols int
	for _, p := range b.piles {
		if p.Empty() {
			switch p.vtable.(type) {
			case *Cell:
				emptyCells++
			case *Tableau:
				if p.Label() == "" && p != pDraggingTo {
					// 'If you are moving into an empty column, then the column you are moving into does not count as empty column.'
					emptyCols++
				}
			}
		}
	}
	// 2^1 == 2, 2^0 == 1, 2^-1 == 0.5
	n := (1 + emptyCells) * util.Pow(2, emptyCols)
	// println(emptyCells, "emptyCells,", emptyCols, "emptyCols,", n, "powerMoves")
	return n
}

// DoingSafeCollect return true if we are doing safe collect
// and the safe ordinal to collect next
func (b *Baize) DoingSafeCollect() (bool, int) {
	if !TheGame.Settings.SafeCollect {
		return false, 0
	}
	if !b.script.SafeCollect() {
		return false, 0
	}
	var fs []*Pile = b.script.Foundations()
	if fs == nil {
		return false, 0
	}
	var f0 *Pile = fs[0]
	if f0 == nil {
		return false, 0
	}
	if f0.Label() != "A" {
		return false, 0 // eg Duchess
	}
	var lowest int = 99
	for _, f := range fs {
		if f.Empty() {
			// it's okay to collect aces and twos to start with
			return true, 2
		}
		var card *Card = f.Peek()
		if card.Ordinal() < lowest {
			lowest = card.Ordinal()
		}
	}
	return true, lowest + 1
}

// collectFromPile is a helper function for Collect2()
func (b *Baize) collectFromPile(pile *Pile) int {
	if pile == nil {
		return 0
	}
	var cardsMoved int = 0
	for _, fp := range b.script.Foundations() {
		for {
			var card *Card = pile.Peek()
			if card == nil {
				return cardsMoved
			}
			ok, _ := fp.vtable.CanAcceptTail([]*Card{card})
			if !ok {
				break // done with this foundation, try another
			}
			if ok, safeOrd := b.DoingSafeCollect(); ok {
				if card.Ordinal() > safeOrd {
					// can't toast here, collect all will create a lot of toasts
					// TheGame.UI.Toast("Glass", fmt.Sprintf("Unsafe to collect %s", card.String()))
					break // done with this foundation, try another
				}
			}
			MoveCard(pile, fp)
			b.AfterUserMove() // does an undoPush()
			b.AfterAfterUserMove()
			cardsMoved += 1
		}
	}
	return cardsMoved
}

// Collect2 should be exactly the same as the user tapping repeatedly on the
// waste, cell, reserve and tableau piles.
// nb there is no collecting to discard piles, they are optional and presence of
// cards in them does not signify a complete game.
// It's called Collect2 because it's the third or fourth rewrite of a
// basic and seemingly simple function.
func (b *Baize) Collect2() {
	for {
		var cardsMoved int = b.collectFromPile(b.script.Waste())
		for _, pile := range b.script.Cells() {
			cardsMoved += b.collectFromPile(pile)
		}
		for _, pile := range b.script.Reserves() {
			cardsMoved += b.collectFromPile(pile)
		}
		for _, pile := range b.script.Tableaux() {
			cardsMoved += b.collectFromPile(pile)
		}
		if cardsMoved == 0 {
			break
		}
	}
	// if ThePreferences.SafeCollect && b.script.SafeCollect() && b.fmoves > 0 {
	// 	TheGame.UI.Toast("Glass", "Not safe to collect card(s)")
	// }
}

func (b *Baize) MaxSlotX() int {
	var maxX int
	for _, p := range b.piles {
		if p.Slot().X > maxX {
			maxX = p.Slot().X
		}
	}
	return maxX
}

// ScaleCards calculates new width/height of cards and margins
// returns true if changes were made
func (b *Baize) ScaleCards() bool {

	// const (
	// 	DefaultRatio = 1.444
	// 	BridgeRatio  = 1.561
	// 	PokerRatio   = 1.39
	// 	OpsoleRatio  = 1.5556 // 3.5/2.25
	// )

	var OldWidth = CardWidth
	var OldHeight = CardHeight

	var maxX int = b.MaxSlotX()

	/*
		71 x 96 = 1:1.352 (Microsoft retro)
		140 x 190 = 1:1.357 (kenney, large)
		64 x 89 = 1:1.390 (official poker size)
		90 x 130 = 1:1.444 (nice looking scalable)
		89 x 137 = 1:1.539 (measured real card)
		57 x 89 = 1:1.561 (official bridge size)
	*/

	// Card padding is 10% of card height/width

	// if ThePreferences.FixedCards {
	// 	CardWidth = ThePreferences.FixedCardWidth
	// 	PilePaddingX = CardWidth / 10
	// 	CardHeight = ThePreferences.FixedCardHeight
	// 	PilePaddingY = CardHeight / 10
	// 	cardsWidth := PilePaddingX + CardWidth*(maxX+2)
	// 	LeftMargin = (b.WindowWidth - cardsWidth) / 2
	// } else {

	// "add" two extra piles and a LeftMargin to make a half-card-width border

	var slotWidth, slotHeight float64
	slotWidth = float64(b.WindowWidth) / float64(maxX+2)
	slotHeight = slotWidth * TheGame.Settings.CardRatio

	PilePaddingX = int(slotWidth / 10)
	CardWidth = int(slotWidth) - PilePaddingX
	PilePaddingY = int(slotHeight / 10)
	CardHeight = int(slotHeight) - PilePaddingY

	TopMargin = ui.ToolbarHeight + CardHeight/3
	LeftMargin = (CardWidth / 2) + PilePaddingX

	// CardDiagonal = math.Sqrt(math.Pow(float64(CardWidth), 2) + math.Pow(float64(CardHeight), 2))
	// }
	CardCornerRadius = float64(CardWidth) / 10.0 // same as lsol

	// if DebugMode {
	// 	if CardWidth != OldWidth || CardHeight != OldHeight {
	// 		log.Println("ScaleCards did something")
	// 	} else {
	// 		log.Println("ScaleCards did nothing")
	// 	}
	// }
	return CardWidth != OldWidth || CardHeight != OldHeight
}

func (b *Baize) PercentComplete() int {
	var pairs, unsorted, percent int
	for _, p := range b.piles {
		if p.Len() > 1 {
			pairs += p.Len() - 1
		}
		unsorted += p.vtable.UnsortedPairs()
	}
	// TheGame.UI.SetMiddle(fmt.Sprintf("%d/%d", pairs-unsorted, pairs))
	percent = (int)(100.0 - util.MapValue(float64(unsorted), 0, float64(pairs), 0.0, 100.0))
	return percent
}

func (b *Baize) Recycles() int {
	return b.recycles
}

func (b *Baize) SetRecycles(recycles int) {
	b.recycles = recycles
	b.setFlag(dirtyPileBackgrounds) // recreate Stock placeholder
}

func (b *Baize) UpdateToolbar() {
	TheGame.UI.EnableWidget("toolbarUndo", len(b.undoStack) > 1)
	TheGame.UI.EnableWidget("toolbarCollect", b.fmoves > 0)
}

func (b *Baize) UpdateStatusbar() {
	if b.script.Stock().Hidden() {
		TheGame.UI.SetStock(-1)
	} else {
		TheGame.UI.SetStock(b.script.Stock().Len())
	}
	if b.script.Waste() == nil {
		TheGame.UI.SetWaste(-1) // previous variant may have had a waste, and this one does not
	} else {
		TheGame.UI.SetWaste(b.script.Waste().Len())
	}
	// if DebugMode {
	// 	TheGame.UI.SetMiddle(fmt.Sprintf("MOVES: %d,%d", b.moves, b.fmoves))
	// }
	TheGame.UI.SetMiddle(fmt.Sprintf("MOVES: %d", len(b.undoStack)-1))
	TheGame.UI.SetPercent(b.PercentComplete())
}

func (b *Baize) UpdateDrawers() {
	TheGame.UI.EnableWidget("restartDeal", len(b.undoStack) > 1)
	TheGame.UI.EnableWidget("gotoBookmark", b.bookmark > 0)
}

func (b *Baize) Conformant() bool {
	for _, p := range b.piles {
		if !p.vtable.Conformant() {
			return false
		}
	}
	return true
}

func (b *Baize) Complete() bool {
	return b.script.Complete()
}

// Layout implements ebiten.Game's Layout.
func (b *Baize) Layout(outsideWidth, outsideHeight int) (int, int) {

	if outsideWidth == 0 || outsideHeight == 0 {
		log.Println("Baize.Layout called with zero dimension")
		return outsideWidth, outsideHeight
	}

	if DebugMode && (outsideWidth != b.WindowWidth || outsideHeight != b.WindowHeight) {
		log.Println("Window resize to", outsideWidth, outsideHeight)
	}

	if outsideWidth != b.WindowWidth {
		b.setFlag(dirtyWindowSize | dirtyCardSizes | dirtyPileBackgrounds | dirtyPilePositions | dirtyCardPositions)
		b.WindowWidth = outsideWidth
	}
	if outsideHeight != b.WindowHeight {
		b.setFlag(dirtyWindowSize | dirtyCardPositions)
		b.WindowHeight = outsideHeight
	}

	if b.dirtyFlags != 0 {
		if b.flagSet(dirtyCardSizes) {
			if b.ScaleCards() {
				if DebugMode {
					log.Printf("ScaleCards %dx%d", CardWidth, CardHeight)
				}
				b.setFlag(dirtyCardImages | dirtyPilePositions | dirtyPileBackgrounds)
			}
			// b.clearFlag(dirtyCardSizes)
		}
		if b.flagSet(dirtyCardImages) {
			CreateCardImages()
			// b.clearFlag(dirtyCardImages)
		}
		if b.flagSet(dirtyPilePositions) {
			for _, p := range b.piles {
				p.SetBaizePos(image.Point{
					X: LeftMargin + (p.Slot().X * (CardWidth + PilePaddingX)),
					Y: TopMargin + (p.Slot().Y * (CardHeight + PilePaddingY)),
				})
			}
			// b.clearFlag(dirtyPilePositions)
		}
		if b.flagSet(dirtyPileBackgrounds) {
			if !(CardWidth == 0 || CardHeight == 0) {
				for _, p := range b.piles {
					if !p.Hidden() {
						p.img = p.vtable.Placeholder()
					}
				}
			}
			// b.clearFlag(dirtyPileBackgrounds)
		}
		if b.flagSet(dirtyWindowSize) {
			TheGame.UI.Layout(outsideWidth, outsideHeight)
			// b.clearFlag(dirtyWindowSize)
		}
		if b.flagSet(dirtyCardPositions) {
			for _, p := range b.piles {
				p.Scrunch()
			}
			// b.clearFlag(dirtyCardPositions)
		}
		b.dirtyFlags = 0
	}

	return outsideWidth, outsideHeight
}

// Update the baize state (transitions, user input)
func (b *Baize) Update() error {

	if b.stroke == nil {
		input.StartStroke(b) // this will set b.stroke when "start" received
	} else {
		b.stroke.Update()
		if b.stroke.IsReleased() || b.stroke.IsCancelled() {
			b.stroke = nil
		}
	}

	// {
	// 	var x, y int = ebiten.CursorPosition()
	// 	// x, y will be 0,0 on mobiles, and before main loop starts
	// 	if !(x == 0 && y == 0) {
	// 		b.hotCard = b.FindLowestCardAt(image.Point{x, y})
	// 	}
	// }

	for _, p := range b.piles {
		p.Update()
	}

	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if inpututil.IsKeyJustReleased(k) {
			Execute(k)
		}
	}

	return nil
}

// Draw renders the baize into the screen
func (b *Baize) Draw(screen *ebiten.Image) {

	screen.Fill(ExtendedColors[TheGame.Settings.BaizeColor])

	for _, p := range b.piles {
		p.Draw(screen)
	}
	for _, p := range b.piles {
		p.DrawStaticCards(screen)
	}
	for _, p := range b.piles {
		p.DrawAnimatingCards(screen)
	}
	for _, p := range b.piles {
		p.DrawDraggingCards(screen)
	}
	// if b.hotCard != nil {
	// 	b.hotCard.Draw(screen)
	// }

	// if DebugMode {
	// var ms runtime.MemStats
	// runtime.ReadMemStats(&ms)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %v, Alloc %v, NumGC %v", ebiten.CurrentTPS(), ms.Alloc, ms.NumGC))
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("%v %v", b.bookmark, len(b.undoStack)))
	// bounds := screen.Bounds()
	// ebitenutil.DebugPrint(screen, bounds.String())
	// }
}
