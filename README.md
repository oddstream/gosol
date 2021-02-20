# Gosol

Don't copy Opsole design (which is shaky and not fully understood)

Choice between using completely flat hierarchy (just Pile class) and depth-1 hierarchy (Pile, then Stock Waste &c).
Flat would mean lots of rules
Depth-1 is kinda supported by Go's struct embedding and interfaces

Classes should be dumb and all the game logic handled by Baize. Eg don't tell a card it's been tapped; tell Baize.

Each variant has
(1) Top level attributes (name, description, wikipedia)
(2) Attributes common to all classes (X, Y, Fan, BuildRule?, MoveRule?) (No, don't use this)
(3) Per-pile attributes (Accept, Deal, SuitFilter, Packs &c)
These are loaded from a reference table/db and stored in the class

Levels of card movement
(1) directed by the engine. Eg put cards in stock after creation.
(2) directed by the engine, animation required. Eg deal, move to foundation, return after unsucessful drag.
(3) directed by the user, animation required.
(4) Dragged by user.

Card transition queue (CTQ)
Create a New() queue
Add {*Card, X, Y} to queue
If card is alreay in queue, update record
Every tick, Update takes head off queue and trigger card's transition