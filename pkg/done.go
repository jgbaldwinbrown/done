package done

import (
	"bufio"
	"fmt"
	"iter"
	"log"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DeletePath(path string) error {
	if FileExists(path) {
		if e := os.Remove(path); e != nil {
			return e
		}
	}
	return nil
}

func IsDone(path string) bool {
	return FileExists(path + ".done")
}

func CreateDone(path string) error {
	_, e := os.Create(path + ".done")
	return e
}
func Donepath(path string) string {
	return path + ".done"
}

func Touch(path string) error {
	w, e := os.Create(path)
	if e != nil {
		return e
	}
	return w.Close()
}

func MakeIfNeeded(donepath string, makef func() error) error {
	if FileExists(donepath) {
		return nil
	}
	if e := makef(); e != nil {
		return e
	}
	return Touch(donepath)
}

func MakeIfNotDone(outpath string, makef func() error) error {
	return MakeIfNeeded(outpath+".done", makef)
}

func FindIncompletes[I iter.Seq[S], S ~string](targets I) iter.Seq[S] {
	return func(yield func(S) bool) {
		for target := range targets {
			if !IsDone(string(target)) && !yield(target) {
				return
			}
		}
	}
}

func FullFindIncompletes() {
	w := bufio.NewWriter(os.Stdout)
	defer func() {
		if e := w.Flush(); e != nil {
			log.Fatal(w)
		}
	}()
	s := bufio.NewScanner(os.Stdin)
	s.Buffer([]byte{}, 1e15)
	for s.Scan() {
		if !IsDone(s.Text()) {
			if _, e := fmt.Fprintf(w, "%v\n", s.Text()); e != nil {
				log.Fatal(e)
			}
		}
	}
}
