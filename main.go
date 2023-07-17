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
	fmt.Println("Enter command followed by desired input (search, dir, searchm, fuzzy): ")
	selec, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	selec = selec[:len(selec)-1]

	if strings.Contains(selec, "searchm") {
		file := strings.Split(string(selec), " ")
		searchM(file[1], "/")
	} else if strings.Contains(selec, "search") {
		file := strings.Split(string(selec), " ")
		search(file[1], "/")
	} else if strings.Contains(selec, "dir") {
		dir := strings.Split(string(selec), " ")[1]
		dirTrav(dir, 0)
	} else if strings.Contains(selec, "fuzzy"){
		TARGET := strings.Split(string(selec), " ")[1]
		colorReset := "\033[0m"
		colorYellow := "\033[33m"

		outputFile, _ := os.Create("./paths.txt")
		defer outputFile.Close()
		recDirTrav("/home", outputFile)

		ofile, err := os.Open("./paths.txt")
		if err != nil {
			fmt.Println(err)
		}

		scn := bufio.NewScanner(ofile)
		for scn.Scan() {
			line := scn.Text()
			path := strings.Split(line, "/")
			path = strings.Fields(strings.Join(path, " "))
		
			for _, word := range path {
				word = strings.Split(word, ".")[0]
				if len(word) < 4 {
					continue
				} else if len(word) > len(TARGET) * 2 { // Word is too big for soundex algorithm
					word1 := word[:len(word) / 2]
					word2 := word[len(word)/2:]
				
					if compare(word1, TARGET) {
						windex := strings.Index(line, word)
						fmt.Println(line[:windex] + string(colorYellow) + (word) + string(colorReset) + line[windex+len(word):])
						//fmt.Println(string(colorYellow), (line), string(colorReset))
					} else if compare(word2, TARGET) {
						windex := strings.Index(line, word)
						fmt.Println(line[:windex] + string(colorYellow) + (word) + string(colorReset) + line[windex+len(word):])
						//fmt.Println(string(colorYellow), (line), string(colorReset))
					}
				} else if len(word) >= len(TARGET){ // Making word isn't too small (search for "face", don't want to return "fs"
					if compare(word, TARGET) {
						windex := strings.Index(line, word)
						fmt.Println(line[:windex] + string(colorYellow) + (word) + string(colorReset) + line[windex+len(word):])
						break
					}
				}
			}
		}

	} else {
		return
	}
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

func dirTrav(root string, indent int){
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
		dirTrav(root + "/" + input, indent + 1)
	} else if input == "q" {
		return
	} else {
		fmt.Println("Opening file")
	}
}

func compare(word1 string, word2 string) bool {
	if calc(word1) == calc(word2) {
		return true
	}
	return false
}

func calc(str string) string {
	var code string
	code += string(str[0])
	
	for i := 1; i < len(str); i++ {
		if strings.Contains(".", string(str[i])) {
			break
		} else if strings.Contains("bfpv", string(str[i])) {
			code += "1"
		} else if strings.Contains("cgjkqsxz", string(str[i])) {
			code += "2"
		} else if strings.Contains("dt", string(str[i])) {
			code += "3"
		} else if strings.Contains("l", string(str[i])) {
			code += "4"
		} else if strings.Contains("mn", string(str[i])) {
			code += "5"
		} else if strings.Contains("r", string(str[i])) {
			code += "6"
		} else {
			code += "0"
		}
		if len(code) == 4 {
			code = removeZeroes(code)
		}
	}
	
	code = removeZeroes(code)

	if len(code) < 4 {
		for i := len(code); i < 4; i++ {
			code += "0"
		}
	} else if len(code) > 4 {
		code = code[:4]
	}

	return code
}

func removeZeroes(str string) string {
	for i := 0; i < len(str); i++ {
		if string(str[i]) == "0" {
			str = str[:i] + str[i + 1:]
			i--;
		}
	}
	return str
}

func recDirTrav(root string, ofile *os.File) {
	files, err := os.ReadDir(root)
	if err != nil {
		return
	}

	for _, file := range files {
		if !strings.Contains(file.Name(), ".") {
			recDirTrav(root + "/" + file.Name(), ofile)
		} else {
			ofile.WriteString(root + "/" + file.Name() + "\n")
		}
	}
}
