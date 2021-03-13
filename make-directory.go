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
// inside. It does remove the directory
// on rollback.
type MakeDirectoryAction struct {
	directory string
	actions   []Action
}

// Creates the specified directory,
// if it does not exist, or raises
// an error. After that, it executes
// all of its inner actions.
func (mda *MakeDirectoryAction) Do(baseDirectory string, dump io.Writer, logRan func(Action)) error {
	full := filepath.Join(baseDirectory, mda.directory)
	if _, err := os.Stat(full); err == nil {
		_, _ = fmt.Fprintln(dump, "The path already exists: "+full)
		return errors.New("the path already exists: " + full)
	} else if err := os.Mkdir(full, 0755); err != nil {
		_, _ = fmt.Fprintln(dump, "Could not create the directory: "+full+" because: "+err.Error())
		return err
	} else {
		_, _ = fmt.Fprintln(dump, "Creating directory: "+full)
		logRan(mda)
		return doTree(full, mda.actions, dump, logRan)
	}
}

// Removes the just created directory.
func (mda *MakeDirectoryAction) Rollback(baseDirectory string, dump io.Writer) {
	full := filepath.Join(baseDirectory, mda.directory)
	_, _ = fmt.Fprintln(dump, "Removing directory: "+full)
	_ = os.Remove(full)
}

// Instantiates a MakeDirectoryAction.
func MakeDirectory(directory string, actions []Action) *MakeDirectoryAction {
	return &MakeDirectoryAction{directory: directory, actions: actions}
}
