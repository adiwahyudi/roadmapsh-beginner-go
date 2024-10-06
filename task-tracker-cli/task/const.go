package task

const (
	Todo       Status = "todo"
	InProgress Status = "in-progress"
	Done       Status = "Done"

	FileName = "tasks.json"
)

func ValidateStatus(status Status) bool {
	switch status {
	case Todo, InProgress, Done:
		return true
	}
	return false
}
