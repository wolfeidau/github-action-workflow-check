package rules

import (
	"context"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/rhysd/actionlint"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/action-workflow-check/internal/ptr"
)

// RuleAction is a rule to check running action in steps of jobs.
// https://docs.github.com/en/actions/learn-github-actions/workflow-syntax-for-github-actions#jobsjob_idstepsuses
type RuleAction struct {
	actionlint.RuleBase
	client *github.Client
}

// NewRuleAction creates new RuleAction instance.
func NewRuleAction(client *github.Client) *RuleAction {
	return &RuleAction{
		RuleBase: actionlint.NewRuleBase("action", "Checks the version of actions released on GitHub"),
		client:   client,
	}
}

// VisitStep is callback when visiting Step node.
func (rule *RuleAction) VisitStep(n *actionlint.Step) error {
	e, ok := n.Exec.(*actionlint.ExecAction)
	if !ok || e.Uses == nil {
		return nil
	}

	if e.Uses.ContainsExpression() {
		// Cannot parse specification made with interpolation. Give up
		return nil
	}

	spec := e.Uses.Value

	if strings.HasPrefix(spec, "./") {
		// Relative to repository root
		return nil
	}

	if strings.HasPrefix(spec, "docker://") {
		return nil
	}

	rule.checkRepoAction(spec, e)

	return nil
}

// Parse {owner}/{repo}@{ref} or {owner}/{repo}/{path}@{ref}
func (rule *RuleAction) checkRepoAction(spec string, exec *actionlint.ExecAction) {
	s := spec
	idx := strings.IndexRune(s, '@')
	if idx == -1 {
		rule.invalidActionFormat(exec.Uses.Pos, spec, "ref is missing")
		return
	}
	ref := s[idx+1:]
	s = s[:idx] // remove {ref}

	idx = strings.IndexRune(s, '/')
	if idx == -1 {
		rule.invalidActionFormat(exec.Uses.Pos, spec, "owner is missing")
		return
	}

	owner := s[:idx]
	s = s[idx+1:] // eat {owner}

	repo := s
	if idx := strings.IndexRune(s, '/'); idx >= 0 {
		repo = s[:idx]
	}

	if owner == "" || repo == "" || ref == "" {
		rule.invalidActionFormat(exec.Uses.Pos, spec, "owner and repo and ref should not be empty")
	}

	rule.Debug("This action skips to check inputs: %s", spec)

	release, _, err := rule.client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get release")
	}

	if release.TagName != nil {
		tag, _, err := rule.client.Git.GetRef(context.Background(), owner, repo, "tags/"+*release.TagName)
		if err != nil {
			log.Error().Err(err).Msg("failed to get release")
		}

		rule.Debug("Ref is %s latest is %s/%s", ref, ptr.ToString(release.TagName), ptr.ToString(tag.Object.SHA))

		if ref != ptr.ToString(tag.Object.SHA) {
			rule.Errorf(exec.Uses.Pos, "update release to latest\n\t%s/%s@%s # %s", owner, repo, ptr.ToString(tag.Object.SHA), ptr.ToString(release.TagName))
		}
	}

}

func (rule *RuleAction) invalidActionFormat(pos *actionlint.Pos, spec string, why string) {
	rule.Errorf(pos, "specifying action %q in invalid format because %s. available formats are \"{owner}/{repo}@{ref}\" or \"{owner}/{repo}/{path}@{ref}\"", spec, why)
}
