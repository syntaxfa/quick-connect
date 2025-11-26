package command

import (
	"os"

	"github.com/spf13/cobra"
)

type Server struct {
	cfg    ServiceConfig
	logger Logger
}

func (s Server) Command(cfg ServiceConfig, logger Logger, trap <-chan os.Signal) *cobra.Command {
	s.cfg = cfg
	s.logger = logger

	run := func(_ *cobra.Command, _ []string) {
		s.run(trap)
	}

	return &cobra.Command{
		Use:   "start",
		Short: "start quick connect in code-level monolith",
		Run:   run,
	}
}

func (s Server) run(_ <-chan os.Signal) {
}
