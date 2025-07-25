package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jinzhu/now"
)

type Status int

const (
	StatusPending Status = iota
	StatusDone
	StatusInProgress
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "todo"
	case StatusInProgress:
		return "in-progress"
	case StatusDone:
		return "done"
	default:
		return "unknown"
	}
}

type Task struct {
	ID             string
	Description    string
	Status         Status
	Notes          []string
	DueDate        time.Time
	CreatedAt      time.Time
	CompletedAt    time.Time
	LastModifiedAt time.Time
}

func (t *Task) UpdateDescription(desc string) {
	t.Description = desc
	t.LastModifiedAt = time.Now()
}

func (t *Task) UpdateDueDate(due time.Time) {
	t.DueDate = due
	t.LastModifiedAt = time.Now()
}

func (t *Task) UpdateStatus(status Status) {
	t.Status = status
	if status == StatusDone {
		t.CompletedAt = time.Now()
	}
	t.LastModifiedAt = time.Now()
}

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	doneStyle     = lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color("8"))
	helpStyle     = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("7"))
)

type state int

const (
	menu state = iota
	add
	addDue
	askNote
	addNote
	list
	viewNotes
	selectTaskNotes
	editNote
)

type model struct {
	cursor        int
	state         state
	choices       []string
	store         *Store
	textInput     textinput.Model
	dueInput      textinput.Model
	noteInput     textinput.Model
	tempDesc      string
	selectedTask  *Task
	numericBuffer string
	selectedNote  int
}

