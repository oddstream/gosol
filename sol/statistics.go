package sol

import (
	"fmt"
	"log"
)

// Statistics is a container for the statistics for all variants
type Statistics struct {
	// PascalCase for JSON
	StatsMap map[string]*VariantStatistics
}

// VariantStatistics holds the statistics for one variant
type VariantStatistics struct {
	// PascalCase for JSON
	Won, Lost, CurrStreak, BestStreak, WorstStreak, SumPercents, BestPercent, BestMoves, WorstMoves, SumMoves int `json:",omitempty"`
	// Won is number of games with 100%
	// Lost is number of games with % less than 100
	// Won + Lost is total number of games played (won or abandoned)
	// SumPercents is a record of games where % < 100
	// average % is (sum of Percents) + (100 * Won) / (Won+Lost)
}

func (stats *VariantStatistics) averagePercent() int {
	played := stats.Won + stats.Lost
	if played > 0 {
		return (stats.SumPercents + (stats.Won * 100)) / played
	}
	return 0
}

func (stats *VariantStatistics) strings(v string) []string {
	strs := []string{}
	if stats.Won+stats.Lost == 0 {
		strs = append(strs, fmt.Sprintf("You have not played %s before", v))
	} else {
		strs = append(strs, fmt.Sprintf("Played: %d", stats.Won+stats.Lost))
		strs = append(strs, fmt.Sprintf("Won: %d", stats.Won))
		strs = append(strs, fmt.Sprintf("Lost: %d", stats.Lost))
		winRate := (stats.Won * 100) / (stats.Won + stats.Lost)
		strs = append(strs, fmt.Sprintf("Win rate: %d%%", winRate))

		avpc := stats.averagePercent()
		if avpc < 100 {
			strs = append(strs, fmt.Sprintf("Average complete: %d%%", avpc))
		}
		if stats.BestPercent < 100 {
			// not yet won a game
			strs = append(strs, "You have yet to win a game")
			strs = append(strs, fmt.Sprintf("Best percent: %d%%", stats.BestPercent))
		} else {
			// won at least one game
			strs = append(strs, fmt.Sprintf("Best number of moves: %d", stats.BestMoves))
			strs = append(strs, fmt.Sprintf("Worst number of moves: %d", stats.WorstMoves))
			strs = append(strs, fmt.Sprintf("Average number of moves: %d", stats.SumMoves/stats.Won))
		}

		if stats.CurrStreak != 0 {
			strs = append(strs, fmt.Sprintf("Current streak: %d", stats.CurrStreak))
		}
		if stats.BestStreak != 0 {
			strs = append(strs, fmt.Sprintf("Best streak: %d", stats.BestStreak))
		}
		if stats.WorstStreak != 0 {
			strs = append(strs, fmt.Sprintf("Worst streak: %d", stats.WorstStreak))
		}
	}
	return strs
}

// NewStatistics creates a new Statistics object (a map)
// and loads the saved statistics into it from file
func NewStatistics() *Statistics {
	s := &Statistics{StatsMap: make(map[string]*VariantStatistics)}
	s.Load()
	return s
}

func (s *Statistics) findVariant(v string) *VariantStatistics {
	stats, ok := s.StatsMap[v]
	if !ok {
		stats = &VariantStatistics{} // everything 0
		s.StatsMap[v] = stats
		// println("statistics has encountered a new variant", v)
	}
	return stats
}

func (s *Statistics) RecordWonGame(v string, moves int) {

	TheUI.Toast("Complete", fmt.Sprintf("Recording completed game of %s", v))

	stats := s.findVariant(v)

	stats.Won = stats.Won + 1

	if stats.CurrStreak < 0 {
		stats.CurrStreak = 1
	} else {
		stats.CurrStreak = stats.CurrStreak + 1
	}
	if stats.CurrStreak > stats.BestStreak {
		stats.BestStreak = stats.CurrStreak
	}

	stats.BestPercent = 100

	if stats.BestMoves == 0 || moves < stats.BestMoves {
		stats.BestMoves = moves
	}
	if stats.WorstMoves == 0 || moves > stats.WorstMoves {
		stats.WorstMoves = moves
	}
	stats.SumMoves += moves

	s.Save()
}

func (s *Statistics) RecordLostGame(v string) {

	percent := TheBaize.PercentComplete()
	if percent == 100 {
		log.Println("*** That's odd, here is a lost game that is 100% complete ***")
	}

	TheUI.Toast("Fail", fmt.Sprintf("Recording lost game of %s, %d%% complete", v, percent))

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

	if percent > stats.BestPercent {
		stats.BestPercent = percent
	}
	stats.SumPercents += percent

	s.Save()
}

func ShowStatisticsDrawer() {
	stats := TheStatistics.findVariant(TheBaize.LongVariantName())
	var strs []string = stats.strings(TheBaize.LongVariantName())
	TheUI.ShowTextDrawer(strs)
}
