package elsafile

import (
	"fmt"
	"os"

	"github.com/risoftinc/elsa/internal/elsafile"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Elsafile in the current directory",
	Long: `Initialize creates a new Elsafile in the current directory.
The Elsafile is similar to a Makefile and allows you to define custom commands
for your project. Commands defined in Elsafile can be run using 'elsa run command'.

Example Elsafile content:
  build:
    go build -o bin/app .
  
  test:
    go test ./...
  
  clean:
    rm -rf bin/
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := createElsafile(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Elsafile: %v\n", err)
			os.Exit(1)
		}
	},
}

func createElsafile() error {
	templateGen := elsafile.NewTemplateGenerator()

	if err := templateGen.CreateDefaultElsafile(); err != nil {
		return err
	}

	fmt.Println(templateGen.GetSuccessMessage())
	return nil
}
