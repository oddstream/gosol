# Gosol

Towards a polymorphic solitaire engine in [Go](https://golang.org/)+[Ebiten](https://ebiten.org/), with help from [fogleman/gg](https://github.com/fogleman/gg).

It's an adaptation of my [Lua](https://www.lua.org/)/[Solar2D](https://solar2d.com/) retained mode engine used in the Android game (which was itself an adaptation of my messy vanilla JavaScript/SVG engine used for the [online game](https://oddstream/games/Solitaire)). The intention is that this version will replace both of those, and provide Linux, Windows, Android and browser-based versions from the same code base. If I had a Mac there'd be iOS and Mac versions, too.

## Variants

It currently knows how to play:

* American Toad (and the original The Toad)
* Australian
* Baker's Dozen (and Baker's Dozen Relaxed)
* Canfield, Storehouse
* Cruel, Ripple Fan
* Duchess
* EasyWin (an easy to win game, for debugging)
* Forty and Eight
* Freecell, Freecell Easy
* Klondike (Draw One, Draw Three and Thoughtful)
* Limited
* Mistress and Mrs Mop
* Simple Simon
* Scorpion, Wasp
* Spider (original, One Suit and Two Suits)
* Spiderette
* Thumb and Pouch
* Will o' the Wisp
* Yukon

Some variants have been tried and discarded as being a bit silly:

* Agnes Sorel
* Giant
* King Albert
* Raglan

(I don't see the point of games that you almost certainly can't win; I like ones that have a 33-66% chance of winnning.)

Some will never make it here because they are just poor games:

* Accordian
* Pyramid (or any card matching variant)

## User interface

It has an intentionally minimal material-style user interface, and some keyboard shortcuts:

* U - undo
* N - new deal
* R - restart deal
* S - save current position ('bookmark')
* L - load a previously saved position (goto bookmark, handy for Freecell and Simple Simon players)
* C - collect cards to the foundations
 
## TODO

* Sounds
* I'd like it to have an inter-user high scores table, but the Google Play games services interface and setup is inpenetrable to me at the moment.

## Live playable version

There's a live WASM version [here](https://oddstream.games/gosol/gosol.html). (Currently an 8 MByte download, I must get that down.)

## Acknowledgements

Original games by Jan Wolter, David Parlett, Paul Alfille, Art Cabral, Albert Morehead, Geoffrey Mott-Smith, Zach Gage and Thomas Warfield.

Retro card back designs by Leslie Kooy and Susan Kare.
