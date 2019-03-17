package main

import (
	"os"

	"github.com/hauke96/kingpin"
	"github.com/hauke96/picl/cmd"
	"github.com/hauke96/picl/pkg"
	"github.com/hauke96/picl/util"
	"github.com/hauke96/sigolo"
)

var (
	version = "0.1"

	app           = kingpin.New("picl", "Maybe the dumbest package manager ever")
	appConfigFile = app.Flag("config", "Specifies the configuration file that should be used. This is \"./picl.conf\" by default.").Short('c').Default("./picl.conf").File()

	installCmd          = app.Command("install", "Installs the given library")
	installOutputFolder = installCmd.Flag("output", "Specifies the output folder where all libraries should be stored.").Short('o').String()
	installUrl          = installCmd.Flag("url", "The base url where picl downloads files from").Short('u').URL()
	installPackageName  = installCmd.Arg("package", "The library to install").Required().String()

	removeCmd         = app.Command("remove", "Uninstalls/removes the given library")
	removePackageName = removeCmd.Arg("package", "The library to remove").Required().String()
)

func configureLogging() {
	sigolo.FormatFunctions[sigolo.LOG_INFO] = sigolo.LogPlain
	//sigolo.LogLevel = sigolo.LOG_DEBUG
}

func configureCliArgs() {
	app.Author("Hauke Stieler")
	app.Version(version)

	app.CustomDescription("Package Name", `This name if the library name including the version you wan't do deal with. The name has the following format:

  my-library@3.5.1

There must be a name and there must be a version. The version is basically the string that is behind the "@" and is not parsed. It just has to exist on the server, but the format "x.y.z" (e.g. 3.5.1) is only recommended.`)

	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')
}

// When a field (e.g. configOutputFolder) is not set yet, but has been specified
// via a command line argument, it'll be set here. Here we also overwrite
// existing values (e.g. the remote URL).
func setConfigFromArguments() {
	if *installOutputFolder != "" {
		configOutputFolder = *installOutputFolder
	}

	if *installUrl != nil {
		configRemoteUrl = *installUrl
	}
}

// The output folder and remote url has to be set when installing a package. This
// function will exit with 1 when one of those is not set.
func handleInvalidInstallConfigs() {
	if configOutputFolder == "" {
		sigolo.Error("Output folder not set")
		os.Exit(1)
	}

	if configRemoteUrl == nil {
		sigolo.Error("Remote url not set")
		os.Exit(1)
	}
}

// The output folder has to be set when removing a package. This function will
// exit with 1 when it's not set.
func handleInvalidRemoveConfigs() {
	if configOutputFolder == "" {
		sigolo.Error("Output folder not set")
		os.Exit(1)
	}
}

func parseCommand() string {
	command, err := app.Parse(os.Args[1:])

	if err != nil {
		sigolo.Error("Error while parsing arguments:\n%s", err)
		os.Exit(1)
	}

	return command
}

func main() {
	configureLogging()
	configureCliArgs()

	command := parseCommand()

	readConfig(*appConfigFile)
	setConfigFromArguments()

	switch command {
	case installCmd.FullCommand():
		handleInvalidInstallConfigs()
		pkg, err := pkg.ParsePackage(*installPackageName)
		util.ExitOnError(err)

		err = cmd.Install(pkg, configOutputFolder, configRemoteUrl)
		util.ExitOnError(err)
	case removeCmd.FullCommand():
		handleInvalidRemoveConfigs()
		pkg, err := pkg.ParsePackage(*removePackageName)
		util.ExitOnError(err)

		err = cmd.Remove(pkg, configOutputFolder)
		util.ExitOnError(err)
	}
}
