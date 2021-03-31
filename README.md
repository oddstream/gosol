# Gosol

Towards a polymorphic solitaire engine in Go+Ebiten, with help from fogleman/gg. It's an adaptation of my Lua/Solar2D retained mode engine used in the Android game (which was itself an adaptation of my messy vanilla JavaScript/SVG engine used for the online game). The intention is that this version will replace both of those.

It currently knows how to play:

* Australian
* Baker's Dozen (and Baker's Dozen Relaxed)
* Canfield (and Storehouse)
* Forty and Eight
* Freecell (and Freecell Easy)
* Klondike (Draw One and Draw Three)
* Limited
* Mistress and Mrs Mop
* Simple Simon
* Scorpion (and Wasp)
* Spider (original, One Suit and Two Suits)
* Spiderette
* Thoughtful
* Thumb and Pouch
* Will o' the Wisp
* Yukon

It uses the Microsoft solitaire retro card set by default, but as these don't scale prettily it also has easy-to-read scalable cards, set via a command line flag `-c=default`. The card width scales to match the window width; the card height is calculated automatically. You can adjust the width:height ratio with `-c=default`, `-c=bridge` or `-c=poker`.

It currently has a minimal user interface, and some keyboard shortcuts:

* U - undo
* N - new deal
* R - restart deal
* S - save current position ('bookmark')
* L - load a previously saved position (goto bookmark, handy for Freecell and Simple Simon players)
* F1 - show the rules
* F2 - change the card back

I'm currently building up a minimal Material Design UI in a separate package for it, keeping an eye on [Ebiten UI](https://ebitenui.github.io/) for ideas.

It currently doesn't do anything pretty after detecting a completed game, or do any fancy highlighting of moveable cards.

That's a lot of 'currently', but this is a work-in-progress.
