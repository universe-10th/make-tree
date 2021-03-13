package make_tree

import "io"

// Actions are the callbacks to use to
// make a particular filesystem tree.
type Action interface {
	// Executes the direct action. It must
	// return an error on failure.
	Do(currentDirectory string, dump io.Writer, logRan func(string, Action)) error
	// Executes the inverse action. It should
	// not return error on failure, but silently
	// forgive.
	Rollback(currentDirectory string, dump io.Writer)
}
