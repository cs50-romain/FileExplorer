package main

import (
	"fmt"
	//"log"
	"os"
	//"github.com/lxn/walk"
	//. "github.com/lxn/declarative"
	"strings"
	"bufio"
)

var STOP bool

func main() {
	fmt.Print("Selection (search, searchM, directory): ")
	selec, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	selec = selec[:len(selec)-1]

	if selec == "search" {
		fmt.Print("Search: ")
		input,_ := bufio.NewReader(os.Stdin).ReadString('\n')
		input = input[:len(input)-1]
		search(input, "/")
	} else if selec == "searchM" {
		fmt.Print("Search (all location): ")
		input,_ := bufio.NewReader(os.Stdin).ReadString('\n')
		input = input[:len(input)-1]
		searchM(input, "/")
	} else if selec == "directory" {
		fmt.Print("Directory: ")
		input,_ := bufio.NewReader(os.Stdin).ReadString('\n')
		input = input[:len(input)-1]
		recDirTrav(input, 0)
	} else {
		return
	}
//	search("main.go", "/")
	//searchM("main.go", "/")
	//fmt.Print("Directory: ")
	//input,_ := bufio.NewReader(os.Stdin).ReadString('\n')
	//input = input[:len(input)-1]
	//recDirTrav(input, 0)
}

// returns all possible locations
func searchM(ffile string, dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		//fmt.Println(file.Name())
		if file.Name() == ffile {
			fmt.Println("Found file: ", file.Name(), " -> ", dir[1:])
			//searchM(ffile, "/")
			return
		}
		if !strings.Contains(file.Name(), ".") {
			searchM(ffile, dir + "/" + file.Name())
		} else {
		}
	}

}

// Returns only one first location found
func search(ffile string, dir string) {
	if STOP {
		return
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		//fmt.Println(file.Name())
		if file.Name() == ffile {
			STOP = true
			fmt.Println("Found file: ", file.Name(), " -> ", dir[1:])
			search(ffile, "/")
		}
		if !strings.Contains(file.Name(), ".") {
			search(ffile, dir + "/" + file.Name())
		} else {
		}
	}

}

func recDirTrav(root string, indent int){
	fmt.Println(root)
	files, err := os.ReadDir(root)
	if err != nil {
		return
	}

	for _, file := range files {
		if !strings.Contains(file.Name(), ".") {
			fmt.Println(strings.Repeat("    ", indent),"/" + file.Name())
			//recDirTrav(root + "/" + file.Name(), indent + 1)
		} else {
			fmt.Println(strings.Repeat("    ", indent),"*", file.Name())
			fmt.Println()
		}
	}
	fmt.Print("Select file/folder (q to quit): ")
	input,_ := bufio.NewReader(os.Stdin).ReadString('\n')
	input = input[:len(input)-1]

	if !strings.Contains(input, ".") {
		recDirTrav(root + "/" + input, indent + 1)
	} else if input == "q" {
		return
	} else {
		fmt.Println("Opening file")
	}

}
