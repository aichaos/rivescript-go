# RiveScript-Go

## Introduction

This is a RiveScript interpreter library written for the Go programming
language. RiveScript is a scripting language for chatterbots, making it easy
to write trigger/response pairs for building up a bot's intelligence.

**This project is still under heavy development and is not ready yet!**

## Development Roadmap

Rough estimation of the current progress on RiveScript-Go:

- [x] `LoadFile`, `LoadDirectory` and parsing RiveScript document (AST)
- [x] Reading parsed data into bot's main memory
- [x] `SortReplies()` to sort the data in the bot's memory
  - [x] Topic inheritence/includes
- [ ] `Reply()` to fetch a reply.
  - [ ] Response tag processing
- [ ] Make a standalone "shell" program for quickly testing the bot
- [ ] Unit tests to make sure it has feature parity with other implementations

## Installation

`go get github.com/aichaos/rivescript-go`

## Usage

```go
package main

import (
    "fmt"
    rivescript "github.com/aichaos/rivescript-go"
)

func main() {
    bot := rivescript.New()

    // Load a directory full of RiveScript documents (.rive files)
    bot.LoadDirectory("eg/brain")

    // Load an individual file.
    bot.LoadFile("brain/testsuite.rive")

    // Sort the replies after loading them!
    bot.SortReplies()

    // Get a reply.
    reply := bot.Reply("local-user", "Hello, bot!")
    fmt.Printf("The bot says: %s", reply)
}
```

## License

```
The MIT License (MIT)

Copyright (c) 2015 Noah Petherbridge

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## See Also

The official RiveScript website, http://www.rivescript.com/
