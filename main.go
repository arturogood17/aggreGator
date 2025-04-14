package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/arturogood17/aggreGator/internal/config"
	"github.com/arturogood17/aggreGator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg     *config.Config
	Queries *database.Queries
}

type command struct {
	name  string
	flags []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
func (c *commands) run(s *state, cmd command) error {
	if f, exists := c.cmds[cmd.name]; !exists {
		return errors.New("this command does not exists")
	} else {
		f(s, cmd)
	}
	return nil
}

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	//Abro conexión a base de datos
	db, err := sql.Open("postgres", config.DbURL)
	if err != nil {
		log.Fatalf("Error open connection to the database - %v", err)
	}
	defer db.Close() //Siempre hay que cerrar la base de datos
	dbQueries := database.New(db)
	s := &state{
		cfg:     &config,
		Queries: dbQueries,
	}
	mapCommands := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	//Register segment
	mapCommands.register("login", handlerLogin)
	mapCommands.register("register", handlerRegister)
	mapCommands.register("reset", handlerReset)
	mapCommands.register("users", handlerAllUsers)
	mapCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	mapCommands.register("agg", handlerFeedFuncs)
	mapCommands.register("feeds", handlerListFeeds)
	mapCommands.register("follow", middlewareLoggedIn(handlerFollowFeed))
	mapCommands.register("following", middlewareLoggedIn(handlerFollowingFeeds))
	mapCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	mapCommands.register("browse", middlewareLoggedIn(handlerBrowsingPosts))

	//Run segment
	if len(os.Args) <= 1 {
		log.Fatal("not enough arguments were provided")
	}
	err = mapCommands.run(s, command{name: os.Args[1], flags: os.Args[2:]})
	if err != nil {
		log.Fatalf("Error running this command: %v. Error value: %v", os.Args[1], err)
	}
}
