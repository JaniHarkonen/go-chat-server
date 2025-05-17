package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/JaniHarkonen/go-chat-server/internal/command"
)

func commandExecutor(server *Server) func(commString *string) error {
	verifyArgument := func(argument *command.Argument, argType int) (bool, *command.Argument) {
		return argument.ArgType() == argType, argument
	}

	var commandHandlers = make(map[string]func(comm *command.Command))

	// Kick user
	commandHandlers["/kick"] = func(c *command.Command) {
		ok, username := verifyArgument(c.GetArgument(0), command.TypeUser)

		if ok {
			client := server.ResolveClient(server.chatManager.FindUserByName(username.AsString()))

			if client != nil {
				client.connection.Close()
			}
		}
	}

	commandHandlers["/mute"] = func(c *command.Command) {
		okUsername, username := verifyArgument(c.GetArgument(0), command.TypeUser)
		okDuration, duration := verifyArgument(c.GetArgument(1), command.TypeNumber)

		if okUsername && okDuration {
			resolvedUser := server.chatManager.FindUserByName(username.AsString())
			vDuration, err := duration.AsInt()

			if err != nil {
				fmt.Println("cannot execute /mute, invalid duration '" + *duration.Data())
				return
			}

			server.chatManager.MuteUser(resolvedUser, time.Duration(vDuration*time.Second.Nanoseconds()))
		}
	}

	return func(commString *string) error {
		if command, err := command.Parse(*commString); err != nil {
			return err
		} else {
			if handler, ok := commandHandlers[*command.Name()]; !ok {
				return errors.New("attempting to execute non-existing command '" + *command.Name() + "'")
			} else {
				handler(command)
				return nil
			}
		}
	}
}
