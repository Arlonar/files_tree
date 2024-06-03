package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func filterDirs(files []os.DirEntry) []os.DirEntry {
	ans := make([]os.DirEntry, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			ans = append(ans, file)
		}
	}
	return ans
}

func dirTreeRecur(path string, prefix string, printFiles bool) (string, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	if !printFiles {
		dir = filterDirs(dir)
	}
	newPrefix := "├───"
	addPrefix := "│\t"
	str := strings.Builder{}
	for index, file := range dir {
		if index == len(dir)-1 {
			newPrefix = "└───"
			addPrefix = "\t"
		}
		str.WriteString(prefix)
		str.WriteString(newPrefix)
		str.WriteString(file.Name())
		info, _ := file.Info()
		if !file.IsDir() {
			if info.Size() > 0 {
				str.WriteString(fmt.Sprintf(" (%db)\n", info.Size()))
			} else {
				str.WriteString(" (empty)\n")
			}
		} else {
			str.WriteString("\n")
			add, err := dirTreeRecur(path+"/"+file.Name(), prefix+addPrefix, printFiles)
			if err != nil {
				return "", err
			}
			str.WriteString(add)
		}
	}
	return str.String(), nil
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	result, err := dirTreeRecur(path, "", printFiles)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(output, result[:len(result)-1])
	if err != nil {
		return err
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
