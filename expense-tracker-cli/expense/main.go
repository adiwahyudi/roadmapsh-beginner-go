package expense

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

func (d *Data) Load(fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			newFile, newFileErr := os.Create(fileName)
			if newFileErr != nil {
				return newFileErr
			}
			newFile.Close()
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	if err = json.Unmarshal(file, d); err != nil {
		return err
	}

	return nil
}

func (d *Data) Store(fileName string) error {
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func (d *Data) Add(description string, category string, amount int) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if len(description) == 0 {
		return errors.New("description cannot be empty")
	}

	if len(category) == 0 {
		return errors.New("category cannot be empty")
	}

	expense := Expense{
		ID:          len(d.Expenses) + 1,
		Category:    category,
		Description: description,
		Amount:      amount,
		Date:        time.Now(),
	}

	d.Expenses = append(d.Expenses, expense)

	todayMonth := int(time.Now().Month())
	if d.Budget == nil {
		d.Budget = make(map[int]Budget)
	}
	totalExpenses := amount
	budgetNow := 0
	if budget, ok := d.Budget[todayMonth]; ok {
		budget.TotalExpenses += amount
		totalExpenses = budget.TotalExpenses
		budgetNow = budget.Amount
	} else {
		d.Budget[todayMonth] = Budget{
			Amount:        -1,
			TotalExpenses: amount,
		}
		budgetNow = -1
	}

	if budgetNow != -1 {
		if totalExpenses > budgetNow {
			fmt.Printf("You have exceeded your budget of $%d\n", budgetNow)
			fmt.Printf("Total expenses: $%d\n", totalExpenses)
		}
	}

	return nil
}

func (d *Data) Update(id int, description string, category string, amount int) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if len(description) == 0 {
		return errors.New("description cannot be empty")
	}

	if len(category) == 0 {
		return errors.New("category cannot be empty")
	}

	d.Expenses[id-1].Category = category
	d.Expenses[id-1].Description = description
	d.Expenses[id-1].Amount = amount

	return nil
}

func (d *Data) Delete(id int) error {
	if id > len(d.Expenses) || id < 1 {
		return errors.New("invalid id")
	}
	d.Expenses = append(d.Expenses[:id-1], d.Expenses[id:]...)
	return nil
}

func (d *Data) List(category string) {
	if len(d.Expenses) == 0 {
		fmt.Println("No expenses, create a new one!")
		return
	}
	for _, expense := range d.Expenses {
		if category != "" && expense.Category != category {
			continue
		}
		fmt.Printf("#%d | %s | %s | %s | $%d\n", expense.ID, expense.Date.Format(time.DateOnly), expense.Category, expense.Description, expense.Amount)
	}
}
func (d *Data) Summary(month int) {
	withMonth := false
	if month == -1 {
		// skip this block
	} else if month < 1 || month > 12 {
		fmt.Println("Invalid month")
		return
	} else {
		withMonth = true
	}

	total := 0
	for _, expense := range d.Expenses {
		if withMonth {
			if int(expense.Date.Month()) == month {
				total += expense.Amount
			}
			continue
		}
		total += expense.Amount
	}
	fmt.Printf("Total expenses: $%d\n", total)
}

func (d *Data) SetBudget(month, amount int) error {
	if month < 1 || month > 12 {
		return errors.New("invalid month")
	}
	if amount <= 0 {
		return errors.New("budget must be greater than 0")
	}

	if d.Budget == nil {
		d.Budget = make(map[int]Budget)
	}

	if budget, ok := d.Budget[month]; ok {
		budget.Amount = amount
	} else {
		d.Budget[month] = Budget{
			Amount:        amount,
			TotalExpenses: 0,
		}
	}
	fmt.Printf("Budget set for month %s (%d): $%d\n", time.Month(month).String(), month, amount)
	return nil
}

func (d *Data) ExportCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	header := []string{"ID", "Date", "Category", "Description", "Amount"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, expense := range d.Expenses {
		//addQuote := func(data string) string {
		//	return fmt.Sprintf(`"%s"`, data)
		//}
		record := []string{
			strconv.Itoa(expense.ID),
			expense.Category,
			expense.Description,
			strconv.Itoa(expense.Amount),
			expense.Date.String(),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
