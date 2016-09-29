# Change History

This documents the history of significant changes to `rivescript-go`.

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
