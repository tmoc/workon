package main

import (
	"fmt"
	"os"
	"os/exec"
)

var projects = map[string]string{}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Projects:\n")
		for k, v := range projects {
			fmt.Println(k + ": " + v)
		}
		return
	}

	project := os.Args[1]

	if projects[project] == "" {
		fmt.Println("No such project.")
		return
	}

	fmt.Println("Working on " + projects[project])

	os.Chdir(projects[project])
	wd, _ := os.Getwd()

	cmd1 := exec.Command("atom", wd)
	cmd2 := exec.Command("open", "-a", "iTerm", wd)
	cmd1.Output()
	cmd2.Output()
}
