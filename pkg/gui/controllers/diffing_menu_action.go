package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type DiffingMenuAction struct {
	c *ControllerCommon
}

func (self *DiffingMenuAction) Call() error {
	names := self.c.Helpers().Diff.CurrentDiffTerminals()

	menuItems := []*types.MenuItem{}
	for _, name := range names {
		name := name
		menuItems = append(menuItems, []*types.MenuItem{
			{
				Label: fmt.Sprintf("%s %s", self.c.Tr.LcDiff, name),
				OnPress: func() error {
					self.c.Modes().Diffing.Ref = name
					// can scope this down based on current view but too lazy right now
					return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				},
			},
		}...)
	}

	menuItems = append(menuItems, []*types.MenuItem{
		{
			Label: self.c.Tr.LcEnterRefToDiff,
			OnPress: func() error {
				return self.c.Prompt(types.PromptOpts{
					Title:               self.c.Tr.LcEnteRefName,
					FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRefsSuggestionsFunc(),
					HandleConfirm: func(response string) error {
						self.c.Modes().Diffing.Ref = strings.TrimSpace(response)
						return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
					},
				})
			},
		},
	}...)

	if self.c.Modes().Diffing.Active() {
		menuItems = append(menuItems, []*types.MenuItem{
			{
				Label: self.c.Tr.LcSwapDiff,
				OnPress: func() error {
					self.c.Modes().Diffing.Reverse = !self.c.Modes().Diffing.Reverse
					return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				},
			},
			{
				Label: self.c.Tr.LcExitDiffMode,
				OnPress: func() error {
					self.c.Modes().Diffing = diffing.New()
					return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				},
			},
		}...)
	}

	return self.c.Menu(types.CreateMenuOptions{Title: self.c.Tr.DiffingMenuTitle, Items: menuItems})
}
