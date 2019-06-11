package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const store = "habits.json"

type habit struct {
	Name       string
	Created    time.Time
	DailyAtoms int
	TotalAtoms int
	Breaks     []time.Time
}

func eraseAndWrite(f *os.File, h *[]habit) {
	if err := f.Truncate(0); err != nil {
		panic(err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		panic(err)
	}

	if *h != nil {
		for _, habit := range *h {
			data, err := json.Marshal(habit)
			if err != nil {
				panic(err)
			}

			data = append(data, byte('\n'))
			if _, err := f.Write(data); err != nil {
				panic(err)
			}
		}
	}
}

func (h *habit) create(f *os.File) {
	data, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}

	if _, err := f.Seek(0, 2); err != nil {
		panic(err)
	}

	if _, err := f.Write(append(data, byte('\n'))); err != nil {
		panic(err)
	}
}

func (h *habit) updateName(f *os.File, newName string) {
	var temp habit
	var habits []habit

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &temp); err != nil {
			panic(err)
		}

		if h.Name == temp.Name {
			temp.Name = newName
			habits = append(habits, temp)
		} else {
			habits = append(habits, temp)
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	eraseAndWrite(f, &habits)

	h.Name = newName
}

func (h *habit) updateDaily(f *os.File, newDaily int) {
	var temp habit
	var habits []habit

	//Ensure we are at start of file in case name is updated prior
	f.Seek(0, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &temp); err != nil {
			panic(err)
		}

		if h.Name == temp.Name {
			temp.DailyAtoms = newDaily

			//Reset total and created: a new daily target is effectively a new habit
			temp.TotalAtoms = 0
			temp.Created = time.Now()
			habits = append(habits, temp)
		} else {
			habits = append(habits, temp)
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	eraseAndWrite(f, &habits)

}

func (h *habit) delete(f *os.File) {
	var temp habit
	var habits []habit

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &temp); err != nil {
			panic(err)
		}

		if h.Name != temp.Name {
			habits = append(habits, temp)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	eraseAndWrite(f, &habits)

}

func (h *habit) increment(f *os.File) {
	var temp habit
	var habits []habit

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &temp); err != nil {
			panic(err)
		}

		if h.Name == temp.Name {
			temp.TotalAtoms++
			habits = append(habits, temp)
		} else {
			habits = append(habits, temp)
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	eraseAndWrite(f, &habits)
}

func (h *habit) decrement(f *os.File) {
	var temp habit
	var habits []habit

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &temp); err != nil {
			panic(err)
		}

		if h.Name == temp.Name {
			if temp.TotalAtoms > 0 {
				temp.TotalAtoms--
			}
			habits = append(habits, temp)
		} else {
			habits = append(habits, temp)
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	eraseAndWrite(f, &habits)
}

func (h *habit) takeBreak(f *os.File) {
	var temp habit
	var habits []habit

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &temp); err != nil {
			panic(err)
		}

		if h.Name == temp.Name {
			temp.Breaks = append(temp.Breaks, time.Now())
			habits = append(habits, temp)
		} else {
			habits = append(habits, temp)
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	eraseAndWrite(f, &habits)
}

func (h *habit) unbreak(f *os.File) {
	var temp habit
	var habits []habit

	today := time.Now()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &temp); err != nil {
			panic(err)
		}

		if h.Name == temp.Name {
			var newBreaks []time.Time
			for _, b := range temp.Breaks {
				if today.Day() != b.Day() && today.Month() != b.Month() && today.Year() != b.Year() {
					newBreaks = append(newBreaks, b)
				}
			}
			temp.Breaks = newBreaks
			habits = append(habits, temp)
		} else {
			habits = append(habits, temp)
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	eraseAndWrite(f, &habits)
}

func list(f *os.File) {
	var temp habit

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)

	fmt.Fprintf(w, "Name:\tCreated:\tDaily Requirement:\tTotal completed:\t\n")

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &temp); err != nil {
			panic(err)
		}

		var strs []string

		strs = append(strs, temp.Name)
		strs = append(strs, temp.Created.Format("02-Jan-2006"))
		strs = append(strs, strconv.Itoa(temp.DailyAtoms))
		strs = append(strs, strconv.Itoa(temp.TotalAtoms))
		strs = append(strs, "\n")

		fmt.Fprintf(w, strings.Join(strs, "\t"))

	}

	if err := w.Flush(); err != nil {
		panic(err)
	}
}

