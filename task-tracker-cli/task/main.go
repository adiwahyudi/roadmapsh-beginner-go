package task

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func (t *Tasks) Load(fileName string) error {
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

	if err = json.Unmarshal(file, t); err != nil {
		return err
	}

	return nil
}

func (t *Tasks) Store(fileName string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func (t *Tasks) List(status *Status) error {
	if len(*t) == 0 {
		fmt.Println("No tasks, create a new one!")
		return nil
	}
	//log.Println(*status)
	for i, task := range *t {
		if status != nil {
			if task.Status == *status {
				fmt.Printf("%d. %s\t\t%s\t\t%v\t\t%v\n", i+1, task.Description, task.Status, task.CreatedAt.Format(time.ANSIC), task.UpdatedAt.Format(time.ANSIC))
				continue
			}
		} else {
			fmt.Printf("%d. %s\t\t%s\t\t%v\t\t%v\n", i+1, task.Description, task.Status, task.CreatedAt.Format(time.ANSIC), task.UpdatedAt.Format(time.ANSIC))
		}
	}

	return nil
}

func (t *Tasks) Add(description string) {
	if len(description) < 0 {
		fmt.Println("Please provide a description.")
		return
	}
	task := Task{
		ID:          len(*t),
		Description: description,
		Status:      Todo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	*t = append(*t, task)

}

func (t *Tasks) Update(id int, description string) {
	if len(description) < 0 {
		fmt.Println("Please provide a description.")
		return
	}
	(*t)[id-1].Description = description
	(*t)[id-1].UpdatedAt = time.Now()
}

func (t *Tasks) Delete(id int) {
	*t = append((*t)[:id-1], (*t)[id:]...)
	for i := range *t {
		(*t)[i].ID = i
	}
}

func (t *Tasks) ChangeStatus(id int, status Status) {
	(*t)[id-1].Status = status
	(*t)[id-1].UpdatedAt = time.Now()
}
