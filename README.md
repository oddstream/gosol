# Gosol

Towards a polymorphic solitaire engine in Go+Ebiten, with help from fogleman/gg. It's an adaptation of my Lua/Solar2D retained mode engine used in the Android game (which was itself an adaptation of my messy vanilla JavaScript/SVG engine used for the online game). The intention is that this version will replace both of those.

It currently only plays a few variants (Australian, Canfield, Klondike (draw 1 and draw 3), Freecell, Limited, Spider1, Spider2). Use a command line flag (eg -v=Limited) to set the variant. There should be around 40 by the time I'm done.

It currently has no user interface other than U - undo, N - new deal, R - restart deal.

It currently doesn't detect a completed game, or do any fancy highlighting of moveable cards, or recording of the player's statistics.

It uses the Microsoft solitaire retro card set by default, but as these don't scale prettily it also has simplified scalable cards.

To start with (this is only my third Go project) I had a fight with using Go's embedded structs and interfaces, but we've resolved our differences and are getting on much better now. I'm finding that, with the help of VS Code's Go extensions, you can write a lot of Go code that works first time, which is nice.
