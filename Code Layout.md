# Code Layout

This project has a handful of `.go` files and it might not always be clear where
to look for each function call.

Here is a walkthrough of the code layout in a more logical fashion.

## RiveScript Module

| File Name        | Purpose and Methods                                                  |
|------------------|----------------------------------------------------------------------|
| `astmap.go`      | Private aliases for `rivescript/ast` structs.                        |
| `brain.go`       | `Reply()` and its implementation.                                    |
| `config.go`      | Config struct and public config methods (e.g. `SetUservar()`).       |
| `debug.go`       | Debugging functions.                                                 |
| `deprecated.go`  | Deprecated methods are moved to this file.                           |
| `doc.go`         | Main module documentation for Go Doc.                                |
| `errors.go`      | Error types used by the RiveScript module.                           |
| `inheritance.go` | Functions related to topic inheritance.                              |
| `loading.go`     | File loading functions (`LoadFile()`, `LoadDirectory()`, `Stream()`) |
| `parser.go`      | Internal implementation of `rivescript/parser`                       |
| `regexp.go`      | Definitions for commonly used regular expressions.                   |
| `rivescript.go`  | `RiveScript` definition, constructor, and `Version()` methods.       |
| `sorting.go`     | `SortReplies()` and its implementation.                              |
| `tags.go`        | Tag processing functions.                                            |
| `utils.go`       | Misc utility functions.                                              |

## Test Files

| File Name       | Purpose                                    |
|-----------------|--------------------------------------------|
| `doc_test.go`   | Example snippets.                          |
| `macro_test.go` | Tests external object macros (JavaScript). |
| `rsts_test.go`  | The RiveScript Test Suite.                 |
