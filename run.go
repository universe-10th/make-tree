package make_tree

import "io"

// doTree invokes all the given actions' Do methods to implement
// the intended changes. On failure, it invoked rollbackTree to
// undo everything.
func doTree(baseDirectory string, actions []Action, dump io.Writer) error {
	ran := []Action{}
	for _, action := range actions {
		if err := action.Do(baseDirectory, dump); err != nil {
			rollbackTree(baseDirectory, actions, dump)
			return err
		} else {
			ran = append(ran, action)
		}
	}
	return nil
}

// rollbackTree invokes all the given actions' Rollback methods to
// implement the intended changes. On failure, it does nothing but
// silently moves forward.
func rollbackTree(baseDirectory string, actions []Action, dump io.Writer) {
	for _, action := range actions {
		action.Rollback(baseDirectory, dump)
	}
}
