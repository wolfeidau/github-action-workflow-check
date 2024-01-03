# action-workflow-check

This is a CLI which checks if the version of an action used in your workflow is up to date, and output changes required. This is to assist with the recommendations relating to pinning third party actions using their git sha in the [Security hardening for GitHub Actions](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions#using-third-party-actions).

This project builds on the [github.com/rhysd/actionlint](https://github.com/rhysd/actionlint) project, using it as a library.

# Usage

```
Usage: action-workflow-check <command>

Flags:
  -h, --help       Show context-sensitive help.
      --debug      Enable debug logging
      --version

Commands:
  scan     Scan the project for GitHub Actions
  login    Login to GitHub to avoid rate limiting
```

**Note:** Given the rate limits for the GitHub api are low, you will probably need to login to GitHub otherwise using `action-workflow-check scan` more than a few times will result in rate limiting.

For example running it on a previous version of [s3iofs](https://github.com/wolfeidau/s3iofs) project.

```
.github/workflows/go.yml:25:15: update release to latest
	actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11 # v4.1.1 [action]
   |
25 |       - uses: actions/checkout@c85c95e3d7251135ab7dc9ce3241c5835cc595a9 # v3.5.3
   |               ^~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
.github/workflows/go.yml:26:15: update release to latest
	actions/setup-go@0c52d547c9bc32b1aa3301fd7a9cb496313a4491 # v5.0.0 [action]
   |
26 |       - uses: actions/setup-go@fac708d6674e30b6ba41289acaab6d4b75aa0753 # v4.0.1
   |               ^~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
.github/workflows/go.yml:34:15: update release to latest
	golangci/golangci-lint-action@3a919529898de77ec3da873e3063ca4b10e7f5cc # v3.7.0 [action]
   |
34 |         uses: golangci/golangci-lint-action@639cd343e1d3b897ff35927a75193d57cfcba299 # v3.6.0
   |               ^~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
3 lint errors found by actionlint
```

# Security

This CLI uses [github.com/zalando/go-keyring](https://github.com/zalando/go-keyring) to store credentials in the OS keychain. 

Authentication is handled by [github.com/cli/oauth](https://github.com/cli/oauth) which uses device flow and an oauth application to login to GitHub.

# License

This application is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au/?utm_source=action-workflow-check).