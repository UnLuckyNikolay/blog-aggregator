package commandhandlers

import (
	"errors"
	"fmt"
)

type command struct {
	name string
	args []string
}

type CommandManager struct {
	list map[string]func(*State, command) error
}

func (c *CommandManager) Initialize() {
	c.list = map[string]func(*State, command) error{}
	c.register("login", handlerLogin)
}

func (c *CommandManager) HandleCommand(s *State, osArgs []string) error {
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

func (c *CommandManager) register(name string, f func(*State, command) error) {
	c.list[name] = f
}

func (c *CommandManager) run(s *State, cmd command) error {
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

func handlerLogin(s *State, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Login command requires 1 argument: username.")
	}

	err := s.Cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Current user set to %s.\n", cmd.args[0])
	return nil
}
