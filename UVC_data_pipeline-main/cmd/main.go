package main

import "os"

func main() {
	app := new(application)
	if err := app.setupApplication(); err != nil {
		app.errorLog.Fatalf("setup of application failed: %v", err)
		// In case the definition of errorLog failed, I will add an os.Exit(),
		// that will not allow the CLI to run if there was an error.
		os.Exit(1)
	}

	app.runTUI()
}
