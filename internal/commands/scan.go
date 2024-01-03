package commands

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/google/go-github/v57/github"
	"github.com/rhysd/actionlint"
	"github.com/wolfeidau/action-workflow-check/internal/rules"
)

func Scan(cliCtx *kong.Context, client *github.Client, projectPath string, debug, all bool) error {
	// The function set at OnRulesCreated is called after rule instances are created. You can
	// add/remove some rules and return the modified slice. This function is called on linting
	// each workflow files.
	o := &actionlint.LinterOptions{
		OnRulesCreated: func(existingRules []actionlint.Rule) []actionlint.Rule {
			if all {
				return append(existingRules, rules.NewRuleAction(client))
			}

			return []actionlint.Rule{rules.NewRuleAction(client)}
		},
		Debug:     debug,
		LogWriter: os.Stdout,
	}

	l, err := actionlint.NewLinter(os.Stderr, o)
	if err != nil {
		return fmt.Errorf("failed to create linter: %w", err)
	}

	// First return value is an array of lint errors found in the workflow file.
	errs, err := l.LintRepository(projectPath)
	if err != nil {
		return fmt.Errorf("failed to lint repository: %w", err)
	}

	// `errs` includes errors like below:
	//
	// testdata/examples/main.yaml:14:9: every step must have its name [step-name]
	//    |
	// 14 |       - uses: actions/checkout@v3
	//    |         ^~~~~
	fmt.Println(len(errs), "lint errors found by actionlint")
	// Output: 1 lint errors found by actionlint

	if len(errs) > 0 {
		return fmt.Errorf("lint errors found by actionlint")
	}

	return nil
}
