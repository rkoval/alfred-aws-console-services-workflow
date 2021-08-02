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

func NewBlankItem(wf *aw.Workflow) *aw.Item {
	item := wf.NewItem("\n").
		Icon(&aw.Icon{Value: "blank"})

	item.Cmd().Subtitle("\n")
	item.Alt().Subtitle("\n")
	item.Ctrl().Subtitle("\n")
	return item
}
