package fetch

import (
	"fmt"
	"time"

	"github.com/UnLuckyNikolay/blog-aggregator/internal/state"
)

func Agg(state *state.State, time_between_req time.Duration) {
	ticker := time.NewTicker(time_between_req)

	fmt.Printf("Fetching new posts every %v.\n\n", time_between_req)

	for ; ; <-ticker.C {
		fetchFeeds(state)
	}
}
