package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jbeausoleil/task-tracker-bubbletea-blue-natural/internal/task"
	"github.com/jinzhu/now"
)

type state int
type addStep int

const (
	menu state = iota
	list
	add
	delete
	mark
	completedToday
	dueThisWeek
)

const (
	addDesc addStep = iota
	addDue
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))         // Bright white
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))         // Bright green
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))                     // Light gray
	doneStyle     = lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color("8")) // Dim gray
	helpStyle     = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("7"))         // Light gray
)

var (
	statusTodoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true) // Bright green
	statusInProgress = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true) // Bright blue
	statusDoneStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Strikethrough(true)
	descriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true) // Bright white
	dueDateStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))            // Yellow
	createdDateStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))            // Cyan
)

type model struct {
	cursor    int
	choices   []string
	state     state
	tasks     []task.Task
	textInput textinput.Model
	dueInput  textinput.Model
	addStep   addStep
	store     *task.Store
	tempDesc  string
}

func initialModel() model {
	descInput := textinput.New()
	descInput.Placeholder = "Type task description"
	descInput.Focus()
	descInput.CharLimit = 128
	descInput.Width = 40

	dueInput := textinput.New()
	dueInput.Placeholder = "YYYY-MM-DD or 'tomorrow'"
	dueInput.CharLimit = 30
	dueInput.Width = 35

	store := task.NewStore("tasks.json")
	tasks, _ := store.Load()

	return model{
		choices: []string{
			"Add Task",
			"List All Tasks",
			"List Completed Tasks Today",
			"List Tasks Due This Week",
			"Mark Task as Done",
			"Delete Task",
			"Quit",
		},
		state:     menu,
		tasks:     tasks,
		textInput: descInput,
		dueInput:  dueInput,
		addStep:   addDesc,
		store:     store,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// parseDueDate supports natural language using jinzhu/now
func parseDueDate(input string) time.Time {
	nowRef := time.Now()
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "", "default":
		return nowRef.Add(48 * time.Hour)
	case "today":
		return nowRef
	case "tomorrow":
		return nowRef.Add(24 * time.Hour)
	case "next week":
		return nowRef.AddDate(0, 0, 7)
	}

	// Fuzzy parsing with jinzhu/now
	if t, err := now.Parse(input); err == nil {
		return t
	}

	// Fallback to YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", input); err == nil {
		return t
	}

	// Default fallback
	return nowRef.Add(48 * time.Hour)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case menu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
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
					m.addStep = addDesc
				case 1:
					m.state = list
				case 2:
					m.state = completedToday
				case 3:
					m.state = dueThisWeek
				case 4:
					m.state = mark
					m.cursor = 0
				case 5:
					m.state = delete
					m.cursor = 0
				case 6:
					return m, tea.Quit
				}
			}
		}

	case add:
		if m.addStep == addDesc {
			m.textInput, cmd = m.textInput.Update(msg)
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case "esc":
					m.state = menu
				case "enter":
					if m.textInput.Value() != "" {
						m.tempDesc = m.textInput.Value()
						m.textInput.SetValue("")
						m.addStep = addDue
						m.dueInput.Focus()
					} else {
						m.state = menu
					}
				}
			}
		} else if m.addStep == addDue {
			m.dueInput, cmd = m.dueInput.Update(msg)
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case "esc":
					m.state = menu
				case "enter":
					dueDate := parseDueDate(m.dueInput.Value())
					newTask := task.Task{
						Description: m.tempDesc,
						Status:      "todo",
						CreatedAt:   time.Now(),
						DueDate:     dueDate,
					}
					m.tasks = append(m.tasks, newTask)
					m.store.Save(m.tasks)
					m.tempDesc = ""
					m.dueInput.SetValue("")
					m.state = menu
				}
			}
		}

	case list, completedToday, dueThisWeek:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.state = menu
			}
		}

	case mark:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.tasks)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.tasks) > 0 {
					m.tasks[m.cursor].Status = "done"
					m.tasks[m.cursor].CompletedAt = time.Now()
					m.store.Save(m.tasks)
				}
			case "esc":
				m.state = menu
			}
		}

	case delete:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.tasks)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.tasks) > 0 {
					m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
					if m.cursor >= len(m.tasks) && m.cursor > 0 {
						m.cursor--
					}
					m.store.Save(m.tasks)
				}
			case "esc":
				m.state = menu
			}
		}
	}
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case list:
		return m.renderTasks(m.tasks, "All Tasks") + helpStyle.Render("\n[↑/↓]: Scroll  [ESC]: Back  [q]: Quit")

	case completedToday:
		todayTasks := task.CompletedToday(m.tasks)
		return m.renderTasks(todayTasks, "Completed Tasks Today") + helpStyle.Render("\n[↑/↓]: Scroll  [ESC]: Back  [q]: Quit")

	case dueThisWeek:
		weekTasks := task.DueThisWeek(m.tasks)
		return m.renderTasks(weekTasks, "Tasks Due This Week") + helpStyle.Render("\n[↑/↓]: Scroll  [ESC]: Back  [q]: Quit")

	case add:
		if m.addStep == addDesc {
			return fmt.Sprintf(
				"%s\n\n%s\n\n(Press Enter to continue or ESC to cancel)\n",
				titleStyle.Render("Add a new task (Description):"),
				m.textInput.View(),
			)
		} else {
			return fmt.Sprintf(
				"%s\n\n%s\n\n(Enter YYYY-MM-DD or natural text: 'tomorrow', 'next week', '3 days from now')\n(Press Enter to save or ESC to cancel)\n",
				titleStyle.Render("Add a new task (Due Date):"),
				m.dueInput.View(),
			)
		}

	case mark:
		return m.renderSelection("Mark Task as Done") + helpStyle.Render("\n[↑/↓]: Navigate  [Enter]: Mark Done  [ESC]: Back")

	case delete:
		return m.renderSelection("Delete Task") + helpStyle.Render("\n[↑/↓]: Navigate  [Enter]: Delete  [ESC]: Back")

	default:
		s := titleStyle.Render("Task Tracker") + "\n\n"
		for i, choice := range m.choices {
			style := normalStyle
			if m.cursor == i {
				style = selectedStyle
			}
			s += style.Render(fmt.Sprintf("%s", choice)) + "\n"
		}
		s += "\n" + helpStyle.Render("[↑/↓]: Navigate  [Enter]: Select  [q]: Quit") + "\n"
		return s
	}
}

func (m model) renderTasks(tasks []task.Task, title string) string {
	if len(tasks) == 0 {
		return "No tasks found. Press ESC to return."
	}

	s := titleStyle.Render(title) + "\n\n"
	for i, t := range tasks {
		var status string
		switch t.Status {
		case "todo":
			status = statusTodoStyle.Render(t.Status)
		case "in-progress":
			status = statusInProgress.Render(t.Status)
		case "done":
			status = statusDoneStyle.Render(t.Status)
		}

		line := fmt.Sprintf("[%d] STATUS: %s | DESC: %s | DUE: %s | CREATED: %s",
			i+1,
			status,
			descriptionStyle.Render(t.Description),
			dueDateStyle.Render(t.DueDate.Format("2006-01-02")),
			createdDateStyle.Render(t.CreatedAt.Format("2006-01-02 15:04")),
		)

		s += line + "\n"
	}
	return s
}

func (m model) renderSelection(title string) string {
	if len(m.tasks) == 0 {
		return fmt.Sprintf("No tasks available. Press ESC to return.")
	}
	s := titleStyle.Render(title) + "\n\n"
	for i, t := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, t.Status, t.Description)
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
