package make_tree

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// UseDirectoryAction creates the
// target directory, but expects it
// to exist, and runs the inner actions
// inside. Conversely, it does not
// remove the directory on rollback.
type MakeDirectoryAction struct {
	directory string
	actions   []Action
}

// Creates the specified directory,
// if it does not exist, or raises
// an error. After that, it executes
// all of its inner actions.
func (mda *MakeDirectoryAction) Do(baseDirectory string, dump io.Writer) error {
	full := filepath.Join(baseDirectory, mda.directory)
	if _, err := os.Stat(full); err == nil {
		return errors.New("the path already exists: " + full)
	} else if err := os.Mkdir(full, 0755); err != nil {
		_, _ = fmt.Fprintln(dump, "Could not create the directory: "+full+" because: "+err.Error())
		return err
	} else {
		_, _ = fmt.Fprintln(dump, "Created directory: "+full)
		return doTree(full, mda.actions, dump)
	}
}

// Removes the just created directory.
func (mda *MakeDirectoryAction) Rollback(baseDirectory string, dump io.Writer) {
	full := filepath.Join(baseDirectory, mda.directory)
	_, _ = fmt.Fprintln(dump, "Removed directory: "+full)
	_ = os.Remove(full)
}
