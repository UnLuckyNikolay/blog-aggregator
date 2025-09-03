package commands

import (
	"context"
	"fmt"

	"github.com/UnLuckyNikolay/blog-aggregator/internal/database"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/state"
)

func middlewareLoggedIn(handler func(s *state.State, cmd command, user database.User) error) func(*state.State, command) error {
	return func(s *state.State, cmd command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Please `login` to use this command. (error: %v)", err)
		}

		return handler(s, cmd, user)
	}
}
