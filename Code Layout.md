# Code Layout

This project has a handful of `.go` files and it might not always be clear where
to look for each function call.

Here is a walkthrough of the code layout in a more logical fashion.

## RiveScript Module

### Constructor and Version Methods

* `rivescript.go`
  * Defines class `RiveScript` and `Version()` method.
  * Internally, also the debug methods like `say()` and `warn()`

### Loading and Parsing Methods

* `loading.go`
  * Public methods:
    * `LoadFile()`
    * `LoadDirectory()`
    * `Stream()`
    * `SortReplies()`
  * Internal methods:
    * `parse()`
* `parser.go`
  * The bulk of the RiveScript source parsing logic is here.
  * Private methods:
    * `parseSource()`
* `ast.go`
  * Defines private structures for representing the "abstract syntax tree"

### Reply Sorting

* `sorting.go`
  * Private methods:
    * `sortTriggerSet()`
    * `sortList()`
    * `sortByWords()`
    * `sortByLength()`
    * `initSortTrack()`

### Public Configuration Methods

* `config.go`
  * Public methods:
    * `SetHandler()` / `RemoveHandler()`
    * `SetSubroutine()` / `RemoveSubroutine()`
    * `SetGlobal()`
    * `SetVariable()` / `GetVariable()`
    * `SetSubstitution()`
    * `SetPerson()`
    * `SetUservar()` / `SetUservars()`
    * `GetUserver()` / `GetUservars()` / `GetAllUservars()`
    * `ClearUservars()` / `ClearAllUservars()`
    * `FreezeUservars()` / `ThawUservars()`
    * `LastMatch()`
    * `CurrentUser()`

### Reply Fetching Methods

* `brain.go`
  * Public methods:
    * `Reply()`
  * Private methods:
    * `getReply()`
* `tags.go`
  * Code for processing tags on both user messages and responses.
  * Private methods:
    * `formatMessage()`
    * `triggerRegexp()`
    * `processTags()`
    * `substitute()`

### Macro Handler Interface

* `macros.go`
  - Defines `MacroInterface`

### Debug Methods

* `debug.go`
  * Public methods:
    * `DumpTopics()`
    * `DumpSorted()`
  * Private methods:
    * `_dumpSorted()`
    * `_dumpSortedList()`

### Documentation

* `doc.go` - Package documentation.
* `doc_test.go` - Examples.

### Inheritance Utility Functions

* `inheritance.go`
  * Private methods:
    * `getTopicTriggers()`
    * `_getTopicTriggers()`
    * `getTopicTree()`

### Miscellaneous

* `regexp.go` - Central location for common regular expressions.
* `structs.go` - Central location of miscellaneous structs.
* `utils.go` - Miscellaneous utility functions.

## Object Macro Language Handlers

* `lang/rivescript_js` - JavaScript support.

## Stand-alone RiveScript Interpreter

* `cmd/rivescript`
