package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/arturogood17/aggreGator/internal/config"

	"github.com/arturogood17/aggreGator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	programS := &state{
		cfg: &cfg,
	}
	comnds := commands{
		registeredcmds: make(map[string]func(*state, command) error),
	}

	db, err := sql.Open("postgres", programS.cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	programS.db = dbQueries

	comnds.register("login", handlerLogin)
	comnds.register("register", registerLogin)
	comnds.register("reset", dbDelete)
	comnds.register("users", getUsers)
	comnds.register("agg", Aggregation)
	comnds.register("addfeed", addFeed)
	comnds.register("feeds", feedList)
	comnds.register("follow", followFeed)
	comnds.register("following", followedList)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	if err := comnds.run(programS, command{name: cmdName, args: cmdArgs}); err != nil {
		log.Fatal(err)
	}
}
