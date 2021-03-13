package sol

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/ui"
	"oddstream.games/gosol/util"
)

// BaizeState is either Virgin, Started or Complete
type BaizeState int

// Virgin, Started or Complete
const (
	Virgin BaizeState = iota
	Started
	Complete
)

// Baize object describes the baize
type Baize struct {
	Piles         []*Pile
	Variant       string
	Seed          int64
	UndoStack     []SaveableBaize
	SavedPosition int
	totalCards    int
	State         BaizeState
	stroke        *input.Stroke
	input         *input.Input
	ui            *ui.UI
	commandTable  map[ebiten.Key]func()
}

// NewBaize is the factory func for the single Baize object
func NewBaize() *Baize {
	// TheUserData may have been injected from command line flags
	log.Printf("%v", TheUserData)
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
		ebiten.KeyF:      TheBaize.ShowPicker,
		ebiten.KeyMenu:   TheBaize.ui.OpenNavDrawer,
		ebiten.KeyEscape: TheBaize.ui.CloseActiveModal,
		ebiten.KeyX:      TheBaize.Exit,
	}
	BuildScalableCardImages() // need to do this after CardWidth,Height set - not in a func init()
	if NoGameLoad || !TheBaize.LoadVariant(TheBaize.Variant) {
		TheBaize.NewVariant(TheBaize.Variant)
	}
	return TheBaize // ugly global-setting kludge
}

// Reset the Baize
func (b *Baize) Reset() {
	b.Piles = nil
	b.UndoStack = nil
	b.SavedPosition = 0
	b.Variant = TheUserData.Variant
	b.Seed = time.Now().UnixNano()
	b.State = Virgin
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
	// 		if !c.Prone() {
	// 			log.Fatal("face up card found in stock")
	// 		}
	// 		if c.owner != stock {
	// 			log.Fatal("card in stock belongs to", c.owner.Class)
	// 		}
	// 	}
	// }

	shuffleCards(stock, b.Seed)

	b.UndoStack = nil // StartGame will deal cards then do initial UndoPush()
	b.SavedPosition = 0
	b.State = Virgin
	b.stroke = nil
}

// StartGame given existing variant and seed
func (b *Baize) StartGame() {
	b.dealCards()
	b.UndoPush()
	TheStatistics.welcomeToast(b.Variant)
}

// RestartGame resets Baize and restarts current variant with same seed
func (b *Baize) RestartGame() {
	// could load first entry on undo stack, start game will push initial state
	if b.State == Started {
		TheStatistics.recordLostGame(b.Variant, b.calcPercentComplete())
	}
	b.Restart()
	b.StartGame()
}

// NewGame resets Baize and restarts current variant with a new seed
func (b *Baize) NewGame() {
	if b.State == Started {
		TheStatistics.recordLostGame(b.Variant, b.calcPercentComplete())
	}
	b.Seed = time.Now().UnixNano()
	b.Restart()
	b.StartGame()
}

// NewVariant resets Baize and starts a new game with a new variant and seed
func (b *Baize) NewVariant(v string) {

	if b.State == Started {
		TheStatistics.recordLostGame(b.Variant, b.calcPercentComplete())
	}
	b.Reset()
	b.Variant = v
	b.Seed = time.Now().UnixNano()

	piles, ok := buildVariantPiles(b.Variant)
	if !ok {
		log.Fatal("unknown variant", b.Variant)
	}
	b.Piles = piles

	// temporary fudge to set window width to center cards on baize
	{
		b.ui.SetTitle(variantDisplayName(b.Variant))

		maxX := 0
		for _, p := range b.Piles {
			if p.X > maxX {
				maxX = p.X
			}
		}
		ebiten.SetWindowSize((maxX+2)*(CardWidth+10), WindowHeight)

		TopMargin = 48 + CardHeight/3
	}

	stock := b.findPilePrefix("Stock")
	if stock == nil {
		log.Fatal("Cannot find stock pile to create cards with")
	}
	createCards(stock)
	b.totalCards = stock.CardCount()
	shuffleCards(stock, b.Seed)

	b.StartGame()
}

