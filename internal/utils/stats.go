package utils

import (
	"sort"
)

func CalculateStatistics(scores []int) (mean, median float64, mode []int) {
	if len(scores) == 0 {
		return 0, 0, []int{}
	}

	sum := 0
	frequency := make(map[int]int)
	maxFreq := 0

	for _, score := range scores {
		sum += score
		frequency[score]++
		if frequency[score] > maxFreq {
			maxFreq = frequency[score]
		}
	}

	// Mean
	mean = float64(sum) / float64(len(scores))

	// Median
	sort.Ints(scores)
	n := len(scores)
	if n%2 == 0 {
		median = float64(scores[n/2-1]+scores[n/2]) / 2.0
	} else {
		median = float64(scores[n/2])
	}

	// Mode
	for score, freq := range frequency {
		if freq == maxFreq {
			mode = append(mode, score)
		}
	}

	return mean, median, mode
}
