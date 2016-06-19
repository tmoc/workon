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

var configFilename = ".workonconfig.json"

func main() {
	app := cli.NewApp()
	app.Name = "workon"
	app.Usage = "Quickly open projects in editor and terminal."
	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "Add project abbreviation and location.",
			Action: func(c *cli.Context) error {
				if len(c.Args()) != 2 {
					fmt.Println("Must provide both project abbreviation and filepath.")
					return nil
				}
				newProject := Project{c.Args()[0], c.Args()[1]}
				var cfg Config
				usr, err := user.Current()
				if err != nil {
					log.Fatal(err)
				}
				fi, err := ioutil.ReadFile(usr.HomeDir + "/" + configFilename)
				if err != nil {
					log.Fatal(err)
				}
				if err := json.Unmarshal(fi, &cfg); err != nil {
					log.Fatal(err)
				}

				for _, v := range cfg.Projects {
					if v.Abbreviation == newProject.Abbreviation {
						fmt.Print("Project \"", newProject.Abbreviation, "\" already exists.\n")
						return nil
					}
				}
				cfg.Projects = append(cfg.Projects, newProject)

				data, err := json.Marshal(cfg)
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(usr.HomeDir+"/"+configFilename, data, 0644)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Print("Added \"", c.Args()[0], "\" located at: ", c.Args()[1], "\n")
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "Remove project from config.",
			Action: func(c *cli.Context) error {
				if len(c.Args()) != 1 {
					fmt.Println("Must provide project abbreviation.")
					return nil
				}
				abbreviation := c.Args()[0]
				var cfg Config
				usr, err := user.Current()
				if err != nil {
					log.Fatal(err)
				}
				fi, err := ioutil.ReadFile(usr.HomeDir + "/" + configFilename)
				if err != nil {
					log.Fatal(err)
				}
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

				data, err := json.Marshal(cfg)
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(usr.HomeDir+"/"+configFilename, data, 0644)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Print("Removed project \"", c.Args().First(), "\".\n")
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
