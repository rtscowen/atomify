package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type habit struct {
	name       string
	created    time.Time
	dailyAtoms int
	totalAtoms int
	breaks     []time.Time
}

func (h *habit) create() {
	fmt.Println(h)
}

func (h *habit) updateName(newName string) {
	fmt.Println(h)
}

func (h *habit) updateDaily(newDaily int) {
	fmt.Println(h)
}

func (h *habit) delete() {
	fmt.Println(h)
}

func (h *habit) increment() {
	fmt.Println(h)
}

func (h *habit) decrement() {
	fmt.Println(h)
}

func (h *habit) takeBreak() {
	fmt.Println(h)
}

func (h *habit) unbreak() {
	fmt.Println(h)
}

func list() {

}

func day() {

}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Usage: atomify <command> [args]")
		fmt.Println("atomify list - lists your currently tracked habits")
		fmt.Println("atomify new [name] - creates a new habit to track")
		return
	}

	//Handle flags for the "update" subcommand
	var newName string
	var newDaily int

	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)

	updateCommand.StringVar(&newName, "r", "", "New name for habit")
	updateCommand.IntVar(&newDaily, "d", 0, "New number of daily atoms")

	switch os.Args[1] {

	case "new":
		if len(os.Args) != 4 {
			fmt.Println("new takes two args [name] [units per day]")
		}

		dailyAtoms, err := strconv.Atoi(os.Args[3])
		if err != nil {
			os.Exit(2)
		}

		h := habit{os.Args[2], time.Now(), dailyAtoms, 0, nil}
		h.create()

	case "remove":
		if len(os.Args) != 3 {
			fmt.Println("remove takes one arg [name]")
		}

		h := habit{name: os.Args[2]}

		h.delete()

	case "list":
		if len(os.Args) != 2 {
			fmt.Println("list takes no arguemnts")
		}

		list()

	case "day":
		if len(os.Args) != 2 {
			fmt.Println("day takes no arguemnts")
		}

		day()

	case "++":
		if len(os.Args) != 3 {
			fmt.Println("++ takes one arg [habit]")
		}

		h := habit{name: os.Args[2]}
		h.increment()

	case "--":
		if len(os.Args) != 3 {
			fmt.Println("-- takes one arg [habit]")
		}

		h := habit{name: os.Args[2]}
		h.decrement()

	case "break":
		if len(os.Args) != 3 {
			fmt.Println("break takes one arg [habit]")
		}

		h := habit{name: os.Args[2]}
		h.takeBreak()

		fmt.Println("Everybody needs a break")

	case "unbreak":
		if len(os.Args) != 3 {
			fmt.Println("unbreak takes one arg [habit]")
		}

		h := habit{name: os.Args[2]}
		h.unbreak()

	case "update":
		updateCommand.Parse(os.Args[2:])

	default:
		fmt.Printf("%q is not a valid command\n", os.Args[1])
		os.Exit(2)
	}

	if updateCommand.Parsed() {
		if len(os.Args) < 4 || len(os.Args) > 5 {
			fmt.Println("atomify update [habit] -r=[new name] -d=[new daily]")
			fmt.Println("At least one of the flags is required")
		}

		h := habit{name: os.Args[3]}

		if newName != "" {
			h.updateName(newName)
		}

		if newDaily != 0 {
			h.updateDaily(newDaily)
		}

		fmt.Println("Habit updated")

	}

}
