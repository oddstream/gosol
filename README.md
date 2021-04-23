# Gosol

Towards a polymorphic solitaire engine in [Go](https://golang.org/)+[Ebiten](https://ebiten.org/), with help from [fogleman/gg](https://github.com/fogleman/gg).

It's an adaptation of my [Lua](https://www.lua.org/)/[Solar2D](https://solar2d.com/) retained mode engine used in the Android game (which was itself an adaptation of my messy vanilla JavaScript/SVG engine used for the [online game](https://oddstream/games/Solitaire)). The intention is that this version will replace both of those, and provide Linux, Windows, Android and browser-based versions from the same code base. If I had a Mac there'd be iOS and Mac versions, too.

## Variants

It currently knows how to play a few games, including:

* American Toad (and the original The Toad)
* Australian
* Baker's Dozen (and Baker's Dozen Relaxed)
* Busy Aces
* Canfield, Storehouse
* Cruel, Ripple Fan
* Duchess
* EasyWin (an easy to win game, for debugging)
* Fortune's Favor
* Forty Thieves (and Forty and Eight, Maria, Limited, Lucas)
* Freecell, Freecell Easy
* Klondike (Draw One, Draw Three and Thoughtful)
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
* N - new deal (resign current game, if started)
* R - restart deal
* S - save current position ('bookmark')
* L - load a previously saved position (goto bookmark, handy for Freecell and Simple Simon players)
* C - collect cards to the foundations

Tapping a card has a few shortcuts:

* Tap the top Stock card to deal one or three cards to the Waste pile, or deal cards to the tableaux if playing a Spider variant
* Tap a card to send it to a Foundation (you can't currently tap a pile in a Spider game to send 13 cards to a Foundation)

## Other features

* Unlimited undo
* Bookmarking positions
* Scalable and retro card designs
* Movable card highlighting (the more useful a move looks, the more the card gets highlighted)
* Statistics (including percent compete and streaks - streaks are great)
* Cards spin prettily when you complete a game
* Slightly randomized sounds

A lot a features have been tried and discarded, in order to keep the game (and player) focused. Weniger aber besser, as [Dieter Rams](https://en.wikipedia.org/wiki/Dieter_Rams) taught us. Just because a feature *can* be implemented, does not mean it *should* be.

## TODO

* Figure out why the touch interface isn't working properly
* I'd like it to have an inter-user high scores table, but the Google Play games services interface and setup is inpenetrable to me at the moment.

## Live playable version

There's a live WASM version [here](https://oddstream.games/gosol/gosol.html). (Currently an 8 MByte download, I must get that down.)

## Acknowledgements

Original games by Jan Wolter, David Parlett, Paul Alfille, Art Cabral, Albert Morehead, Geoffrey Mott-Smith, Zach Gage and Thomas Warfield.

Retro card back designs by Leslie Kooy and Susan Kare.
