package sol

import (
	"fmt"
	"strings"
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

// NewStatistics creates a new Statistics object (a map)
// and loads the saved statistics into it from file
func NewStatistics() *Statistics {
	s := &Statistics{StatsMap: make(map[string]*VariantStatistics)}
	s.Load()
	return s
}

// averagePercent is a helper function
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
		strs = append(strs, strings.ToUpper(v))
		strs = append(strs, fmt.Sprintf("Played: %d", stats.Won+stats.Lost))
		strs = append(strs, fmt.Sprintf("Won: %d", stats.Won))
		strs = append(strs, fmt.Sprintf("Lost: %d", stats.Lost))
		winRate := (stats.Won * 100) / (stats.Won + stats.Lost)
		strs = append(strs, fmt.Sprintf("Win rate: %d%%", winRate))
		strs = append(strs, " ")

		avpc := stats.averagePercent()
		if avpc < 100 {
			strs = append(strs, fmt.Sprintf("Average incomplete: %d%%", avpc))
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

func (s *Statistics) strings() []string {
	var strs []string = []string{}
	var numPlayed, numWon, numLost int
	for _, vs := range s.StatsMap {
		numPlayed += vs.Won + vs.Lost
		numWon += vs.Won
		numLost += vs.Lost
	}
	strs = append(strs, fmt.Sprintf("Played: %d", numPlayed))
	strs = append(strs, fmt.Sprintf("Won: %d", numWon))
	strs = append(strs, fmt.Sprintf("Lost: %d", numLost))
	winRate := (numWon * 100) / (numPlayed)
	strs = append(strs, fmt.Sprintf("Win rate: %d%%", winRate))
	return strs
}

func (s *Statistics) findVariant(v string) *VariantStatistics {
	vstats, ok := s.StatsMap[v]
	if !ok {
		vstats = &VariantStatistics{} // everything 0
		s.StatsMap[v] = vstats
		// println("statistics has encountered a new variant", v)
	}
	return vstats
}

func (s *Statistics) Played(v string) int {
	vstats := s.findVariant(v)
	return vstats.Won + vstats.Lost
}

func (s *Statistics) RecordWonGame(v string, moves int) string {

	vstats := s.findVariant(v)

	vstats.Won = vstats.Won + 1

	if vstats.CurrStreak < 0 {
		vstats.CurrStreak = 1
	} else {
		vstats.CurrStreak = vstats.CurrStreak + 1
	}
	if vstats.CurrStreak > vstats.BestStreak {
		vstats.BestStreak = vstats.CurrStreak
	}

	vstats.BestPercent = 100

	if vstats.BestMoves == 0 || moves < vstats.BestMoves {
		vstats.BestMoves = moves
	}
	if vstats.WorstMoves == 0 || moves > vstats.WorstMoves {
		vstats.WorstMoves = moves
	}
	vstats.SumMoves += moves

	s.Save()

	return fmt.Sprintf("Recording completed game of %s", v)
}

func (s *Statistics) RecordLostGame(v string, percent int) string {

	vstats := s.findVariant(v)

	vstats.Lost = vstats.Lost + 1
	// don't see that currStreak can ever be zero
	if vstats.CurrStreak > 0 {
		vstats.CurrStreak = -1
	} else {
		vstats.CurrStreak = vstats.CurrStreak - 1
	}
	if vstats.CurrStreak < vstats.WorstStreak {
		vstats.WorstStreak = vstats.CurrStreak
	}

	if percent > vstats.BestPercent {
		vstats.BestPercent = percent
	}
	vstats.SumPercents += percent

	s.Save()

	return fmt.Sprintf("Recording lost game of %s, %d%% complete", v, percent)
}

func ShowStatisticsDrawer() {
	vstats := TheGame.Statistics.findVariant(TheGame.Baize.variant)
	var strs []string = vstats.strings(TheGame.Baize.variant)
	strs = append(strs, " ") // n.b. can't use empty string
	strs = append(strs, "ALL VARIANTS")
	strs = append(strs, TheGame.Statistics.strings()...)
	TheGame.UI.ShowTextDrawer(strs)
}
