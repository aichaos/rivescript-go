# RiveScript-Go

## Introduction

This is a RiveScript interpreter library written for the Go programming
language. RiveScript is a scripting language for chatterbots, making it easy
to write trigger/response pairs for building up a bot's intelligence.

**This project is currently in Beta status.**

## Documentation

* RiveScript Library: <http://godoc.org/github.com/aichaos/rivescript-go>
* RiveScript Stand-alone Interpreter: <http://godoc.org/github.com/aichaos/rivescript-go/cmd/rivescript>
* JavaScript Object Macros: <http://godoc.org/github.com/aichaos/rivescript-go/lang/rivescript_js>

## Installation

`go get github.com/aichaos/rivescript-go`

For the stand-alone binary for testing a RiveScript bot:

`go install github.com/aichaos/rivescript-go/cmd/rivescript`

## Usage

The distribution of RiveScript includes an interactive shell for testing your
RiveScript bot. Run it with the path to a folder on disk that contains your
RiveScript documents. Example:

```bash
# (Linux)
$ rivescript eg/brain

# (Windows)
> rivescript.exe eg/brain
```

See `rivescript --help` for options it accepts, including debug mode and UTF-8
mode.

When used as a library for writing your own chatbot, the synopsis is as follows:

```go
package main

import (
    "fmt"
    rivescript "github.com/aichaos/rivescript-go"
)

func main() {
    bot := rivescript.New()

    // Load a directory full of RiveScript documents (.rive files)
    err := bot.LoadDirectory("eg/brain")
    if err != nil {
      fmt.Printf("Error loading from directory: %s", err)
    }

    // Load an individual file.
    err = bot.LoadFile("brain/testsuite.rive")
    if err != nil {
      fmt.Printf("Error loading from file: %s", err)
    }

    // Sort the replies after loading them!
    bot.SortReplies()

    // Get a reply.
    reply := bot.Reply("local-user", "Hello, bot!")
    fmt.Printf("The bot says: %s", reply)
}
```

## UTF-8 Support

UTF-8 support in RiveScript is considered an experimental feature. It is
disabled by default. Enable it by setting `RiveScript.UTF8 = true`.

By default (without UTF-8 mode on), triggers may only contain basic ASCII
characters (no foreign characters), and the user's message is stripped of all
characters except letters, numbers and spaces. This means that, for example,
you can't capture a user's e-mail address in a RiveScript reply, because of
the @ and . characters.

When UTF-8 mode is enabled, these restrictions are lifted. Triggers are only
limited to not contain certain metacharacters like the backslash, and the
user's message is only stripped of backslashes and HTML angled brackets
(to protect from obvious XSS if you use RiveScript in a web application).
Additionally, common punctuation characters are stripped out, with the default
set being `/[.,!?;:]/g`. This can be overridden by providing a new `Regexp`
object as the `RiveScript.UnicodePunctuation` attribute. Example:

```go
// Make a new bot with UTF-8 mode enabled.
bot := rivescript.New()
bot.UTF8 = true

// Override the punctuation characters that get stripped from the
// user's message.
bot.UnicodePunctuation = regexp.MustCompile(`[.,!?;:]`);
```

The `<star>` tags in RiveScript will capture the user's "raw" input, so you can
write replies to get the user's e-mail address or store foreign characters in
their name.

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
