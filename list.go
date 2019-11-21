package recursivelist

import (
	"os"
	"path/filepath"
	"strings"
)

/*
RecursiveList Lists files of a directory recurrently using channels.
The output paths will be appeared with relative path.

directory string - the directory path to list files
*/
func RecursiveList(directory string) (chan string, chan error, chan bool, error) {
	return recursiveList(directory, "/")
}

func recursiveList(directory, root string) (chan string, chan error, chan bool, error) {
	out := make(chan string)
	errs := make(chan error)
	done := make(chan bool)

	dir, err := filepath.Abs(directory)
	if err != nil {
		return nil, nil, nil, err
	}

	matches, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return nil, nil, nil, err
	}
	go insertFiles(matches, root, out, errs, done)

	return out, errs, done, nil
}

func insertFiles(matches []string, root string, out chan string, errs chan error, done chan bool) {
	if len(matches) == 0 {
		// If it's a directory without file then give just directory name to out.
		out <- root
	} else {
		for _, file := range matches {
			_, fname := filepath.Split(file)
			if strings.HasPrefix(fname, ".") {
				continue
			}

			stat, err := os.Stat(file)
			if err != nil {
				errs <- err
			}

			switch mode := stat.Mode(); {
			case mode.IsRegular():
				out <- filepath.Join(root, fname)
			case mode.IsDir():
				ifiles, ierrs, idone := recursiveList(file, filepath.Join(root, fname))
			L:
				for {
					select {
					case ifile := <-ifiles:
						out <- ifile
					case ierr := <-ierrs:
						errs <- ierr
					case <-idone:
						break L
					}
				}
			}
		}
	}
	done <- true
}
