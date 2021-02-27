/*
history:
2016/06/08 v1
20/1106 print only file name by default, add options -r, -s, -t, -m, -l

GoFmt GoBuildNull GoBuild GoRelease GoRun
*/

package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Version string

	Recursive bool
	ShowTime  bool
	ShowSize  bool
	ShowMode  bool
)

func seps(i int, e int) string {
	ee := int(math.Pow(10, float64(e)))
	if i < ee {
		return fmt.Sprintf("%d", i%ee)
	} else {
		f := fmt.Sprintf("0%dd", e)
		return fmt.Sprintf("%s.%"+f, seps(i/ee, e), i%ee)
	}
}

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

func printinfo(path string, info os.FileInfo) error {
	var err error

	s := path

	var finfo os.FileInfo = info
	if (info.Mode() & os.ModeSymlink) != 0 {
		finfo, err = os.Lstat(path)
		if err != nil {
			return err
		}
		var linkpath string
		linkpath, err = os.Readlink(path)
		if err != nil {
			return err
		}
		s += "@" + linkpath + "@"
	}

	if info.Mode().IsDir() {
		s += string(os.PathSeparator)
	}

	if ShowMode {
		s += fmt.Sprintf("\t%04o", finfo.Mode()&os.ModePerm)
	}
	if ShowSize {
		s += fmt.Sprintf("\t%s.bytes", seps(int(finfo.Size()), 3))
	}
	if ShowTime {
		s += fmt.Sprintf("\t%s", finfo.ModTime().UTC().Format("06/0102/15:04"))
	}

	fmt.Println(s)
	return nil
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

func init() {
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Println(Version)
		os.Exit(0)
	}
}

func main() {
	var err error

	switch os.Args[0] {
	case "lr":
		Recursive = true
	case "ll":
		ShowMode = true
		ShowTime = true
		ShowSize = true
	case "llr":
		Recursive = true
		ShowMode = true
		ShowTime = true
		ShowSize = true
	}

	paths := os.Args[1:]
	for _, p := range paths {
		if !strings.HasPrefix(p, "-") {
			break
		}
		switch p {
		case "-r":
			Recursive = true
		case "-m":
			ShowMode = true
		case "-s":
			ShowSize = true
		case "-t":
			ShowTime = true
		case "-l":
			ShowMode = true
			ShowTime = true
			ShowSize = true
		default:
			log("invalid option `%s`", p)
			os.Exit(1)
		}
		paths = paths[1:]
	}

	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	for _, p := range paths {
		err = do(p)
		if err != nil {
			log("%v", err)
			os.Exit(1)
		}
	}
}

func do(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	path = filepath.Clean(path)

	if Recursive {
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
