package expense

import "time"

type Expense struct {
	ID          int       `json:"id"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Date        time.Time `json:"date"`
}

type Budget struct {
	Amount        int `json:"amount"`
	TotalExpenses int `json:"totalExpenses"`
}

/*
	Budget.Amount [LEGEND]
	-1: No budget set
	>0 : Budget set
*/

type Data struct {
	Expenses []Expense      `json:"expenses"`
	Budget   map[int]Budget `json:"budget"`
}
