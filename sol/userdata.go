package sol

// UserData contains the settings and preferences for the user
type UserData struct {
	// Capitals to emit to json
	Game             string
	Variant          string
	CardBackPattern  string
	CardBackColor    string
	CardStyle        string
	HighlightMovable bool
	PowerMoves       bool
	MuteSounds       bool
}
