package util

import (
	"bufio"
	"fmt"
	"os"
)

func GenerateFile(data []byte, dryRun bool) error {
	var f *os.File
	defer f.Close()
	if dryRun {
		f = os.Stdout
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		filepath := fmt.Sprintf("%s/.vimspector.json", wd)
		f, err = os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
	}

	w := bufio.NewWriter(f)
	_, err := w.Write(data)
	if err != nil {
		return err
	}
	return w.Flush()
}
