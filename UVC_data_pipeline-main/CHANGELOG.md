## v0.9.2
* [minor] Expand sysadmin documentation.

## v0.9.1
* Add sysadmin documentation to the `/docs` folder, explaining in detail how to use this program.

## v0.9.0
* Add the half-done function to work with LinkedIn data structures in the data base.
* BUG FIX! Important: Increase the timeout of the functions that connect to the database.
	- The timeout was too small, it did not allow MongoDB time to perform the operations.
	- BUG: This BUG use to happen in all previous versions. 

## v0.8.0
* Change repository folder structure to a more orthodox Go repo structure, with the code starting in the root directory.
* Add dependabot GitHub Action.

## v0.7.1
* Fix GoReleaser YAML file. GoReleaser was not performing a build and release.

## v0.7.0 
* Use `urfave/cli` to handle the CLI.
* Enable the use of a proxy to extract data from CB. 
* Use GoReleaser to handle new releases.

## v0.6.0
* Clean the repo by migrating the _Cloud deployment_ scripts to a new repo. This improves and makes versioning and releases, easier.
* The main script in charge of pushing documents to the db remotely and getting CB data has not been modified.

## v0.5.0
* Be able to connect to a MongoDB remote instance (with username:password URI authentication, SSL has not been configured properly yet).
* One can insert multiple documents from a local computer to a remote MOngoDB instance which has been setup to permit remote connections.

## v0.4.0
* Add proper error handling in all methods and functions.
* Return context-aware error messages.

## v0.3.0
* Automatically handle the configuration of pf in the cloud deployment with a Go script.

## v0.2.1
* Add package `viddy` (watch command written in Go) to the default cloud deployment.

## v0.2.0
* Create a new read-only user in the MongoDB instance. This user will be used as a client with MongoDB Compass, so that there is no risk of modifing data in a production server while using Compass.
* Add the `-o` and `-h` flags to the script to setup the cloud VM. With the `-o` flag mongod opens a port to the internet, so that remote clients can access the DB, i.e. through MongoDB Compass.
* Add error handling functions in the cloud deployment script.

## v0.1.1
* Handle the sensitive information to get auth tokens from CB API through a .env file, remove the sensitive information from the repository.

## production0_v0.1.0
* Uses the `-fetchFile` flag to insert documents from a file into the db in the server.
* If no flag is used, it tries to extract data from the CB API, and stores it in a file created with a unique random name.
