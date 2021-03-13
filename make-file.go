package make_tree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// MakeFileAction creates the file,
// which must not exist beforehand,
// and runs the inner actions inside.
// It does remove the file on rollback.
type MakeFileAction struct {
	file   string
	action func(file *os.File) error
}

// Do Ensures the specified file
// exists as a file and then
// executes the inner actions.
func (mfa *MakeFileAction) Do(baseDirectory string, dump io.Writer, logRan func(Action)) error {
	full := filepath.Join(baseDirectory, mfa.file)
	if _, err := os.Stat(full); err == nil {
		_, _ = fmt.Fprintln(dump, "The file already exists: "+full)
		return err
	} else {
		_, _ = fmt.Fprintln(dump, "Creating file: "+full)
		if file, err := os.OpenFile(full, os.O_RDWR|os.O_CREATE, 0644); err != nil {
			_, _ = fmt.Fprintln(dump, "Could not create: "+full+" because: "+err.Error())
			return err
		} else {
			defer func() {
				file.Close()
				logRan(mfa)
			}()
			return mfa.action(file)
		}
	}
}

// Rollback does not have any particular
// implementation for this type.
func (mfa *MakeFileAction) Rollback(baseDirectory string, dump io.Writer) {
	full := filepath.Join(baseDirectory, mfa.file)
	_, _ = fmt.Fprintln(dump, "Removing file: "+full)
	_ = os.Remove(full)
}

// Instantiates a MakeFileAction.
func MakeFile(file string, action func(file *os.File) error) *MakeFileAction {
	return &MakeFileAction{file: file, action: action}
}
