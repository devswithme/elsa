package elsafile

import (
	"fmt"
	"os"

	"github.com/risoftinc/elsa/constants"
)

// TemplateGenerator handles creation of default Elsafile templates
type TemplateGenerator struct{}

// NewTemplateGenerator creates a new TemplateGenerator instance
func NewTemplateGenerator() *TemplateGenerator {
	return &TemplateGenerator{}
}

// CreateDefaultElsafile creates a default Elsafile in the current directory
func (tg *TemplateGenerator) CreateDefaultElsafile() error {
	// Check if Elsafile already exists
	if _, err := os.Stat(constants.DefaultElsafileName); err == nil {
		return fmt.Errorf(constants.ErrElsafileAlreadyExists)
	}

	// Create default Elsafile content
	content := tg.GetDefaultTemplate()

	// Write Elsafile
	if err := os.WriteFile(constants.DefaultElsafileName, []byte(content), constants.DefaultFilePermissions); err != nil {
		return fmt.Errorf(constants.ErrFailedToWriteFile, err)
	}

	return nil
}

// GetDefaultTemplate returns the default Elsafile template content
func (tg *TemplateGenerator) GetDefaultTemplate() string {
	return constants.DefaultTemplateHeader + "\n\n" +
		constants.DefaultBuildCommand + "\n\n" +
		constants.DefaultTestCommand + "\n\n" +
		constants.DefaultCleanCommand + "\n\n" +
		constants.DefaultDepsCommand + "\n\n" +
		constants.DefaultRunCommand + "\n\n" +
		constants.DefaultFmtCommand + "\n"
}

// CreateCustomElsafile creates a custom Elsafile with specified content
func (tg *TemplateGenerator) CreateCustomElsafile(content string) error {
	// Check if Elsafile already exists
	if _, err := os.Stat(constants.DefaultElsafileName); err == nil {
		return fmt.Errorf(constants.ErrElsafileAlreadyExists)
	}

	// Write custom Elsafile
	if err := os.WriteFile(constants.DefaultElsafileName, []byte(content), constants.DefaultFilePermissions); err != nil {
		return fmt.Errorf(constants.ErrFailedToWriteFile, err)
	}

	return nil
}

// GetSuccessMessage returns the success message after creating Elsafile
func (tg *TemplateGenerator) GetSuccessMessage() string {
	return fmt.Sprintf(`%s %s
%s You can now define custom commands in the Elsafile
%s Run commands using: elsa run command_name

Example:
  elsa run build    # Run the build command
  elsa run test     # Run the test command`,
		constants.SuccessEmoji, constants.MsgElsafileCreatedSuccess,
		constants.PencilEmoji,
		constants.RocketEmoji)
}
