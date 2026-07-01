package main

import (
	"aggregator/internal/config"
	"aggregator/internal/database"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)
	s := &state{db: dbQueries, cfg: &cfg}

	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))

	args := os.Args
	if len(args) < 2 {
		fmt.Println("error: not enough arguments, usage: gator <command> [args]")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	if err := cmds.run(s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
