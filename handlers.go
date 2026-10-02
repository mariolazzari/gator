package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mariolazzari/gator/internal/database"
)

func checkArgs(cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	return nil
}

func loginHandler(s *state, cmd command) error {
	// check param
	if err := checkArgs(cmd); err != nil {
		return err
	}
	name := cmd.Args[0]

	// get user from database by name
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user %q does not exist", name)
	}

	// save current user
	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User '%s' switched successfully\n", name)

	return nil
}

func registerHandler(s *state, cmd command) error {
	// check params
	if err := checkArgs(cmd); err != nil {
		return err
	}
	name := cmd.Args[0]
	now := time.Now()

	// insert user into database
	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	// save current user
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("User '%s' switched successfully\n", name)

	return nil
}
