package sol

import (
	"fmt"

	"oddstream.games/gosol/util"
)

// Statistics is a container for the statistics for all variants
type Statistics struct {
	// PascalCase for JSON
	StatsMap map[string]*VariantStatistics
}

// VariantStatistics holds the statistics for one variant
type VariantStatistics struct {
	// PascalCase for JSON
	Won, Lost, BestWinningMoves, WorstWinningMoves, CurrStreak, BestStreak, WorstStreak, BestPercent int
}

// NewStatistics creates a new Statistics object
func NewStatistics() *Statistics {
	s := &Statistics{StatsMap: make(map[string]*VariantStatistics)}
	s.Load()
	return s
}

func (s *Statistics) findVariant(v string) *VariantStatistics {
	stats, ok := s.StatsMap[v]
	if !ok {
		stats = &VariantStatistics{BestWinningMoves: 9999} // everything else is 0
		s.StatsMap[v] = stats
		println("statistics has encountered a new variant", v)
	}
	return stats
}

func (s *Statistics) recordWonGame(v string, numberOfMoves int) {

	stats := s.findVariant(v)

	stats.Won = stats.Won + 1
	stats.BestPercent = 100
	if numberOfMoves < stats.BestWinningMoves {
		stats.BestWinningMoves = numberOfMoves
	}
	if numberOfMoves > stats.WorstWinningMoves {
		stats.WorstWinningMoves = numberOfMoves
	}
	if stats.CurrStreak < 0 {
		stats.CurrStreak = 1
	} else {
		stats.CurrStreak = stats.CurrStreak + 1
	}
	if stats.CurrStreak > stats.BestStreak {
		stats.BestStreak = stats.CurrStreak
	}

	println("Statistics recording won game")
	s.Save()
}

func (s *Statistics) recordLostGame(v string, percent int) {

	stats := s.findVariant(v)

	stats.Lost = stats.Lost + 1
	// don't see that currStreak can ever be zero
	if stats.CurrStreak > 0 {
		stats.CurrStreak = -1
	} else {
		stats.CurrStreak = stats.CurrStreak - 1
	}
	if stats.CurrStreak < stats.WorstStreak {
		stats.WorstStreak = stats.CurrStreak
	}
	if stats.BestPercent < 100 {
		if percent > stats.BestPercent {
			stats.BestPercent = percent
		}
	}
	println("Statistics recording lost game")
	s.Save()
}

func (s *Statistics) welcomeToast(v string) {
	toasts := []string{}

	stats, ok := s.StatsMap[v]
	if !ok || stats.Won+stats.Lost == 0 {
		toasts = append(toasts, fmt.Sprintf("You have not played %s before", v))
	} else {
		if stats.BestPercent == 0 {
			toasts = append(toasts, fmt.Sprintf("You have yet to score anything in %s", util.Pluralize("attempt", stats.Lost)))
		} else if stats.BestPercent < 100 {
			toasts = append(toasts, fmt.Sprintf("Your best score is %d%% in %s", stats.BestPercent, util.Pluralize("attempt", stats.Lost)))
		} else {
			toasts = append(toasts,
				fmt.Sprintf("You have won %s, and lost %s (%d%%)",
					util.Pluralize("game", stats.Won),
					util.Pluralize("game", stats.Lost),
					((stats.Won*100)/(stats.Won+stats.Lost))))
			if stats.CurrStreak > 0 {
				toasts = append(toasts, fmt.Sprintf("You are on a winning streak of %s", util.Pluralize("game", stats.CurrStreak)))
			}
			if stats.CurrStreak < 0 {
				toasts = append(toasts, fmt.Sprintf("You are on a losing streak of %s", util.Pluralize("game", util.Abs(stats.CurrStreak))))
			}
		}
	}

	for _, t := range toasts {
		TheUI.Toast(t)
	}
}

func (s *Statistics) wonToast(v string, moves int) {
	toasts := []string{}

	stats, ok := s.StatsMap[v]
	if !ok || stats.Won+stats.Lost == 0 {
		toasts = append(toasts, fmt.Sprintf("You have not played %s before", v))
	} else {
		toasts = append(toasts, fmt.Sprintf("%s complete in %d moves", v, moves))
		toasts = append(toasts,
			fmt.Sprintf("You have won %s, and lost %s (%d%%)",
				util.Pluralize("game", stats.Won),
				util.Pluralize("game", stats.Lost),
				((stats.Won*100)/(stats.Won+stats.Lost))))
	}

	for _, t := range toasts {
		TheUI.Toast(t)
	}
}

func (s *Statistics) ShowStatistics() {
	var toasts = []string{}

	stats, ok := s.StatsMap[TheBaize.Variant]
	if !ok || stats.Won+stats.Lost == 0 {
		toasts = append(toasts, fmt.Sprintf("You have not played %s before", TheBaize.Variant))
	} else {
		toasts = append(toasts, fmt.Sprintf("You have made %s in this game, which is %d%% complete", util.Pluralize("move", len(TheBaize.UndoStack)-1), TheBaize.percentComplete))

		if stats.BestPercent == 0 {
			toasts = append(toasts, fmt.Sprintf("You have yet to score anything in %s", util.Pluralize("attempt", stats.Lost)))
		} else if stats.BestPercent < 100 {
			toasts = append(toasts, fmt.Sprintf("Your best score is %d%% in %s", stats.BestPercent, util.Pluralize("attempt", stats.Lost)))
		} else {
			toasts = append(toasts,
				fmt.Sprintf("You have won %s, and lost %s (%d%%)",
					util.Pluralize("game", stats.Won),
					util.Pluralize("game", stats.Lost),
					((stats.Won*100)/(stats.Won+stats.Lost))))
			if stats.BestWinningMoves != stats.WorstWinningMoves {
				toasts = append(toasts, fmt.Sprintf("Your best number of winning moves is %d, average %d, your worst is %d",
					stats.BestWinningMoves,
					(stats.BestWinningMoves+stats.WorstWinningMoves)/2,
					stats.WorstWinningMoves))
			}
			if stats.CurrStreak > 0 {
				toasts = append(toasts, fmt.Sprintf("You are on a winning streak of %s", util.Pluralize("game", stats.CurrStreak)))
			}
			if stats.CurrStreak < 0 {
				toasts = append(toasts, fmt.Sprintf("You are on a losing streak of %s", util.Pluralize("game", util.Abs(stats.CurrStreak))))
			}
			toasts = append(toasts, fmt.Sprintf("Your best streak is %d, your worst is %d", stats.BestStreak, stats.WorstStreak))
		}
	}
	TheUI.ShowTextDrawer(toasts)
}
