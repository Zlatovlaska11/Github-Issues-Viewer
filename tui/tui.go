package tui

import (
	"fmt"
	"os"
	"zlatolas/projectManager/dataSchemes"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	page         int
	repoInput    textinput.Model
	titleInput   textinput.Model
	descInput    textinput.Model
	labelInput   textinput.Model
	issueTitle   string
	issueDesc    string
	issueLabels  string
	currentInput int

	table table.Model
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color("240"))

var username string;

const (
	page1 = iota
	page2
	page3
)

var Reset = "\033[0m"
var Red = "\033[31m"

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "up":
			if m.currentInput > 0 {
				m.currentInput--
				m.updateFocus() // Update focus when navigating
			}
		case "down":
			if m.currentInput < 3 { // 0-based index for 3 inputs
				m.currentInput++
				m.updateFocus() // Update focus when navigating
			}
		case "tab":
			// Switch between pages
			if m.page == page1 {
				m.page = page2
			} else {
				m.page = page1
			}
		case "q":
			if m.page == page3 {
				m.page = page1
			}
		case "enter":
			if m.page != 0 {
				m.page = page3
			} else {
				dataSchemes.CreateIssue(m.titleInput.Value(), m.descInput.Value(), m.labelInput.Value(), m.repoInput.Value(), username)
				m.titleInput.Reset()
				m.descInput.Reset()
				m.labelInput.Reset()
				m.repoInput.Reset()
			}
			// TODO add some sending logic
		}
	}

	// Update the input fields if on page 1
	var cmd tea.Cmd
	if m.page == page1 {
		switch m.currentInput {
		case 0:
			m.repoInput, cmd = m.repoInput.Update(msg)
		case 1:
			m.titleInput, cmd = m.titleInput.Update(msg)
		case 2:
			m.descInput, cmd = m.descInput.Update(msg)
		case 3:
			m.labelInput, cmd = m.labelInput.Update(msg)
		}
		return m, cmd
	}
	m.table, cmd = m.table.Update(msg)

	return m, nil
}

// Update the focus based on current input field
func (m *model) updateFocus() {
	switch m.currentInput {
	case 0:
    m.titleInput.Blur()
    m.descInput.Blur()
    m.labelInput.Blur()
    m.repoInput.Focus()
	case 1:
    m.titleInput.Focus()
    m.repoInput.Blur()
    m.descInput.Blur()
    m.labelInput.Blur()
	case 2:
    m.titleInput.Blur()
    m.repoInput.Blur()
    m.descInput.Focus()
    m.labelInput.Blur()
  case 3:
    m.titleInput.Blur()
    m.repoInput.Blur()
    m.descInput.Blur()
    m.labelInput.Focus()
	}
}

func (m model) View() string {
	var view string

	// Tab-like view at the top
	view += "== Tabs ==\n"
	if m.page == page1 {
		view += "-> Page 1\nPage 2\n"
	} else {
		view += "Page 1\n-> Page 2\n"
	}

	// Main view content based on current page
	switch m.page {
	case page1:
		// Page 1 with the GitHub Issue form
		view += "\nPage 1: Create a GitHub Issue\n\n"
		view += fmt.Sprintf("Title: %s\n", m.titleInput.View())
		view += fmt.Sprintf("Description: %s\n", m.descInput.View())
		view += fmt.Sprintf("Labels: %s\n", m.labelInput.View())
		view += fmt.Sprintf("[Submit]")
		view += "\nPress 'Tab' to go to the summary page Press 'esc' to quit."
	case page2:
		view += m.table.View() + "\nPress q to quit."
	case page3:
		view += fmt.Sprintf(Red+"Issue number:"+Reset+" %s\n", m.table.SelectedRow()[0])
		view += fmt.Sprintf(Red+"Description:"+Reset+" %s\n", m.table.SelectedRow()[1])
		view += fmt.Sprintf(Red+"Status:"+Reset+" %s\n", m.table.SelectedRow()[2])
		view += fmt.Sprintf(Red+"Assignee:"+Reset+" %s\n", m.table.SelectedRow()[3])
	}

	return view
}

func InitTui(repo string, user string) {
	// Initialize form inputs for GitHub issue creation
	repoInput := textinput.New()
	repoInput.Placeholder = repo
	repoInput.Focus()

	titleInput := textinput.New()
	titleInput.Placeholder = "Issue Title"

	descInput := textinput.New()
	descInput.Placeholder = "Issue Description"

	labelInput := textinput.New()
	labelInput.Placeholder = "Issue Labels (comma separated)"

	columns := []table.Column{
		{Title: "Number", Width: 4},
		{Title: "Title", Width: 10},
		{Title: "State", Width: 10},
		{Title: "Asignee", Width: 10},
	}

	issues := dataSchemes.GetIssues(repo, user)

  username = user

	rows := dataSchemes.ParseIssues(issues)

	rws := []table.Row{}
	for _, issue := range rows.Issues {
		row := table.Row{
			fmt.Sprintf("#%d", issue.Number),
			issue.Title,
			issue.State,
			issue.Asignee,
		}
		rws = append(rws, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rws),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	// Initialize the program with a model containing the form inputs
	p := tea.NewProgram(model{
		page:         page1,
		table:        t,
    repoInput:    repoInput,
		titleInput:   titleInput,
		descInput:    descInput,
		labelInput:   labelInput,
		currentInput: 0,
	})

	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
