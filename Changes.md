# Change History

This documents the history of significant changes to `rivescript-go`.

## v0.1.1 - TBD

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

### Changes

* Separate the unit tests into multiple files and put them in the `rivescript`
  package instead of `rivescript_test`; this enables test code coverage
  reporting (we're at 72.1% coverage!)
* Handle module configuration at the root package instead of in the `src`
  package. This enabled getting rid of the `rivescript/config` package and
  making the public API more sane.
* Add more documentation and examples to the Go doc.
* Fix `@Redirects` not working sometimes when tags like `<bot>` insert capital
  letters (bug #1)
* Fix an incorrect regexp that makes wildcards inside of optionals, like `[_]`,
  not matchable in `<star>` tags. For example, with `+ my favorite [_] is *`
  and a message of "my favorite color is red", `<star1>` would be "red" because
  the optional makes its wildcard non-capturing (bug #15)
* Fix the `<star>` tag handling to support star numbers greater than `<star9>`:
  you can use as many star numbers as will be captured by your trigger (bug #16)

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
