package sol

import (
	"fmt"
	"image"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"oddstream.games/gosol/input"
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
	Variant         string
	Seed            int64
	UndoStack       []SaveableBaize
	SavedPosition   int
	totalCards      int
	percentComplete int
	movableCards    int
	State           BaizeStateType
	stroke          *input.Stroke
	input           *input.Input
	ui              *ui.UI
	stock           *Pile
	commandTable    map[ebiten.Key]func()
	DragOffsetY     int
	DragOffsetBaseY int
	OldWindowWidth  int
}

// NewBaize is the factory func for the single Baize object
func NewBaize() *Baize {
	// TheUserData may have been injected from command line flags
	// log.Printf("%v", TheUserData)

	// bug lurking here; scalables start at 71x96, which is the size needed for CardBackPicker
	CreateScalables() // sets global TheCIP (sorry)

	TheBaize = &Baize{Variant: TheUserData.Variant, Seed: time.Now().UnixNano()}
	TheBaize.input = input.NewInput()
	TheBaize.input.Add(TheBaize) // TheBaize.NotifyCallback() will receive input event notifications
	TheBaize.ui = ui.New(TheBaize.input)
	TheBaize.commandTable = map[ebiten.Key]func(){
		ebiten.KeyN:      TheBaize.NewGame,
		ebiten.KeyR:      TheBaize.RestartGame,
		ebiten.KeyU:      TheBaize.Undo,
		ebiten.KeyS:      TheBaize.SavePosition,
		ebiten.KeyL:      TheBaize.LoadPosition,
		ebiten.KeyC:      TheBaize.Collect,
		ebiten.KeyF1:     TheBaize.ShowRules,
		ebiten.KeyF:      TheBaize.ShowVariantPicker,
		ebiten.KeyF2:     TheBaize.ShowCardBackPicker,
		ebiten.KeyF3:     TheBaize.ShowSettingsDrawer,
		ebiten.KeyI:      TheBaize.ShowInfo,
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

	for _, p := range b.Piles {
		p.Reset() // stock needs resetting, too
		if p == b.stock {
			continue
		}
		if p.CardCount() > 0 {
			b.MoveCards(p.Cards[0], b.stock)
		}
		// if p.CardCount() != 0 {
		// 	log.Fatal(p.Class, " still contains ", p.CardCount(), " cards")
		// }
	}

	if DebugMode {
		println("cards recalled to stock, now contains", b.stock.CardCount(), "cards")
		for _, c := range b.stock.Cards {
			if !c.Prone() {
				log.Fatal("face up card found in stock")
			}
			if c.owner != b.stock {
				log.Fatal("card in stock belongs to", c.owner.Class)
			}
		}
	}

	shuffleCards(b.stock, b.Seed)

	b.UndoStack = nil
	b.SavedPosition = 0
	b.State = Virgin
	b.stroke = nil
}

// NewGame resets Baize and restarts current variant with a new seed
func (b *Baize) NewGame() {
	if b.State == Started {
		TheStatistics.recordLostGame(b.Variant, b.percentComplete)
	}
	b.Seed = time.Now().UnixNano()
	b.RecallCardsToStock()
	b.dealCards()
	b.UndoPush()
	b.percentComplete = b.calcPercentComplete()
	TheStatistics.welcomeToast(b.Variant)
}

// NewVariant resets Baize and starts a new game with a new variant and seed
func (b *Baize) NewVariant(v string) {

	if b.State == Started {
		TheStatistics.recordLostGame(b.Variant, b.percentComplete)
	}

	// reset the baize
	b.Piles = b.Piles[:0]
	b.UndoStack = nil
	b.SavedPosition = 0
	b.State = Virgin
	b.stroke = nil

	// switch to new variant
	if _, exists := Variants[v]; !exists {
		v = "Klondike"
	}
	b.Variant = v
	TheUserData.Variant = v
	b.BuildVariant(v)
	b.ui.SetTitle(v)
	b.Seed = time.Now().UnixNano()

	b.OldWindowWidth = 0 // force a rescale
	b.Scale()            // now we know number of piles and can discover window width; scale the cards before creating them

	b.stock = b.findPilePrefix("Stock")
	if b.stock == nil {
		log.Fatal("Cannot find stock pile to create cards with")
	}
	createCards(b.stock)
	b.totalCards = b.stock.CardCount()
	shuffleCards(b.stock, b.Seed)

	b.dealCards()
	b.UndoPush()
	b.percentComplete = b.calcPercentComplete()
	TheStatistics.welcomeToast(b.Variant)
}

// LoadVariant tries to load a game from json resets Baize and continues an old game
func (b *Baize) LoadVariant(v string) bool {

	if _, exists := Variants[v]; !exists {
		v = "Klondike"
	}
	b.Variant = v
	TheUserData.Variant = v
	b.BuildVariant(v)
	b.ui.SetTitle(v)

	if !b.Load(v) {
		return false
	}

	sav, ok := b.UndoPop() // removes extra pushed state
	if !ok {
		log.Panic("error popping extra state from undo stack")
	}

	b.OldWindowWidth = 0 // force a rescale
	b.Scale()            // now we know number of piles and can discover window width; scale the cards before creating them

	b.stock = b.findPilePrefix("Stock")
	if b.stock == nil {
		log.Fatal("Cannot find stock pile to create cards with")
	}
	createCards(b.stock)
	b.totalCards = b.stock.CardCount()

	b.UpdateFromSaveable(sav)

	b.UndoPush()
	b.percentComplete = b.calcPercentComplete()
	TheStatistics.welcomeToast(b.Variant)

	if b.State == Complete {
		b.ui.ShowFAB("star", ebiten.KeyN)
	}

	return true
}

func (b *Baize) ShowInfo() {
	TheStatistics.welcomeToast(b.Variant)
	b.ui.Toast(fmt.Sprintf("You have made %s in this game, which is %d%% complete", util.Pluralize("move", len(b.UndoStack)-1), b.percentComplete))
	if !b.stock.Hidden() {
		b.ui.Toast(fmt.Sprintf("The stock contains %s", util.Pluralize("card", b.stock.CardCount())))
	}
}

func (b *Baize) dealCards() {
	for _, p := range b.Piles {
		deal := p.GetStringAttribute("Deal")
		if deal == "" {
			continue
		}
		for _, d := range deal {
			switch d {
			case 'u':
				c := b.stock.Pop() // this will flip card up
				if c == nil {
					log.Fatal("out of cards during deal u from ", deal)
				}
				if c.Prone() {
					log.Fatal("popped a face down card from stock")
				}
				p.Push(c)
			case 'd':
				c := b.stock.Pop() // this will flip card up
				if c == nil {
					log.Fatal("out of cards during deal d from ", deal)
				}
				c.FlipDown()
				p.Push(c)
			case '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D':
				idx, ok := findCard(b.stock.Cards, d)
				if ok {
					c := b.stock.Extract(idx)
					p.Push(c)
				} else {
					log.Fatal("cannot find", d, "during deal from ", deal)
				}
			default:
				println("unknown rune in Deal", string(d))
			}
		}
	}

	for _, p := range b.Piles {
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
			println(b.stock.CardCount(), "cards remaining in hidden stock")
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

func (b *Baize) countPiles(cls string) (total, count int) {
	total = 0
	count = 0
	for _, p := range b.Piles {
		if p.Class == cls {
			total++
			if p.CardCount() == 0 {
				count++
			}
		}
	}
	return
}

// PileTapped is called when a pile has been tapped
func (b *Baize) PileTapped(pTapped *Pile) {

	// this method is in Baize because it needs access to Baize.findPile()
	switch pTapped.Class {
	case "Stock":
		if pTapped.localRecycles > 0 {
			waste := b.findPile("Waste")
			if waste == nil || len(waste.Cards) == 0 {
				return
			}
			for len(waste.Cards) > 0 {
				c := waste.Pop()
				b.stock.Push(c) // this will flip card down
			}
			pTapped.SetRecycles(pTapped.localRecycles - 1)
			b.AfterUserMove()
		}
	case "StockCruel":
		/*
		   https://politaire.com/help/cruel

		   The redeal procedure begins by picking up all cards on the tableau.
		   The cards from the tableau are collected, one column at a time, starting with the left-most column,
		   picking up the cards in each column in top to bottom order.
		   Then, without shuffling, the cards are dealt out again, starting with the first card picked up,
		   and dealing the cards in the same order as they were picked up.
		*/
		if pTapped.localRecycles > 0 {
			for _, pTab := range b.Piles {
				if pTab.Class == "Tableau" {
					if pTab.CardCount() > 0 {
						// don't use Pop/Push, cards will be in wrong order
						b.MoveCards(pTab.Cards[0], b.stock)
					}
					pTab.scrunchPercentage = 100 // TODO ugly
				}
			}
			if b.stock.CardCount() == 0 {
				break
			}

			// use an intermediate stack to make sure cards are redealt in proper order
			tmp := make([]*Card, 0, cap(b.stock.Cards))
			for b.stock.CardCount() > 0 {
				tmp = append(tmp, b.stock.Pop())
			}

			for _, pTab := range b.Piles {
				if pTab.Class == "Tableau" {
					deal := pTab.GetStringAttribute("Deal")
					for _, d := range deal {
						var c *Card
						if len(tmp) > 0 {
							c = tmp[len(tmp)-1]
							tmp = tmp[:len(tmp)-1]
						} else {
							goto FinishedDealing
						}
						switch d {
						case 'u':
							c.FlipUp()
						case 'd':
							c.FlipDown()
						}
						pTab.Push(c)
					}
				}
			}
		FinishedDealing:
			pTapped.SetRecycles(pTapped.localRecycles - 1)
			b.AfterUserMove()
		}
	}

}

// CardTapped is called when a card has been tapped
func (b *Baize) CardTapped(c *Card) {

	// println("card",c.ID.String(),"tapped")

	if c.Transitioning() || c.Flipping() {
		println("cannot tap an animating card")
		return
	}

	pSrc := c.owner

	// can only tap top card
	// TODO might be playing Spider &c and trying to send a conformant pile to Foundation
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
		cardsToMove, ok := pSrc.GetIntAttribute("CardsToMove")
		if !ok || cardsToMove == 0 {
			cardsToMove = 1
		}
		for _, p := range b.Piles {
			if targetClass == p.Class {
				// println("found a", p.Class)
				if p.CanAcceptCard(c) {
					// println(p.Class, "can accept", c.ID.String())
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
			if b.stock.CardCount() > empty {
				TheBaize.ui.Toast("All tableaux spaces must be filled before dealing a new row")
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
				c.SetProne(false)
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
			// if p.Class == "Foundation" && p.buildFlags&8==8 {
			// 	// fake a drag
			// 	if c.owner.StartDrag(c) {
			// 		if b.CanAcceptTail(p, c.owner.Tail) {
			// 			b.MoveCards(c, p)
			// 			moved = true
			// 		}
			// 		p.CancelDrag(c)
			// 	}
			// 	if moved {
			// 		break
			// 	}
			// }
		}
	default:
		println("clueless when tapping on a", pSrc.Class, "card")
	}

	if moved {
		b.AfterUserMove()
	} else {
		if c != nil {
			c.Shake()
		}
	}

	// TODO else test other piles to see if this card is accepted?
}

// MoveCards from one pile to another, always from card downwards (inclusive)
func (b *Baize) MoveCards(c *Card, dst *Pile) {

	src := c.owner
	moveFrom := len(src.Cards)

	// find the index of the first card we will move
	for i, sc := range src.Cards {
		if sc == c {
			moveFrom = i
			break
		}
	}

	if moveFrom == len(src.Cards) {
		log.Panic("MoveCards could not find first card in source")
	}

	oldSrcLen := len(src.Cards)

	tmp := make([]*Card, 0, cap(src.Cards))

	// pop the tail off the source and push onto temp stack
	for i := len(src.Cards) - 1; i >= moveFrom; i-- {
		tmp = append(tmp, src.Pop())
	}

	// pop all cards off the temp stack and onto the destination
	for len(tmp) > 0 {
		dc := tmp[len(tmp)-1]
		tmp = tmp[:len(tmp)-1]
		dst.Push(dc)
	}

	if oldSrcLen == len(src.Cards) {
		println("nothing happened in MoveCards")
		return
	}

	// flip up an exposed source card
	if !strings.HasPrefix(src.Class, "Stock") {
		if tc := src.Peek(); tc != nil {
			tc.FlipUp()
		}
	}
	src.ScrunchCards()
	dst.ScrunchCards()

}

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
		if p.CardCount() == 0 {
			if aff := p.GetStringAttribute("AutoFillFrom"); aff != "" {
				affPiles := strings.Split(aff, ",")
				for _, srcPile := range affPiles {
					if src := b.findPile(srcPile); src != nil {
						if c := src.Peek(); c != nil {
							b.MoveCards(c, p)
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
			TheStatistics.recordWonGame(b.Variant, len(b.UndoStack)-1)
			TheStatistics.wonToast(b.Variant, len(b.UndoStack)-1)
			b.ui.ShowFAB("star", ebiten.KeyN)
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
	} else {
		log.Println("not pushing to undo because checksums match")
	}

	//

	b.percentComplete = b.calcPercentComplete()
}

func (b *Baize) largestIntersection(c *Card) *Pile {
	var largest int = 0
	var pile *Pile = nil
	cx0, cy0, cx1, cy1 := c.BaizeRect()
	for _, p := range b.Piles {
		px0, py0, px1, py1 := p.FannedBaizeRect()
		i := util.OverlapArea(cx0, cy0, cx1, cy1, px0, py0, px1, py1)
		if i > largest {
			largest = i
			pile = p
		}
	}
	return pile
}

func (b *Baize) calcPercentComplete() int {
	var foundations, sorted, unsorted int
	for _, p := range b.Piles {
		if p.Class == "Foundation" {
			foundations++
		}
		for i := 0; i < len(p.Cards)-1; i++ {
			c1 := p.Cards[i]
			c2 := p.Cards[i+1]
			if isCardPairConformant(p.Build, p.Flags, c1, c2) {
				sorted++
			} else {
				unsorted++
			}
		}
	}
	return int(util.MapValue(float64(sorted)-float64(unsorted)+float64(foundations), float64(-b.totalCards), float64(b.totalCards), 0, 100))
}

// StartDrag return true if the Baize can be dragged (vscrolled)
func (b *Baize) StartDrag() bool {
	return true
}

// DragBy move (vscroll) the Baize by dragging it
func (b *Baize) DragBy(dx, dy int) {
	b.DragOffsetY = b.DragOffsetBaseY + dy
	if b.DragOffsetY > 0 {
		b.DragOffsetY = 0 // DragOffsetY should only ever be 0 or -ve
	}
}

// StopDrag stop dragging the Baize
func (b *Baize) StopDrag() {
	// remember the amount of drag so the next drag starts from here
	b.DragOffsetBaseY = b.DragOffsetY
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (b *Baize) NotifyCallback(event interface{}) {
	switch v := event.(type) { // type switch https://tour.golang.org/methods/16
	case image.Point:
		// println("image.Point (tap)", v.X, v.Y)
		// a tap outside any open ui drawer (ie on the baize) closes the drawer
		if con := b.ui.VisibleDrawer(); con != nil {
			if !util.InRect(v.X, v.Y, con.Rect) {
				con.Hide()
			}
		} else {
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
				b.CardTapped(c)
			} else {
				if p := b.findPileAt(v.X, v.Y); p != nil {
					b.PileTapped(p)
				}
			}
		}
	case ebiten.Key:
		// println("ebiten.Key", v)
		if fn, ok := b.commandTable[v]; ok {
			b.ui.HideActiveDrawer()
			b.ui.HideFAB()
			fn()
		}
	case ui.ChangeRequest:
		// a widget has sent a change request
		b.ui.HideActiveDrawer()
		b.ui.HideFAB()
		switch v.ChangeRequested {
		case "Variant":
			newVariant := v.Data
			if newVariant == "" {
				println("unknown variant", v.Data)
				break
			}
			if newVariant != b.Variant {
				b.Save()
				if !TheBaize.LoadVariant(newVariant) {
					b.NewVariant(newVariant)
				}
			}
		case "CardBack":
			if TheUserData.CardStyle == "retro" {
				TheUserData.CardBackPattern = v.Data
				CardBackImage = TheCIP.BackImage(TheUserData.CardBackPattern)
			} else {
				TheUserData.CardBackColor = v.Data
				CardBackImage = TheCIP.BackImage(TheUserData.CardBackColor)
			}
		case "Highlight":
			TheUserData.HighlightMovable, _ = strconv.ParseBool(v.Data)
			// println("TheUserData.HighlightMovable :=", TheUserData.HighlightMovable)
			if TheUserData.HighlightMovable {
				b.HighlightMovable()
			} else {
				for _, p := range b.Piles {
					for _, c := range p.Cards {
						c.SetMovable(false)
					}
				}
			}
		case "Power moves":
			TheUserData.PowerMoves, _ = strconv.ParseBool(v.Data)
			if TheUserData.HighlightMovable {
				b.HighlightMovable()
			} else {
				for _, p := range b.Piles {
					for _, c := range p.Cards {
						c.SetMovable(false)
					}
				}
			}
		case "Retro cards":
			retro, _ := strconv.ParseBool(v.Data)
			if retro {
				TheUserData.CardStyle = "retro"
			} else {
				TheUserData.CardStyle = "default"
			}
			b.OldWindowWidth = 0 // force a rescale
			b.Scale()
		default:
			println("unknown change request", v.ChangeRequested, v.Data)
		}
	case input.StrokeEvent:
		// if v.Event != "move" {
		// 	println("stroke event", v.Event, v.X, v.Y)
		// }
		switch v.Event {
		case "start":
			b.stroke = v.Stroke
			if con := b.ui.VisibleDrawer(); con != nil {
				if util.InRect(v.X, v.Y, con.Rect) && con.StartDrag() {
					b.stroke.SetDraggedObject(con)
				} else {
					v.Stroke.Cancel()
				}
			} else {
				c := b.findCardAt(v.X, v.Y)
				if c != nil {
					b.stroke.SetDraggedObject(c)
					if !c.owner.StartDrag(c) {
						println("cancel stroke because drag not allowed")
						v.Stroke.Cancel()
					}
				} else {
					if b.StartDrag() {
						// println("starting baize drag")
						b.stroke.SetDraggedObject(b)
					} else {
						println("cancel stroke because not over a card")
						v.Stroke.Cancel()
					}
				}
			}
		case "move":
			if v.Stroke.DraggedObject() == nil {
				println("move stroke with nil dragged object")
				break
			}
			switch v.Stroke.DraggedObject().(type) { // type switch
			case *Card:
				c := v.Stroke.DraggedObject().(*Card)
				c.owner.DragTailBy(v.Stroke.PositionDiff())
			case ui.Container:
				con := v.Stroke.DraggedObject().(ui.Container)
				con.DragBy(v.Stroke.PositionDiff())
			case *Baize:
				// println("baize drag")
				b2 := v.Stroke.DraggedObject().(*Baize)
				if b2 != b {
					println("baize drag - something has gone terribly wrong")
				}
				b2.DragBy(v.Stroke.PositionDiff())
			default:
				println("unknown move dragging object")
			}
		case "stop":
			if v.Stroke.DraggedObject() == nil {
				println("stop stroke with nil dragged object")
				break
			}
			switch v.Stroke.DraggedObject().(type) { // type switch
			case *Card:
				c := v.Stroke.DraggedObject().(*Card)
				p := b.largestIntersection(c)
				if p == nil || p == c.owner {
					c.owner.CancelDrag(c)
				} else {
					if p.CanAcceptTail(c.owner.Tail, true) {
						c.owner.StopDrag(c)
						b.MoveCards(c, p)
						b.AfterUserMove()
					} else {
						c.owner.CancelDrag(c)
					}
				}
			case ui.Container:
				con := v.Stroke.DraggedObject().(ui.Container)
				con.StopDrag()
			case *Baize:
				// println("stop baize drag")
				b2 := v.Stroke.DraggedObject().(*Baize)
				if b2 != b {
					println("baize drag - something has gone terribly wrong")
				}
				b2.StopDrag()
			default:
				println("unknown stop dragging object")
			}
		default:
			println("unknown stroke event", v.Event)
			// case "cancel":
			// 	// c := v.Stroke.DraggingObject().(*Card)
			// 	c := v.Stroke.DraggedCard()
			// 	if c != nil {
			// 		c.owner.CancelDrag(c)
			// 	}
			// 	b.stroke = nil
		}
	default:
		println("unknown notification received", v)
	}
}

// ScaleCards calculates new width/height of cards and margins
func (b *Baize) ScaleCards() {

	windowWidth, _ := ebiten.WindowSize()

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
	switch TheUserData.CardStyle {
	default:
		slotWidth := float64(windowWidth) / float64(maxX+2)
		PilePaddingX = int(slotWidth / 10)
		CardWidth = int(slotWidth) - PilePaddingX
		slotHeight := slotWidth * 1.444
		PilePaddingY = int(slotHeight / 10)
		CardHeight = int(slotHeight) - PilePaddingY
		LeftMargin = (CardWidth / 2) + PilePaddingX
	case "poker":
		slotWidth := float64(windowWidth) / float64(maxX+2)
		PilePaddingX = int(slotWidth / 10)
		CardWidth = int(slotWidth) - PilePaddingX
		slotHeight := slotWidth * 1.39
		PilePaddingY = int(slotHeight / 10)
		CardHeight = int(slotHeight) - PilePaddingY
		LeftMargin = (CardWidth / 2) + PilePaddingX
	case "bridge":
		slotWidth := float64(windowWidth) / float64(maxX+2)
		PilePaddingX = int(slotWidth / 10)
		CardWidth = int(slotWidth) - PilePaddingX
		slotHeight := slotWidth * 1.561
		PilePaddingY = int(slotHeight / 10)
		CardHeight = int(slotHeight) - PilePaddingY
		LeftMargin = (CardWidth / 2) + PilePaddingX
	case "retro":
		CardWidth = 71
		PilePaddingX = 7
		CardHeight = 96
		PilePaddingY = 10
		cardsWidth := int(PilePositionType(PilePaddingX+CardWidth) * (maxX + 1)) // add 1 for half width card margin
		LeftMargin = (windowWidth - cardsWidth) / 2
	}
	log.Printf("scaled card size %s %dx%d", TheUserData.CardStyle, CardWidth, CardHeight)

	TopMargin = 48 + CardHeight/3

}

// Scale resizes piles, cards (inc shadow image), fonts and then repositions piles and cards
func (b *Baize) Scale() {

	// on startup, b.OldWindowWidth will be 0 so scalables will be built
	w, _ := ebiten.WindowSize()
	if w == b.OldWindowWidth {
		return
	}

	b.ScaleCards()

	CreateScalables()

	for _, p := range b.Piles {
		p.CreateBackgroundImage()
		if p.CardCount() == 0 {
			continue
		}
		// because we're about to use copy(), tmp must have a length
		var tmp = make([]*Card, len(p.Cards), cap(p.Cards)) // https://github.com/golang/go/wiki/SliceTricks#copy
		// len(tmp) == len(p.Cards)
		copy(tmp, p.Cards)
		p.Cards = p.Cards[:0] // keep the underlying array, slice the slice to zero length
		for _, c := range tmp {
			c.RefreshFaceImage()
			p.Push(c)
		}
	}

	b.OldWindowWidth = w
}

// Layout implements ebiten.Game's Layout.
func (b *Baize) Layout(outsideWidth, outsideHeight int) (int, int) {

	b.Scale()

	b.ui.Layout(outsideWidth, outsideHeight)

	return outsideWidth, outsideHeight

}

// Update the baize state (transitions, user input)
func (b *Baize) Update() error {

	b.input.Update() // detect mouse taps and keyboard input

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
		ebitenutil.DebugPrint(screen, fmt.Sprintf("NumGC %v, Undo %d, State %d, Complete %d%%, Movable %d", ms.NumGC, len(b.UndoStack), b.State, b.percentComplete, b.movableCards))
	}
}

// Exit this app
func (b *Baize) Exit() {
	if !NoGameSave {
		b.Save()
	}
	TheUserData.Save()
	os.Exit(0)
}
