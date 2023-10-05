package argparse

import (
	"fmt"
	"os"
	"strings"
)

type argparse struct {
	args        []Argument
	argMap      map[string]interface{}
	name        string
	description string
	details     string
}

// Create a new argparse object and return its pointer
func NewParser(name, description, details string) *argparse {
	var args []Argument
	parser := argparse{name: name, description: description, details: details, args: args}
	parser.AddArg("-h", "--help", "show this help message and exit", "string", "", false)
	return &parser
}

// AddArg new option, ex AddArg('-c', '--count', 'number of x', ArgType.STRING)
func (a *argparse) AddArg(alias, flag, description, argType, defaultValue string, required bool) {
	// Check for empty values
	if ok := newArgParamCheck(alias, flag, argType); !ok {
		os.Exit(0)
	}
	argTypeObj := getArgTypeFromString(argType)

	arg := NewArgument(argTypeObj, flag, alias, description, defaultValue, required)
	a.args = append(a.args, arg)
}

// newArgParamCheck examine each new argument, ensure they have value
func newArgParamCheck(alias, flag, argType string) bool {
	if strings.TrimSpace(alias) == "" {
		return false
	} else if strings.TrimSpace(flag) == "" {
		return false
	} else if strings.TrimSpace(argType) == "" {
		return false
	}
	return true
}

func (a *argparse) getRequiredArgs() []Argument {
	var required []Argument

	for i, _ := range a.args {
		arg := a.args[i]

		if arg.isRequired() {
			required = append(required, arg)
		}
	}
	return required
}

func (a *argparse) getOptionalArgs() []Argument {
	var optional []Argument

	for i, _ := range a.args {
		arg := a.args[i]

		if !arg.isRequired() {
			optional = append(optional, arg)
		}
	}
	return optional
}

// helpRequested iterates over user arguments to determine whether -h or --help has been passed
func (a *argparse) helpRequested() bool {
	for _, arg := range os.Args[1:] {
		argLower := strings.ToLower(strings.Trim(arg, "- \t"))
		if strings.EqualFold(argLower, "h") || strings.EqualFold(argLower, "help") {
			return true
		}
	}
	return false
}

func (a *argparse) getProgramName() string {
	programPath := os.Args[0]
	start := strings.LastIndex(programPath, "/") + 1
	programName := programPath[start:]
	return programName
}

// Print the banner for the user
func (a *argparse) printBanner() {

	// Set help parameterized help banner
	banner := `
usage: $name [-h --help] $required $optional

$description

options:
$options

$details`
	// Grab the program name for usage line
	banner = strings.Replace(banner, "$name", a.name, 1)

	// Replace the description
	banner = strings.Replace(banner, "$description", a.description, 1)

	// Replace the details
	banner = strings.Replace(banner, "$details", a.details, 1)

	// Gather optional parameters
	optional := a.getOptionalArgs()

	// Gather required parameters
	required := a.getRequiredArgs()

	// determine length of options line
	maxParam, maxDescription := getOptionsLength(optional, required)

	// Add optional parameters to banner
	optionalString := ""
	optionsString := ""
	for i, _ := range optional {
		arg := optional[i]
		optionalString += fmt.Sprintf(" [%s %s] ", arg.alias, arg.flag)
		left := padRight(fmt.Sprintf("   %s %s ", arg.alias, arg.flag), maxParam)
		right := padLeft(arg.description, maxDescription)
		optionsString += fmt.Sprintf("%s%s\n", left, right)
	}
	banner = strings.Replace(banner, "$optional", strings.TrimSpace(optionalString), 1)

	// Add required parameters to banner
	requiredString := ""
	for i, _ := range required {
		arg := required[i]
		requiredString += fmt.Sprintf(" {%s %s} ", arg.alias, arg.flag)
		left := padRight(fmt.Sprintf("   %s %s ", arg.alias, arg.flag), maxParam)
		right := padLeft(arg.description, maxDescription)
		optionsString += fmt.Sprintf("%s%s\n", left, right)
	}
	banner = strings.Replace(banner, "$required", strings.TrimSpace(requiredString), 1)
	banner = strings.Replace(banner, "$options", optionsString, 1)

	fmt.Println(banner)
	os.Exit(1)
}

// Parse parses arguments, return a mapping of each
func (a *argparse) Parse() {
	// Check if -h or --help in parameters, print help banner
	if help := a.helpRequested(); help {
		a.printBanner()
		os.Exit(1)
	}

	// Get a list of all required parameters
	required := a.getRequiredArgs()

	// Arguments passed by the user, excluding path to executable
	userArgs := os.Args[1:]

	// If the number of arguments passed does not equal the total number of required fields
	if len(required) > len(userArgs) {
		err := "[!] "
		for _, arg := range required {
			err += fmt.Sprintf("\"%s\", ", strings.Trim(arg.flag, "-"))
		}
		err += " is a required field"
		fmt.Println(fmt.Errorf(err))
		os.Exit(0)
	}

	// Map arguments to parameters
	a.argMap = a.mapArgs()
}

// Iterate over the user input, assign parameter names to the real value of the argument passed by the user
func (a *argparse) mapArgs() map[string]interface{} {
	var arg *Argument

	argsMapped := make(map[string]interface{})

	// Get the arguments passed
	userArgs := os.Args[1:]

	// Iterate over each index of the user arguments, assign a key -> value pair to the map for each
	for index, value := range userArgs {
		if string(value[0]) == "-" {
			arg = a.getArg(value)
			if arg == nil {
				continue
			}
			arg.value = userArgs[index+1]
			argsMapped[arg.name] = arg.getRealValue()
		}
	}

	return argsMapped
}

func (a *argparse) getArg(param string) *Argument {
	for i, _ := range a.args {
		arg := a.args[i]

		if strings.EqualFold(arg.alias, param) || strings.EqualFold(arg.flag, param) {
			return &arg
		}
	}
	return nil
}

// Helper Functions

// getOptionsLength get the maximum length of a parameter and description, used for formatting help output
func getOptionsLength(optional, required []Argument) (int, int) {
	maxDescriptionLength := 0
	maxParamsLength := 0

	for i, _ := range optional {
		arg := optional[i]
		argLength := len(arg.description)
		if argLength > maxDescriptionLength {
			maxDescriptionLength = argLength
		}
		argLength = len(arg.flag) + len(arg.alias)
		if argLength > maxParamsLength {
			maxParamsLength = argLength
		}
	}

	for i, _ := range required {
		arg := required[i]
		argLength := len(arg.description)
		if argLength > maxDescriptionLength {
			maxDescriptionLength = argLength
		}
		argLength = len(arg.flag) + len(arg.alias)
		if argLength > maxParamsLength {
			maxParamsLength = argLength
		}
	}

	return maxParamsLength + 5, maxDescriptionLength + 10
}

// padLeft add space to the left of a value up to a maximum length
func padLeft(val string, length int) string {
	out := ""
	delta := length - len(val)
	for i := 0; i < delta; i++ {
		out += " "
	}
	out += val
	return out
}

// padRight add space to the right of a value up to a maximum length
func padRight(val string, length int) string {
	out := val
	delta := length - len(val)
	for i := 0; i < delta; i++ {
		out += " "
	}
	return out
}
