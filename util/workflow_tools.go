package util

import (
	aw "github.com/deanishe/awgo"
)

func NewURLItem(wf *aw.Workflow, title string) *aw.Item {
	item := wf.NewItem(title).
		Valid(true).
		Var("action", "open-url")

	item.Cmd().Subtitle("Copy URL to clipboard").
		Var("action", "copy-to-clipboard")

	return item
}
