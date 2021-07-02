package sol

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/ui"
	"oddstream.games/gosol/util"
)

// PilePositionType is the position, in logical coords, of a Pile on the Baize
// the leftmost Pile will have X=0, the topmost will have Y=0
// Piles can be hidden by setting X or Y to negative values
type PilePositionType int

// BaizeStateType is either Virgin, Started or Complete
type BaizeStateType int

// Virgin, Started or Complete
const (
	Virgin BaizeStateType = iota
	Started
	Complete
)

// Baize object describes the baize
type Baize struct {
	Piles           []*Pile
	stock           *Pile
	Variant         string
	UndoStack       []SaveableBaize
	SavedPosition   int
	totalCards      int
	percentComplete int
	movableCards    int
	State           BaizeStateType
	stroke          *input.Stroke
	ui              *ui.UI
	commandTable    map[ebiten.Key]func()
	DragOffsetX     int // value of current horz drag
	DragOffsetBaseX int // value last-used horz drag
	DragOffsetY     int // value of current vertical drag
	DragOffsetBaseY int // value last-used vertical drag
	WindowWidth     int // the most recent window width given to Layout
	OldWindowWidth  int // the window width last used to scale baize and cards
	WindowHeight    int // the most recent window height given to Layout
	OldWindowHeight int // the window height last used to scale baize and cards
}

//--+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8

// NewBaize is the factory func for the single Baize object
func NewBaize() *Baize {
	// TheUserData may have been injected from command line flags
	// log.Printf("%v", TheUserData)

	// bug lurking here; scalables start at 71x96, which is the size needed for CardBackPicker
	CreateCardImages() // sets global TheCIP (sorry)

	TheBaize = &Baize{Variant: ThePreferences.Variant}
	TheBaize.ui = ui.New(TheBaize.Execute)
	TheBaize.commandTable = map[ebiten.Key]func(){
		ebiten.KeyN:      TheBaize.NewGame,
		ebiten.KeyR:      TheBaize.RestartGame,
		ebiten.KeyU:      TheBaize.Undo,
		ebiten.KeyS:      TheBaize.SavePosition,
		ebiten.KeyL:      TheBaize.LoadPosition,
		ebiten.KeyC:      TheBaize.Collect,
		ebiten.KeyF:      TheBaize.ShowVariantPicker,
		ebiten.KeyF1:     TheBaize.ShowRules,
		ebiten.KeyF2:     TheBaize.ShowCardBackPicker,
		ebiten.KeyF3:     TheBaize.ShowSettingsDrawer,
		ebiten.KeyF4:     TheStatistics.ShowStatistics,
		ebiten.KeyF5:     TheBaize.StartSpinning,
		ebiten.KeyF6:     TheBaize.StopSpinning,
		ebiten.KeyF7:     func() { TheBaize.ui.ShowFAB("star", ebiten.KeyN) },
		ebiten.KeyF8:     func() { TheBaize.ui.HideFAB() },
		ebiten.KeyMenu:   TheBaize.ui.ToggleNavDrawer,
		ebiten.KeyEscape: TheBaize.ui.HideActiveDrawer,
		ebiten.KeyX:      TheBaize.Exit,
	}

	if NoGameLoad || !TheBaize.LoadVariant(TheBaize.Variant) {
		TheBaize.NewVariant(TheBaize.Variant)
	}

	return TheBaize // ugly global-setting kludge
}

// RecallCardsToStock without changing variant or seed
func (b *Baize) RecallCardsToStock() {

	for i, p := range b.Piles {
		p.Reset() // stock needs resetting, too
		if i == 0 {
			continue // 0th pile is the stock
		}
		if !p.Empty() {
			b.stock.MoveCards(p.Cards[0])
		}
		// if !p.Empty() {
		// 	log.Fatal(p.Class, " still contains ", p.CardCount(), " cards")
		// }
	}

	// if DebugMode {
	// 	println("cards recalled to stock, now contains", b.stock.CardCount(), "cards")
	// 	for _, c := range b.stock.Cards {
	// 		if !c.Prone() {
	// 			log.Fatal("face up card found in stock")
	// 		}
	// 		if c.owner != b.stock {
	// 			log.Fatal("card in stock belongs to", c.owner.Class)
	// 		}
	// 	}
	// }
}

