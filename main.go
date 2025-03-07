package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . -f")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	// debug data
	// path := "testdata"
	// printFiles := true
	//debug data end
	//fmt.Println(path, printFiles)
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, status bool) (err error) {

	if !status {
		return err
	}

	err = ExplorationDirs(out, path, "")

	if err != nil {
		panic(err.Error())
	}

	err = nil

	return err
}

func ExplorationDirs(out io.Writer, path string, prefix string) (err error) {

	entries, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("В директории ошибка")
		return
	}

	err = SortNameDir(entries)
	if err != nil {
		return
	}

	var infoFileSize string

	for _, file := range entries {

		currentPrefix := prefix

		fileInfo, err := file.Info()
		if err != nil {
			return err
		}

		infoFileSize = CheckFileSize(fileInfo)

		if file.IsDir() {

			if file == entries[len(entries)-1] {
				fmt.Fprint(out, currentPrefix)
				fmt.Fprintln(out, "└───"+file.Name())
				currentPrefix += "\t"
			} else {
				fmt.Fprintln(out, currentPrefix+"├───"+file.Name())
				currentPrefix += "│" + "\t"
			}
			path += "\\" + file.Name()

			ExplorationDirs(out, path, currentPrefix)
			path = filepath.Dir(path)
		} else {
			// обычным файлам не нужны линии если они находятся внутри последней папки массива
			if file == entries[len(entries)-1] {
				fmt.Fprint(out, currentPrefix+"└───"+file.Name())
				fmt.Fprintln(out, " ("+infoFileSize+")")
				currentPrefix += "\t"
			} else {
				fmt.Fprint(out, currentPrefix+"├───"+file.Name())
				fmt.Fprintln(out, " ("+infoFileSize+")")

				if prefix != "" {
					currentPrefix += "│" + "\t"
				}
			}

		}

	}

	err = nil
	return err
}

func SortNameDir(names []os.DirEntry) (err error) {
	sort.Slice(names, func(i int, j int) bool {
		return names[i].Name() < names[j].Name()
	})
	return err
}

func CheckFileSize(info fs.FileInfo) (a string) {
	if info.Size() == 0 {
		return "empty"
	} else {
		return fmt.Sprint(info.Size()) + "b"
	}
}
