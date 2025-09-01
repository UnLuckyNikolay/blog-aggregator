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

type command struct {
	name string
	args []string
}

type CommandManager struct {
	list map[string]func(*state.State, command) error
}

func (c *CommandManager) Initialize() {
	c.list = map[string]func(*state.State, command) error{}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
}

func (c *CommandManager) HandleCommand(s *state.State, osArgs []string) error {
	if len(osArgs) < 2 {
		return fmt.Errorf("Not enough arguments.")
	}

	cmd := command{
		name: osArgs[1],
		args: []string{},
	}
	if len(osArgs) > 2 {
		cmd.args = osArgs[2:]
	}

	err := c.run(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandManager) register(name string, f func(*state.State, command) error) {
	c.list[name] = f
}

func (c *CommandManager) run(s *state.State, cmd command) error {
	var err error

	_, ok := c.list[cmd.name]
	if !ok {
		return fmt.Errorf("Command `%s`not found.", cmd.name)
	}

	err = c.list[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state.State, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Login command requires 1 argument: username.")
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
		return errors.New("Register command requires 1 argument: username.")
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

	fmt.Printf("User `%s` successfully registered.\n", name)
	fmt.Println(user)
	err = s.Cfg.SetUser(name)
	if err != nil {
		return err
	}

	return nil
}

func handlerReset(s *state.State, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("Reset command requires no arguments.")
	}

	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func handlerUsers(s *state.State, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("Users command requires no arguments.")
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
		return errors.New("Agg command requires no arguments.")
	}

	err := rssfeed.Agg(s)
	if err != nil {
		return err
	}

	return nil
}
