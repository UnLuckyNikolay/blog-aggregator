package rssfeed

import (
	"context"
	"fmt"

	"github.com/UnLuckyNikolay/blog-aggregator/internal/state"
)

func Agg(state *state.State) error {
	feed, err := FetchFeed(state, context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err

	}
	fmt.Println(feed)

	return nil
}
