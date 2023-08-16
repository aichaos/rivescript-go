# Change History

This documents the history of significant changes to `rivescript-go`.

## v0.4.0 - Aug 15, 2023

This update will modernize the Go port of RiveScript bringing some of the
newer features that were available on the JavaScript or Python ports.

### The ?Keyword Command

This update adds support for the newer `?Keyword` command in RiveScript.

This command works around a Unicode matching bug that affected the
Go, JavaScript and Python ports of RiveScript. For example:

```rivescript
// You wanted this trigger to match if "你好" appears anywhere
// in a user's message, but it wasn't working before.
+ [*] 你好 [*]
- 你好!

// Now: use the ?Keyword command in place of +Trigger
? 你好
- 你好!
```

The optional wildcard `[*]` syntax didn't work when paired with Unicode
symbols because the regular expression that `[*]` is turned into (which
involves "word boundary" or `\b` characters) only worked with ASCII latin
characters. The ?Keyword command works around this by translating into a
+Trigger that tries _every_ combination of literal word, wildcard on either
or both sides, and optional wildcard, to ensure that your keyword trigger
will indeed match your keyword _anywhere_ in the user's message.

### CaseSensitive User Message Support

By default RiveScript would always lowercase the user's message as it
comes in. If this is undesirable and you'd like to preserve their _actual_
capitalization (when it gets captured by wildcards and comes out in their
`<star>` tags), provide the new CaseSensitive boolean to your RiveScript
constructor:

```go
bot := rivescript.New(&rivescript.Config{
    CaseSensitive: true,
})
```

The built-in `rivescript` command line program adds a `-case` parameter
to enable this option for testing:

```bash
rivescript -case /path/to/brain
```

### JavaScript Object Macros

The official JavaScript object macro support library has a couple of
exciting new updates.

