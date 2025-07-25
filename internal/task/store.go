package task

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

// Store holds all tasks and manages persistence.
type Store struct {
	FilePath string
	Tasks    []Task
}

// NewStore initializes the store and loads tasks from file.
func NewStore(path string) *Store {
	s := &Store{FilePath: path}
	s.load()
	return s
}

func (s *Store) load() {
	data, err := os.ReadFile(s.FilePath)
	if os.IsNotExist(err) {
		s.Tasks = []Task{}
		return
	}
	if err != nil {
		fmt.Println("failed to load tasks:", err)
		s.Tasks = []Task{}
		return
	}
	_ = json.Unmarshal(data, &s.Tasks)
}

// Save persists the tasks to the file.
func (s *Store) Save() {
	data, err := json.MarshalIndent(s.Tasks, "", "  ")
	if err != nil {
		fmt.Println("failed to Save tasks:", err)
		return
	}
	_ = os.WriteFile(s.FilePath, data, 0644)
}

func (s *Store) AddTask(desc string, due time.Time) Task {
	task := Task{
		ID:          uuid.NewString(),
		Description: desc,
		Status:      StatusPending,
		DueDate:     due,
		CreatedAt:   time.Now(),
		Notes:       []string{},
	}
	s.Tasks = append(s.Tasks, task)
	s.Save()
	return task
}

func (s *Store) MarkTaskDone(id string) {
	for i, t := range s.Tasks {
		if t.ID == id {
			s.Tasks[i].Status = StatusDone
			s.Tasks[i].CompletedAt = time.Now()
			break
		}
	}
	s.Save()
}

func (s *Store) AddNoteToTask(id, note string) {
	for i, t := range s.Tasks {
		if t.ID == id {
			s.Tasks[i].Notes = append(s.Tasks[i].Notes, note)
			break
		}
	}
	s.Save()
}

func (s *Store) DeleteNoteFromTask(id string) {
	for i, t := range s.Tasks {
		if t.ID == id && len(s.Tasks[i].Notes) > 0 {
			s.Tasks[i].Notes = s.Tasks[i].Notes[:len(s.Tasks[i].Notes)-1]
			break
		}
	}
	s.Save()
}

// DeleteTask removes a task by its ID and saves the updated store.
func (s *Store) DeleteTask(id string) {
	for i, t := range s.Tasks {
		if t.ID == id {
			s.Tasks = append(s.Tasks[:i], s.Tasks[i+1:]...)
			break
		}
	}
	s.Save()
}

// DeleteNoteFromTaskAtIndex deletes a note from a specific task at the given index.
func (s *Store) DeleteNoteFromTaskAtIndex(taskID string, noteIndex int) {
	for ti, t := range s.Tasks {
		if t.ID == taskID {
			if noteIndex >= 0 && noteIndex < len(t.Notes) {
				s.Tasks[ti].Notes = append(t.Notes[:noteIndex], t.Notes[noteIndex+1:]...)
			}
			break
		}
	}
	s.Save()
}

// UpdateNoteAtIndex updates a note's content at a given index for a specific task.
func (s *Store) UpdateNoteAtIndex(taskID string, noteIndex int, newContent string) {
	for ti, t := range s.Tasks {
		if t.ID == taskID {
			if noteIndex >= 0 && noteIndex < len(t.Notes) {
				s.Tasks[ti].Notes[noteIndex] = newContent
			}
			break
		}
	}
	s.Save()
}

// GetTaskByID retrieves a task by its ID.
func (s *Store) GetTaskByID(id string) (Task, bool) {
	for _, t := range s.Tasks {
		if t.ID == id {
			return t, true
		}
	}
	return Task{}, false
}
