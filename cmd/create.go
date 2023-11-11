package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ross96D/gohtmx-stack/cmd/download"
	"github.com/ross96D/gohtmx-stack/output"
	"github.com/ross96D/gohtmx-stack/program"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Start creating your structure for your web app",
	Run:   create,
}

type State struct {
	models      []tea.Model
	proyectName string
}

func create(cmd *cobra.Command, args []string) {
	var state State
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	out := output.Output{
		BasePath: wd,
	}
	state.models = []tea.Model{
		program.NewTextInput("gohtmx/test"),
		// program.Menu{
		// 	Options: []program.MenuItem{
		// 		{
		// 			Text: "Use default htmx",
		// 			Action: func() error {
		// 				println("Copying default htmx")
		// 				out.DefaultHtmx = true
		// 				return nil
		// 			},
		// 		},
		// 		{
		// 			Text: "Download lastest version of htmx",
		// 			Action: func() error {
		// 				println("Downloading lastest version htmx")
		// 				out.DefaultHtmx = false
		// 				return nil
		// 			},
		// 		},
		// 	},
		// },
		program.Menu{
			Options: []program.MenuItem{
				{
					Text: "Use tailwind",
					Action: func() error {
						println("Adding to package-json tailwind")
						return nil
					},
				},
				{
					Text: "Dont use tailwind",
					Action: func() error {
						println("Skipping tailwind")
						return nil
					},
				},
			},
		},
	}
	for i := 0; i < len(state.models); i++ {
		p := tea.NewProgram(state.models[i])
		m, err := p.Run()
		if err != nil {
			panic(err)
		}
		state.models[i] = m
	}
	println("Building proyect")
	var htmx []byte
	if htmx, err = download.LatestHtmx(); err != nil {
		fmt.Printf("An error occur while downloading the latest version of htmx: %s\n", err.Error())
		os.Exit(1)
	}
	out.Htmx = htmx
	for i := 0; i < len(state.models); i++ {
		m := state.models[i]
		switch v := m.(type) {
		case program.Menu:
			err := v.Options[v.SelectedIndex].Action()
			if err != nil {
				fmt.Printf("An error happens on the action %s: %s", v.Options[v.SelectedIndex].Text, err.Error())
			}
		case program.TextInput:
			out.ProyectName = v.TextInput.Value()
			fmt.Println("Name of the proyect is", state.proyectName)
		}
	}
	err = out.Build()
	if err != nil {
		fmt.Printf("Error building: %v", err)
	}
}