func day(f *os.File) {
	var temp habit

	fmt.Println("Todo list: ")

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &temp); err != nil {
			panic(err)
		}

		//If today is a break day for this habit, it is not on the todo list
		if temp.Breaks != nil {
			lastBreak := temp.Breaks[len(temp.Breaks)-1]
			if lastBreak.Day() != time.Now().Day() && lastBreak.Month() != time.Now().Month() && lastBreak.Year() != time.Now().Year() {
				temp.printTodo()
			}
		} else {
			temp.printTodo()
		}
	}
}

func (h *habit) printTodo() {
	daysBetween := (time.Now().Day() - h.Created.Day()) + 1

	expectedAtoms := h.DailyAtoms * daysBetween

	if h.TotalAtoms < expectedAtoms {
		fmt.Println("[+] ", h.Name, " x ", expectedAtoms-h.TotalAtoms)
	}
}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Usage: atomify <command> [args]")
		fmt.Println("atomify list - lists your currently tracked habits")
		fmt.Println("atomify new [name] - creates a new habit to track")
		return
	}

	f, err := os.OpenFile(store, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	//Handle flags for the "update" subcommand
	var newName string
	var newDaily int

	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)
	updateCommand.StringVar(&newName, "r", "", "New name for habit")
	updateCommand.IntVar(&newDaily, "d", 0, "New number of daily atoms")

	switch os.Args[1] {

	case "new":
		if len(os.Args) != 4 {
			fmt.Println("new takes two args [Name] [units per day]")
			return
		}

		dailyAtoms, err := strconv.Atoi(os.Args[3])
		if err != nil {
			panic(err)
		}

		h := habit{os.Args[2], time.Now(), dailyAtoms, 0, nil}
		h.create(f)

	case "remove":
		if len(os.Args) != 3 {
			fmt.Println("remove takes one arg [Name]")
			return
		}

		h := habit{Name: os.Args[2]}

		h.delete(f)

	case "list":
		if len(os.Args) != 2 {
			fmt.Println("list takes no arguemnts")
			return
		}

		list(f)

	case "day":
		if len(os.Args) != 2 {
			fmt.Println("day takes no arguemnts")
			return
		}

		day(f)

	case "++":
		if len(os.Args) != 3 {
			fmt.Println("++ takes one arg [habit]")
			return
		}

		h := habit{Name: os.Args[2]}
		h.increment(f)

	case "--":
		if len(os.Args) != 3 {
			fmt.Println("-- takes one arg [habit]")
			return
		}

		h := habit{Name: os.Args[2]}
		h.decrement(f)

	case "break":
		if len(os.Args) != 3 {
			fmt.Println("break takes one arg [habit]")
			return
		}

		h := habit{Name: os.Args[2]}
		h.takeBreak(f)

		fmt.Println("Everybody needs a break")

	case "unbreak":
		if len(os.Args) != 3 {
			fmt.Println("unbreak takes one arg [habit]")
			return
		}

		h := habit{Name: os.Args[2]}
		h.unbreak(f)

	case "update":
		updateCommand.Parse(os.Args[3:])

	default:
		fmt.Printf("%q is not a valid command\n", os.Args[1])
		return
	}

	if updateCommand.Parsed() {
		if len(os.Args) < 4 || len(os.Args) > 5 {
			fmt.Println("atomify update [habit] -r=[new name] -d=[new daily]")
			fmt.Println("At least one of the flags is required")
			return
		}

		h := habit{Name: os.Args[2]}

		if newName != "" {
			h.updateName(f, newName)
		}

		if newDaily != 0 {
			if !(newDaily > 0) {
				fmt.Println("Daily expectation must be greater than 0")
			}
			h.updateDaily(f, newDaily)
		}

		fmt.Println("Habit updated")

	}
}
