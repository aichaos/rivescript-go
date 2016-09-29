# Change History

This documents the history of significant changes to `rivescript-go`.

## v0.0.3 - Sept 28, 2016

This update was all about restructuring the internal source code to make certain
internal modules exportable for third party users (e.g. the parser) and to
reduce clutter from the root of the git repo.

* **Important Breaking Change:** The package name of the base RiveScript
  module has been updated. Fixing this is easy:

```diff
- import rivescript "github.com/aichaos/rivescript-go"
+ import "github.com/aichaos/rivescript-go/rivescript"
```

* Massive restructuring of the internal source code:
  * Moved all source code *out* of the root of the git repo to reduce clutter.
    All the core library source code is now under the `rivescript/` subpackage.
  * Added public facing Parser submodule:
    [rivescript-go/parser](https://github.com/aichaos/rivescript-go/tree/master/parser).
    It enables third party developers to write applications that simply parse
    RiveScript code and getting an abstract syntax tree from it.
  * Moved exported object macro helpers to
    [rivescript-go/macro](https://github.com/aichaos/rivescript-go/tree/master/macro).