// LoadVariant tries to load a game from json resets Baize and continues an old game
func (b *Baize) LoadVariant(v string) bool {

	if !b.Load(v) {
		return false
	}

	sav, ok := b.UndoPop() // removes extra pushed state
	if !ok {
		log.Fatal("error popping extra state from undo stack")
	}

	piles, ok := buildVariantPiles(b.Variant)
	if !ok {
		log.Fatal("unknown variant", b.Variant)
	}
	b.Piles = piles

	// temporary fudge to set window width to center cards on baize
	{
		b.ui.SetTitle(variantDisplayName(b.Variant))

		maxX := 0
		for _, p := range b.Piles {
			if p.X > maxX {
				maxX = p.X
			}
		}
		ebiten.SetWindowSize((maxX+2)*(CardWidth+10), WindowHeight)

		TopMargin = 48 + CardHeight/3
	}

	stock := b.findPilePrefix("Stock")
	if stock == nil {
		log.Fatal("Cannot find stock pile to create cards with")
	}
	createCards(stock)
	b.totalCards = stock.CardCount()

	b.UpdateFromSaveable(sav)
	b.UndoPush()
	TheStatistics.welcomeToast(b.Variant)

	return true
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
				if c == nil {
					log.Fatal("out of cards during deal")
				}
				if c.Prone() {
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
			case '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D':
				idx, ok := findCard(stock.Cards, d)
				if ok {
					c := stock.Extract(idx)
					p.Push(c)
				} else {
					log.Fatal("cannot find", d, "during deal")
				}
			default:
				println("unknown rune in Deal", d)
			}
		}
	}

	for _, p := range b.Piles {
		bury, ok := p.GetIntAttribute("Bury")
		if ok {
			p.BuryCards(bury)
		}
		disinter, ok := p.GetIntAttribute("Disinter")
		if ok {
			p.DisinterCards(disinter)
		}
		if p.Class == "Foundation" && p.CardCount() == 1 {
			afp := p.GetBoolAttribute("AcceptFirstPush")
			if afp {
				ord := p.Peek().Ordinal()
				for _, fp := range b.Piles {
					if fp.Class == "Foundation" {
						fp.SetAccept(ord)
					}
				}
			}
		}
	}

	if DebugMode {
		if stock.Y < 0 {
			println(stock.CardCount(), "cards remaining in hidden stock")
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
			if p.CardCount() == 0 {
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

	// println("card",c.ID.String(),"tapped")

	if c.Animating() {
		println("cannot tap an animating card")
		return
	}

	pSrc := c.owner

	// can only tap top card
	// TODO might be playing Spider &c
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
			stock := b.findPilePrefix("Stock")
			if stock.CardCount() > empty {
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
			// if p.Class == "FoundationSpider" {
			// 	// fake a drag
			// 	if c.owner.StartDrag(c) {
			// 		if p.CanAcceptTail(b.Piles, c.owner.Tail) {
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
	tmp := make([]*Card, 0, cap(src.Cards))

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
		println("nothing happened in MoveCards")
		return
	}
	// flip up an exposed source card
	if !strings.HasPrefix(src.Class, "Stock") {
		tc := src.Peek()
		if tc != nil {
			tc.FlipUp()
		}
	}
	src.ScrunchCards()
	dst.ScrunchCards()

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

// AfterUserMove runs after the user has made a move;
// perform any auto moves (on behalf of user), test for game complete, push state onto undo stack
func (b *Baize) AfterUserMove() {

	b.AutoMoves()

	//

	switch b.State {
	case Virgin:
		TheStatistics.startGame(b.Variant)
		b.State = Started
		b.ui.Toast(fmt.Sprintf("%s started", b.Variant))
	case Started:
		if b.Complete() {
			b.ui.Toast(fmt.Sprintf("%s complete in %d moves", b.Variant, len(b.UndoStack)-1))
			b.State = Complete
			TheStatistics.recordWonGame(b.Variant, len(b.UndoStack)-1)
		}
	case Complete:
		println("what are we doing here?")
	}

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

func (b *Baize) calcPercentComplete() int {
	var count int
	for _, p := range b.Piles {
		if strings.HasPrefix(p.Class, "Foundation") {
			count += p.CardCount()
		}
	}
	return int(math.Round(float64(count) / float64(b.totalCards) * 100))
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (b *Baize) NotifyCallback(event interface{}) {
	switch v := event.(type) { // type switch https://tour.golang.org/methods/16
	case image.Point:
		// println("image.Point (tap)", v.X, v.Y)
		if b.ui.ActiveModal() {
			if !util.InRect(v.X, v.Y, b.ui.ActiveRect) {
				b.ui.CloseActiveModal()
			}
		} else {
			c := b.findCardAt(v.X, v.Y)
			if b.stroke != nil {
				println("cancel stroke because tap")
				c.owner.CancelDrag(c) // if we have a stroke we must have a card
				b.stroke.Cancel()
			}
			if c != nil {
				b.CardTapped(c)
			} else {
				p := b.findPileAt(v.X, v.Y)
				if p != nil {
					b.PileTapped(p)
				}
			}
		}
	case ebiten.Key:
		// println("ebiten.Key", v)
		fn, ok := b.commandTable[v]
		if ok {
			if b.ui.ActiveModal() {
				b.ui.CloseActiveModal()
			}
			fn()
		}
	case string:
		newVariant := findVariantFromDisplayName(v)
		if b.ui.ActiveModal() {
			b.ui.CloseActiveModal()
		}
		if newVariant == "" {
			println("unknown variant", v)
		} else {
			TheUserData.Variant = v
			b.NewVariant(v)
		}
	case input.StrokeEvent:
		// if v.Event != "move" {
		// 	println("stroke event", v.Event, v.X, v.Y)
		// }
		switch v.Event {
		case "start":
			b.stroke = v.Stroke
			if b.ui.ActiveModal() {
				con := b.ui.FindContainerAt(v.X, v.Y)
				if con != nil {
					if con.StartDrag() {
						b.stroke.SetDraggedObject(con)
					} else {
						v.Stroke.Cancel()
					}
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
					println("cancel stroke because not over a card")
					v.Stroke.Cancel()
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
					if p.CanAcceptTail(b.Piles, c.owner.Tail) {
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

// Layout implements ebiten.Game's Layout.
func (b *Baize) Layout(outsideWidth, outsideHeight int) (int, int) {

	for _, p := range b.Piles {
		p.Layout(outsideWidth, outsideHeight)
	}

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
		p.DrawCards(screen)
	}
	for _, p := range b.Piles {
		p.DrawAnimatingCards(screen)
	}
	if DebugMode {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("NumGC %v, Undo %d, State %d, Percent %d", ms.NumGC, len(b.UndoStack), b.State, b.calcPercentComplete()))
	}

	b.ui.Draw(screen)
}

// Exit this app
func (b *Baize) Exit() {
	if !NoGameSave {
		b.Save()
	}
	TheUserData.Save()
	os.Exit(0)
}
