package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mariolazzari/gator/internal/database"
)

func checkArgs(cmd command, n int) error {
	if len(cmd.Args) != n {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	return nil
}

func loginHandler(s *state, cmd command) error {
	// check param
	if err := checkArgs(cmd, 1); err != nil {
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
	if err := checkArgs(cmd, 1); err != nil {
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

func resetHandler(s *state, cmd command) error {

	// reset users table
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't deelte users: %w", err)
	}

	// reset current user
	err = s.cfg.SetUser("")
	if err != nil {
		return fmt.Errorf("couldn't reser current user: %w", err)
	}

	fmt.Printf("users table successfully truncated")

	return nil
}

func usersHandler(s *state, cmd command) error {

	// reset users table
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't load users: %w", err)
	}

	for _, user := range users {
		label := user.Name
		if label == s.cfg.CurrentUserName {
			label = fmt.Sprintf("%s (current)", label)
		}
		fmt.Println(label)
	}

	return nil
}

func aggHandler(s *state, cmd command) error {

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't load users: %w", err)
	}

	fmt.Println(*feed)

	return nil
}

func addFeedHandler(s *state, cmd command) error {
	// check params
	if err := checkArgs(cmd, 2); err != nil {
		return err
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	now := time.Now()

	ctx := context.Background()

	// get current user
	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't load current user: %w", err)
	}

	// save feed to database
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		Name:      name,
		Url:       url,
		CreatedAt: now,
		UpdatedAt: now,
	})

	fmt.Println(feed)

	return nil
}
