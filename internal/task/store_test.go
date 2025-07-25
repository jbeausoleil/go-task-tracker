package task

import (
	"os"
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {
	tmpFile := "test_tasks.json"
	defer os.Remove(tmpFile)

	store := NewStore(tmpFile)
	task := store.AddTask("Test Task", time.Now().Add(24*time.Hour))

	if task.Description != "Test Task" {
		t.Errorf("expected 'Test Task', got %s", task.Description)
	}

	if len(store.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(store.Tasks))
	}
}

func TestMarkTaskDone(t *testing.T) {
	tmpFile := "test_tasks.json"
	defer os.Remove(tmpFile)

	store := NewStore(tmpFile)
	task := store.AddTask("Test Done Task", time.Now().Add(24*time.Hour))

	store.MarkTaskDone(task.ID)

	if store.Tasks[0].Status != StatusDone {
		t.Errorf("expected StatusDone, got %s", store.Tasks[0].Status)
	}
}

func TestAddNoteToTask(t *testing.T) {
	tmpFile := "test_tasks.json"
	defer os.Remove(tmpFile)

	store := NewStore(tmpFile)
	task := store.AddTask("Test Note Task", time.Now().Add(24*time.Hour))

	store.AddNoteToTask(task.ID, "This is a note")

	if len(store.Tasks[0].Notes) != 1 {
		t.Errorf("expected 1 note, got %d", len(store.Tasks[0].Notes))
	}
}

func TestDeleteNoteFromTask(t *testing.T) {
	tmpFile := "test_tasks.json"
	defer os.Remove(tmpFile)

	store := NewStore(tmpFile)
	task := store.AddTask("Test Delete Note", time.Now().Add(24*time.Hour))

	store.AddNoteToTask(task.ID, "Note 1")
	store.AddNoteToTask(task.ID, "Note 2")
	store.DeleteNoteFromTask(task.ID)

	if len(store.Tasks[0].Notes) != 1 {
		t.Errorf("expected 1 note, got %d", len(store.Tasks[0].Notes))
	}
}
