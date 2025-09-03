package commands

import (
	"fmt"

	"github.com/UnLuckyNikolay/blog-aggregator/internal/state"
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
	c.register("agg", handlerAgg)
	c.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	c.register("feeds", middlewareLoggedIn(handlerFeeds))
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("following", middlewareLoggedIn(handlerFollowing))
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
		return fmt.Errorf("Command `%s` not found.", cmd.name)
	}

	err = c.list[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}
