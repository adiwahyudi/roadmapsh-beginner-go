package main

import (
	"fmt"
	"os"
	"strconv"
	"task-tracker-cli/task"
)

func main() {
	fmt.Println("Welcome! This is a simple task tracker CLI.")
	tasks := &task.Tasks{}
	if err := tasks.Load(task.FileName); err != nil {
		fmt.Println(err)
		return
	}

	switch os.Args[1] {
	case "list":
		if len(os.Args) > 2 {
			status := task.Status(os.Args[2])
			if isValidStatus := task.ValidateStatus(status); !isValidStatus {
				fmt.Printf(`Invalid status '%s'`, status)
			}
			if err := tasks.List(&status); err != nil {
				fmt.Println(err)
			}
			return
		}
		if err := tasks.List(nil); err != nil {
			fmt.Println(err)
		}

	case "add":
		if len(os.Args) > 2 {
			description := os.Args[2]
			tasks.Add(description)
			if err := tasks.Store(task.FileName); err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Task added successfully (ID: %d)", len(*tasks))
		}
		fmt.Println("Please provide the description.")

	case "update":
		if len(os.Args) > 3 {
			id := getTaskID(os.Args[2])
			description := os.Args[3]
			tasks.Update(id, description)
			if err := tasks.Store(task.FileName); err != nil {
				fmt.Println(err)
			}
			fmt.Println("Task successfully updated")
			return
		}
		fmt.Println("Please provide task ID and description.")

	case "delete":
		if len(os.Args) > 2 {
			id := getTaskID(os.Args[2])
			tasks.Delete(id)
			if err := tasks.Store(task.FileName); err != nil {
				fmt.Println(err)
			}
			fmt.Println("Task successfully deleted")
			return
		}
		fmt.Println("Please provide a task ID.")

	case "mark-in-progress":
		if len(os.Args) > 2 {
			id := getTaskID(os.Args[2])
			tasks.ChangeStatus(id, task.InProgress)
			if err := tasks.Store(task.FileName); err != nil {
				fmt.Println(err)
			}
			fmt.Println("Task successfully mark as in-progress")
			return
		}
		fmt.Println("Please provide a task ID.")

	case "mark-done":
		if len(os.Args) > 2 {
			id := getTaskID(os.Args[2])
			tasks.ChangeStatus(id, task.Done)
			if err := tasks.Store(task.FileName); err != nil {
				fmt.Println(err)
			}
			fmt.Println("Task successfully mark as done")
			return
		}
		fmt.Println("Please provide a task ID.")

	default:
		fmt.Println("Invalid command.")
	}

}

func getTaskID(args string) int {
	idString := args
	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println("Please provide a valid task ID.")
	}
	return id
}
