package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// runTUI, run the TUI (Terminal User Interface), handled by the CLI package.
func (app *application) runTUI() {
	app.setupCLI()

	if err := app.tui.Run(os.Args); err != nil {
		app.errorLog.Fatal(err)
	}
}

// setupCLI, configure and initialize all commands, flags and options of
// the TUI.
func (app *application) setupCLI() {
	app.tui = &cli.App{
		Name:  "cbExtractor",
		Usage: "Extract data from Crunchbase.",
		// This options enables short flag abbreviations to be merged into a
		// single flag with a '-' prefix.
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "extract",
				Usage: "Extract data from Crunchbase.",
				// Flags for extract command.
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "no-proxy",
						Usage: "Do not use a proxy while establishing a connection with Crunchbase.",
					},
				},
				Action: func(cCtx *cli.Context) error {
					noProxyFlag := cCtx.Bool("no-proxy")
					// Perform the required setup and configuration, e.g.
					// configuring the *http.Client with or without a proxy,
					// or handling the authentication with the CB API.
					if err := app.setupExtractCommands(noProxyFlag); err != nil {
						err = fmt.Errorf("setup for 'extract' command failed: %w", err)
						app.errorLog.Printf("extracting data from CB API failed: %v", err)
						return cli.Exit(err, 1)

					}

					if err := app.extractCBData(); err != nil {
						err = fmt.Errorf("error while executing 'extractCBData' command: %w", err)
						app.errorLog.Printf("extracting data from CB API failed: %v", err)
						return cli.Exit(err, 1)
					}
					return nil
				},
			},
			&cli.Command{
				Name:  "db",
				Usage: "Perform operations in the database.",
				// Help will not be another subcommand apart from 'insert'. It
				// will still be a flag '-h'.
				// HideHelpCommand: true,
				Subcommands: []*cli.Command{
					&cli.Command{
						Name:  "insert",
						Usage: "Insert a `file` into a database.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "file",
								Aliases:  []string{"f"},
								Required: true,
								Usage:    "`PATH` to the file that will be inserted into the database.",
							},
							&cli.StringFlag{
								Name:     "remote",
								Aliases:  []string{"r"},
								Required: true,
								Usage:    "`IP` address of remote server hosting the MongoDB instance.",
							},
						},
						Action: func(cCtx *cli.Context) error {
							// If the "file" flag was not set properly, the
							// application should exit.
							if cCtx.String("file") == "" {
								err := fmt.Errorf("--file flag missing: no path for a file was given.")
								return cli.Exit(err, 1)
							}

							// Perform the required setup and configuration.
							// Pass the given IP for the remote db to connect to
							// to the setup method.
							if err := app.setupDBCommands(cCtx.String("remote")); err != nil {
								err = fmt.Errorf("setup for 'db' command failed: %w", err)
								app.errorLog.Print(err)
								return cli.Exit(err, 1)
							}

							if err := app.insertDB(cCtx.String("file")); err != nil {
								err = fmt.Errorf("error while executing 'insertDB' command, file %s could not be inserted into the database: %w", cCtx.String("file"), err)
								app.errorLog.Print(err)
								return cli.Exit(err, 1)
							}

							app.infoLog.Printf("The documents in the external file %s were correctly inserted into the db.", cCtx.String("file"))
							return nil
						},
					},
				},
			},
			&cli.Command{
				Name:  "linkedin",
				Usage: "LinkedIn operations.",
				// Help will not be another subcommand apart from 'insert'. It
				// will still be a flag '-h'.
				// HideHelpCommand: true,
				Subcommands: []*cli.Command{
					&cli.Command{
						Name:  "present",
						Usage: "Is company in collection.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "remote",
								Aliases:  []string{"r"},
								Required: true,
								Usage:    "`IP` address of remote server hosting the MongoDB instance.",
							},
							&cli.StringFlag{
								Name:     "uuid",
								Aliases:  []string{"u"},
								Required: true,
								Usage:    "UUID of the company.",
							},
						},
						Action: func(cCtx *cli.Context) error {
							// Perform the required setup and configuration.
							// Pass the given IP for the remote db to connect to
							// to the setup method.
							if err := app.setupDBCommands(cCtx.String("remote")); err != nil {
								err = fmt.Errorf("setup for 'db' command failed: %w", err)
								app.errorLog.Print(err)
								return cli.Exit(err, 1)
							}

							if err := app.companyIsInColl(cCtx.String("uuid")); err != nil {
								err = fmt.Errorf("error while executing 'present' command: %w", err)
								app.errorLog.Print(err)
								return cli.Exit(err, 1)
							}
							return nil
						},
					},
					&cli.Command{
						Name:  "list",
						Usage: "List companies after a certain date.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "remote",
								Aliases:  []string{"r"},
								Required: true,
								Usage:    "`IP` address of remote server hosting the MongoDB instance.",
							},
							&cli.StringFlag{
								Name:     "date",
								Aliases:  []string{"d"},
								Required: true,
								Usage:    "Date after which companies should be listed. Format: '2010-Feb-02'.",
							},
						},
						Action: func(cCtx *cli.Context) error {
							// Perform the required setup and configuration.
							// Pass the given IP for the remote db to connect to
							// to the setup method.
							if err := app.setupDBCommands(cCtx.String("remote")); err != nil {
								err = fmt.Errorf("setup for 'db' command failed: %w", err)
								app.errorLog.Print(err)
								return cli.Exit(err, 1)
							}

							if err := app.findCompaniesTimestamp(cCtx.String("date")); err != nil {
								err = fmt.Errorf("error while executing 'list' command: %w", err)
								app.errorLog.Print(err)
								return cli.Exit(err, 1)
							}
							return nil
						},
					},
					&cli.Command{
						Name:  "update",
						Usage: "Update the LinkedIn URLs of targets.",
						Subcommands: []*cli.Command{
							&cli.Command{
								Name:  "companies",
								Usage: "The companies are the targets of the update.",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "remote",
										Aliases:  []string{"r"},
										Required: true,
										Usage:    "`IP` address of remote server hosting the MongoDB instance.",
									},
									&cli.StringFlag{
										Name:     "date",
										Aliases:  []string{"d"},
										Required: true,
										Usage:    "Date after which companies should be listed. Format: '2010-Feb-02'.",
									},
								},

								Action: func(cCtx *cli.Context) error {
									// Perform the required setup and configuration.
									// Pass the given IP for the remote db.
									if err := app.setupDBCommands(cCtx.String("remote")); err != nil {
										err = fmt.Errorf("setup for 'db' command failed: %w", err)
										app.errorLog.Print(err)
										return cli.Exit(err, 1)
									}

									if err := app.updateUniqueCompanies(cCtx.String("date")); err != nil {
										err = fmt.Errorf("error while executing 'update companies' command: %w", err)
										app.errorLog.Print(err)
										return cli.Exit(err, 1)
									}
									return nil
								},
							},
						},
					},
				},
			},
		},
	}
}
