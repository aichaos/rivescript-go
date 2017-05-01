# Contributing

Interested in contributing to RiveScript? Great!

First, check the general guidelines for RiveScript and its primary implementations
found at <https://www.rivescript.com/contributing> - in particular, understand
the goals and scope of the RiveScript language and the style guide for the
Go implementation (briefly: use `gofmt`).

## Quick Start

Fork, then clone the git repo:

```bash
$ git clone git@github.com:your-username/rivescript-go
```

If you are an experienced Go developer, you can clone the repo into your
standard `$GOPATH`. If you are new to Go or don't want to deal with the
`$GOPATH`, you can use the commands in the Makefile; these create a "private"
Go path inside the repo folder so you can simply clone the repo and get up and
running in no time.

After cloning, run these Make commands to get your dev environment set up:

```bash
$ make setup
$ make build
```

See the README.md for more Make commands.

## Submitting Code Changes

Run `make fmt` or `gofmt` to clean up your source code before submitting a
pull request. Also verify that `make test` works and that all the unit tests
pass.

Push to your fork and [submit a pull request](https://github.com/aichaos/rivescript-go/compare/).

At this point you're waiting on me. I'm usually pretty quick to comment on pull
requests (within a few days) and I may suggest some changes, improvements or
alternatives.

Some things that will increase the chance that your pull request is accepted:

* Follow the style guide at <https://www.rivescript.com/contributing>
* Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).
