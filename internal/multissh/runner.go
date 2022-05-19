package multissh

// Run : controller for application workflow
func Run() {
	// Parse command line arguments and store them in commandParams
	commandParams := Cli()

	// Feed commandParams into the configurer
	// Receive an ClientConfig object to use in each connection
	config := newClientConfig(commandParams)

	// for each of the nodes, use the ClientConfig to connect and run the command
	Distribute(config, commandParams)
}
