# Gosol

Don't copy Opsole design (which is shaky and not fully understood)

Choice between using completely flat hierarchy (just Pile class) and depth-1 hierarchy (Pile, then Stock Waste &c).
Flat would mean lots of rules
Depth-1 is kinda supported by Go's struct embedding and interfaces

Classes should be dumb and all the game logic handled by Baize. Eg don't tell a card it's been tapped; tell Baize.

Each variant has
(1) Top level attributes (name, description, wikipedia)
(2) Attributes common to all classes (X, Y, Fan, BuildRule?, MoveRule?)
(3) Per-class attributes (Accept, Deal, SuitFilter, Packs &c)
These are loaded from a reference table/db and stored in the class
