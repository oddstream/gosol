package sol

// UserData contains the level the user is on
type UserData struct {
	// Capitals to emit to json
	Copyright       string
	Game            string
	Variant         string
	CardBack        string
	CardStyle       string
	BackColor       string
	CompletedLevels int
}

// UserDataIO performsLoad/Save of UserData objects
// type UserDataIO interface {
// 	Load(*UserData)
// 	Save(*UserData)
// }
