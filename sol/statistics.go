package sol

import (
	"fmt"

	"oddstream.games/gosol/sound"
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
	Won, Lost, CurrStreak, BestStreak, WorstStreak, SumPercents, BestPercent int `json:",omitempty"`
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

func (stats *VariantStatistics) generalToasts() []string {

	v := TheBaize.LongVariantName()

	toasts := []string{}
	toasts = append(toasts,
		fmt.Sprintf("You have played %s %s (won %d, lost %d)", v, util.Pluralize("time", stats.Won+stats.Lost), stats.Won, stats.Lost))
	// fmt.Sprintf("You have won %s, and lost %s (%d%%)",
	// 	util.Pluralize("game", stats.Won),
	// 	util.Pluralize("game", stats.Lost),
	// 	((stats.Won*100)/(stats.Won+stats.Lost))))

	avpc := stats.averagePercent()
	if avpc > 0 && avpc < 100 {
		toasts = append(toasts, fmt.Sprintf("Your average score is %d%%", avpc))
	}

	if stats.CurrStreak > 1 {
		toasts = append(toasts, fmt.Sprintf("You are on a winning streak of %s", util.Pluralize("game", stats.CurrStreak)))
	}
	if stats.CurrStreak < 1 {
		toasts = append(toasts, fmt.Sprintf("You are on a losing streak of %s", util.Pluralize("game", util.Abs(stats.CurrStreak))))
	}

	return toasts
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
		stats = &VariantStatistics{} // everything 0
		s.StatsMap[v] = stats
		println("statistics has encountered a new variant", v)
	}
	return stats
}

func (s *Statistics) RecordWonGame(v string) {

	sound.Play("Complete")
	TheUI.Toast(fmt.Sprintf("Recording completed game of %s", v))

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

	toasts := stats.generalToasts()
	for _, t := range toasts {
		TheUI.Toast(t)
	}

	s.Save()
}

func (s *Statistics) RecordLostGame(v string) {

	percent := TheBaize.PercentComplete()
	if percent == 100 {
		println("*** That's odd, here is a lost game that is 100% complete ***")
	}

	TheUI.Toast(fmt.Sprintf("Recording lost game of %s, %d%% complete", v, percent))

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

func (s *Statistics) WelcomeToast(v string) {

	toasts := []string{}

	stats, ok := s.StatsMap[v]
	if !ok || stats.Won+stats.Lost == 0 {
		toasts = append(toasts, fmt.Sprintf("You have not played %s before", v))
	} else {
		avpc := stats.averagePercent()

		if stats.Won == 0 {
			toasts = append(toasts, fmt.Sprintf("You have yet to win a game of %s in %s", v, util.Pluralize("attempt", stats.Lost)))
			if stats.BestPercent > 0 && stats.BestPercent != avpc {
				toasts = append(toasts, fmt.Sprintf("Your best score is %d%%, your average score is %d%%", stats.BestPercent, avpc))
			}
		} else {
			toasts = stats.generalToasts()
		}
	}

	for _, t := range toasts {
		TheUI.Toast(t)
	}
}
