package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aichaos/rivescript-go"
	"github.com/aichaos/rivescript-go/lang/javascript"
)

// Bot is a global RiveScript instance to share between requests, so that the
// bot's replies only need to be parsed and sorted one time.
var Bot *rivescript.RiveScript
var forgetful bool

func main() {
	// Command line arguments.
	var (
		port  = flag.Int("port", 8000, "Port to listen on (default 8000)")
		host  = flag.String("host", "0.0.0.0", "Interface to listen on.")
		debug = flag.Bool("debug", false, "Enable debug mode for RiveScript.")
		utf8  = flag.Bool("utf8", true, "Enable UTF-8 mode")
	)
	flag.BoolVar(&forgetful, "forgetful", false,
		"Do not store user variables in server memory between requests.",
	)
	flag.Parse()

	// Set up the RiveScript bot.
	Bot = rivescript.New(&rivescript.Config{
		Debug: *debug,
		UTF8:  *utf8,
	})
	Bot.SetHandler("javascript", javascript.New(Bot))
	Bot.LoadDirectory("../brain")
	Bot.SortReplies()

	http.HandleFunc("/", LogMiddleware(IndexHandler))
	http.HandleFunc("/reply", LogMiddleware(ReplyHandler))

	addr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Server listening at http://%s/\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Request describes the JSON arguments to the API.
type Request struct {
	Username string            `json:"username"`
	Message  string            `json:"message"`
	Vars     map[string]string `json:"vars"`
}

// Response describes the JSON output from the API.
type Response struct {
	Status string            `json:"status"` // 'ok' or 'error'
	Error  string            `json:"error,omitempty"`
	Reply  string            `json:"reply,omitempty"`
	Vars   map[string]string `json:"vars,omitempty"`
}

// ReplyHandler is the JSON endpoint for the RiveScript bot.
func ReplyHandler(w http.ResponseWriter, r *http.Request) {
	// Only POST allowed.
	if r.Method != "POST" {
		writeError(w, "This endpoint only works with POST requests.", http.StatusMethodNotAllowed)
		return
	}

	// Get the request information.
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		writeError(w, "Content-Type of the request should be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Get JSON parameters.
	var params Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// The username is required.
	if params.Username == "" {
		writeError(w, "username is required", http.StatusBadRequest)
		return
	}

	// Let RiveScript know all the user vars of the client.
	for k, v := range params.Vars {
		Bot.SetUservar(params.Username, k, v)
	}

	// Get a reply from the bot.
	reply, err := Bot.Reply(params.Username, params.Message)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve all user variables from the bot.
	var vars map[string]string
	userdata, err := Bot.GetUservars(params.Username)
	if err == nil {
		vars = userdata.Variables
	}

	// Are we being forgetful?
	if forgetful {
		Bot.ClearUservars(params.Username)
	}

	// Prepare the JSON response.
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	response := Response{
		Status: "ok",
		Error:  "",
		Reply:  reply,
		Vars:   vars,
	}

	out, _ := json.MarshalIndent(response, "", "  ")
	w.Write(out)
}

// IndexHandler is the default page handler and just shows a `curl` example.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "index.html")
}

// LogMiddleware does basic logging to the console for HTTP requests.
func LogMiddleware(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log line looks like:
		// [127.0.0.1] POST /reply HTTP/1.1
		log.Printf("[%s] %s %s %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			r.Proto,
		)
		fn(w, r)
	}
}

// writeError handles sending JSON errors to the client.
func writeError(w http.ResponseWriter, message string, code int) {
	// Prepare the error JSON.
	response, err := json.MarshalIndent(Response{
		Status: "error",
		Error:  message,
	}, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send it.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("[ERROR] %s\n", err)
	}
}
