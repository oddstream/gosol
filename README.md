# Gosol

There's a live WASM version [here](https://oddstream.games/gomps5/gosol.html).

Towards a polymorphic solitaire engine in [Go](https://golang.org/)+[Ebiten](https://ebiten.org/), with help from [fogleman/gg](https://github.com/fogleman/gg).

It's an adaptation of my [Lua](https://www.lua.org/)/[Solar2D](https://solar2d.com/) retained mode engine used in the Android game (which was itself an adaptation of my messy vanilla JavaScript/SVG engine used for the [online game](https://oddstream/games/Solitaire)). The intention is that this version will replace both of those, and provide Linux, Windows, Android and browser-based versions from the same code base. If I had a Mac there would be iOS and Mac versions, too.

## Variants

It currently knows how to play:

* Aces and Kings
* Algerian
* Alhambra
* American Toad (also The Toad)
* Australian
* Baker's Dozen (also Baker's Dozen Relaxed)
* Bisley
* Bristol (also Dover)
* Canfield (also Acme, Storehouse)
* Cruel, Ripple Fan
* Duchess
* EasyWin (an easy to win game, for debugging)
* Fortune's Favor
* Forty Thieves (also Busy Aces, Fortune's Favor, Forty and Eight, Josephine, Maria, Limited, Lucas, Red and Black)
* Freecell (also Eight Off, Freecell Easy)
* Klondike (also Draw One, Draw Three, Batsford, Double Klondike, Gargantua, Thumb and Pouch, Thoughtful)
* La Belle Lucie (and Trefoil, The Fan)
* Mistress and Mrs Mop
* Scorpion, Wasp
* Simple Simon
* Spider (also Beetle, Spiderette, Spider One Suit, Spider Two Suits, Will o' the Wisp)
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

## Other features

* Unlimited undo, without penalty. Also, you can restart a deal without penalty.
* Bookmarking positions (really good for games like Freecell or Simple Simon)
* Scalable or fixed-size cards
* Statistics (including percent complete and streaks - streaks are great)
* Cards spin and flutter when you complete a game, so you feel rewarded and happy
* Slightly randomized sounds
* Automatic saving of game in progress
* A dragable baize; if cards spill out of view to the bottom or right of the screen, just drag the baize to move them into view

## Deliberate minimalism

A lot a features have been tried and discarded, in order to keep the game (and player) focused. Weniger aber besser, as [Dieter Rams](https://en.wikipedia.org/wiki/Dieter_Rams) taught us. Design is all about saying "no", as Steve Jobs preached. Just because a feature *can* be implemented, does not mean it *should* be.

So taken out were:

* Movable card highlighting (the more useful a move looks, the more the card gets highlighted)
* Turn-offable one tap interface

## Configurability is the root of all evil

Every configuration option in a program is a place where the program is too stupid to figure out for itself what the user really wants, and should be considered a failure of both the program and the programmer who implemented it.

## FAQ

### What makes this different from the other solitaire implementations?

This solitaire is all about [Flow](https://en.wikipedia.org/wiki/Flow_(psychology)).
Anything that distracts from your interaction with the flow of the game has been either been tried and removed or not included.
Crucially, the games can be played by single-clicking the card you wish to move, and the software figures out where you want the card to go
(mostly to the foundation if possible, and if not, the biggest tableau). If you don't like where the card goes,
just try clicking it again or dragging it.

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

### Keyboard shortcuts?

* U - undo
* N - new deal (resign current game, if started)
* R - restart deal
* S - save current position ('bookmark')
* L - load/return to a previously saved position
* C - collect cards to the foundations

### What about scores?

Nope, the software doesn't keep an arbitary score. Too confusing.
Just the number of moves made, number of wins, and your winning streak (streaks are great).
A game isn't counted until you move a card.
Thereafter, if you ask for a new deal, that counts as a loss.

You can cheat the score system by restarting a deal and then asking for a new deal.

### But you can cheat

You can when playing with actual cards, too. Cheat if you like; I'm not your mother.

### What about a timer?

Nope, there isn't one of those. Too stressful.
Solitaire is also called *patience*; it's hard to feel patient when you're pressured by a clock.

### What's with the settings?

#### Fixed cards

The size of the cards is fixed, preventing them from scaling.

Otherwise, the size of the cards is changed dynamically so the cards fill the width of the screen. In some variants, this can cause the
cards to disappear off the bottom, in which case you can (a) drag the baize, (b) switch to fixed or retro cards, or (c) change
the size of the window (if not running on a mobile device).

#### Card back

Changes the card back to either a different color (when using scalable or fixed size cards) or design (when using retro cards).

#### Single tap

Enabling this allows a single tap on a card to move it, for example from the stock pile to the waste pile, or to the
fullest tableaux pile, or to a foundation pile. You you want the card to go to a specific place, then you have to drag it there.

You can always tap on the card on top of a Stock pile.

#### Free drag

By default, the game only allows you to drag cards that are allowed to be dragged by the current variant's rules.
When this setting is checked, you can drag any cards anywhere, anytime; if the move isn't allowed, the game will send the cards
back, and maybe give you a message saying why.

#### Highlights

This highlights cards that *can* be moved by you. The software uses different levels of highlighting, depending on it's simple assessment
of how useful it thinks that move may be. For example, a card that can be moved to a foundation pile is highlighted the most.

#### Power moves

Some variants (eg Freecell or Forty Thieves) only allow you to move one card at a time. Moving several cards between piles requires
you to move them, one at a time, via an empty pile or cell. Enabling power moves automates this, allowing multi-card moves between piles.
The number of cards you can move is calculated from the number of empty piles and cells (if any).

#### Mirror baize

Mirrors the card piles on the baize from right to left, because not everyone is right handed, or likes the stock to be on the left of
the screen when they are right handed.

#### Mute sounds

So you can, for example, listen to an audio book while playing.

### Is the game rigged?

No. The cards are shuffled randomly using a Fisher-Yates shuffle
driven by a Park-Miller pseudo random number generator,
which is in itself seeded by a random number. This mechanism was tested and analysed to make sure it produced an even distribution of shuffled cards.

There are 80658175170943878571660636856403766975289505440883277824000000000000
possible deals of a pack of 52 playing cards; you're never going to play the same game twice, nor indeed play the same game
that anyone else ever has, or ever will.

### Any hints and tips?

* For games that start with face down cards (like Klondike or Yukon) the priority is to get the face down cards turned over.
* For games that start with a block of cards in the tableaux and only allow single cards to be moved (like Forty Thieves), the priority is usually to open up some space (create empty tableaux piles) to allow you to juggle cards around.
* For Forty Thieves-style games, the *other* priority is to minimize the number of cards in the waste pile.
* For puzzle-type games (like Freecell, Simple Simon and Mistress Mop), take your time and think ahead.
* For games with reshuffles (like Cruel and Perseverance) you need to anticipate the effects of the reshuffle.
* Use undo and bookmark. Undo isn't cheating; it's improvising, adapting and overcoming.

## TODO

* Get it working on Android.
* Get the size of the executable and WASM down.
* Use SVG for the card graphics. The obvious SVG renderer (github.com/srwiley/oksvg, rasterx) doesn't render the wikimedia SVG files properly.
* I'd like it to have an inter-user high scores table, but the Google Play games services interface and setup is inpenetrable to me at the moment.
* Give up and rewrite the whole thing in [Defold](https://www.defold.com) or C/[raylib](https://www.raylib.com).

## Acknowledgements

Original games by Jan Wolter, David Parlett, Paul Alfille, Art Cabral, Albert Morehead, Geoffrey Mott-Smith, Zach Gage and Thomas Warfield.

Retro card back designs by [Leslie Kooy](https://www.lesliekooy.com/), jazz cup pattern by Gina Ekiss.

Sounds by [kenney.nl](https://www.kenney.nl/assets) and Victor Vashenko
