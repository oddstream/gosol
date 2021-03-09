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
	Played, Won, BestWinningMoves, WorstWinningMoves, CurrStreak, BestStreak, WorstStreak, BestPercent int
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
	stats.Played = stats.Played + 1
	s.Save()
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

// func (b *Baize) recordStatistics() {
// 	switch b.State {
// 	case Started:
// 		TheStatistics.recordLostGame(b.Variant, b.calcPercentComplete())
// 	case Complete:
// 		TheStatistics.recordWonGame(b.Variant, len(b.UndoStack)-1)
// 	}
// }

func (s *Statistics) welcomeToast(v string) {
	stats, ok := s.StatsMap[v]
	if !ok {
		TheBaize.ui.Toast(fmt.Sprintf("You have not played %s before", v))
		return
		// log.Fatal("welcomeToast unknown variant ", v)
	}
	if stats.Played == 0 {
		TheBaize.ui.Toast(fmt.Sprintf("You have not played %s before", v))
	} else {
		TheBaize.ui.Toast(fmt.Sprintf("You have started %d games of %s", stats.Played, v))
	}
	if stats.BestPercent < 100 {
		TheBaize.ui.Toast(fmt.Sprintf("Your best score is %d%%", stats.BestPercent))
	} else {
		TheBaize.ui.Toast(fmt.Sprintf("You have won %d games", stats.Won))
	}
	if stats.BestPercent == 100 {
		if stats.CurrStreak > 0 {
			TheBaize.ui.Toast(fmt.Sprintf("You are one a winning streak of %d games", stats.CurrStreak))
		}
		if stats.CurrStreak < 0 {
			TheBaize.ui.Toast(fmt.Sprintf("You are on a losing streak of %d games", util.Abs(stats.CurrStreak)))
		}
	}
}
