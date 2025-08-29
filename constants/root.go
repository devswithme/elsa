package constants

// Root command constants
const (
	// Command usage
	RootUse = "elsa"

	// Command descriptions
	RootShort = "Elsa - Engineer's Little Smart Assistant"

	// Help text
	HelpText = `-h, --help      help for elsa
  -v, --version   version for elsa`

	// Usage text
	UsageText = `Usage:
  elsa [flags]
  elsa [command]

Available Commands:`

	// Suggestion text
	SuggestionText = "💡 Did you mean one of these commands?"

	// Version template
	VersionTemplate = "ELSA v%s (CLI)\ngo version %s %s\nLearn more at: https://risoftinc.com\n"

	// Banner template
	BannerTemplate = `Developer productivity toolkit for Go.
      ___           ___       ___           ___     
     /\  \         /\__\     /\  \         /\  \    
    /::\  \       /:/  /    /::\  \       /::\  \   
   /:/\:\  \     /:/  /    /:/\ \  \     /:/\:\  \  
  /::\~\:\  \   /:/  /    _\:\~\ \  \   /::\~\:\  \ 
 /:/\:\ \:\__\ /:/__/    /\ \:\ \ \__\ /:/\:\ \:\__\
 \:\~\:\ \/__/ \:\  \    \:\ \:\ \/__/ \/__\:\/:/  /
  \:\ \:\__\    \:\  \    \:\ \:\__\        \::/  / 
   \:\ \/__/     \:\  \    \:\/:/  /        /:/  /  
    \:\__\        \:\__\    \::/  /        /:/  /   
     \/__/         \/__/     \/__/         \/__/    V %s
(migration, scaffolding, project runner and task automation)`
)
