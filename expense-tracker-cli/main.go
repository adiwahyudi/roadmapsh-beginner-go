package main

import (
	"expense-tracker-cli/expense"
	"flag"
	"fmt"
	"os"
)

const (
	fileNameJSON = "expenses.json"
	fileNameCSV  = "expenses.csv"
)

func main() {
	fmt.Println("Expense Tracker CLI")
	data := &expense.Data{}

	if err := data.Load(fileNameJSON); err != nil {
		fmt.Println(err)
		return
	}

	subcommand := os.Args[1]
	subcommandFlags := flag.NewFlagSet(subcommand, flag.ExitOnError)

	// Define flag
	id := subcommandFlags.Int("id", 0, "Expense ID")
	description := subcommandFlags.String("description", "", "Expense description")
	category := subcommandFlags.String("category", "", "Expense category")
	filter := subcommandFlags.String("filter", "", "Filter expenses")
	amount := subcommandFlags.Int("amount", 0, "Expense amount")
	month := subcommandFlags.Int("month", -1, "Month for summary")
	budget := subcommandFlags.Int("budget", 0, "Maximum budget")

	if err := subcommandFlags.Parse(os.Args[2:]); err != nil {
		fmt.Println(err)
		return
	}

	// Menu
	switch subcommand {
	case "add":
		if err := data.Add(*description, *category, *amount); err != nil {
			fmt.Println(err)
			return
		}
		if err := data.Store(fileNameJSON); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Expense added successfully (ID: %d)", len(data.Expenses))
	case "list":
		data.List(*filter)
	case "update":
		if err := data.Update(*id, *description, *category, *amount); err != nil {
			fmt.Println(err)
			return
		}
		if err := data.Store(fileNameJSON); err != nil {
			fmt.Println(err)
			return
		}
	case "delete":
		if err := data.Delete(*id); err != nil {
			fmt.Println(err)
			return
		}
		if err := data.Store(fileNameJSON); err != nil {
			fmt.Println(err)
			return
		}
	case "summary":
		data.Summary(*month)
	case "budget":
		if err := data.SetBudget(*month, *budget); err != nil {
			fmt.Println(err)
			return
		}
		if err := data.Store(fileNameJSON); err != nil {
			fmt.Println(err)
			return
		}
	case "export":
		if err := data.ExportCSV(fileNameCSV); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Exported to %s", fileNameCSV)

	default:
		fmt.Println("Invalid command")
	}

}
