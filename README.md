# Minimal Polymorphic Solitaire in Go

Towards a polymorphic solitaire engine in [Go](https://golang.org/)+[Ebiten](https://ebiten.org/), with help from [fogleman/gg](https://github.com/fogleman/gg) (both of which are highly recommended), with game variants run by (user supplied) scripts.

![Screenshot](https://github.com/oddstream/gosol/blob/7152668f4b5053a1d438981e9d4564624616da6a/screenshots/Simple%20Simon.png)

It's tested on Linux, Windows and in a web browser. If you have go installed, you should be able to run it on Linux or Windows by cloning this repo and then `go run .` in the cloned directory. Or, install it using `go install github.com/oddstream/gosol@latest`.

There is a live playable WASM version [here](https://oddstream.games/gosol/gosol.html) (sorry about the large initial download).

It's created because I *have* to write software, and for my own personal enjoyment. It's skewed toward puzzle-type games, because they're the ones I mostly play. It's definitely not for profit and will never contain ads.

## Variants

It currently knows how to play:

* Agnes Bernauer
* Australian
* Baker's Dozen
* Canfield (also Storehouse, American Toad, Duchess)
* Easy (an easy to win game, for debugging)
* Forty Thieves (also Sixty Thieves, Busy Aces, Forty and Eight, Josephine, Maria, Limited, Lucas, Red and Black, Rank and File, Number Ten)
* Freecell (also Freecell Easy, Blind Freecell, Eight Off, Seahaven Towers)
* Klondike (also Gargantua, Klondike Draw Three, Triple Klondike, Thoughtful)
* Penguin
* Scorpion (also Wasp)
* Simple Simon
* Spider (also Spider One Suit, Spider Two Suits)
* Usk
* Whitehead
* Westcliff (Classic, American and Easthaven)
* Yukon (also Yukon Cells)

Variants are added when the whim takes me, or when some aspect of the engine needs testing/extending, or when someone asks.

Some variants have been tried and discarded as being a bit silly, or just too hard:

* Agnes Sorel
* Giant
* King Albert
* Raglan

Some will never make it here because they are just poor games:

* Accordian
* Golf
* Pyramid (or any card matching variant)

![Screenshot](https://github.com/oddstream/gosol/blob/7152668f4b5053a1d438981e9d4564624616da6a/screenshots/Australian.png)

## Other features

* Permissive card moves. If you want to move a card from here to there, go ahead and do it. If that move is not allowed by the current rules, the game will put the cards back *and explain why that move is not allowed*.
* Unlimited undo, without penalty. Also, you can restart a deal without penalty.
* Bookmarking positions (really good for puzzle-style games like Freecell or Simple Simon).
* Scalable cards. Change the size and shape of the window to make the cards fit.
* One-tap interface. Tapping on a card or cards tries to move them to a foundation, or to a suitable tableau pile.
* Cards in traditional red and black (best for games like Klondike or Yukon where cards are sorted into alternating colors), or in four colors (for games where cards are sorted by suit, like Australian or Spider).
* Every game has a link to it's Wikipedia page.
* Statistics (including percent complete and streaks; percent is good for games that are not often won, and streaks are good for games that are).
* Cards spin and flutter when you complete a game, so you feel rewarded and happy.
* Slightly randomized sounds.
* Automatic saving of game in progress.
* A dragable baize; if cards spill out of view to the bottom or right of the screen, just drag the baize to move them into view.
* A 'discard' pile type so that Spideresque games can be implemented as they are described in the textbooks (other software reuses Foundation piles).

## Deliberate minimalism

A lot a features have been tried and discarded, in order to keep the game (and player) focused. Weniger aber besser, as [Dieter Rams](https://en.wikipedia.org/wiki/Dieter_Rams) taught us. Design is all about saying "no", as Steve Jobs preached. Just because a feature *can* be implemented, does not mean it *should* be.

Configurability is the root of all evil, someone said. Every configuration option in a program is a place where the program is too stupid to figure out for itself what the user really wants, and should be considered a failure of both the program and the programmer who implemented it. So, there's one card face, one color palette, one card animation speed, and so on.

![Screenshot](https://github.com/oddstream/gosol/blob/7152668f4b5053a1d438981e9d4564624616da6a/screenshots/American%20Toad.png)

## FAQ

### What makes this different from the other solitaire implementations?

This solitaire is all about [Flow](https://en.wikipedia.org/wiki/Flow_(psychology)).

Anything that distracts from your interaction with the flow of the game has been either been tried and removed or not included.

Crucially, the games can be played by single-clicking the card you wish to move, and the software figures out where you want the card to go (mostly to the foundation if possible, and if not, the biggest tableau, or an empty cell). If you don't like where the card goes, just try clicking it again or dragging it.

Also, I'm trying to make games authentic, by taking the rules from reputable sources and implementing them exactly.

### Why are the graphics so basic?

Anything that distracts from your interaction with the flow of the game,
or the ability to scan a deck of cards,
has either been tried and removed, or not included.
This includes:
fancy card designs (front and back),
changing the screen/baize background,
keeping an arbitrary score,
distracting graphics on the screen.

The user interface tries to stick to the Material Design guidelines, and so is minimal and tactile.
I looked at a lot of the other solitaire websites and apps out there, and think how distracting some of them are. Features seem to have been added because the developers thought they were cool; they never seem to have stopped to consider that just because they *could* implement a feature, that they *should*.

### Sometimes the cards are really huge or really tiny

Either resize your browser/desktop window (if using scalable cards) or change the settings to fixed size cards.

### The rules for a variation are wrong

There's no ISO or ANSI or FIDE-like governing body for solitaire; so there's no standard set of rules.
Other implementations vary in how they interpret each variant.
For example, some variants of American Toad build the tableau down by suit, some by alternate color.
So, rather than just making this stuff up, I've tried to find a well researched set of rules for each variant and stick to them, leaning heavily on Jan Wolter (RIP, and thanks for all the fish), David Parlett and Thomas Warfield. Where possible, I've implemented the games from the book "The Complete Book of Solitaire and Patience Games" by Albert Morehead and Geoffrey Mott-Smith.

### Keyboard shortcuts?

* C - collect cards to the foundations
* B - bookmark current position; Ctrl+B - return position to last bookmark
* H - hint/help - show movable cards
* N - new deal (resign current game, if started)
* R - restart deal
* U - undo

### What about scores?

Nope, the software doesn't keep an arbitary score. Too confusing.
Just the number of wins, the average 'completeness percentage' and your winning streak (streaks are great).
A game isn't counted until you move a card.
Thereafter, if you ask for a new deal or switch to a different variant, that counts as a loss.

You can cheat the score system by restarting a deal and then asking for a new deal.

'Completeness percentage' is calculated from the number of unsorted pairs of cards in all the piles.

### But you can cheat

You can when playing with actual cards, too. Cheat if you like; I'm not your mother.

### What about a timer?

Nope, there isn't one of those. Too stressful.
Solitaire is also called *patience*; it's hard to feel patient when you're pressured by a clock.

### You can't move cards off a foundation pile

Nope. Reading the "original" rules for a lot of the games seem to explicitly forbid this, so there's a complete ban on moving cards off a foundation pile.

You can always use undo if you get stuck or change you mind about a move.

### What's with the settings?

#### Fixed cards

With this checked, the size of the cards is fixed, preventing them from scaling.

Otherwise, the size of the cards is changed dynamically so the cards fill the width of the screen. In some variants, this can cause the
cards to disappear off the bottom, in which case you can drag the baize, or change the size of the window (if not running on a mobile device).

#### Power moves

Some variants (eg Freecell or Forty Thieves) only allow you to move one card at a time. Moving several cards between piles requires
you to move them, one at a time, via an empty pile or cell. Enabling power moves automates this, allowing multi-card moves between piles.
The number of cards you can move is calculated from the number of empty piles and cells (if any).

#### Colorful cards

Depending on the variant, enabling this draws the cards in four colors, rather than the usual black and red.

#### Mirror baize

Mirrors the card piles on the baize from right to left, because not everyone is right handed, or likes the stock to be on the left of the screen when they are right handed.

#### Mute sounds

So you can, for example, listen to an audio book while playing.

#### Auto collect

Enabling this will cause cards to be moved to the Foundation piles after every move you make.

#### Safe collect

In games like Klondike that build tableau cards in alternating colors, you can sometimes get into trouble by moving cards to the foundations too soon. With this option turned on, the titlebar collect button will only move cards to the foundation piles when it is safe to do so.

### Is the game rigged?

No. The cards are shuffled randomly using a Fisher-Yates shuffle driven by a Park-Miller pseudo random number generator, which is in itself seeded by a random number. This mechanism was tested and analysed to make sure it produced an even distribution of shuffled cards.

There are 80658175170943878571660636856403766975289505440883277824000000000000 possible deals of a pack of 52 playing cards; you're probably more likely to grow wings than play the same game twice, or play the same game that anyone else ever has, or ever will.

### Any hints and tips?

* For games that start with face down cards (like Klondike or Yukon) the priority is to get the face down cards turned over. This is more important than moving cards to the foundations.
* For games that start with a block of cards in the tableau and only allow single cards to be moved (like Forty Thieves), the priority is usually to open up some space (create empty tableaux piles) to allow you to juggle cards around.
* For Forty Thieves-style games, the *other* priority is to minimize the number of cards in the waste pile.
* For puzzle-type games (like Baker's Dozen, Freecell, Penguin, Simple Simon), take your time and think ahead.
* For games with reshuffles (like Cruel and Perseverance) you need to anticipate the effects of the reshuffle.
* Focus on sorting the cards in the tableaux, rather than moving cards to the foundations. Only move cards to the foundations when you *have* to.
* Use undo and bookmark. Undo isn't cheating; it's improvising, adapting and overcoming.

## Where are the preferences, statistics and saved games stored?

On Linux, you'll find them as `.json` files in a folder called `~/.config/oddstream.games/gosol`.

On Windows, you'll find them as `.json` files in a folder called `C:\Users\<username>\AppData\Roaming\oddstream.games\gosol`.

## Terminology and conventions

* A PILE of cards

* A CONFORMANT series of cards in a pile is called a SEQUENCE

* A set of cards is called a PACK (not 'deck')

* Suits are listed in alphabetic order: Club, Diamond, Heart, Spade

* Cards changing between face down and face up is called FLIPPING

* The user never moves or flips a face down card, only the dealer can

* Cards cannot be played from the foundation

* Cell, Foundation and Waste piles only hold face up cards

* Stock only has face down cards

* A game is EASY when the deal has been 'fixed', usually by moving Aces to the foundations, or shuffling Kings or Aces in the tableaux

* A game is RELAXED when any restriction on what card can be moved to an empty tableaux pile has been removed

![Screenshot](https://github.com/oddstream/gosol/blob/7152668f4b5053a1d438981e9d4564624616da6a/screenshots/Klondike.png)

## The seven different types of piles

### Stock

All games have a stock pile, because this is where the cards are created and start their life.

In some games, like Freecell, the stock pile is off screen (invisible). In others, like Klondike, it's on screen (usually at the top left corner) and tapping the top card will cause one card to be flipped up and moved to a waste pile. In other games, like Spider, tapping the top card will cause cards to be moved to each of the tableau piles.

All cards in the stock are always face down. You can't move a card to the stock pile. There is only ever one stock pile.

### Tableau

Tableau piles are where the main building in the game happens. The player tries to move the cards around the tableau and other piles, so that the cards in each tableau pile are sorted into some game-specific order. For example, in Klondike and Freecell, the tableau cards start in some random order, and must be sorted into increasing rank and alternating color.

Normally, only the top card of a tableau pile may be moved. This can be relaxed to allow consecutive sequences of conformant cards at the top of the pile to be moved as a unit. Some games appear to allow the latter while actually only allowing one card at a time to be moved; Freecell is a prime example of this. In reality, Freecell allows 'power moves' to hide the one-card-only rule.

Sometimes, there is a constraint on which card may be placed onto an empty tableau, for example in Klondike, and empty tableau can only contain a King.

Some cards in the tableau pile may start life face down; the game will automatically turn the cards up when they are exposed.

### Foundation

Foundation piles are where the player is trying to move the cards to, so that the game is completed.

The cards in each foundation usually start with an Ace, and build up, always the same suit. A foundation pile is full (complete) when it contains 13 cards.

Only one card at a time can be moved to a foundation. Cards cannot be taken off a foundation.

### Discard

Discard piles aren't usually found in other solitaire implementations.

Discard piles are like foundation piles, except that only a complete set of 13 cards can be moved at once.

Moving completed sets of cards to a discard is optional, and is usally done to create space in the tableaux. You do not have to move cards to a discard pile to complete a game.

### Waste

A waste pile can store any number of cards, all face up. You can only move one card at a time to a waste pile, and that card must come from the stock pile. There is only ever one waste pile.

In some games (like Klondike) cards in the waste pile can be recycled back to the stock pile, by tapping on an empty stock pile. The game may restrict the number of times this can happen.

### Cell

A cell is either empty, or it can contain one card of any type. Cell cards are always face up, and available for play to tableau or foundation piles. Cells are used as temporary holding areas for cards.

### Reserve

A reserve pile contains a series of cards, usually all face down with only the top card face up and available for play to a foundation, tableau or cell pile.

Only one card at a time may be moved from a reserve, and cards can never be moved to a reserve pile.

## TODO

* Split the code into front and back end, and add a universal solver.
* Reduce the size of the executable (using [UPX](https://upx.github.io/)?) and WASM files.
* Scripted game variants, possibly using [GopherLua](https://github.com/yuin/gopher-lua), or a Tcl-style little language.
* ~~The LÖVE+Lua version contains several things that are implemented better, so I'm in the process of copying the designs back to this version.~~
* ~~Get it working on Android (agggh! help!).~~
* ~~I'd like it to have an inter-user high scores table, but the Google Play games services interface and setup is inpenetrable to me at the moment.~~
* (Don't) replace Ebiten with [Fyne](https://fyne.io). That's a complete rewrite because Ebiten is [immediate mode](https://en.wikipedia.org/wiki/Immediate_mode_(computer_graphics)), and Fyne is [retained mode](https://en.wikipedia.org/wiki/Retained_mode). I prefer immediate mode.
* Give up and rewrite the whole thing in Go+Fyne, [Godot](https://godotengine.org/), [Defold](https://www.defold.com), Dart+Flutter, Java+libGDX, Kotlin+Korge, Haxe, Rust, Tcl/Tk, Pascal, Wren, Clojure, or something else. I agonize over this, and have made several false starts, but keep coming back to C, Go or Lua.

## History

First, there was a Javascript version that used SVG graphics and ran in web browsers. Game variants were configured using static lookup tables, which I thought was a good idea at the time. It's still available via the https://oddstream.games website, but I don't recommend it. The cards are not pretty and some of the rules are incorrect.

Second, there was a version in Lua, using the Solar2d game engine, that made it to the Google Play store.  Game variants were configured using static lookup tables, which I still thought was a good idea.

Third, there was a version in Go, using the Ebiten game engine, with help from gg/fogleman. The design was internally a mess (you can't write Java in Go!), and the cards didn't scale, so this was abandoned. Game variants were configured using static lookup tables, which was starting to become a source of clumsiness and code smells.

Fourth, there is a version in C that uses the Raylib game engine and uses Lua to script the game variants. The design was good, but has problems with scaling cards.

Fifth, there was this version in Go, using the Ebiten game engine, with help from gg/fogleman. The design is way better than the original attempt in Go, allows the option for scripting games, and has sharp graphics.

Sixth, there was a complete rewrite in Lua + the LÖVE game engine. It replaced the second version, runs on Android/Linux/Windows, and is available in the Google Play Store. It's similar to the Go version, and the two get together to swap features and design ideas.

## Acknowledgements

Original games by Jan Wolter, David Parlett, Paul Alfille, Art Cabral, Albert Morehead, Geoffrey Mott-Smith, Zach Gage and Thomas Warfield.

Sounds by [kenney.nl](https://www.kenney.nl/assets) and Victor Vashenko.
