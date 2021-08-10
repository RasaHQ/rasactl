package cmd

import (
	"bufio"
	"os"
	"strings"
	"syscall"

	"github.com/kyokomi/emoji"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// HandleSignals receives a signal from the channel and runs an action depends on the type of the signal.
func HandleSignals(sigs chan os.Signal) {
	signal := <-sigs
	runOnClose(signal)
}

func runOnClose(signal os.Signal) {
	emoji.Println("Bye :wave:")

	switch signal {
	case os.Interrupt:
		os.Exit(130)
	case syscall.SIGTERM:
		os.Exit(143)
	default:
		os.Exit(0)
	}
}

func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return errors.Errorf(
			"%q accepts no arguments\n\nUsage:  %s",
			cmd.CommandPath(),
			cmd.UseLine(),
		)
	}
	return nil
}

func examples(s string) string {
	trimmedText := strings.TrimSpace(s)
	if trimmedText == "" {
		return ""
	}

	const indent = `  `
	inLines := strings.Split(trimmedText, "\n")
	outLines := make([]string, 0, len(inLines))

	for _, line := range inLines {
		outLines = append(outLines, indent+strings.TrimSpace(line))
	}

	return strings.Join(outLines, "\n")
}

func getRasaXPasswordStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return line, err
	}
	return strings.TrimSuffix(line, "\n"), nil
}