// NewGame resets Baize and restarts current variant with a new seed
func (b *Baize) NewGame() {

	if DebugMode {
		defer util.Duration(time.Now(), "Baize.NewGame")
	}

	if b.State == Started {
		TheStatistics.recordLostGame(b.Variant, b.percentComplete)
	}
	b.RecallCardsToStock()
	b.Shuffle()

	// reset the Baize
	b.UndoStack = nil
	b.SavedPosition = 0
	b.State = Virgin
	b.stroke = nil

	b.dealCards()
	b.UndoPush()
	TheStatistics.welcomeToast(b.Variant)
}

// NewVariant resets Baize and starts a new game with a new variant and seed
func (b *Baize) NewVariant(v string) {

	if DebugMode {
		defer util.Duration(time.Now(), "Baize.NewVariant")
	}

	if b.State == Started {
		TheStatistics.recordLostGame(b.Variant, b.percentComplete)
	}

	// reset Baize
	b.Piles = b.Piles[:0]
	b.UndoStack = nil
	b.SavedPosition = 0
	b.State = Virgin
	b.stroke = nil

	// switch to new variant
	b.Variant = findVariant(v) // v is now invalid because AKA
	ThePreferences.Variant = b.Variant
	b.BuildVariant(b.Variant)
	b.ui.SetTitle(b.Variant)

	b.OldWindowWidth, b.OldWindowHeight = 0, 0 // force a rescale
	b.Scale()                                  // now we know number of piles and can discover window width; scale the cards before creating them

	b.CreateStock()
	b.Shuffle()

	b.dealCards()
	b.UndoPush()
	TheStatistics.welcomeToast(b.Variant)
}

// LoadVariant tries to load a game from json resets Baize and continues an old game
func (b *Baize) LoadVariant(v string) bool {

	undoStack := LoadUndoStack(v)
	if undoStack == nil {
		return false
	}

	if DebugMode {
		defer util.Duration(time.Now(), "Baize.LoadVariant")
	}

	b.Variant = findVariant(v) // v is now invalid because AKA
	ThePreferences.Variant = b.Variant
	b.BuildVariant(b.Variant)
	b.ui.SetTitle(b.Variant)
	b.UndoStack = undoStack

	sav, ok := b.UndoPop() // removes extra pushed state
	if !ok {
		log.Panic("error popping extra state from undo stack")
	}

	b.OldWindowWidth, b.OldWindowHeight = 0, 0 // force a rescale
	b.Scale()                                  // now we know number of piles and can discover window width; scale the cards before creating them

	b.CreateStock()

	b.UpdateFromSaveable(sav)

	b.UndoPush()
	TheStatistics.welcomeToast(b.Variant)

	if b.State == Complete {
		b.ui.ShowFAB("star", ebiten.KeyN)
	}

	return true
}

