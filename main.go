package main

import (
	"fmt"
	"os"

	"github.com/arturogood17/aggreGator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	} else {
		S := state{
			cfg: &cfg,
		}
		comnds := commands{
			cmds: make(map[string]func(*state, command) error),
		}
		comnds.register("login", handlerLogin)
		if len(os.Args) < 2 {
			fmt.Println("invalid input")
			os.Exit(1)
		}
		comnd := command{
			name: os.Args[1],
			args: os.Args[2:],
		}
		if err := comnds.run(&S, comnd); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
