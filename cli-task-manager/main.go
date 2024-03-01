package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

type CLI struct {
	DB   *bolt.DB
	Name string
}

type Task struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		printHelp("./doc/help.txt")
	}

	db, err := bolt.Open("tasks.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	cli := CLI{DB: db, Name: "tasks"}

	help := isHelpCommand(args)

	switch args[0] {
	case "add":
		if help {
			printHelp("./doc/add.txt")
		}
		txt := strings.TrimSpace(strings.Join(args[1:], " "))
		if txt == "" {
			log.Fatal("Invalid task.")
		}
		err = cli.add(txt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Added \"%v\" to your task list.\n", txt)

	case "list":
		if help {
			printHelp("./doc/list.txt")
		}
		list, err := cli.list(nil)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("You have the following tasks:\n\n")
		for _, v := range list {
			fmt.Printf("%v, %v, %v, %v\n", v.ID, v.Text, v.Done, v.UpdatedAt)
		}
		fmt.Println()

	case "do":
		if help {
			printHelp("./doc/do.txt")
		}
		idStr := args[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		err = cli.do(id)
		if err != nil {
			log.Fatal(err)
		}
	case "rm":
		if help {
			printHelp("./doc/do.txt")
		}
		idStr := args[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		err = cli.rm(id)
		if err != nil {
			log.Fatal(err)
		}
	case "completed":
		if help {
			printHelp("./doc/completed.txt")
		}

		oneDay := 24 * time.Hour
		today := time.Now().Truncate(oneDay)

		list, err := cli.list(func(item Task) bool {
			a := item.UpdatedAt.Truncate(oneDay)
			return item.Done && a == today
		})

		if err != nil {
			log.Fatal(err)
		}
		if len(list) == 0 {
			fmt.Printf("You have not finished any Tasks today.\n\n")
		} else {
			fmt.Printf("You have finished the following tasks today:\n")
			for _, v := range list {
				fmt.Printf("- %v\n", v.Text)
			}
			fmt.Println()
		}

	default:
		log.Fatal("Invalid command.")
	}

}
