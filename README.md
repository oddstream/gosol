# Gosol

Towards a polymorphic solitaire engine in Go+Ebiten, with help from fogleman/gg. It's an adaptation of my Lua/Solar2D retained mode engine used in the Android game (which was itself an adaptation of my messy vanilla JavaScript/SVG engine used for the online game). The intention is that this version will replace both of those.

It currently knows how to play:

* American Toad (and the original The Toad)
* Australian
* Baker's Dozen (and Baker's Dozen Relaxed)
* Canfield (and Storehouse)
* Duchess
* EasyWin (an easy to win game, for debugging)
* Forty and Eight
* Freecell (and Freecell Easy)
* Klondike (Draw One, Draw Three and Thoughtful)
* Limited
* Mistress and Mrs Mop
* Simple Simon
* Scorpion (and Wasp)
* Spider (original, One Suit and Two Suits)
* Spiderette
* Thumb and Pouch
* Will o' the Wisp
* Yukon

Some variants have been tried and discarded as being a bit silly:

* Giant
* King Albert
* Raglan

Some will never make it here because they a just bad games:

* Pyramid

It uses the Microsoft solitaire retro card set by default, but as these don't scale prettily it also has easy-to-read scalable cards.

It currently has a minimal user interface, and some keyboard shortcuts:

* U - undo
* N - new deal
* R - restart deal
* S - save current position ('bookmark')
* L - load a previously saved position (goto bookmark, handy for Freecell and Simple Simon players)

It currently doesn't do anything pretty after detecting a completed game.

That's a lot of 'currently', but this is a work-in-progress.
