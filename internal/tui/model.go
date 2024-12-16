package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/meliadamian17/bootstrapme/internal/config"
	"github.com/meliadamian17/bootstrapme/internal/runner"
)

type step int

const (
	stepSelectLanguage step = iota
	stepSelectFramework
	stepEnterUsername
	stepEnterProjectName
	stepBootstrapping
	stepDone
)

type quitMsg struct{}

type Model struct {
	presetsByLang    map[string][]config.Preset
	languages        []table.Row
	frameworks       []table.Row
	selectedLanguage string
	SelectedPreset   config.Preset
	currentStep      step

	langTable      table.Model
	frameworkTable table.Model
	spinner        spinner.Model
	logs           string

	projectNameInput textinput.Model
	usernameInput    textinput.Model
	runner           *runner.Runner
	logChanClosed    bool
	doneTimerFired   bool
}

func NewModel(presetsByLang map[string][]config.Preset) Model {
	var langRows []table.Row
	for lang, ps := range presetsByLang {
		desc := fmt.Sprintf("%d presets available", len(ps))
		langRows = append(langRows, table.Row{lang, desc})
	}

	langCols := []table.Column{
		{Title: "Language", Width: 20},
		{Title: "Description", Width: 40},
	}
	langTable := table.New(
		table.WithColumns(langCols),
		table.WithRows(langRows),
		table.WithFocused(true),
	)

	langTable.SetStyles(TableStyle)

	s := spinner.New()
	s.Spinner = spinner.Dot

	projectTI := textinput.New()
	projectTI.Placeholder = "myproject"
	projectTI.Focus()

	usernameTI := textinput.New()
	usernameTI.Placeholder = "username"

	m := Model{
		presetsByLang:    presetsByLang,
		languages:        langRows,
		langTable:        langTable,
		spinner:          s,
		currentStep:      stepSelectLanguage,
		projectNameInput: projectTI,
		usernameInput:    usernameTI,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if wasQuitKey(msg) {
		return m, tea.Quit
	}

	if wasBackKey(msg) {
		switch m.currentStep {
		case stepSelectFramework:
			m.currentStep = stepSelectLanguage
			return m, nil
		case stepEnterUsername:
			m.currentStep = stepSelectFramework
			return m, nil
		case stepEnterProjectName:
			if m.selectedLanguage == "go" {
				m.currentStep = stepEnterUsername
			} else {
				m.currentStep = stepSelectFramework
			}
			return m, nil
		}
	}

	switch m.currentStep {
	case stepSelectLanguage:
		var cmd tea.Cmd
		m.langTable, cmd = m.langTable.Update(msg)
		if wasEnterPressed(msg) {
			idx := m.langTable.Cursor()
			if idx >= 0 && idx < len(m.languages) {
				m.selectedLanguage = m.languages[idx][0]
				m = m.loadFrameworks(m.selectedLanguage)
			}
		}
		return m, cmd

	case stepSelectFramework:
		var cmd tea.Cmd
		m.frameworkTable, cmd = m.frameworkTable.Update(msg)
		if wasEnterPressed(msg) {
			idx := m.frameworkTable.Cursor()
			if idx >= 0 && idx < len(m.frameworks) {
				selectedFramework := m.frameworks[idx][0]
				fps := m.presetsByLang[m.selectedLanguage]
				var chosen config.Preset
				for _, p := range fps {
					if p.Name == selectedFramework {
						chosen = p
						break
					}
				}
				m.SelectedPreset = chosen
				if m.selectedLanguage == "go" {
					m.currentStep = stepEnterUsername
					m.usernameInput.Focus()
				} else {
					m.currentStep = stepEnterProjectName
					m.projectNameInput.Focus()
				}
			}
		}
		return m, cmd

	case stepEnterUsername:
		var cmd tea.Cmd
		m.usernameInput, cmd = m.usernameInput.Update(msg)
		if wasEnterPressed(msg) {
			username := m.usernameInput.Value()
			if username == "" {
				username = "username"
			}
			if m.SelectedPreset.Variables == nil {
				m.SelectedPreset.Variables = make(map[string]string)
			}
			m.SelectedPreset.Variables["username"] = username

			m.currentStep = stepEnterProjectName
			m.projectNameInput.Focus()
		}
		return m, cmd

	case stepEnterProjectName:
		var cmd tea.Cmd
		m.projectNameInput, cmd = m.projectNameInput.Update(msg)
		if wasEnterPressed(msg) {
			projectName := m.projectNameInput.Value()
			if projectName == "" {
				projectName = "myproject"
			}
			m.currentStep = stepBootstrapping
			m.startBootstrapping(projectName)
			return m, m.tickCmd()
		}
		return m, cmd

	case stepBootstrapping:
		if _, ok := msg.(time.Time); ok {
			// On each tick, read lines
			for {
				if m.logChanClosed {
					m.currentStep = stepDone
					return m, m.autoCloseCmd()
				}
				select {
				case line, ok := <-m.runner.LogChan:
					if !ok {
						m.logChanClosed = true
					} else {
						// Color lines
						var styledLine string
						if strings.HasPrefix(line, "ERR:") || strings.Contains(line, "Error") {
							styledLine = LogErrorStyle.Render(line)
						} else {
							styledLine = LogInfoStyle.Render(line)
						}
						m.logs += styledLine + "\n"
					}
				default:
					// no more lines this tick
					return m, tea.Batch(m.tickCmd(), m.spinner.Tick)
				}
			}
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case stepDone:
		if _, ok := msg.(quitMsg); ok {
			return m, tea.Quit
		}
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	title := TitleStyle.Render("BootstrapMe") + "\n"

	switch m.currentStep {
	case stepSelectLanguage:
		return BorderStyle.Render(
			title + "\n" +
				TitleStyle.Render("Select a Language") + "\n\n" +
				m.langTable.View() + "\n(Use arrows, Enter to select, ctrl+c to quit)",
		)

	case stepSelectFramework:
		return BorderStyle.Render(
			title + "\n" +
				TitleStyle.Render(
					fmt.Sprintf(
						"Select a Framework/Preset for %s",
						m.selectedLanguage,
					),
				) + "\n\n" +
				m.frameworkTable.View() + "\n(Use arrows, Enter to select, ^ to go back, ctrl+c to quit)",
		)

	case stepEnterUsername:
		return BorderStyle.Render(
			title + "\n" +
				TitleStyle.Render("Enter Username:") + "\n\n" +
				m.usernameInput.View() + "\n\n(Enter to confirm, ^ to go back, ctrl+c to quit)",
		)

	case stepEnterProjectName:
		return BorderStyle.Render(
			title + "\n" +
				TitleStyle.Render("Enter Project Name:") + "\n\n" +
				m.projectNameInput.View() + "\n\n(Enter to confirm, ^ to go back, ctrl+c to quit)",
		)

	case stepBootstrapping:
		return BorderStyle.Render(
			title + m.spinner.View() + "\n" +
				m.logs +
				"\n(Press ctrl+c to quit early)",
		)

	case stepDone:
		return BorderStyle.Render(
			title + "\n" +
				DoneStyle.Render("Project created successfully!\n\nLogs:\n") +
				m.logs +
				"\n\nClosing in 2 seconds... (Press ctrl+c to quit now)",
		)
	}
	return ""
}

func (m *Model) startBootstrapping(projectName string) {
	if m.SelectedPreset.Variables == nil {
		m.SelectedPreset.Variables = make(map[string]string)
	}
	m.SelectedPreset.Variables["project_name"] = projectName

	presetInfo := runner.PresetInfo{
		Name:            m.SelectedPreset.Name,
		Files:           convertFiles(m.SelectedPreset.Files),
		PostInstallCmds: m.SelectedPreset.PostInstallCmds,
		Variables:       m.SelectedPreset.Variables,
	}
	r := &runner.Runner{
		ProjectName: projectName,
		Preset:      presetInfo,
		LogChan:     make(chan string, 1024),
	}
	m.runner = r
	r.RunPresetAsync()
}

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return t })
}

