/*
history:
2016/06/08 v1

GoFmt GoBuildNull GoBuild GoRelease GoRun
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func log(msg string, args ...interface{}) {
	const Beat = time.Duration(24) * time.Hour / 1000
	tzBiel := time.FixedZone("Biel", 60*60)
	t := time.Now().In(tzBiel)
	ty := t.Sub(time.Date(t.Year(), 1, 1, 0, 0, 0, 0, tzBiel))
	td := t.Sub(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, tzBiel))
	ts := fmt.Sprintf(
		"%d/%d@%d",
		t.Year()%1000,
		int(ty/(time.Duration(24)*time.Hour))+1,
		int(td/Beat),
	)
	fmt.Fprintf(os.Stderr, ts+" "+msg+"\n", args...)
}

func printinfo(path string, info os.FileInfo) {
	if info.Mode().IsDir() {
		path = path + string(os.PathSeparator)
	}

	fmt.Printf(
		"%s\t%04o\t%db\n",
		//"%s %04o %s %db\n",
		path,
		info.Mode()&os.ModePerm,
		//info.ModTime().UTC().Format("06/0102/15:04"),
		info.Size(),
	)
}

func fls(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	printinfo(path, info)

	return nil
}

func ls(path string) error {
	var err error

	pathstat, err := os.Lstat(path)
	if err != nil {
		return err
	}

	if pathstat.Mode().IsDir() {
		ff, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}

		for _, filestat := range ff {
			filepath := filepath.Join(path, filestat.Name())
			printinfo(filepath, filestat)
		}
	} else {
		printinfo(path, pathstat)
	}

	return nil
}

func main() {
	var err error

	var recursive bool
	if os.Args[0] == "lr" {
		recursive = true
	}

	paths := os.Args[1:]
	if len(os.Args) > 1 && os.Args[1] == "-r" {
		recursive = true
		paths = paths[1:]
	}

	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	for _, path := range paths {
		err = do(path, recursive)
		if err != nil {
			log("%v", err)
			os.Exit(1)
		}
	}
}

func do(path string, recursive bool) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	path = filepath.Clean(path)

	if recursive {
		err = filepath.Walk(path, fls)
		if err != nil {
			return err
		}
	} else {
		err = ls(path)
		if err != nil {
			return err
		}
	}

	return nil
}
