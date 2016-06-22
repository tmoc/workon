package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"

	"github.com/urfave/cli"
)

type Project struct {
	Abbreviation string `json:"abbreviation"`
	Filepath     string `json:"filepath"`
}

type Config struct {
	EditorApp   string    `json:"editorApp"`
	TerminalApp string    `json:"terminalApp"`
	Projects    []Project `json:"projects"`
}

var configFilename = ".workon"

func main() {
	app := cli.NewApp()
	app.Name = "workon"
	app.Usage = "Quickly open projects in editor and terminal."
	app.Commands = []cli.Command{
		{
			Name:  "setup",
			Usage: "Create new workon configuration file.",
			Action: func(c *cli.Context) error {
				usr, err := user.Current()
				if err != nil {
					log.Fatal(err)
				}
				newCfg, err := json.Marshal(Config{"atom", "terminal", nil})
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(usr.HomeDir+"/"+configFilename, newCfg, 0644)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Created new configuration file in", usr.HomeDir)
				return nil
			},
		},
		{
			Name:  "add",
			Usage: "Add project abbreviation and location.",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if len(args) != 2 {
					fmt.Println("Must provide both project abbreviation and filepath.")
					return nil
				}
				abbreviation, filepath := args[0], args[1]

				usr, err := user.Current()
				if err != nil {
					log.Fatal(err)
				}
				file, err := ioutil.ReadFile(usr.HomeDir + "/" + configFilename)
				if err != nil {
					log.Fatal(err)
				}
				var cfg Config
				if err := json.Unmarshal(file, &cfg); err != nil {
					log.Fatal(err)
				}

				for _, v := range cfg.Projects {
					if v.Abbreviation == abbreviation {
						fmt.Print("Project \"", abbreviation, "\" already exists.\n")
						return nil
					}
				}
				cfg.Projects = append(cfg.Projects, Project{abbreviation, filepath})

				cfgJson, err := json.Marshal(cfg)
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(usr.HomeDir+"/"+configFilename, cfgJson, 0644)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Print("Added \"", abbreviation, "\" located at: ", filepath, "\n")
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "Remove project from config.",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if len(args) != 1 {
					fmt.Println("Must provide project abbreviation.")
					return nil
				}
				abbreviation := args[0]

				usr, err := user.Current()
				if err != nil {
					log.Fatal(err)
				}
				fi, err := ioutil.ReadFile(usr.HomeDir + "/" + configFilename)
				if err != nil {
					log.Fatal(err)
				}
				var cfg Config
				if err := json.Unmarshal(fi, &cfg); err != nil {
					log.Fatal(err)
				}

				found := false
				for i, v := range cfg.Projects {
					if v.Abbreviation == abbreviation {
						cfg.Projects = append(cfg.Projects[:i], cfg.Projects[i+1:]...)
						found = true
						break
					}
				}
				if found == false {
					fmt.Print("Project \"", abbreviation, "\" does not exist.\n")
					return nil
				}

				cfgJson, err := json.Marshal(cfg)
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(usr.HomeDir+"/"+configFilename, cfgJson, 0644)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Print("Removed project \"", abbreviation, "\".\n")
				return nil
			},
		},
		{
			Name:  "editor",
			Usage: "Set editor to open project with.",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if len(args) != 1 {
					fmt.Println("Must provide editor name.")
					return nil
				}
				editor := args[0]

				usr, err := user.Current()
				if err != nil {
					log.Fatal(err)
				}
				file, err := ioutil.ReadFile(usr.HomeDir + "/" + configFilename)
				if err != nil {
					log.Fatal(err)
				}
				var cfg Config
				if err := json.Unmarshal(file, &cfg); err != nil {
					log.Fatal(err)
				}

				cfg.EditorApp = editor

				cfgJson, err := json.Marshal(cfg)
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(usr.HomeDir+"/"+configFilename, cfgJson, 0644)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Set editor to", editor)
				return nil
			},
		},
		{
			Name:  "terminal",
			Usage: "Set terminal to open project with.",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if len(args) != 1 {
					fmt.Println("Must provide terminal name.")
					return nil
				}
				terminal := args[0]

				usr, err := user.Current()
				if err != nil {
					log.Fatal(err)
				}
				file, err := ioutil.ReadFile(usr.HomeDir + "/" + configFilename)
				if err != nil {
					log.Fatal(err)
				}
				var cfg Config
				if err := json.Unmarshal(file, &cfg); err != nil {
					log.Fatal(err)
				}

				cfg.TerminalApp = terminal

				cfgJson, err := json.Marshal(cfg)
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(usr.HomeDir+"/"+configFilename, cfgJson, 0644)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Set terminal to", terminal)
				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		fi, err := ioutil.ReadFile(usr.HomeDir + "/" + configFilename)
		if err != nil {
			log.Fatal(err)
		}

		var cfg Config
		if err := json.Unmarshal(fi, &cfg); err != nil {
			log.Fatal(err)
		}

		if len(c.Args()) == 0 {
			fmt.Println("Projects:")
			for _, v := range cfg.Projects {
				fmt.Println(v.Abbreviation + ": " + v.Filepath)
			}
			return nil
		}

		abbreviation := c.Args().First()
		found := false
		var foundProject Project
		for _, v := range cfg.Projects {
			if v.Abbreviation == abbreviation {
				found = true
				foundProject = v
				break
			}
		}
		if found == false {
			fmt.Print("Project \"", abbreviation, "\" does not exist.\n")
			return nil
		}

		fmt.Println("Working on " + foundProject.Abbreviation)
		os.Chdir(foundProject.Filepath)
		wd, _ := os.Getwd()
		exec.Command(cfg.EditorApp, wd).Output()
		exec.Command("open", "-a", cfg.TerminalApp, wd).Output()
		return nil
	}

	app.Run(os.Args)
}
