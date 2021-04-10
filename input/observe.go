package input

type (
	// Observable https://gist.github.com/patrickmn/1549985
	Observable interface {
		Add(Observer)
		Notify(StrokeEvent)
		Remove(Observer)
	}

	// Observer https://gist.github.com/patrickmn/1549985
	Observer interface {
		NotifyCallback(StrokeEvent)
	}
)
