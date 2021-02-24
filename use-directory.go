package make_tree

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// UseDirectoryAction does not create
// the directory, but expects it to
// exist, and runs the inner actions
// inside. Conversely, it does not
// remove the directory on rollback.
type UseDirectoryAction struct {
	directory string
	actions   []Action
}

// Do Ensures the specified directory
// exists as a directory and then executes
// the inner actions.
func (uda *UseDirectoryAction) Do(baseDirectory string, dump io.Writer) error {
	full := filepath.Join(baseDirectory, uda.directory)
	if info, err := os.Stat(full); err != nil {
		return err
	} else if !info.IsDir() {
		return errors.New("the path must be a directory: " + full)
	} else {
		return doTree(full, uda.actions, dump)
	}
}

// Rollback does not have any particular
// implementation for this type.
func (uda *UseDirectoryAction) Rollback(baseDirectory string, dump io.Writer) {}

// Instantiates a UseDirectoryAction.
func UseDirectory(directory string, actions []Action) *UseDirectoryAction {
	return &UseDirectoryAction{directory: directory, actions: actions}
}
