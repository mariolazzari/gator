package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/mariolazzari/gator/internal/config"
	"github.com/mariolazzari/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Open db connection
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error DB connection: %v", err)
	}
	dbQueries := database.New(db)

	appState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", loginHandler)
	cmds.register("register", registerHandler)
	cmds.register("reset", resetHandler)
	cmds.register("users", usersHandler)
	cmds.register("agg", aggHandler)
	cmds.register("addfeed", middlewareLoggedIn(addFeedHandler))
	cmds.register("feeds", feedsHandler)
	cmds.register("follow", middlewareLoggedIn(followHandler))
	cmds.register("following", middlewareLoggedIn(followingHandler))
	cmds.register("unfollow", middlewareLoggedIn(unfollowHandler))
	cmds.register("browse", browseHandler)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	err = cmds.run(appState, command{Name: os.Args[1], Args: os.Args[2:]})
	if err != nil {
		log.Fatal(err)
	}
}
