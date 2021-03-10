package input

type (
	// Observable https://gist.github.com/patrickmn/1549985
	Observable interface {
		Add(observer Observer)
		Notify(event interface{})
		Remove(event interface{})
	}

	// Observer https://gist.github.com/patrickmn/1549985
	Observer interface {
		NotifyCallback(event interface{})
	}
)
