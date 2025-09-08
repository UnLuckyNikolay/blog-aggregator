package commands

import (
	"context"
	"fmt"

	"github.com/UnLuckyNikolay/blog-aggregator/internal/database"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/state"
)

// middlewareLoggedIn is used to wrap a handler function.
// Requires the user to be logged in to run the function.
func middlewareLoggedIn(handler func(s *state.State, cmd command, user database.User) error) func(*state.State, command) error {
	return func(s *state.State, cmd command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Please `login` to use this command. (error: %v)", err)
		}

		return handler(s, cmd, user)
	}
}
