# Gosol

Towards a polymorphic solitaire engine in Go+Ebiten, with help from fogleman/gg. It's an adaptation of my Lua/Solar2D retained mode engine used in the Android game (which was itself an adaptation of my messy vanilla JavaScript/SVG engine used for the online game). The intention is that this version will replace both of those.

It currently only plays a few variants (Australian, Baker's Dozen, Canfield, Klondike (draw 1 and draw 3), Freecell, Limited, Scorpion, Spider1, Spider2, Storehouse Canfield, Wasp, Yukon). Use a command line flag (eg `-v=Limited`) to set the variant. There should be around 40 by the time I'm done.

It currently has a minimal user interface other than a some toasts, a toolbar and some keyboard shortcuts:

* U - undo
* N - new deal
* R - restart deal
* S - save current position
* L - load a previously saved position (handy for the Freecell and Simple Simon players)

I'm currently building up a Material Design UI in a separate package for it, keeping an eye on [Ebiten UI](https://ebitenui.github.io/) for ideas.

It currently doesn't do anything pretty after detecting a completed game, or do any fancy highlighting of moveable cards, or recording of the player's statistics.

That's a lot of 'currently', but this is a work-in-progress.

It uses the Microsoft solitaire retro card set by default, but as these don't scale prettily it also has simplified scalable cards, set via a command line flag `-c=scalable`, and `-cardwidth=` to set the card width; `-cardwidth=90` looks good. The card height is calculated automatically.
