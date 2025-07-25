package task

import (
    "encoding/json"
    "os"
    "time"
)

type Task struct {
    Description string    `json:"description"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    DueDate     time.Time `json:"due_date"`
    CompletedAt time.Time `json:"completed_at,omitempty"`
}

type Store struct {
    filePath string
}

func NewStore(path string) *Store {
    return &Store{filePath: path}
}

func (s *Store) Load() ([]Task, error) {
    if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
        return []Task{}, nil
    }
    data, err := os.ReadFile(s.filePath)
    if err != nil {
        return nil, err
    }
    var tasks []Task
    if err := json.Unmarshal(data, &tasks); err != nil {
        return nil, err
    }
    return tasks, nil
}

func (s *Store) Save(tasks []Task) error {
    data, err := json.MarshalIndent(tasks, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(s.filePath, data, 0644)
}

func CompletedToday(tasks []Task) []Task {
    today := time.Now().Truncate(24 * time.Hour)
    var result []Task
    for _, t := range tasks {
        if t.Status == "done" && t.CompletedAt.After(today) {
            result = append(result, t)
        }
    }
    return result
}

func DueThisWeek(tasks []Task) []Task {
    now := time.Now()
    _, currentWeek := now.ISOWeek()
    var result []Task
    for _, t := range tasks {
        _, dueWeek := t.DueDate.ISOWeek()
        if dueWeek == currentWeek {
            result = append(result, t)
        }
    }
    return result
}