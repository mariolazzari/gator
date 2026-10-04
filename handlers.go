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
	if err := checkArgs(cmd, 1); err != nil {
		return err
	}

	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user %q does not exist", name)
	}

	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User '%s' switched successfully\n", name)

	return nil
}

func registerHandler(s *state, cmd command) error {
	if err := checkArgs(cmd, 1); err != nil {
		return err
	}

	name := cmd.Args[0]
	now := time.Now()

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User '%s' switched successfully\n", name)

	return nil
}

func resetHandler(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}

	if err := s.cfg.SetUser(""); err != nil {
		return fmt.Errorf("couldn't reset current user: %w", err)
	}

	fmt.Println("users table successfully truncated")

	return nil
}

func usersHandler(s *state, cmd command) error {
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
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	fmt.Println(*feed)

	return nil
}

func addFeedHandler(s *state, cmd command) error {
	if err := checkArgs(cmd, 2); err != nil {
		return err
	}

	name := cmd.Args[0]
	url := cmd.Args[1]
	now := time.Now()
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't load current user: %w", err)
	}

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		Name:      name,
		Url:       url,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println(feed)

	return nil
}

func feedsHandler(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't load feeds: %w", err)
	}

	for _, feed := range feeds {
		label := fmt.Sprintf("%s %s %s", feed.Name, feed.Url, feed.User)
		fmt.Println(label)
	}

	return nil
}

func followHandler(s *state, cmd command) error {
	if err := checkArgs(cmd, 1); err != nil {
		return err
	}

	url := cmd.Args[0]
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't load current user: %w", err)
	}

	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't load feed by url: %w", err)
	}

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Printf("%s %s\n", feed.Name, user.Name)

	return nil
}

func followingHandler(s *state, cmd command) error {
	ctx := context.Background()

	follows, err := s.db.GetFeedFollowsForUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't load feeds followed by user: %w", err)
	}

	for _, f := range follows {
		fmt.Printf("%s %s\n", f.FeedName, f.UserName)
	}

	return nil
}
