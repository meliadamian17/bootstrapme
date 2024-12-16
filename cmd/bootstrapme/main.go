package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/meliadamian17/bootstrapme/internal/config"
	"github.com/meliadamian17/bootstrapme/internal/tui"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	presetsByLang, err := config.LoadAllPresets()
	if err != nil {
		log.Fatalf("Failed to load presets: %v", err)
	}

	if len(presetsByLang) == 0 {
		log.Println("No presets found. Please add YAML configs in ~/.config/bootstrapme.")
		os.Exit(1)
	}

	initialModel := tui.NewModel(presetsByLang)
	p := tea.NewProgram(initialModel)
	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}

	selectedPreset := finalModel.(tui.Model).SelectedPreset
	if selectedPreset.Name == "" {
		log.Println("No preset selected, exiting.")
		return
	}

}
