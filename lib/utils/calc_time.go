package utils

import (
	"fmt"
	"time"
)

// CalcAvgTime gets called when the program is about to exit
func CalcAvgTime(durations []time.Duration) {
	if len(durations) > 0 {
		var sum time.Duration
		for _, d := range durations {
			sum += d
		}
		fmt.Printf(
			"Average GET time: %v\n",
			time.Duration(int(sum)/len(durations)),
		)
	}
}
