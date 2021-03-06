package description

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gojektech/proctor/daemon"
	"github.com/gojektech/proctor/io"
	"github.com/gojektech/proctor/proc"
	"github.com/spf13/cobra"
)

func NewCmd(printer io.Printer, proctorEngineClient daemon.Client) *cobra.Command {
	return &cobra.Command{
		Use:     "describe",
		Short:   "Describe a proc, list help for variables and constants",
		Long:    "In order to execute a proc, you need to provide certain variables. Describe command helps you with those variables and their meanings/convention/usage, etc.",
		Example: "proctor describe proc-one\nproctor describe proc-two",
		Aliases: []string{"help"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				printer.Println("Incorrect command. See `proctor describe --help` for usage", color.FgRed)
				return
			}

			procList, err := proctorEngineClient.ListProcs()
			if err != nil {
				printer.Println(err.Error(), color.FgRed)
				return
			}

			userProvidedProcName := args[0]
			desiredProc := proc.Metadata{}
			for _, proc := range procList {
				if userProvidedProcName == proc.Name {
					desiredProc = proc
				}
			}
			if len(desiredProc.Name) == 0 {
				printer.Println(fmt.Sprintf("Proctor doesn't support Proc `%s`\nRun `proctor list` to view supported Procs", userProvidedProcName), color.FgRed)
				return
			}

			printer.Println(fmt.Sprintf("%-40s %-100s", "Description", desiredProc.Description), color.Reset)
			printer.Println(fmt.Sprintf("%-40s %-100s", "Contributors", desiredProc.Contributors), color.Reset)
			printer.Println(fmt.Sprintf("%-40s %-100s", "Organization", desiredProc.Organization), color.Reset)

			printer.Println("\nArgs", color.FgMagenta)
			for _, arg := range desiredProc.EnvVars.Args {
				printer.Println(fmt.Sprintf("%-40s %-100s", arg.Name, arg.Description), color.Reset)
			}

			printer.Println(fmt.Sprintf("\nTo %s, run:\nproctor execute %s ARG_ONE=foo ARG_TWO=bar", userProvidedProcName, userProvidedProcName), color.FgGreen)
		},
	}
}