func (m Model) autoCloseCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
		return quitMsg{}
	})
}

func (m Model) loadFrameworks(language string) Model {
	fps := m.presetsByLang[language]
	var fwRows []table.Row
	for _, p := range fps {
		fwRows = append(fwRows, table.Row{p.Name, p.Description})
	}

	fwCols := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Description", Width: 40},
	}

	fwTable := table.New(
		table.WithColumns(fwCols),
		table.WithRows(fwRows),
		table.WithFocused(true),
	)

	fwTable.SetStyles(TableStyle)

	newM := m
	newM.frameworks = fwRows
	newM.frameworkTable = fwTable
	newM.currentStep = stepSelectFramework
	return newM
}

func wasEnterPressed(msg tea.Msg) bool {
	km, ok := msg.(tea.KeyMsg)
	return ok && km.Type == 13 // Enter
}

func wasQuitKey(msg tea.Msg) bool {
	km, ok := msg.(tea.KeyMsg)
	return ok && (km.String() == "ctrl+c")
}

func wasBackKey(msg tea.Msg) bool {
	km, ok := msg.(tea.KeyMsg)
	return ok && (km.String() == "^")
}

func convertFiles(files []config.FileSpec) []runner.FileInfo {
	res := make([]runner.FileInfo, len(files))
	for i, f := range files {
		res[i] = runner.FileInfo{Path: f.Path, Content: f.Content}
	}
	return res
}