func InitialModel() model {
	descInput := textinput.New()
	descInput.Placeholder = "Task description"
	descInput.Focus()
	descInput.Width = 40

	dueInput := textinput.New()
	dueInput.Placeholder = "Due date (YYYY-MM-DD or 'tomorrow')"
	dueInput.Width = 35

	noteInput := textinput.New()
	noteInput.Placeholder = "Note content"
	noteInput.Width = 40

	store := NewStore("tasks.json")

	return model{
		choices:   []string{"Add Task", "List Tasks", "View Notes", "Quit"},
		state:     menu,
		store:     store,
		textInput: descInput,
		dueInput:  dueInput,
		noteInput: noteInput,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case menu:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				switch m.cursor {
				case 0:
					m.state = add
					m.textInput.Focus()
				case 1:
					m.state = list
					m.cursor = 0
				case 2:
					m.state = selectTaskNotes
					m.cursor = 0
				case 3:
					return m, tea.Quit
				}
			}
		}
	case add:
		m.textInput, cmd = m.textInput.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				m.tempDesc = m.textInput.Value()
				m.textInput.SetValue("")
				m.state = addDue
				m.dueInput.Focus()
			case "esc":
				m.textInput.SetValue("")
				m.state = menu
			}
		}
	case addDue:
		m.dueInput, cmd = m.dueInput.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				dueDate := parseDueDate(m.dueInput.Value())
				newTask := m.store.AddTask(m.tempDesc, dueDate)
				m.selectedTask = &newTask
				m.dueInput.SetValue("")
				m.state = askNote
			case "esc":
				m.dueInput.SetValue("")
				m.state = menu
			}
		}
	case askNote:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch strings.ToLower(key.String()) {
			case "y":
				m.state = addNote
				m.noteInput.Focus()
			case "n", "esc":
				m.state = menu
			}
		}
	case list:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "esc":
				m.state = menu
				m.cursor = 0
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.store.Tasks)-1 {
					m.cursor++
				}
			case "d":
				if len(m.store.Tasks) > 0 {
					m.store.MarkTaskDone(m.store.Tasks[m.cursor].ID)
				}
			case "i":
				if len(m.store.Tasks) > 0 {
					currentTask := &m.store.Tasks[m.cursor]
					if currentTask.Status != StatusInProgress {
						currentTask.Status = StatusInProgress
					} else {
						currentTask.Status = StatusPending
					}
					m.store.Save()
				}
			case "x":
				if len(m.store.Tasks) > 0 {
					m.store.DeleteTask(m.store.Tasks[m.cursor].ID)
					if m.cursor >= len(m.store.Tasks) && m.cursor > 0 {
						m.cursor--
					}
				}
			case "enter":
				if len(m.store.Tasks) > 0 {
					m.selectedTask = &m.store.Tasks[m.cursor]
					m.selectedNote = 0
					m.state = viewNotes
				}
			}
		}
	case selectTaskNotes:
		if key, ok := msg.(tea.KeyMsg); ok {
			if num, err := strconv.Atoi(key.String()); err == nil {
				idx := num - 1
				if idx >= 0 && idx < len(m.store.Tasks) {
					m.selectedTask = &m.store.Tasks[idx]
					m.selectedNote = 0
					m.state = viewNotes
				}
			}
			switch key.String() {
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.store.Tasks)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.store.Tasks) > 0 {
					m.selectedTask = &m.store.Tasks[m.cursor]
					m.selectedNote = 0
					m.state = viewNotes
				}
			case "esc":
				m.state = menu
			}
		}
	case addNote:
		m.noteInput, cmd = m.noteInput.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				if m.noteInput.Value() != "" && m.selectedTask != nil {
					m.store.AddNoteToTask(m.selectedTask.ID, m.noteInput.Value())
					m.store.Save()

					// Refresh selectedTask from updated store
					for i := range m.store.Tasks {
						if m.store.Tasks[i].ID == m.selectedTask.ID {
							m.selectedTask = &m.store.Tasks[i]
							break
						}
					}

					addedNote := m.noteInput.Value()
					if len(addedNote) > 20 {
						addedNote = addedNote[:20] + "..."
					}
					m.noteInput.SetValue("")
					return m, tea.Printf("Note added: %s", addedNote)
				}
			case "esc":
				m.noteInput.SetValue("")
				if m.selectedTask != nil {
					m.state = viewNotes
				} else {
					m.state = menu
					m.cursor = 0
				}
			}
		}
	case viewNotes:
		// Always refresh selectedTask from updated store before rendering notes
		for i := range m.store.Tasks {
			if m.store.Tasks[i].ID == m.selectedTask.ID {
				m.selectedTask = &m.store.Tasks[i]
				break
			}
		}
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "esc":
				m.state = list
				m.selectedNote = 0
				m.cursor = 0
			case "a":
				m.state = addNote
				m.noteInput.Focus()
			case "x":
				if m.selectedTask != nil && len(m.selectedTask.Notes) > 0 && m.selectedNote >= 0 && m.selectedNote < len(m.selectedTask.Notes) {
					m.store.DeleteNoteFromTaskAtIndex(m.selectedTask.ID, m.selectedNote)
					m.store.Save()
					if m.selectedNote >= len(m.selectedTask.Notes) && m.selectedNote > 0 {
						m.selectedNote--
					}
				}
			case "up":
				if m.selectedTask != nil && m.selectedNote > 0 {
					m.selectedNote--
				}
			case "down":
				if m.selectedTask != nil && m.selectedNote < len(m.selectedTask.Notes)-1 {
					m.selectedNote++
				}
			case "e":
				if m.selectedTask != nil && len(m.selectedTask.Notes) > 0 && m.selectedNote >= 0 && m.selectedNote < len(m.selectedTask.Notes) {
					m.state = editNote
					m.noteInput.SetValue(m.selectedTask.Notes[m.selectedNote])
					m.noteInput.Focus()
				}
			default:
				// Allow selecting note by number keys
				if num, err := strconv.Atoi(key.String()); err == nil {
					idx := num - 1
					if m.selectedTask != nil && idx >= 0 && idx < len(m.selectedTask.Notes) {
						m.selectedNote = idx
					}
				}
			}
		}
	case editNote:
		m.noteInput, cmd = m.noteInput.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				if m.noteInput.Value() != "" && m.selectedTask != nil && m.selectedNote >= 0 && m.selectedNote < len(m.selectedTask.Notes) {
					m.store.UpdateNoteAtIndex(m.selectedTask.ID, m.selectedNote, m.noteInput.Value())
					m.store.Save()
					m.noteInput.SetValue("")
					m.state = viewNotes
				}
			case "esc":
				m.noteInput.SetValue("")
				if m.selectedTask != nil {
					m.state = viewNotes
				} else {
					m.state = menu
					m.cursor = 0
				}
			}
		}
	}
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case menu:
		s := titleStyle.Render("Task Tracker") + "\n\n"
		for i, choice := range m.choices {
			style := normalStyle
			if m.cursor == i {
				style = selectedStyle
			}
			s += style.Render(choice) + "\n"
		}
		return s + "\n" + helpStyle.Render("[↑/↓] Navigate • [Enter] Select • [q] Quit")
	case add:
		return titleStyle.Render("Add Task Description:") + "\n" + m.textInput.View() + "\n" + helpStyle.Render("[Enter] Save • [ESC] Cancel")
	case addDue:
		return titleStyle.Render("Add Task Due Date:") + "\n" + m.dueInput.View() + "\n" + helpStyle.Render("[Enter] Save • [ESC] Cancel")
	case askNote:
		return titleStyle.Render("Would you like to add a note? (y/n)")
	case list:
		return m.renderTaskList() + "\n" + helpStyle.Render("[↑/↓] Navigate • [d] Mark Done • [i] Mark In-Progress • [x] Delete Task • [Enter] View Notes • [ESC] Back")
	case selectTaskNotes:
		return titleStyle.Render("Select Task for Notes:") + "\n" +
			m.renderTaskList() + "\n" +
			helpStyle.Render("[↑/↓] Navigate • [#] Type Task Number • [Enter] Open Notes • [ESC] Back")
	case addNote:
		var existing string
		if m.selectedTask != nil {
			// Always refresh selectedTask from store
			for i := range m.store.Tasks {
				if m.store.Tasks[i].ID == m.selectedTask.ID {
					m.selectedTask = &m.store.Tasks[i]
					break
				}
			}

			existing = titleStyle.Render("Existing Notes for "+m.selectedTask.Description+":") + "\n"
			if len(m.selectedTask.Notes) > 0 {
				for i, note := range m.selectedTask.Notes {
					line := fmt.Sprintf("%d. %s", i+1, note)
					existing += line + "\n"
				}
			}
			// else {
			// 	existing += "No existing notes.\n"
			// }
		}
		if m.state == askNote {
			return existing
		}
		return existing + "\n" + titleStyle.Render("New Note:") + "\n" + m.noteInput.View() + "\n" + helpStyle.Render("[Enter] Save • [ESC] Cancel")
	case viewNotes:
		return m.renderNotes()
	case editNote:
		return titleStyle.Render("Edit Note:") + "\n" + m.noteInput.View()
	default:
		return "Unknown state"
	}
}

