package multissh

import (
	"fmt"
	"os"
	"os/user"

	"github.com/akamensky/argparse"
)

// CommandParameters : Valid parameters to be passed to the main function
type CommandParameters struct {
	Nodes          []string
	Command        string
	Username       string
	Password       string
	PrivateKeyPath string
	// AddSudo        bool
	// BlockOutput    bool
	// MatchWidth     bool
	Debug bool
}

// ParseArgs : Primary parser for arguments
func Cli() CommandParameters {

	// Create the parser
	parser := argparse.NewParser("multissh", "Run a command on any number of nodes simultaneously")

	// Main criteria for running the application
	nodesArg := parser.String("n", "nodes", &argparse.Options{Required: true, Help: "Comma-delimited list of hostnames"})
	usernameArg := parser.String("u", "username", &argparse.Options{Required: false, Help: "Username to connect with"})
	passwordArg := parser.String("p", "password", &argparse.Options{Required: false, Help: "Password for the provided username"})
	privateKeyArg := parser.String("k", "private-key", &argparse.Options{Required: false, Help: "Path to your ssh private key"})
	commandArg := parser.String("c", "command", &argparse.Options{Required: true, Help: "Command to run"})
	debugFlag := parser.Flag("d", "debug", &argparse.Options{Required: false, Help: "Show debug output"})
    formatArg := parser.String("f", "format", &argparse.Options{Required: false, Help: "output formats: human, json", Default: "human"})

	// Optional modes
	//flag.BoolVar(&as, "add-sudo", false, "[Optional] Default: false - 'sudo' will be added to the beginning of every command")
	//flag.BoolVar(&bo, "block-output", false, "[Optional] Default: false - output will wait until all commands are finished, then display all output from a host at once")
	//flag.BoolVar(&mw, "match-width", false, "[Optional] Default: false - match the width of the hostnames so output is easier to compare")

	// parse the flags into their vars
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

    // Set up logging
	InitLogger(*debugFlag, *formatArg)

	// Process username
	username := getUsername(*usernameArg)

	// Process password
	password := getPassword(*passwordArg)

	// Process private key path
	privateKeyPath, err := getPrivateKeyPath(*privateKeyArg)
	if err != nil {
		logger("args").Error(err.Error())
		os.Exit(1)
	}

	// Process the nodes
	nodeList := getNodes(*nodesArg)
	if err != nil {
		// logger("args").Error(err.Error())
		// logger("args").Error("Invalid Syntax for node list")
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	return CommandParameters{
		Nodes:          nodeList,
		Command:        *commandArg,
		Username:       username,
		Password:       password,
		PrivateKeyPath: privateKeyPath,
		Debug:          *debugFlag,
	}

}

func getUsername(u string) string {
	// priority goes to the passed argument, gotten from cli args
	if u != "" {
		logger("args").Debug("username passed via args")
		return u
	}

	// next try MULTISSH_USERNAME env
	env_username := os.Getenv("MULTISSH_USERNAME")
	if env_username != "" {
		logger("args").Debug("using environment variable MULTISSH_USERNAME for username")
		return env_username
	}

	// fall back to the current user
	user, err := user.Current()
	if err != nil {
		logger("args").Error(err.Error())
		os.Exit(1)
	}
	logger("args").Debug(fmt.Sprintf("fallback to current user: %s", user.Username))
	return user.Username
}

func getPassword(p string) string {
	// priority goes to the passed argument, gotten from cli args
	if p != "" {
		logger("args").Debug("password passed via args")
		return p
	}

	// next try MULTISSH_PASSWORD env
	env_password := os.Getenv("MULTISSH_PASSWORD")
	if env_password != "" {
		logger("args").Debug("using environment variable MULTISSH_PASSWORD for password")
		return env_password
	}

	// return an empty string otherwise,
	// hopefully ssh key works
	logger("args").Debug("no password provided")
	return ""
}

func getNodes(n string) []string {
	// Validate and format the node list
	nl := getNodeList(n)
	return nl
}

func getPrivateKeyPath(k string) (string, error) {
	// priority goes to the passed argument, gotten from cli args
	if k == "" {
		logger("args").Debug("no private key provided")
		return "", nil
	}
	if _, err := os.Stat(k); err != nil {
		return "", fmt.Errorf("file not found: %s", k)
	}
	return k, nil
}
