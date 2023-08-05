package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubCommitsHelper struct {
	c *HelperCommon

	refreshHelper *RefreshHelper
	setSubCommits func([]*models.Commit)
}

func NewSubCommitsHelper(
	c *HelperCommon,
	refreshHelper *RefreshHelper,
	setSubCommits func([]*models.Commit),
) *SubCommitsHelper {
	return &SubCommitsHelper{
		c:             c,
		refreshHelper: refreshHelper,
		setSubCommits: setSubCommits,
	}
}

func (self *SubCommitsHelper) ViewSubCommits(ref types.Ref, refForSymmetricDifference string, context types.Context, showBranchHeads bool) error {
	commits, err := self.c.Git().Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                   true,
			FilterPath:              self.c.Modes().Filtering.GetPath(),
			IncludeRebaseCommits:    false,
			RefName:                 ref.FullRefName(),
			RefToShowDivergenceFrom: refForSymmetricDifference,
		},
	)
	if err != nil {
		return err
	}

	self.setSubCommits(commits)
	self.refreshHelper.RefreshAuthors(commits)

	subCommitsContext := self.c.Contexts().SubCommits
	subCommitsContext.SetSelectedLineIdx(0)
	subCommitsContext.SetParentContext(context)
	subCommitsContext.SetWindowName(context.GetWindowName())
	subCommitsContext.SetTitleRef(ref.Description())
	subCommitsContext.SetRef(ref)
	subCommitsContext.SetLimitCommits(true)
	subCommitsContext.SetShowBranchHeads(showBranchHeads)
	subCommitsContext.ClearSearchString()
	subCommitsContext.GetView().ClearSearch()

	err = self.c.PostRefreshUpdate(self.c.Contexts().SubCommits)
	if err != nil {
		return err
	}

	return self.c.PushContext(self.c.Contexts().SubCommits)
}
