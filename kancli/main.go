package main

import (
	"fmt"
	"os"

	bubbles_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 3

const (
	todo status = iota
	inProgress
	done
)

/* STYLING */
var (
	containerStyle = lipgloss.NewStyle().Padding(2, 3)
	columnStyle    = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle   = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

var StatusName = map[status]string{
	todo:       "Todo",
	inProgress: "In Progress",
	done:       "Done",
}

/* CUSTOM ITEM */
type Task struct {
	status      status
	title       string
	description string
}

// Implement the list.item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

// main model
type Model struct {
	focused  status
	lists    []bubbles_list.Model
	loaded   bool
	err      error
	quitting bool
}

func (m *Model) FocusNext() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

func (m *Model) FocusPrev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) initLists(width int, height int) {
	defaultList := bubbles_list.New([]bubbles_list.Item{}, bubbles_list.NewDefaultDelegate(), width, height)

	defaultList.SetShowHelp(false)

	m.focused = todo
	m.lists = []bubbles_list.Model{defaultList, defaultList, defaultList}

	// Set default title
	m.lists[todo].Title = "To Do"
	m.lists[inProgress].Title = "In Progress"
	m.lists[done].Title = "Done"

	// Set default items
	m.lists[todo].SetItems([]bubbles_list.Item{
		Task{status: todo, title: "Buy Milk from the store!", description: "Tesco is running a sale..."},
		Task{status: todo, title: "Look at clouds", description: "They are so fluffy!"},
		Task{status: todo, title: "Call Mom", description: "It's her birthday!"},
		Task{status: todo, title: "Go for a run", description: "It's been a while..."},
		Task{status: todo, title: "Read a book", description: "I have a few in mind..."},
	})

	m.lists[inProgress].SetItems([]bubbles_list.Item{
		Task{status: inProgress, title: "Finish the report", description: "Due on Friday..."},
		Task{status: inProgress, title: "Write a blog post", description: "About the new feature..."},
		Task{status: inProgress, title: "Work on the presentation", description: "For the meeting..."},
	})

	m.lists[done].SetItems([]bubbles_list.Item{
		Task{status: done, title: "Buy Milk from the store!", description: "Tesco is running a sale..."},
		Task{status: done, title: "Look at clouds", description: "They are so fluffy!"},
		Task{status: done, title: "Call Mom", description: "It's her birthday!"},
		Task{status: done, title: "Go for a run", description: "It's been a while..."},
	})

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			m.FocusPrev()
		case "right", "l":
			m.FocusNext()
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	// based on focus update the title colour
	m.lists[m.focused].Styles.Title.Background(lipgloss.Color("60"))
	m.lists[m.focused].Styles.Title.Foreground(lipgloss.Color("0"))

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		todoView := m.lists[todo].View()
		inProgressView := m.lists[todo].View()
		doneView := m.lists[done].View()
		switch m.focused {
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				inProgressView,
				doneView,
			)
		case inProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				todoView,
				focusedStyle.Render(inProgressView),
				doneView,
			)
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				todoView,
				inProgressView,
				focusedStyle.Render(doneView),
			)
		}
	}

	return "Loading..."
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
