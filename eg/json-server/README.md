# JSON Server

This example demonstrates embedding RiveScript in a Go web app, accessible via
a JSON endpoint.

## Run the Example

Run one of these in a terminal:

```bash
# Quick run
go run main.go

# Build and run
go build -o json-server main.go
./json-server [options] [path/to/brain]
```

Then you can visit the web server at <http://localhost:8000/> where you can
find an example `curl` command to run from a terminal, and an in-browser demo
that makes an ajax request to the endpoint.

From another terminal, you can use `curl` to test a JSON endpoint for the
chatbot. Or, you can use your favorite REST client.

```bash
curl -X POST -H 'Content-Type: application/json' \
  -d '{"username": "kirsle", "message": "Hello, robot"}' \
  http://localhost:8000/reply
```

### Options

The JSON server accepts the following command line options.

```
json-server [-host=string -port=int -debug -utf8 -forgetful -help] [path]
```

#### Server Options

* `-host string`

  The interface to listen on (default `"0.0.0.0"`)

* `-port int`

  The port number to bind to (default `8000`)

#### RiveScript Options

* `-debug`

  Enable debug mode within RiveScript (default `false`)

* `-utf8`

  Enable UTF-8 mode within RiveScript (default `true`)

* `-forgetful`

  Do not store user variables in server memory between requests (default
  `false`). See [User Variables](#user-variables) for more information about
  how user variables are dealt with in this program.

* `path`

  Specify a path on disk where RiveScript source files (`*.rive`) can be found.
  The default is `../brain`, or `/eg/brain` relative to the git root
  of rivescript-go.

## API Documentation

### POST /reply

Post a JSON message (`Content-Type: application/json`) to this endpoint to get
a response from the chatbot.

Request payload follows this format (all types are strings):

```javascript
{
  "username": "demo",       // Unique user ID (for user variables in the bot)
  "message": "Hello robot", // The message to send.
  "vars": {                 // Optional user variables to include.
    "name": "Demo User"
  }
}
```

The only **required** parameter is the `username`. A missing or blank `message`
would be handled by the chatbot's fall-back `*` trigger.

The response follows this format (all types are strings):

```javascript
// On successful outputs.
{
  "status": "ok",
  "reply": "Hello human.",
  "vars": {                 // All user variables the bot has for that user.
    "topic": "random",
    "name": "Demo User"
  }
}

// On errors.
{
  "status": "error",
  "error": "username is required"
}
```

The only key guaranteed to be in the response is `status`. Other keys are
excluded when empty.

## User Variables

The server keeps a shared RiveScript instance in memory for the lifetime of
the program. When the server exits, the user variables are lost.

The REST client that consumes this API *should* always send the full set of
user vars that it knows about on each request. This is the safest way to keep
consistent state for the end user. However, the client does not need to provide
these variables; the server will temporarily use its own and send its current
state to the client with each response.

A client that cares about long-term consistency of user variables should take
the `vars` returned by the server and store them somewhere, and send them back
to the server on the next request. This way the server could be rebooted
between requests and the bot won't forget the user's name, because the client
always sends its variables to the server.

To ensure that the server does not keep user variables around after the
request, you can provide the `-forgetful` command line option to the program.
This will clear the user's variables at the end of every request, forcing the
REST client to manage them on their end.

## Disclaimer for Deployment

This code is only intended for demonstration purposes, but as a Go web server
it can be used in a production environment. To that end, you should be aware of
some security and performance considerations:

* **The API is non-authenticated.** If the server is publicly accessible, then
  anybody on the Internet can interact with it, providing *any* data for the
  `username`, `message` and `vars` fields.

  If this is a problem (for example, if you're implementing some sort of User
  Access Control within RiveScript keyed off the user's username, or `<id>`),
  then you should bind the server to a non-public interface, for example by
  using the command line option: `-host localhost`

  You could put a reverse proxy like Nginx in front of the server to provide
  authentication on public interfaces if needed.

* **The server remembers users by default.** RiveScript stores user variables in
  memory by default, and this server doesn't change that behavior. This may be
  a memory leak concern if your bot interacts with large amounts of distinct
  usernames, or stores a ton of user variables per user.

  To prevent the server from holding onto user variables in memory, use the
  `-forgetful` command line option. The bot will then clear its user variables
  after every request.

  See [User Variables](#user-variables) for tips on how the client should
  then keep track of variables on its end rather than depend on the server.

## License

This example is released under the same license as rivescript-go itself.
