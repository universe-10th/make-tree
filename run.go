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
// silently moves forward. The actions are executed in reverse order.
func rollbackTree(baseDirectory string, actions []Action, dump io.Writer) {
	last := len(actions) - 1
	for index, _ := range actions {
		actions[last-index].Rollback(baseDirectory, dump)
	}
}

// MakeTree takes a list of actions and executes them one by one
// but recursively on a given base directory.
func MakeTree(baseDirectory string, actions []Action, dump io.Writer) error {
	return doTree(baseDirectory, actions, dump)
}
