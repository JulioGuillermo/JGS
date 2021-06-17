package main

import (
	"fmt"
	"github.com/JulioGuillermo/jgs/JGScript"
	"os"
	"strings"
)

func showHelp() {
	fmt.Println("JG Script version 2.0")
	fmt.Println("Usage: jgs [options] [file]")
	fmt.Println()
	fmt.Println("      -c --compile      Compile the code to a binary AST.")
	fmt.Println()
	fmt.Println("      -r --recursive    Used to compile a directory recursive.")
	fmt.Println()
	fmt.Println("      -i --interactive  Run the interactive mode.")
	fmt.Println()
	fmt.Println("      -h --help         Show this help.")
	os.Exit(0)
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
        interactive := false
		compile := false
		recursive := false
		file := ""
		out := ""
		for _, a := range args {
			if strings.HasPrefix(a, "-") {
				switch a {
				case "-h":
					fallthrough
				case "--help":
					showHelp()
				case "-c":
					compile = true
				case "--compile":
					compile = true
				case "-r":
					recursive = true
				case "--recursive":
					recursive = true
				case "-i":
					interactive = true
				case "--interactive":
					interactive = true
				default:
					fmt.Println("Unkown argument: ", a)
					os.Exit(-1)
				}
			} else {
				if file == "" {
					file = a
				} else if out == "" {
					out = a
				} else {
					fmt.Println("Unkown argument: ", a)
					os.Exit(-1)
				}
			}
		}
		if compile {
			if file == "" {
				file = "."
				recursive = true
			}
			JGScript.Compile(file, out, recursive)
        } else if interactive {
            JGScript.RunInterpreter()
		} else {
			JGScript.RunFile(file)
		}
	} else {
		showHelp()
	}
}
