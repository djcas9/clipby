package main

import (
	"os"
	"os/signal"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/Unknwon/com"
	"github.com/mitchellh/go-mruby"
	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	Name    = "Clipby"
	Version = "0.1.0"
)

var (
	Debug = kingpin.Flag("debug", "Enable debug mode.").Bool()
	Quiet = kingpin.Flag("quiet", "Remove all output logging.").Short('q').Bool()

	MainPath   = ""
	PluginPath = ""

	OutputChan = make(chan CBType)
	CBChan     = make(chan string)
	DoneChan   = make(chan bool)
)

func init() {
	kingpin.Version(Version)
	kingpin.Parse()

	log.SetOutput(os.Stderr)

	if *Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if *Quiet {
		log.SetLevel(log.FatalLevel)
	}

}

func main() {
	log.Infof("Initializing %s Version: %s.", Name, Version)

	home, err := com.HomeDir()

	if err != nil {
		log.Fatal(err)
	}

	MainPath := path.Join(home, ".clipby")
	PluginPath := path.Join(home, ".clipby", "plugins")

	os.Mkdir(MainPath, 0777)
	os.Mkdir(PluginPath, 0777)

	log.Info("Loading plugins.")

	vm := VM{}

	vm.Mrb = mruby.NewMrb()
	defer vm.Close()

	vm.Init()

	plugins, err := GetPlugins(PluginPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, plugin := range plugins {
		log.Info("\tLOAD: ", plugin.Name)
		vm.Load(plugin)
	}

	if len(plugins) <= 0 {
		log.Warn("No plugins found.")
	}

	// Get ruby class names
	vm.PluginNames()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log.Info("Monitoring System Clipboard.")
	go ClipBoardStart()

	for {
		select {
		case sig := <-sigChan:
			if sig.String() == "interrupt" {
				log.Debug("Closing done channel.")
				close(DoneChan)
			}
		case data := <-OutputChan:

			// fmt.Println("GOT DATA:", data)
			go vm.Run(data)

		case <-DoneChan:
			log.Info("Cleaning up...")
			os.Exit(1)
		}
	}
}
