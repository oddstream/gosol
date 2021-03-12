package sol

import (
	"fmt"
	"log"

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

func (s *Statistics) startGame(v string) {
	stats, ok := s.StatsMap[v]
	if !ok {
		stats = &VariantStatistics{BestWinningMoves: 9999} // everything else is 0
		s.StatsMap[v] = stats
	}
}

func (s *Statistics) recordWonGame(v string, numberOfMoves int) {
	stats, ok := s.StatsMap[v]
	if !ok {
		log.Fatal("recordWonGame unknown variant ", v)
	}

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
	stats, ok := s.StatsMap[v]
	if !ok {
		log.Fatal("recordLostGame unknown variant ", v)
	}

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
	displayName := variantDisplayName(v)
	stats, ok := s.StatsMap[v]
	if !ok {
		TheBaize.ui.Toast(fmt.Sprintf("You have not played %s before", displayName))
		return
		// log.Fatal("welcomeToast unknown variant ", v)
	}
	if stats.Won+stats.Lost == 0 {
		TheBaize.ui.Toast(fmt.Sprintf("You have not played %s before", displayName))
	} else {
		// TheBaize.ui.Toast(fmt.Sprintf("You have played %s of %s", util.Pluralize("game", stats.Won+stats.Lost), displayName))

		if stats.BestPercent == 0 {
			TheBaize.ui.Toast(fmt.Sprintf("You have yet to score anything"))
		} else if stats.BestPercent < 100 {
			TheBaize.ui.Toast(fmt.Sprintf("Your best score is %d%%", stats.BestPercent))
		} else {
			TheBaize.ui.Toast(
				fmt.Sprintf("You have won %s, and lost %s (%d%%)",
					util.Pluralize("game", stats.Won),
					util.Pluralize("game", stats.Lost),
					((stats.Won * 100) / (stats.Won + stats.Lost))))
		}
		if stats.BestPercent == 100 {
			if stats.CurrStreak > 0 {
				TheBaize.ui.Toast(fmt.Sprintf("You are on a winning streak of %s", util.Pluralize("game", stats.CurrStreak)))
			}
			if stats.CurrStreak < 0 {
				TheBaize.ui.Toast(fmt.Sprintf("You are on a losing streak of %s", util.Pluralize("game", util.Abs(stats.CurrStreak))))
			}
		}
	}
}

// ShowStatistics open statistics window
func (s *Statistics) ShowStatistics(v string) {

	displayName := variantDisplayName(v)
	stats, ok := s.StatsMap[v]
	if !ok {
		TheBaize.ui.AppendTextToWindow(fmt.Sprintf("You have not played %s before", displayName))
		return
		// log.Fatal("welcomeToast unknown variant ", v)
	}
	if stats.Won+stats.Lost == 0 {
		TheBaize.ui.AppendTextToWindow(fmt.Sprintf("You have not played %s before", displayName))
	} else {
		if stats.BestPercent == 0 {
			TheBaize.ui.AppendTextToWindow(fmt.Sprintf("You have yet to score anything"))
		} else if stats.BestPercent < 100 {
			TheBaize.ui.AppendTextToWindow(fmt.Sprintf("Your best score is %d%%", stats.BestPercent))
		} else {
			TheBaize.ui.AppendTextToWindow(
				fmt.Sprintf("You have won %s, and lost %s (%d%%)",
					util.Pluralize("game", stats.Won),
					util.Pluralize("game", stats.Lost),
					((stats.Won * 100) / (stats.Won + stats.Lost))))
		}
		if stats.BestPercent == 100 {
			if stats.CurrStreak > 0 {
				TheBaize.ui.AppendTextToWindow(fmt.Sprintf("You are on a winning streak of %s", util.Pluralize("game", stats.CurrStreak)))
			}
			if stats.CurrStreak < 0 {
				TheBaize.ui.AppendTextToWindow(fmt.Sprintf("You are on a losing streak of %s", util.Pluralize("game", util.Abs(stats.CurrStreak))))
			}
		}
	}
}

// ShowStatistics open statistics window
func (b *Baize) ShowStatistics() {
	b.ui.OpenWindow(b.input, "Statistics")
	TheStatistics.ShowStatistics(b.Variant)
}
