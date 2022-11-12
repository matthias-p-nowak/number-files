package main

import (
	"flag"
	"fmt"

	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	prefix string
	start  uint64 = 0
)

func scanNames(path string, info os.FileInfo, err error) error {
	if info.Mode().IsRegular() {
		// dir := filepath.Dir(path)
		ext := filepath.Ext(path)
		base := filepath.Base(path)
		base = strings.TrimSuffix(base, ext)
		parts := strings.Split(base, "-")
		// fmt.Printf("got path %s dir=%s base=%s ext=%s\n", path, dir, base, ext)
		if len(parts) > 1 {
			if parts[0] == prefix {
				// fmt.Printf("got one match \n")
				num, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return err
				}
				if num > start {
					start = num
				}
			}
		}
	}
	return nil
}
func rename(path string, info os.FileInfo, err error) error {
	if info.Mode().IsRegular() {
		dir := filepath.Dir(path)
		ext := filepath.Ext(path)
		start++
		fn := fmt.Sprintf("%s-%d%s", prefix, start, ext)
		newPath := filepath.Join(dir, fn)
		fmt.Printf("renaming %s -> %s\n", path, newPath)
		os.Rename(path, newPath)
	}
	return nil
}
func main() {
	fmt.Println("numbering files")
	defer fmt.Println("all done")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	args := flag.Args()
	fmt.Printf("#%v\n", args)
	if len(args) < 2 {
		fmt.Println("use with <prefix> <directory>")
		os.Exit(1)
	}
	prefix = args[0]
	dir := args[1]
	fmt.Printf("prefix='%s' dir='%s'\n", prefix, dir)
	err := filepath.Walk(dir, scanNames)
	if err != nil {
		fmt.Errorf("walk failed ", err)
	}
	err = filepath.Walk(dir, rename)
}
