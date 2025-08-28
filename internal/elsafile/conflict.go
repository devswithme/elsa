package elsafile

import (
	"reflect"
)

// HasConflict checks if a command name conflicts with built-in commands
func (em *Manager) HasConflict(name string) bool {
	// If we have root command, get built-in commands dynamically
	if em.rootCommand != nil {
		// Use reflection to call Commands() method on cobra.Command
		rootCmdValue := reflect.ValueOf(em.rootCommand)
		if rootCmdValue.Kind() == reflect.Ptr {
			rootCmdValue = rootCmdValue.Elem()
		}

		// Try to call Commands() method
		commandsMethod := rootCmdValue.MethodByName("Commands")
		if commandsMethod.IsValid() {
			commandsResult := commandsMethod.Call(nil)
			if len(commandsResult) > 0 {
				commands := commandsResult[0]
				if commands.Kind() == reflect.Slice {
					for i := 0; i < commands.Len(); i++ {
						cmd := commands.Index(i)
						if cmd.Kind() == reflect.Ptr {
							cmd = cmd.Elem()
						}

						// Try to get Name() and Hidden() methods
						nameMethod := cmd.MethodByName("Name")
						hiddenMethod := cmd.MethodByName("Hidden")

						if nameMethod.IsValid() && hiddenMethod.IsValid() {
							nameResult := nameMethod.Call(nil)
							hiddenResult := hiddenMethod.Call(nil)

							if len(nameResult) > 0 && len(hiddenResult) > 0 {
								cmdName := nameResult[0].String()
								isHidden := hiddenResult[0].Bool()

								if cmdName == name && !isHidden {
									return true
								}
							}
						}
					}
				}
			}
		}
		return false
	}

	// Fallback to static list if no root command available
	builtinCommands := []string{
		"init", "run", "list", "exec", "migrate", "watch", "help", "version",
	}

	for _, builtin := range builtinCommands {
		if name == builtin {
			return true
		}
	}

	return false
}
