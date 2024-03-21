package main

import (
    "fmt"
    "os"
    tea "github.com/charmbracelet/bubbletea"
)

type model struct{
	choices []string
	cursor int
	selected map[int]struct{}
}

func initialModel() model{
	return model{
		choices: []string{"Buy carrots", "Buy Celery", "Buy Goats"},
		cursor: 0,
		selected: make(map[int]struct{}),
	}
}

// now is the init, we're not returning any cmd
func (m model) Init() tea.Cmd{
	return nil
}

// update function, as part of the elm architecture
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	switch msg := msg.(type){
		case tea.KeyMsg:
			switch msg.String() {
				case "q":
					return m, tea.Quit

				case "up", "k":
					if m.cursor > 0{
						m.cursor--
					}

				case "down", "j":
					if m.cursor < len(m.choices) - 1 {
						m.cursor++
					}

				case "enter", " ":
					_, ok := m.selected[m.cursor]
					if ok{
						delete(m.selected, m.cursor)
					} else {
						m.selected[m.cursor] = struct{}{}
					}
			}
	}

	return m, nil

}

func (m model) View() string {
	s := "What are you doing today? \n\n"
	for i, choice := range m.choices {
		// is the cursor pointing a this choice?
		cursor := " "
		if m.cursor == i{
			cursor = ">"
		}

		// is this choice selected!
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main(){
	p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