func (b *Baize) dealCards() {
	sound.Play("Fan")
	for _, p := range b.Piles {
		deal := p.GetStringAttribute("Deal")
		if deal == "" {
			continue
		}
		cards := parseCardsFromDeal(b.stock, deal)
		for _, c := range cards {
			p.Push(c)
		}
		if bury, ok := p.GetIntAttribute("Bury"); ok {
			p.BuryCards(bury)
		}
		if disinter, ok := p.GetIntAttribute("Disinter"); ok {
			p.DisinterCards(disinter)
		}
	}

	b.AutoMoves()

	if DebugMode {
		if b.stock.Hidden() {
			log.Println(b.stock.CardCount(), "cards remaining in hidden stock")
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

// func (b *Baize) findPilePrefix(cls string) *Pile {
// 	for _, p := range b.Piles {
// 		if strings.HasPrefix(p.Class, cls) {
// 			return p
// 		}
// 	}
// 	return nil
// }

// findPileAt finds the pile under the mouse click or touch
func (b *Baize) findPileAt(x, y int) *Pile {
	for _, p := range b.Piles {
		if util.InRect(x, y, p.FannedScreenRect) {
			return p
		}
	}
	return nil
}

// findCardAt finds the tile under the mouse click or touch
func (b *Baize) findCardAt(x, y int) *Card {
	// go backwards, for King Albert's overlapping reserve piles
	for j := len(b.Piles) - 1; j >= 0; j-- {
		p := b.Piles[j]
		for i := p.CardCount() - 1; i >= 0; i-- {
			c := p.Cards[i]
			if util.InRect(x, y, c.ScreenRect) {
				return c
			}
		}
	}
	return nil
}

/* retired by Driver
// CardTapped is called when a card has been tapped
func (b *Baize) CardTapped(c *Card) {

	// println("card",c.ID.String(),"tapped")

	// if c.Transitioning() || c.Flipping() {
	// 	log.Println("cannot tap an animating card ", c.String())
	// 	return
	// }

	if !ThePreferences.SingleTap {
		return
	}

	var sourcePile *Pile = c.owner
	var tail []*Card = sourcePile.makeTail(c)
	var targetClass string = sourcePile.GetStringAttribute("Target")
	var anyCardsMoved bool = false

	switch sourcePile.Class {
	case "Stock":
		// Tap on a Stock card to send one or more cards to Waste, Golf
		if targetClass == "" {
			targetClass = "Waste"
		}
		cardsToMove, ok := sourcePile.GetIntAttribute("CardsToMove")
		if !ok || cardsToMove == 0 {
			cardsToMove = 1
		}
		for _, p := range b.Piles {
			if targetClass == p.Class {
				// println("found a", p.Class)
				// Waste, Golf should always accept a card from Stock, should don't need to
				// check if targetClass can accept Stock card
				// if p.CanAcceptCard(c) {
				// println(p.Class, "can accept", c.ID.String())
				for cardsToMove > 0 && c != nil {
					cardsToMove--
					p.MoveCards(c)
					c = sourcePile.Peek()
					anyCardsMoved = true
				}
				// }
				if c == nil {
					break
				}
			}
		}
	case "StockSpider":
		var tabCards, tabPiles, emptyTabPiles int
		{ // scope for p *Pile (in case of silent Go variable shadowing)
			for _, p := range b.Piles {
				if p.Class == "Tableau" {
					tabPiles++
					if p.Empty() {
						emptyTabPiles++
					} else {
						tabCards += p.CardCount()
					}
				}
			}
		}
		if tabCards >= tabPiles && emptyTabPiles > 0 {
			if b.stock.CardCount() > emptyTabPiles {
				TheBaize.ui.Toast("All empty tableaux must be filled before dealing a new row")
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
				p.MoveCards(c)
				anyCardsMoved = true
				c.SetProne(false)
				c = sourcePile.Peek()
			}
			if c == nil {
				break
			}
		}
	case "Tableau":
		// try to move a conformant run of 13 cards to a spider foundation
		if sourcePile.Spider() && len(tail) == 13 && isTailConformant(sourcePile.Build, sourcePile.Flags, tail) {
			for _, p := range b.Piles {
				if p.Class == "Foundation" && p.Spider() {
					if p.CanAcceptTail(tail, false) {
						p.MoveCards(c)
						anyCardsMoved = true
						break
					}
				}
			}
		}
		fallthrough
	case "Waste", "Cell", "Reserve", "Golf":
		// try to move a top or single tapped card to a Foundation
		if !anyCardsMoved && c == sourcePile.Peek() {
			for _, p := range b.Piles {
				if p.Class == "Foundation" {
					if p.CanAcceptTail([]*Card{c}, false) {
						p.MoveCards(c)
						anyCardsMoved = true
						break
					}
				}
			}
		}
		// else try to move card to the longest Tableau or Golf
		if !anyCardsMoved {
			if len(tail) == 0 {
				log.Panic("CardTapped empty tail")
			}
			var pLongest *Pile
			for _, p := range b.Piles {
				if p == c.owner {
					continue
				}
				if !(p.Class == "Tableau" || p.Class == "Golf") {
					continue
				}
				if p.CanAcceptTail(tail, false) {
					if pLongest == nil || p.CardCount() > pLongest.CardCount() {
						pLongest = p
					}
				}
			}
			if pLongest != nil {
				pLongest.MoveCards(c)
				anyCardsMoved = true
			}
		}
	case "Foundation":
		TheBaize.ui.Toast("You cannot move cards from a Foundation")
	default:
		println("clueless when tapping on a", sourcePile.Class, "card")
	}

	if anyCardsMoved {
		b.AfterUserMove()
	} else {
		if c != nil {
			c.Shake()
		}
	}

}
*/

// AutoMoves performs post user-moves
func (b *Baize) AutoMoves() {

	for _, p := range b.Piles {
		if p.Class == "Foundation" && p.CardCount() == 1 {
			if afp := p.GetStringAttribute("AcceptFirstPush"); afp != "" {
				ord := p.Peek().Ordinal()
				for _, fp := range b.Piles {
					if fp.Class == "Foundation" {
						fp.SetAccept(ord)
					}
				}
			}
		}
	}

	for _, p := range b.Piles {
		if p.Empty() {
			if aff := p.GetStringAttribute("AutoFillFrom"); aff != "" {
				affPiles := strings.Split(aff, ",")
				for _, srcPile := range affPiles {
					if src := b.findPile(srcPile); src != nil {
						if c := src.Peek(); c != nil {
							p.MoveCards(c)
							break
						}
					}
				}
			}
		}
	}

}

// AfterUserMove runs after the user has made a move;
// perform any auto moves (on behalf of user), test for game complete, push state onto undo stack
func (b *Baize) AfterUserMove() {

	b.AutoMoves()

	//

	switch b.State {
	case Virgin:
		b.State = Started
		b.ui.Toast(fmt.Sprintf("%s started", b.Variant))
	case Started:
		if b.Complete() {
			b.State = Complete
			sound.Play("Complete")
			TheStatistics.recordWonGame(b.Variant, len(b.UndoStack)-1)
			TheStatistics.wonToast(b.Variant, len(b.UndoStack)-1)
			b.ui.ShowFAB("star", ebiten.KeyN)
			b.StartSpinning()
		} else if b.Conformant() {
			b.ui.ShowFAB("done_all", ebiten.KeyC)
		} else {
			b.ui.HideFAB()
		}
	case Complete:
		log.Println("what are we doing here?")
	}

	//

	var oldChecksum, newChecksum uint32
	var ok bool

	//

	if len(b.UndoStack) == 0 {
		log.Panic("undo stack is empty in AfterUserMove()")
	} else {
		oldChecksum, ok = b.UndoPeekChecksum()
		if !ok {
			log.Panic("error peeking undo stack checksum")
		}
	}
	newChecksum = b.Checksum()
	// println(oldChecksum, newChecksum)
	if oldChecksum != newChecksum {
		b.UndoPush()
		if b.State == Started && b.movableCards == 0 && b.stock.localRecycles == 0 {
			b.ui.Toast("No movable cards")
		}
	} else {
		log.Println("not pushing to undo because checksums match")
	}

}

func (b *Baize) largestIntersection(c *Card) *Pile {
	var largestArea int = 0
	var pile *Pile = nil
	cx0, cy0, cx1, cy1 := c.BaizeRect()
	for _, p := range b.Piles {
		if p == c.owner {
			continue
		}
		px0, py0, px1, py1 := p.FannedBaizeRect()
		area := util.OverlapArea(cx0, cy0, cx1, cy1, px0, py0, px1, py1)
		// can't test for AcceptTail here as it would filter out warning toasts
		if area > largestArea /*&& p.CanAcceptTail(c.owner.Tail, false)*/ {
			largestArea = area
			pile = p
		}
	}
	return pile
}

func (b *Baize) calcPercentComplete() int {
	var sorted, unsorted int
	for _, p := range b.Piles {
		sorted, unsorted = p.CountSortedAndUnsorted(sorted, unsorted)
	}
	return int(util.MapValue(float64(sorted)-float64(unsorted), float64(-b.totalCards), float64(b.totalCards), 0, 100))
}

// StartDrag return true if the Baize can be dragged (vscrolled)
func (b *Baize) StartDrag() bool {
	return true
}

// DragBy move (vscroll) the Baize by dragging it
func (b *Baize) DragBy(dx, dy int) {
	b.DragOffsetX = b.DragOffsetBaseX + dx
	if b.DragOffsetX > 0 {
		b.DragOffsetX = 0 // DragOffsetX should only ever be 0 or -ve
	}
	b.DragOffsetY = b.DragOffsetBaseY + dy
	if b.DragOffsetY > 0 {
		b.DragOffsetY = 0 // DragOffsetY should only ever be 0 or -ve
	}
}

// StopDrag stop dragging the Baize
func (b *Baize) StopDrag() {
	// remember the amount of drag so the next drag starts from here
	b.DragOffsetBaseX = b.DragOffsetX
	b.DragOffsetBaseY = b.DragOffsetY
}

// StartSpinning tells all the cards to start spinning
func (b *Baize) StartSpinning() {
	for _, p := range b.Piles {
		p.ApplyToCards((*Card).StartSpinning)
	}
}

// StopSpinning tells all the cards to stop spinning and return to their upright position
// debug only
func (b *Baize) StopSpinning() {
	for _, p := range b.Piles {
		if p.Empty() {
			continue
		}
		p.RepushAllCards()
	}
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (b *Baize) NotifyCallback(v input.StrokeEvent) {
	switch v.Event {
	case "start":
		// try UI Container > Card > Pile > Baize
		b.stroke = v.Stroke
		if con := b.ui.FindContainerAt(v.X, v.Y); con != nil {
			// println("found container")
			if con.StartDrag(b.stroke) {
				b.stroke.SetDraggedObject(con)
			} else {
				v.Stroke.Cancel()
			}
		} else {
			if c := b.findCardAt(v.X, v.Y); c != nil {
				b.stroke.SetDraggedObject(c)
				if !c.owner.StartDrag(c) {
					// log.Println("cancel stroke because drag not allowed")
					v.Stroke.Cancel()
				}
				// sound.Play("Place")
			} else {
				if p := b.findPileAt(v.X, v.Y); p != nil {
					// we can't really drag piles, but nevertheless...
					b.stroke.SetDraggedObject(p)
				} else {
					if b.StartDrag() {
						// println("starting baize drag")
						b.stroke.SetDraggedObject(b)
					} else {
						// log.Println("cancel stroke because not over a card")
						v.Stroke.Cancel()
					}
				}
			}
		}
	case "move":
		if v.Stroke.DraggedObject() == nil {
			println("*** move stroke with nil dragged object ***")
			break
		}
		switch v.Stroke.DraggedObject().(type) { // type switch
		case ui.Container:
			con := v.Stroke.DraggedObject().(ui.Container)
			con.DragBy(v.Stroke.PositionDiff())
		case *Card:
			c := v.Stroke.DraggedObject().(*Card)
			c.owner.DragTailBy(v.Stroke.PositionDiff())
		case *Pile:
			// do nothing
		case *Baize:
			// println("baize drag")
			b2 := v.Stroke.DraggedObject().(*Baize)
			if b2 != b {
				println("baize drag - something has gone terribly wrong")
			}
			b2.DragBy(v.Stroke.PositionDiff())
		default:
			println("*** unknown move dragging object ***")
		}
	case "stop":
		if v.Stroke.DraggedObject() == nil {
			println("*** stop stroke with nil dragged object ***")
			break
		}
		switch v.Stroke.DraggedObject().(type) { // type switch
		case ui.Container:
			con := v.Stroke.DraggedObject().(ui.Container)
			con.StopDrag()
		case *Card:
			c := v.Stroke.DraggedObject().(*Card)
			p := b.largestIntersection(c)
			if p == nil || p == c.owner {
				c.owner.CancelDrag(c)
			} else {
				if len(c.owner.Tail) == 0 {
					println("*** stop dragging card - empty tail ***")
					c.owner.CancelDrag(c)
				} else {
					if ok, err := p.driver.CanAcceptTail(c.owner.Tail); ok {
						c.owner.StopDrag(c) // this makes the Tail = nil
						if c.owner == b.stock && p.Class == "Waste" {
							// special case: dragging a card from Stock to Waste in Canfield, Klondike (Draw Three)
							cardsToMove, ok := c.owner.GetIntAttribute("CardsToMove")
							if !ok || cardsToMove == 0 {
								cardsToMove = 1
							}
							for cardsToMove > 0 && !b.stock.Empty() {
								cardsToMove--
								p.MoveCards(b.stock.Peek()) // this reassigns c.owner to p
							}
						} else {
							p.MoveCards(c)
						}
						b.AfterUserMove()
					} else {
						c.owner.CancelDrag(c)
						if err != nil {
							TheBaize.ui.Toast(err.Error())
						}
					}
				}
			}
		case *Pile:
			// p := v.Stroke.DraggedObject().(*Pile)
			// println("stop dragging pile", p.Class)
			// do nothing
		case *Baize:
			// println("stop dragging baize")
			b2 := v.Stroke.DraggedObject().(*Baize)
			if b2 != b {
				println("baize drag - something has gone terribly wrong")
			}
			b2.StopDrag()
		default:
			println("*** stop dragging unknown object ***")
		}
	case "cancel":
		if v.Stroke.DraggedObject() == nil {
			println("*** cancel stroke with nil dragged object ***")
			break
		}
		switch v.Stroke.DraggedObject().(type) { // type switch
		case ui.Container:
			con := v.Stroke.DraggedObject().(ui.Container)
			con.StopDrag()
		case *Card:
			c := v.Stroke.DraggedObject().(*Card)
			c.owner.CancelDrag(c)
		case *Pile:
			// p := v.Stroke.DraggedObject().(*Pile)
			// println("stop dragging pile", p.Class)
			// do nothing
		case *Baize:
			// println("stop dragging baize")
			b2 := v.Stroke.DraggedObject().(*Baize)
			if b2 != b {
				println("baize drag - something has gone terribly wrong")
			}
			b2.StopDrag()
		default:
			println("*** cancel dragging unknown object ***")
		}
	case "tap":
		// println("Baize.NotifyCallback() tap", v.X, v.Y)
		// a tap outside any open ui drawer (ie on the baize) closes the drawer
		if con := b.ui.VisibleDrawer(); con != nil && !util.InRect(v.X, v.Y, con.Rect) {
			con.Hide()
		} else if con := b.ui.FindContainerAt(v.X, v.Y); con == nil {
			// not a tap on a UI container, so must be on a pile or a card
			c := b.findCardAt(v.X, v.Y)
			// we've received a tap, so cancel any stroke that has started
			if b.stroke != nil {
				// println("cancel stroke because tap")
				if c != nil {
					c.owner.CancelDrag(c)
				}
				b.stroke.Cancel()
			}
			if c != nil {
				if ThePreferences.SingleTap {
					_, err := c.owner.driver.CardTapped(c)
					if err != nil {
						TheBaize.ui.Toast(err.Error())
					}
				}
			} else {
				if p := b.findPileAt(v.X, v.Y); p != nil {
					p.driver.Tapped()
				}
			}
		}
	default:
		println("*** unknown stroke event", v.Event)
		// case "cancel":
		// 	// c := v.Stroke.DraggingObject().(*Card)
		// 	c := v.Stroke.DraggedCard()
		// 	if c != nil {
		// 		c.owner.CancelDrag(c)
		// 	}
		// 	b.stroke = nil
	}
}

// ScaleCards calculates new width/height of cards and margins
func (b *Baize) ScaleCards() {

	// const (
	// 	DefaultRatio = 1.444
	// 	BridgeRatio  = 1.561
	// 	PokerRatio   = 1.39
	// 	OpsoleRatio  = 1.5556 // 3.5/2.25
	// )

	var maxX PilePositionType
	for _, p := range b.Piles {
		if p.X > maxX {
			maxX = p.X
		}
	}

	// "add" two extra piles and a LeftMargin to make a half-card-width border

	/*
		71 x 96 = 1:1.352 (Microsoft retro)
		64 x 89 = 1:1.390 (official poker size)
		90 x 130 = 1:1.444 (nice looking scalable)
		89 x 137 = 1:1.539 (measured real card)
		57 x 89 = 1:1.561 (official bridge size)
	*/

	// Card gap is 10% of card width
	/*
		switch TheUserData.CardStyle {
		default:
			slotWidth := float64(b.WindowWidth) / float64(maxX+2)
			PilePaddingX = int(slotWidth / 10)
			CardWidth = int(slotWidth) - PilePaddingX
			slotHeight := slotWidth * 1.5
			PilePaddingY = int(slotHeight / 10)
			CardHeight = int(slotHeight) - PilePaddingY
			LeftMargin = (CardWidth / 2) + PilePaddingX
		case "fixed":
			CardWidth = 70
			PilePaddingX = 7
			CardHeight = 70 * 1.5 // 105
			PilePaddingY = 10
			cardsWidth := int(PilePositionType(PilePaddingX+CardWidth) * (maxX + 1)) // add 1 for half width card margin
			LeftMargin = (b.WindowWidth - cardsWidth) / 2
		case "retro":
			CardWidth = 71
			PilePaddingX = 7
			CardHeight = 96
			PilePaddingY = 10
			cardsWidth := int(PilePositionType(PilePaddingX+CardWidth) * (maxX + 1)) // add 1 for half width card margin
			LeftMargin = (b.WindowWidth - cardsWidth) / 2
		}
		log.Printf("%s card size %dx%d", TheUserData.CardStyle, CardWidth, CardHeight)
	*/
	if ThePreferences.FixedCards {
		// the retro cards dimensions in the front and back sprite sheets are width 71, height 96
		// but that makes Modern cards that look too wide
		if ThePreferences.RetroCards {
			CardWidth = 71
			PilePaddingX = 7
			CardHeight = 96
			PilePaddingY = 10
		} else {
			CardWidth = 70
			PilePaddingX = 7
			CardHeight = 70 * 1.5 // 105
			PilePaddingY = 10
		}
		cardsWidth := int(PilePositionType(PilePaddingX+CardWidth) * (maxX + 1)) // add 1 for half width card margin
		LeftMargin = (b.WindowWidth - cardsWidth) / 2
	} else {
		slotWidth := float64(b.WindowWidth) / float64(maxX+2)
		PilePaddingX = int(slotWidth / 10)
		CardWidth = int(slotWidth) - PilePaddingX
		slotHeight := slotWidth * 1.5
		PilePaddingY = int(slotHeight / 10)
		CardHeight = int(slotHeight) - PilePaddingY
		LeftMargin = (CardWidth / 2) + PilePaddingX
	}
	TopMargin = 48 + CardHeight/3

}

// Scale resizes piles, cards (inc shadow image), fonts and then repositions piles and cards
func (b *Baize) Scale() {

	// on startup, b.OldWindowWidth will be 0 so scalables will be built
	if b.WindowWidth == b.OldWindowWidth && b.WindowHeight == b.OldWindowHeight {
		return
	}

	b.ScaleCards()

	CreateCardImages()

	for _, p := range b.Piles {
		p.CreateBackgroundImage()
		if p.Empty() {
			continue
		}
		p.ApplyToCards((*Card).RefreshFaceImage)
		p.RepushAllCards()
	}

	b.OldWindowWidth, b.OldWindowHeight = b.WindowWidth, b.WindowHeight
}

// Layout implements ebiten.Game's Layout.
func (b *Baize) Layout(outsideWidth, outsideHeight int) (int, int) {

	b.WindowWidth, b.WindowHeight = outsideWidth, outsideHeight

	b.Scale()

	b.ui.Layout(outsideWidth, outsideHeight)

	return outsideWidth, outsideHeight

}

// Update the baize state (transitions, user input)
func (b *Baize) Update() error {

	if b.stroke == nil {
		input.StartStroke(b) // this will set b.stroke when "start" received
	} else {
		b.stroke.Update()
		if b.stroke == nil || b.stroke.IsReleased() || b.stroke.IsCancelled() {
			b.stroke = nil
		}
	}

	for _, p := range b.Piles {
		p.Update()
	}

	b.ui.Update()

	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if inpututil.IsKeyJustReleased(k) {
			b.Execute(k)
		}
	}

	// if _, yoff := ebiten.Wheel(); yoff != 0 {
	// 	b.DragBy(0, int(yoff*24))
	// }

	return nil
}

// Draw renders the baize into the screen
func (b *Baize) Draw(screen *ebiten.Image) {

	screen.Fill(colorBaize)

	for _, p := range b.Piles {
		p.Draw(screen)
	}
	for _, p := range b.Piles {
		p.DrawStaticCards(screen)
	}
	for _, p := range b.Piles {
		p.DrawTransitioningCards(screen)
	}
	for _, p := range b.Piles {
		p.DrawFlippingCards(screen)
	}

	b.ui.Draw(screen)

	if DebugMode {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("NumGC %v, State %d, Movable %d", ms.NumGC, b.State, b.movableCards))
	}
}

// Exit this app
func (b *Baize) Exit() {

	if !NoGameSave {
		b.Save()
	}

	if runtime.GOARCH != "wasm" {
		ThePreferences.WindowX, ThePreferences.WindowY = ebiten.WindowPosition()
		ThePreferences.WindowWidth, ThePreferences.WindowHeight = ebiten.WindowSize()
		ThePreferences.WindowMaximized = ebiten.IsWindowMaximized()
	}

	ThePreferences.Save()

	if runtime.GOARCH != "wasm" {
		os.Exit(0)
	}

}
