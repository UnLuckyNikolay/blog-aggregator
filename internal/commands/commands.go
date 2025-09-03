package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/UnLuckyNikolay/blog-aggregator/internal/database"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/rssfeed"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/state"
	"github.com/google/uuid"
)

func handlerLogin(s *state.State, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Command `login` requires 1 argument: username.")
	}
	name := cmd.args[0]

	_, err := s.Db.GetUser(context.Background(), name)
	if err != nil {
		log.Fatal(1)
	}

	err = s.Cfg.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Printf("Current user set to %s.\n", name)
	return nil
}

func handlerRegister(s *state.State, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Command `register` requires 1 argument: username.")
	}
	name := cmd.args[0]

	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User `%s` successfully registered.\n", user.Name)
	//fmt.Println(user)
	err = s.Cfg.SetUser(name)
	if err != nil {
		return err
	}

	return nil
}

func handlerReset(s *state.State, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("Command `reset` requires no arguments.")
	}

	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("The tables have been reset.")

	return nil
}

func handlerUsers(s *state.State, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("Command `users` requires no arguments.")
	}

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		var currentCheck string
		if user.Name == s.Cfg.CurrentUserName {
			currentCheck = " (current)"
		} else {
			currentCheck = ""
		}
		fmt.Printf("* %s%s\n", user.Name, currentCheck)
	}

	return nil
}

func handlerAgg(s *state.State, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("Command `agg` requires no arguments.")
	}

	err := rssfeed.Agg(s)
	if err != nil {
		return err
	}

	return nil
}

func handlerAddfeed(s *state.State, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("Command `addfeed` requires 2 arguments: name, url.")
	}

	feed, err := s.Db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully added new feed `%s` (%s)\n", feed.Name, feed.Url)

	err = handlerFollow(s, command{
		name: "follow",
		args: []string{cmd.args[1]},
	}, user)
	if err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *state.State, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return errors.New("Command `feeds` requires no arguments.")
	}

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("> RSS Name: %s\n", feed.Name)
		fmt.Printf("  URL: %s\n", feed.Url)
		fmt.Printf("  Added by user: %s\n", feed.UserName)
	}

	return nil
}

func handlerFollow(s *state.State, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("Command `follow` requires 1 argument: url.")
	}

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedUrl:   cmd.args[0],
	})
	if err != nil {
		return err
	}

	fmt.Printf("User %s now follows '%s' feed.\n", user.Name, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state.State, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return errors.New("Command `following` requires 0 arguments.")
	}

	feedList, err := s.Db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	if len(feedList) == 0 {
		fmt.Printf("%s currently doesn't follow any feeds.\n", user.Name)
	} else {
		fmt.Printf("%s's followed feeds:\n", user.Name)
		for _, feed := range feedList {
			fmt.Printf("> %s\n", feed)
		}
	}

	return nil
}

func handlerUnfollow(s *state.State, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("Command `unfollow` requires 1 argument: feed url.")
	}

	deletedFeed, err := s.Db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID:  user.ID,
		FeedUrl: cmd.args[0],
	})
	if err != nil {
		return err
	}

	fmt.Printf("User %s unfollowed `%s` feed.\n", user.Name, deletedFeed.Name)

	return nil
}
