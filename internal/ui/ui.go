package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gitswitch/internal/config"
)

var (
	subtle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	success   = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
	warning   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	muted     = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	label     = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	cursor     = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	selectedBg = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("255"))

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true).
			MarginBottom(1)

	tagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			MarginRight(1)
)

type PickResult struct {
	Profile  config.Profile
	Canceled bool
}

type pickerModel struct {
	profiles []config.Profile
	cursor   int
	chosen   *config.Profile
	canceled bool
	repoName string
}

func newPicker(profiles []config.Profile, repoName string) pickerModel {
	return pickerModel{profiles: profiles, repoName: repoName}
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.canceled = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
		case "enter", " ":
			p := m.profiles[m.cursor]
			m.chosen = &p
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m pickerModel) View() string {
	var b strings.Builder

	header := "  gitswitch"
	if m.repoName != "" {
		header += subtle.Render("  /" + m.repoName)
	}
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(muted.Render("  pick a profile  ·  ↑↓ navigate  ·  enter select  ·  q quit"))
	b.WriteString("\n\n")

	for i, p := range m.profiles {
		isSelected := i == m.cursor

		var line strings.Builder

		if isSelected {
			line.WriteString(cursor.Render("▸ "))
		} else {
			line.WriteString("  ")
		}

		nameStr := p.Name
		if isSelected {
			nameStr = highlight.Render(nameStr)
		}
		line.WriteString(nameStr)
		line.WriteString("  ")
		line.WriteString(muted.Render(p.Email))

		if p.GitHubUser != "" {
			line.WriteString("  ")
			line.WriteString(tagStyle.Render("@" + p.GitHubUser))
		}
		if p.SSHKey != "" {
			key := p.SSHKey
			if idx := strings.LastIndex(key, "/"); idx >= 0 {
				key = key[idx+1:]
			}
			line.WriteString(tagStyle.Render("ssh:" + key))
		}

		row := "  " + line.String()
		if isSelected {
			row = selectedBg.Render(row)
		}
		b.WriteString(row)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	return b.String()
}

func RunPicker(profiles []config.Profile, repoName string) PickResult {
	m := newPicker(profiles, repoName)
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return PickResult{Canceled: true}
	}
	fm := final.(pickerModel)
	if fm.canceled || fm.chosen == nil {
		return PickResult{Canceled: true}
	}
	return PickResult{Profile: *fm.chosen}
}

type formField struct {
	label    string
	value    string
	required bool
	hint     string
}

type formModel struct {
	fields  []formField
	cursor  int
	editing bool
	done    bool
	abort   bool
	input   string
}

func newForm() formModel {
	return formModel{
		fields: []formField{
			{label: "Name", required: true, hint: "e.g. John Doe"},
			{label: "Email", required: true, hint: "e.g. john@company.com"},
			{label: "SSH Key", required: true, hint: "e.g. ~/.ssh/id_ed25519_work"},
			{label: "GitHub Username", required: false, hint: "e.g. johndoe  (optional)"},
			{label: "GitHub URL", required: false, hint: "e.g. https://github.com/johndoe  (optional)"},
			{label: "Notes", required: false, hint: "e.g. Work laptop profile  (optional)"},
		},
		editing: true,
	}
}

func (m formModel) Init() tea.Cmd { return nil }

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.abort = true
			return m, tea.Quit

		case "esc":
			m.input = ""

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		case "enter":
			f := m.fields[m.cursor]
			val := strings.TrimSpace(m.input)
			if val == "" && f.required {
				return m, nil
			}
			m.fields[m.cursor].value = val
			m.input = ""
			if m.cursor < len(m.fields)-1 {
				m.cursor++
			} else {
				m.done = true
				return m, tea.Quit
			}

		default:
			if msg.Type == tea.KeyRunes {
				m.input += string(msg.Runes)
			}
		}
	}
	return m, nil
}

func (m formModel) View() string {
	var b strings.Builder

	b.WriteString(headerStyle.Render("  gitswitch") + "  " + muted.Render("new profile") + "\n\n")

	for i, f := range m.fields {
		isCurrent := i == m.cursor
		isFilledIn := f.value != ""

		var line strings.Builder
		if isCurrent {
			line.WriteString(cursor.Render("▸ "))
		} else {
			line.WriteString("  ")
		}

		lbl := label.Render(fmt.Sprintf("%-20s", f.label))
		line.WriteString(lbl)

		if isCurrent {
			display := m.input
			if display == "" {
				display = muted.Render(f.hint)
			} else {
				display = highlight.Render(display) + subtle.Render("█")
			}
			line.WriteString(display)
		} else if isFilledIn {
			line.WriteString(success.Render("✓ ") + f.value)
		} else {
			line.WriteString(muted.Render("—"))
		}

		b.WriteString("  " + line.String() + "\n")
	}

	b.WriteString("\n")
	b.WriteString(muted.Render("  enter confirm field  ·  esc clear  ·  ctrl+c cancel"))
	b.WriteString("\n")
	return b.String()
}

type FormResult struct {
	Profile config.Profile
	Aborted bool
}

func RunForm() FormResult {
	m := newForm()
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return FormResult{Aborted: true}
	}
	fm := final.(formModel)
	if fm.abort || !fm.done {
		return FormResult{Aborted: true}
	}

	get := func(i int) string { return fm.fields[i].value }

	return FormResult{
		Profile: config.Profile{
			Name:       get(0),
			Email:      get(1),
			SSHKey:     get(2),
			GitHubUser: get(3),
			GitHubURL:  get(4),
			Notes:      get(5),
		},
	}
}

func PrintProfiles(profiles []config.Profile) {
	if len(profiles) == 0 {
		fmt.Println(muted.Render("  no profiles found — run `gitswitch new` to create one"))
		return
	}

	fmt.Println(headerStyle.Render("  gitswitch") + muted.Render("  profiles"))
	fmt.Println()

	for i, p := range profiles {
		idx := subtle.Render(fmt.Sprintf("  %2d  ", i+1))
		name := highlight.Render(p.Name)
		email := muted.Render(p.Email)

		fmt.Printf("%s%s  %s\n", idx, name, email)

		if p.GitHubUser != "" {
			fmt.Printf("       %s %s\n", label.Render("github"), p.GitHubUser)
		}
		if p.GitHubURL != "" {
			fmt.Printf("       %s %s\n", label.Render("url   "), muted.Render(p.GitHubURL))
		}
		if p.SSHKey != "" {
			fmt.Printf("       %s %s\n", label.Render("ssh   "), p.SSHKey)
		}
		if p.Notes != "" {
			fmt.Printf("       %s %s\n", label.Render("notes "), subtle.Render(p.Notes))
		}
		fmt.Println()
	}

	fmt.Printf("  %s\n\n", warning.Render(fmt.Sprintf("%d profile(s)", len(profiles))))
}

func PrintApplied(p config.Profile, repoName string) {
	fmt.Printf("\n  %s applied %s\n",
		success.Render("✓"),
		highlight.Render(p.Name),
	)
	if repoName != "" {
		fmt.Printf("  %s\n", muted.Render("→ "+repoName))
	}
	fmt.Printf("  %s %s\n", label.Render("name "), p.Name)
	fmt.Printf("  %s %s\n", label.Render("email"), p.Email)
	if p.SSHKey != "" {
		fmt.Printf("  %s %s\n", label.Render("ssh  "), p.SSHKey)
	}
	fmt.Println()
}
