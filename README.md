# Gosol

Towards a polymorphic solitaire engine in Go+Ebiten, with help from fogleman/gg. It's an adaptation of my Lua/Solar2D retained mode engine used in the Android game (which was itself an adaptation of my messy vanilla JavaScript/SVG engine used for the online game). The intention is that this version will replace both of those.

It currently only plays a few variants (Klondike, Freecell, Limited, Spider1, Spider2). Use a command line flag (eg -v=Limited) to set the variant.

It currently has no user interface other than U - undo, N - new deal, R - restart deal.

It currently doesn't detect a completed game, or do any fancy highlighting of moveable cards, or automoving of cards.

Had a fight with using Go's embedded structs and interfaces, but we've resolved our differences and are getting on much better now.
