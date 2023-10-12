package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

// setupExtractCommands, configures http.Client and CB cookies for 'extract'
// commands. If the 'noProxy' parameter is true, then no proxy connection is
// established to extract data from CB.
func (app *application) setupExtractCommands(noProxy bool) error {
	if err := app.configureClient(noProxy); err != nil {
		return fmt.Errorf("error while configuring the HTTP client: %w", err)
	}
	if err := app.handleAuthentication(); err != nil {
		return fmt.Errorf("authentication with CB API failed: %w", err)
	}
	return nil
}

// configureClient, configures the http.Client used by the application, so
// that if a flag was set or not, it proxies its traffic through a proxy.
func (app *application) configureClient(noProxy bool) error {
	// Create an HTTP Client to use in multiple API calls, and to handle the
	// cookie jar.
	app.client = new(http.Client)
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("unable to create a new cookiejar: %w", err)
	}
	// If no errors happened during the initialization of the cookie jar assign
	// its value to the app's cookiejar.
	app.client.Jar = jar

	// If noProxy is false (it was not set as a flag) then a proxy should be
	// configured.
	if !noProxy { // Configure a proxy.
		app.infoLog.Print("[PROXY] Running an extraction with a proxy. Configuring proxy connection.")
		transport, err := parseProxyURI()
		if err != nil {
			return fmt.Errorf("error while configuring proxy: %w", err)
		}
		app.client.Transport = transport
	}
	if noProxy { // Do not configure a proxy.
		app.infoLog.Print("[NO PROXY] Executing extraction WITHOUT a proxy.")
	}

	return nil
}

// handleAuthentication, handles getting session cookies by sending a login
// request to the Crunchbase API.
// If it returns an error, exit the application, fatal error.
func (app *application) handleAuthentication() error {
	// Parameter for app.login = false, so that login() does not store the
	// cookies on a persistent file.
	if err := app.login(false); err != nil {
		// If login fails, exit program, fatal error.
		return fmt.Errorf("unable to login (authenticate) into Crunchbase API: %w", err)
	}
	app.infoLog.Print("Received new cookies from Crunchbase API.")
	return nil
}

// handleAuthenticationPersistentCookies, handles getting session cookies by
// sending a login request to the Crunchbase API. If first tries to load an
// already present persistent file with cookies. It also stores new cookies,
// if it needs to ask for new cookies to the CB API.
// If it returns an error, exit the application, fatal error.
func (app *application) handleAuthenticationPersistentCookies() error {
	err := app.loadCookies()
	if err != nil {
		err = fmt.Errorf("unable to load cookies from file: %w", err)
		// Print an error here, just in case there is no cookies file.
		app.errorLog.Print(err)
		// If loadCookies fails, try to get new cookies through login().
		// loadCookies() can fail, if for example, there is no file in
		// the local repository that contains the current cookies.
		err = app.login(true)
		if err != nil {
			// If login fails, exit program, fatal error.
			return fmt.Errorf("unable to login (authenticate) into Crunchbase API: %w", err)
		}
		app.infoLog.Print("Received new cookies from Crunchbase API.")
	}
	return nil
}
