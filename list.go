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
func RecursiveList(directory string) (chan string, chan error, chan bool) {
	return recursiveList(directory, "/")
}

func recursiveList(directory, root string) (chan string, chan error, chan bool) {
	out := make(chan string)
	errs := make(chan error)
	done := make(chan bool)

	dir, err := filepath.Abs(directory)
	if err != nil {
		errs <- err
	}

	matches, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		errs <- err
	}

	go func() {
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
				ifiles, ierrs, idone := listRecurrent(file, filepath.Join(root, fname))
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
		done <- true
	}()

	return out, errs, done
}
