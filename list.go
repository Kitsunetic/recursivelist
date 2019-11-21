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

	go func() {
		dir, err := filepath.Abs(directory)
		if err != nil {
			errs <- err
			return
		}

		files, err := filepath.Glob(filepath.Join(dir, "*"))
		if err != nil {
			errs <- err
			return
		}
		go insertFiles(files, root, out, errs, done)
	}()

	return out, errs, done
}

func insertFiles(files []string, root string, out chan string, errs chan error, done chan bool) {
	if len(files) == 0 {
		// If it's a directory without file then give just directory name to out.
		out <- root
	} else {
	L1:
		for _, file := range files {
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
				if err != nil {
					errs <- err
					continue L1
				}
			L2:
				for {
					select {
					case ifile := <-ifiles:
						out <- ifile
					case ierr := <-ierrs:
						errs <- ierr
					case <-idone:
						break L2
					}
				}
			}
		}
	}
	done <- true
}
