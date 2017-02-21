# Redis Sessions for RiveScript

[![GoDoc](https://godoc.org/github.com/aichaos/rivescript-go/sessions/redis?status.svg)](https://godoc.org/github.com/aichaos/rivescript-go/sessions/redis)

This package provides support for using a [Redis cache](https://redis.io/) to
store user variables for RiveScript.

```bash
go get github.com/aichaos/rivescript-go/sessions/redis
```

## Quick Start

```go
package main

import (
    "fmt"

    rivescript "github.com/aichaos/rivescript-go"
    "github.com/aichaos/rivescript-go/sessions/redis"
    goRedis "gopkg.in/redis.v5"
)

func main() {
    // Verbose example with ALL options spelled out. All the settings are
    // optional, and their default values are shown here.
    bot := rivescript.New(&rivescript.Config{
        // Initialize the Redis session manager here.
        SessionManager: redis.New(&redis.Config{
            // The prefix is added before all the usernames in the Redis cache.
            // For a username of 'alice' it would go into 'rivescript/alice'
            Prefix: "rivescript/",

            // The prefix used to store 'frozen' copies of user variables. The
            // default takes the form "frozen:<prefix>" using your Prefix,
            // so this field is doubly optional unless you wanna customize it.
            FrozenPrefix: "frozen:rivescript/",

            // If you need to configure the underlying Redis instance, you can
            // pass its options along here.
            Redis: &goRedis.Options{
                Addr: "localhost:6379",
                DB:   0,
            },
        }),
    })

    // A minimal version of the above that uses all the default options.
    bot = rivescript.New(&rivescript.Config{
        SessionManager: redis.New(nil),
    })

    bot.LoadDirectory("eg/brain")
    bot.SortReplies()

    // And go on as normal.
    reply, err := bot.Reply("soandso", "hello bot")
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    } else {
        fmt.Printf("Reply: %s\n", reply)
    }
}
```

## Testing

Running these unit tests requires a local Redis server to be running. In the
future I'll look into mocking the server.

## License

Released under the same terms as RiveScript itself (MIT license).
