package bookmark

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

// Repository is needed for an Action
type Repository interface {
	Save(bookmark Bookmark) error
	// Find returns a Bookmark or an error
	Find() (Bookmark, error)
	// Delete returns true if successfully deleted, else error
	Delete() (bool, error)
}

type Action struct {
	repository     Repository
	activeBookmark Bookmark // should result always be a bookmark?
	err            error
}

func NewAction(repository Repository) *Action {
	return &Action{
		repository: repository,
	}
}

func (action *Action) Save(bookmark Bookmark) *Action {
	if err := action.repository.Save(bookmark); err != nil {
		action.err = fmt.Errorf("failed to save bookmark: %w", err)
		return action
	}
	return action
}

func (action *Action) Find() *Action {
	bookmark, err := action.repository.Find()
	if err != nil {
		action.err = err
		return action
	}
	action.activeBookmark = bookmark
	return action
}

func (action *Action) Delete() *Action {
	_, err := action.repository.Delete()
	if err != nil {
		action.err = err
		return action
	}
	return action
}

func (action *Action) Open() *Action {
	openLink(action.activeBookmark.Url)
	return action
}

// And will simply return Action if
// the previous action succeeded, else
// it will log the error and abort
func (action *Action) And() *Action {
	if action.err != nil {
		log.Fatalf("failed to run next action, previous failed with err %v\n", action.err)
	}
	return action
}

func openLink(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
