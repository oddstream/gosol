# Gosol

Towards a polymorphic solitaire engine in [Go](https://golang.org/)+[Ebiten](https://ebiten.org/), with help from [fogleman/gg](https://github.com/fogleman/gg).

It's an adaptation of my [Lua](https://www.lua.org/)/[Solar2D](https://solar2d.com/) retained mode engine used in the Android game (which was itself an adaptation of my messy vanilla JavaScript/SVG engine used for the [online game](https://oddstream/games/Solitaire)). The intention is that this version will replace both of those, and provide Linux, Windows, Android and browser-based versions from the same code base. If I had a Mac there would be iOS and Mac versions, too.

## Variants

It currently knows how to play about 40 games, including:

* American Toad (and the original The Toad)
* Australian
* Baker's Dozen (and Baker's Dozen Relaxed)
* Bisley
* Canfield, Storehouse
* Cruel, Ripple Fan
* Duchess
* EasyWin (an easy to win game, for debugging)
* Fortune's Favor
* Forty Thieves (and Busy Aces, Fortune's Favor, Forty and Eight, Josephine, Maria, Limited, Lucas, Red and Black)
* Freecell, Freecell Easy
* Klondike (and Draw One, Draw Three, Double Klondike, Gargantua, Thumb and Pouch, Thoughtful)
* Mistress and Mrs Mop
* Scorpion, Wasp
* Simple Simon
* Spider (and Spiderette, Spider One Suit, Spider Two Suits, Will o' the Wisp)
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
* L - load/return to a previously saved position
* C - collect cards to the foundations

Tapping a card has a few shortcuts:

* Tap the top Stock card to deal one or three cards to the Waste pile, or deal cards to the tableaux if playing a Spider variant
* Tap a card to send it to a Foundation (you can't currently tap a pile in a Spider game to send 13 cards to a Foundation)

## Other features

* Unlimited undo, without penalty. Also, you can restart a deal without penalty.
* Bookmarking positions (really good for games like Freecell or Simple Simon)
* Scalable and retro card designs
* Movable card highlighting (the more useful a move looks, the more the card gets highlighted)
* Statistics (including percent complete and streaks - streaks are great)
* Cards spin and flutter when you complete a game, so you feel rewarded and happy
* Slightly randomized sounds
* Automatic saving of games in progress

A lot a features have been tried and discarded, in order to keep the game (and player) focused. Weniger aber besser, as [Dieter Rams](https://en.wikipedia.org/wiki/Dieter_Rams) taught us. Just because a feature *can* be implemented, does not mean it *should* be.

## FAQ

### What makes this different from the other solitaire implementations?

This solitaire is all about [Flow](https://en.wikipedia.org/wiki/Flow_(psychology)). 
Anything that distracts from your interaction with the flow of the game has been either been tried and removed or not included. 
Crucially, the games can be played by single-clicking the card you wish to move, and the software figures out where you want the card to go 
(mostly to the foundation if possible, and if not, the biggest tableau). If you don't like where the card goes, 
just try clicking it again or drag it.

Also, I'm trying to make games authentic, by taking the rules from reputable sources
and implementing them exactly.

### Why are the graphics so basic?

Anything that distracts from your interaction with the flow of the game, 
or the ability to scan a deck of cards,
has either been tried and removed, or not included. 
This includes: 
fancy card designs (front and back), 
changing the screen/baize background, 
animating card flips, 
keeping an arbitrary score, 
distracting graphics on the screen, 
and forcing you to drag the cards to move them.

The user interface tries to stick to the Material Design guidelines, and so is minimal and tactile.
I looked at a lot of the other solitaire websites out there, and think how distracting some of them are. 
Features seem to have been added because the developers thought they were cool; 
they never seem to have stopped to 
consider that just because they *could* implement a feature, that they *should*.

### Sometimes the cards are really huge or really tiny

Either resize your browser/desktop window (if using scalable cards) or change the settings to fixed size cards.

### The rules for a variation are wrong

There's no ISO or ANSI or FIDE-like governing body for solitaire; so there's no standard set of rules.
Other implementations vary in how they interpret each variant.
For example, some variants of American Toad build the tableau down by suit, some by alternate color.
So, rather than just making this stuff up, I've tried to find a well researched set of rules for each variant and stick to them,
leaning heavily on Jan Wolter (RIP, and thanks for all the fish), David Parlett and Thomas Warfield.
Where possible, I've implemented the games from the book 
"The Complete Book of Solitaire and Patience Games" by Albert Morehead and Geoffrey Mott-Smith.

### But you can cheat!

You can when playing with actual cards, too. Cheat if you like; I'm not your mother.

### What about scores?

Nope, the software doesn't keep an arbitary score. Too confusing. 
Just the number of moves made, number of wins, and your winning streak (streaks are great).
A game isn't counted at all until you move a card. 
Thereafter, if you ask for a new deal, that counts as a loss.

### What about a timer?

Nope, there isn't one of those. Too stressful.
Solitaire is also called patience; it's hard to feel patient when you're pressured by a clock.

### What's with the settings?

#### Scaled cards

The size of the cards is changed dymaically so the cards fill the width of the screen. In some variants, this can cause the
cards to disappear off the bottom, in which case you can (a) drag the baize, (b) switch to fixed or retro cards, or (c) change
the size of the window (if not running on a mobile device).

#### Fixed cards

#### Retro cards

#### Card back...

#### Single tap

Enabling this allows a single tap on a card to move it, for example from the stock pile to the waste pile, or to the
fullest tableaux pile, or to a foundation pile. You you want the card to go to a specific place, then you have to drag it there.

#### Highlights

This highlights cards that *can* be moved by you. The software uses diffenent levels of highlighting, depending on it's simple assessment
of how useful it thinks that move may be. For example, a card that can be moved to a foundation pile is highlighted the most.

#### Power moves

Some variants (eg Freecell or Forty Thieves) only allow you to move one card at a time.

#### Mute sounds

So you can, say, listen to an audio book while playing.

### Is the game rigged?

No. The cards are shuffled randomly using a Fisher-Yates shuffle 
driven by a Park-Miller pseudo random number generator, 
which is in itself seeded by a random number. This mechanism was tested and analysed to make sure it produced an even distribution of shuffled cards. 

There are 80658175170943878571660636856403766975289505440883277824000000000000
possible deals of a pack of 52 playing cards; you're never going to play the same game twice, nor indeed play the same game 
that anyone else ever has, or ever will.     

## TODO

* Get it working on Android 
* I'd like it to have an inter-user high scores table, but the Google Play games services interface and setup is inpenetrable to me at the moment.

## Live playable version

There's a live WASM version [here](https://oddstream.games/gosol/gosol.html). (Currently an 12.8 MByte download, I must get that down.)

## Acknowledgements

Original games by Jan Wolter, David Parlett, Paul Alfille, Art Cabral, Albert Morehead, Geoffrey Mott-Smith, Zach Gage and Thomas Warfield.

Retro card back designs by Leslie Kooy and Susan Kare.
