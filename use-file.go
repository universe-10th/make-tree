package make_tree

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// UseFileAction does not create the
// file, but expects it to exist, and
// runs the inner actions inside.
// Conversely, it does not remove
// the file on rollback.
type UseFileAction struct {
	file   string
	action func(file *os.File) error
	append bool
}

// Do Ensures the specified file
// exists as a file and then
// executes the inner actions.
func (ufa *UseFileAction) Do(baseDirectory string, dump io.Writer, logRan func(string, Action)) error {
	full := filepath.Join(baseDirectory, ufa.file)
	if info, err := os.Stat(full); err != nil {
		_, _ = fmt.Fprintln(dump, "The file does not exist: "+full)
		return err
	} else if info.IsDir() {
		_, _ = fmt.Fprintln(dump, "The path must not be a directory: "+full)
		return errors.New("the path must not be a directory: " + full)
	} else {
		_, _ = fmt.Fprintln(dump, "Using file: "+full)
		mode := os.O_RDWR
		if ufa.append {
			mode |= os.O_APPEND
		}
		if file, err := os.OpenFile(full, mode, 0666); err != nil {
			_, _ = fmt.Fprintln(dump, "Could not open file for edition: "+full+" because: "+err.Error())
			return err
		} else {
			defer func() {
				file.Close()
				logRan(baseDirectory, ufa)
			}()
			return ufa.action(file)
		}
	}
}

// Rollback does not have any particular
// implementation for this type.
func (ufa *UseFileAction) Rollback(baseDirectory string, dump io.Writer) {}

// Instantiates a UseFileAction.
func UseFile(file string, append bool, action func(file *os.File) error) *UseFileAction {
	return &UseFileAction{file: file, append: append, action: action}
}
