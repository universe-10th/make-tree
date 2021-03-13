package make_tree

import "io"

// doTree invokes all the given actions' Do methods to implement
// the intended changes. On failure, it invoked rollbackTree to
// undo everything.
func doTree(baseDirectory string, actions []Action, dump io.Writer, logRan func(Action)) error {
	for _, action := range actions {
		if err := action.Do(baseDirectory, dump, logRan); err != nil {
			return err
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
	var ran []Action

	logRan := func(action Action) {
		ran = append(ran, action)
	}

	if err := doTree(baseDirectory, actions, dump, logRan); err != nil {
		rollbackTree(baseDirectory, ran, dump)
		return err
	} else {
		return nil
	}
}