Firstly, the JavaScript engine has been replaced from
[github.com/robertkrimen/otto](https://github.com/robertkrimen/otto) with
[github.com/dop251/goja](https://github.com/dop251/goja) which should
provide a better quality of life for writing JavaScript functions for
your bot. Goja supports many of the modern ES6+ features including the
let and const keyword and arrow functions.

Additionally, the Goja runtime will be exposed at the `.VM` accessor
for the JavaScriptHandler class, so you can directly play with it and
set global variables or Go functions that can be called from your
JavaScript macros in your bot, allowing for much greater extensibility.

### Other Changes

* Add shellword parsing for `<call>` tags: you can pass "quoted strings"
  in which will go in as one 'word' in the `args` array, instead of the
  arguments being literally split by space characters. This brings the
  Go port of RiveScript in line with the JavaScript port which has been
  doing this for a while.
* Fix the rsts_test.go to properly load from the RSTS (RiveScript Test
  Suite) git submodule.

## v0.3.1 - Aug 20, 2021

This release simply adds a `go.mod` to this project so that it gets along well
with modern Go projects.

## v0.3.0 - Apr 30, 2017

This update brings some long-needed restructuring to the source layout of
`rivescript-go`. Briefly: it moves all source files from the `src/` subpackage
into the root package namespace, and removes the wrapper shim functions (their
documentation was then moved to the actual implementation functions).

### API Breaking Changes

* The `github.com/aichaos/rivescript-go/src` subpackage has been removed, and
  all of the things you used to import from there can now be found in the root
  package instead. Most notably, the `RiveScript` struct needed to be imported
  (again) from the `src` subpackage for use with Go object macros.

  To update source code where you used Go object macros:

  ```diff
    import (
        "github.com/aichaos/rivescript-go"
  -     rss "github.com/aichaos/rivescript-go/src"
    )

    func main() {
        bot = rivescript.New(nil)

  -     subroutine := func(rs *rss.RiveScript, args []string) string {
  +     subroutine := func(rs *rivescript.RiveScript, args []string) string {
            return "Hello world"
        }

        bot.SetSubroutine("hello", subroutine)
    }
  ```
* `rivescript.Version` is now a string constant (replacing `VERSION`). The
  instance method `Version()` has been removed.

## Changes

* All RiveScript unit tests have been removed in favor of those from the
  [RiveScript Test Suite](https://github.com/aichaos/rsts). The test file
  `rsts_test.go` implements the Go test runner, and the `rsts` repo was added
  as a Git submodule.
* The Git commit hash is now encoded into the front-end command line client,
  printed along with the version number in the welcome banner.

## v0.2.0 - Feb 7, 2017

This update focuses on bug fixes and code reorganization.

### API Breaking Changes

* `rivescript.New()` and the `Config` struct have been refactored. `Config`
  now comes from the `rivescript` package directly rather than needing to
  import from `rivescript/config`.

  For your code, this means you can remove the `aichaos/rivescript-go/config`
  import and change the `config.Config` name to `rivescript.Config`:

  ```go
  import "github.com/aichaos/rivescript-go"

  func main() {
      // Example defining the struct to override defaults.
      bot := rivescript.New(&rivescript.Config{Debug: true})

      // For the old `config.Basic()` that provided default settings, just
      // pass in a nil Config object.
      bot = rivescript.New(nil)

      // For the old `config.UTF8()` helper function that provided a Config with
      // UTF-8 mode enabled, instead call rivescript.WithUTF8()
      bot = rivescript.New(rivescript.WithUTF8())
  }
  ```
* `Reply()`, `SortReplies()` and `CurrentUser()` now return an `error` value
  in addition to what they already returned.

### Changes

* Add ANSI colors to the RiveScript shell (`cmd/rivescript`); they can be
  disabled with the `-nocolor` command line option.
* Add new commands to the RiveScript shell:
  * `/debug [true|false]` to toggle the debug mode (`/debug` will print
    the current setting of debug mode).
  * `/dump <topics|sorted>` to print the internal data structures for the
    topics and sorted trigger sets, respectively.
* Separate the unit tests into multiple files and put them in the `rivescript`
  package instead of `rivescript_test`; this enables test code coverage
  reporting (we're at 72.1% coverage!)
* Handle module configuration at the root package instead of in the `src`
  package. This enabled getting rid of the `rivescript/config` package and
  making the public API more sane.
* Code cleanup via `go vet`
* Add more documentation and examples to the Go doc.
* Fix `@Redirects` not working sometimes when tags like `<bot>` insert capital
  letters (bug #1)
* Fix an incorrect regexp that makes wildcards inside of optionals, like `[_]`,
  not matchable in `<star>` tags. For example, with `+ my favorite [_] is *`
  and a message of "my favorite color is red", `<star1>` would be "red" because
  the optional makes its wildcard non-capturing (bug #15)
* Fix the `<star>` tag handling to support star numbers greater than `<star9>`:
  you can use as many star numbers as will be captured by your trigger (bug #16)
* Fix a probable bug within inheritance/includes: some parts of the code were
  looking in one location for them, another in the other, so they probably
  didn't work perfectly before.
* Fix `RemoveHandler()` to make it remove all known object macros that used that
  handler, which protects against a possible null pointer exception.
* Fix `LoadDirectory()` to return an error when doesn't find any RiveScript
  source files to load, which helps protect against the common error that you
  gave it the wrong directory.
* New unit tests: object macros.
* An internal optimization that allowed for cleaning up a redundant storage
  location for triggers that have `%Previous` commands (PR #20)

## v0.1.0 - Dec 11, 2016

This update changes some function prototypes in the API which breaks backward
compatibility with existing code.

* **API Breaking Changes:**
  * `rivescript.New()` now takes a `*config.Config` struct to configure the
    instance. This is now the preferred way to configure debug mode, strict
    mode, UTF-8, etc. rather than functions like `SetUTF8()`.

    For RiveScript's default settings, you can do `rivescript.New(config.Basic())`
    or `rivescript.New(nil)`. For UTF-8 support, `rivescript.New(config.UTF8())`
    is a convenient config template to use.
  * `GetDepth()` and `SetDepth()` now use a `uint` instead of an `int`. But
    these functions are deprecated anyway.
  * `GetUservars()` and `GetAllUservars()` return `*sessions.UserData` objects
    instead of `map[string]string` for the user data.
  * `ThawUservars()` now takes a `sessions.ThawAction` instead of a string to
    specify the action. Valid values are `Thaw`, `Discard`, or `Keep`
    (constants from the `sessions` package).
* **Deprecated Functions:**
  * Configuration functions (getters and setters). Use the `Config` struct
    when calling `rivescript.New(*config.Config)` instead:
    * `SetDebug()`, `SetUTF8()`, `SetDepth()`, `SetStrict()`
    * `GetDebug()`, `GetUTF8()`, `GetDepth()`, `GetStrict()`
* **Changes:**
  * Add support for pluggable session stores for user variables. The default
    one still keeps user variables in memory, but you can specify your own
    implementation instead.

    The interface for a `SessionManager` is in the `sessions` package. The
    default in-memory manager is in `sessions/memory`. By implementing your own
    session manager, you can change where RiveScript keeps track of user
    variables, e.g. to put them in a database or cache.
  * Make the library thread-safe with regards to getting/setting user variables
    while answering a message. The default in-memory session manager implements
    a mutex for accessing user variables.

## v0.0.3 - Sept 28, 2016

This update was all about restructuring the internal source code to make certain
internal modules exportable for third party users (e.g. the parser) and to
reduce clutter from the root of the git repo.

* Massive restructuring of the internal source code:
  * Tidied up the root of the git repo by moving *all* of the implementation
    code into the `src/` subdirectory and making the root RiveScript module a
    very lightweight API wrapper around the code. *Note: do not import the
    src package directly. Only go through the public API at the root module.*
  * Added public facing Parser submodule:
    [rivescript-go/parser](https://github.com/aichaos/rivescript-go/tree/master/parser).
    It enables third party developers to write applications that simply parse
    RiveScript code and getting an abstract syntax tree from it.
  * Moved exported object macro helpers to
    [rivescript-go/macro](https://github.com/aichaos/rivescript-go/tree/master/macro).
