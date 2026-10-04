package main

import (
	"context"
)

type Handler = func(*state, command) error

func middlewareLoggedIn(handler Handler) Handler {
	return func(s *state, cmd command) error {
		_, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd)
	}
}