func (m model) renderTaskList() string {
	if len(m.store.Tasks) == 0 {
		return "No tasks found."
	}

	// Determine maximum widths for each column
	maxIndexWidth := len(fmt.Sprintf("[%d]", len(m.store.Tasks)))
	maxDescWidth := len("DESCRIPTION")
	maxNotesWidth := len("NOTES")
	maxStatusWidth := len("STATUS")

	for _, t := range m.store.Tasks {
		if len(t.Description) > maxDescWidth {
			maxDescWidth = len(t.Description)
		}
		notesCount := fmt.Sprintf("%d", len(t.Notes))
		if len(notesCount) > maxNotesWidth {
			maxNotesWidth = len(notesCount)
		}
		statusText := t.Status.String()
		if len(statusText) > maxStatusWidth {
			maxStatusWidth = len(statusText)
		}
	}

	// Prepare header with proper alignment
	header := fmt.Sprintf("%-*s  %-*s  %-10s  %-*s  %-*s",
		maxIndexWidth, "#",
		maxDescWidth, "DESCRIPTION",
		"DUE DATE",
		maxStatusWidth, "STATUS",
		maxNotesWidth, "NOTES",
	)
	s := header + "\n" + strings.Repeat("-", len(header)) + "\n"

	// Prepare rows
	for i, t := range m.store.Tasks {
		statusText := "todo"
		if t.Status == StatusDone {
			statusText = "done"
		} else if t.Status == StatusInProgress {
			statusText = "in-progress"
		}
		line := fmt.Sprintf("%-*s  %-*s  %-10s  %-*s  %*d",
			maxIndexWidth, fmt.Sprintf("[%d]", i+1),
			maxDescWidth, t.Description,
			t.DueDate.Format("2006-01-02"),
			maxStatusWidth, statusText,
			maxNotesWidth, len(t.Notes),
		)
		if t.Status == StatusDone {
			line = doneStyle.Render(line)
		}
		if i == m.cursor {
			line = selectedStyle.Render(line)
		}
		s += line + "\n"
	}

	return s
}

func (m model) renderNotes() string {
	if m.selectedTask == nil {
		return "No task selected."
	}
	s := titleStyle.Render("Notes for "+m.selectedTask.Description) + "\n"
	if len(m.selectedTask.Notes) == 0 {
		s += "No notes yet. Press 'a' to add.\n"
	} else {
		for i, note := range m.selectedTask.Notes {
			line := fmt.Sprintf("%d. %s", i+1, note)
			if i == m.selectedNote {
				line = selectedStyle.Render(line)
			}
			s += line + "\n"
		}
	}
	return s + "\n" + helpStyle.Render("[↑/↓] Navigate • [a] Add Note • [e] Edit • [x] Delete • [ESC] Back")
}

func parseDueDate(input string) time.Time {
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" || input == "default" {
		return time.Now().Add(48 * time.Hour)
	}
	if input == "tomorrow" {
		return time.Now().Add(24 * time.Hour)
	}
	if input == "next week" {
		return time.Now().AddDate(0, 0, 7)
	}
	if t, err := time.Parse("2006-01-02", input); err == nil {
		return t
	}
	if t, err := now.Parse(input); err == nil {
		return t
	}
	return time.Now().Add(48 * time.Hour)
}
