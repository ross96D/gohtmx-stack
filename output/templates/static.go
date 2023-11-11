package templates

const Serve = `package cmd

import (
	"%s/handlers"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "Use for start the server",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	server := handlers.NewServer()

	fmt.Printf("Starting server on %%s\n", server.Addr)
	go printLink(server)
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Errorf("server could not start %%w", err))
	}
}

func printLink(s *http.Server) {
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("visit page on http://127.0.0.1%%s\n", s.Addr)
}
`
