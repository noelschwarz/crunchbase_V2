# UVC Sysadmin Guide
The following document is a very detailed guide on how to manage the UVC cloud resources used for the UVC data pipeline.

Working as a sysadmin on the system using a Unix operating system (Mac OS X, Linux or FreeBSD) makes your life much easier, since lots of the tools needed to manage the cloud systems are natively installed on the system, or an installation of them is very uncomplicated.

## Table of contents
<!-- vim-markdown-toc GFM -->

* [Software requirements to follow along](#software-requirements-to-follow-along)
* [What is the UVC data pipeline?](#what-is-the-uvc-data-pipeline)
	- [What are the data sources that collect data (as of February 2023)?](#what-are-the-data-sources-that-collect-data-as-of-february-2023)
* [FAQs data collection](#faqs-data-collection)
	- [Which commands should I run to collect the data?](#which-commands-should-i-run-to-collect-the-data)
	- [Restarting a MongoDB database](#restarting-a-mongodb-database)

<!-- vim-markdown-toc -->

## Software requirements to follow along
The following software is needed/recommended to work with the cloud resources:

* A terminal program
	- Required to run [SSH](https://en.wikipedia.org/wiki/Secure_Shell) to establish a connection with the servers.
	- Recommended: default terminal program of any Unix operating system.
	- It should also be possible to work with Windows [Google how to establish an SSH connection with Windows (it probably requires using Putty)].
	- Further remarks: in the following code snippets meant to be executed with a terminal program if a `$` [dollar sign] appears before the command, it means that the user does NOT require sudo rights to run the code in the terminal. 
	On the other hand, if a `#` [pound sign] appears, the command should be executed as root.
		- **VERY IMPORTANT**: DO NOT, under any circumstances, write either the dollar or pound sign, when executing a command in the terminal, they are only symbolically present in the documentation to inform you about the priviledges needed to run a command.

* SSH/OpenSSH
	- This program actually establishes the SSH connection.
	- It should be in the default installation of all Unix operating systems.

* The Go programming language
	- Required to compile and run the codebase, after cloning/downloading it from GitHub.
	- It is recommended to use Go version +1.19 (previous versions might not even compile the current codebase, so it is STRONGLY advised against using previous versions to compile the codebase).
	- If Go is not present in your Mac, run the following code in your terminal to install it (it requires you to have already install the `brew` Mac package manager):
		- `# brew install go`
		- After `brew` finishes installing go, run `$ go version` to check what version of Go has been installed in your system.
	- Further remarks: if it would be too cumbersome to install Go locally in your development laptop, you can generate compiled versions of the codebase for basically all platforms using GitHub Actions and Goreleaser.

## What is the UVC data pipeline?

### What are the data sources that collect data (as of February 2023)?
As of February 2023, we are able to collect data reliably from Crunchbase.

## FAQs data collection
### Which commands should I run to collect the data?
1. Install the latest version of the UVC data collector software:
	- If you have configured your computer to have access with your SSH keys to the UVC's repos in GitHub, you can directly clone the GitHub repository to your computer by running the following code snippet in your terminal (the repo will be cloned to the directory were you currently are in your terminal session):

	```
	$ git clone git@github.com:uvc-partners/UVC_data_pipeline.git
	```
	- Otherwise, if you do not have SSH key access to the repository, visit `https://github.com/uvc-partners/UVC_data_pipeline` and download a zip file with the codebase.
	- The GitHub account that you use to login to the repository should have access to the repository. 
	You can find the credentials for the official UVC GitHub account on the `Credentials` file on the UVC Data Team Sharepoint folder.
2. Go to the base directory of the repository that was previously cloned, and compile the codebase by executing:

```
$ go build -o cbExtractor.bin ./cmd 
```

If no errors occured while compiling the codebase, you will now have an executable file named `cbExtractor.bin` at the root folder of the repository. 
You can check if the file was correctly stored at the root directory of the repository by executing the following command:

```
$ ls
```

`ls` will print all the files in the current directory.

You will use the executable that was compiled for the following tasks:
1. Extracting data from Crunchbase
2. Inserting the extracted data from Crunchbase into the cloud databases.

After you have once built this executable, you do not have to repeat steps 1 and 2 next time you want to extract data or insert the data to a database. 
Just use the same executable file that was compiled the first time.

3. Check that you have the `.env` file with the **secrets** in the same directory in which the executable file is.
The `.env` file contains all the passwords and/or sensitive data to access the databases and to extract data from Crunchbase, that should not be leaked to third parties through a repository.

If this file is missing, you won't be able to either extract data from Crunchbase or insert data into the databases.

4. Run the executable to extract data from Crunchbase by running the following command in the same folder where the executable and the `.env` file are located:

```
$ ./cbExtractor.bin extract --no-proxy
```

**Important details:**

* You should execute this file when connected to either a residential or office IP, if you try to execute this file from within a cloud instance, cloud server, cloud VM, etc. Crunchbase will flag you as a bot and you will not be able to extract any data at all.
* The extraction process can take up to an hour or more, to fully extract the thousands of start-ups available through the Crunchbase query used.
* Do not close your terminal window while the extraction process is taking place, the program will tell you how much progress it has made. 
When it is ready with the extraction, it will also tell you that through a text message in the terminal session.
* If you want to know more about the different options and subcommands available through the executable, you can always provide the executable with the `-h` or `--help` flags.
For example, `./cbExtractor.bin --help` will print a help menu on the terminal presenting all available subcommands, like `extract` and `insert`.
* If you have to abruptly cancel an ongoing data extraction before it is done, you can type `Ctrl + C` in the terminal window where the data extraction is taking place.
This will cancel the ongoing process.

5. After a successful extraction of Crunchbase data, you will now have a file in the folder where the executable is, named something along the lines of *CBData_xxxxx* where the *xxxx* are a random string of numbers.
You can now insert this data to the cloud MongoDB databases that host the data by running the following command:

```
$ ./cbExtractor.bin db insert --file <DATA_FILE> --remote <IP_DATABASE>
```

You should replace `<DATA_FILE>` with the path to the file that will be inserted into the database, in this case `./CBData_xxxx`.

`<IP_DATABASE>` should be the IP address of the remote server hosting the MongoDB instance.

**Remarks**

* I normally insert the data right away to the `production1` and `staging1` servers.
After running the command that inserts the data, I check with MongoDB Compass, if the collections within MongoDB have been updated.
Simply check if there are any documents with a timestamp equal to the timestamp when the extraction took place.
* In the unlikely event that a connection cannot be established with the server and the command to insert the documents fails, you should connect through SSH with the server and restart MongoDB.

### Restarting a MongoDB database
1. Open your terminal or command line interface.
2. You should already have in your system the SSH key that allows you to connect with the server.
The key should be stored at `~/.ssh/`
3. Type the following command to connect to the FreeBSD server via SSH:

```
$ ssh root@serverIP
```

Replace "serverIP" with the actual IP address of the server.

**Important remarks** 

* Sometimes, especially the first time that you try to connect to the server through SSH, the system does not recognize which SSH key in your system should be used to connect to the server.
You will get an error while connecting to the server.
In that case run the same ssh command but add the `-i` flag, that tells SSH specifically which SSH key to use:

```
$ ssh root@serverIP -i <PATH_SSH_KEY>
```

Replace `<PATH_SSH_KEY>` with the path to the SSH key in your system, like for example, `~/.ssh/UVC_Key`, if you named your key `UVC_Key` and stored it at `~/.ssh`

4. If everything went well with the configuration of your SSH key, you should now have access to the remote server you are trying to access.

First check the status of the MongoDB instance running in the server by executing the following command:

```
$ service mongod status
```

If the MongoDB instance is working properly, you should see the following prompt:

```
mongod is running as pid 81562.
```

The number after "pid" will always be different, it is the _Process ID_, so each time the MongoDB process is restarted it gets a new number.


If you still want to restart the MongoDB service, type the following command:

```
$ service mongod restart
```

5. Wait for the service to restart, and then type `exit` to log out of the server.

```
exit
```
