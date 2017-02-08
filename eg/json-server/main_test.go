package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/aichaos/rivescript-go"
)

func init() {
	Bot = rivescript.New(rivescript.WithUTF8())
	Bot.Stream(`
		+ hello bot
		- Hello human.

		+ my name is *
		- <set name=<formal>>Nice to meet you, <get name>.

		+ what is my name
		- Your name is <get name>.

		+ i am # years old
		- <set age=<star>>I will remember you are <get age> years old.

		+ how old am i
		- You are <get age> years old.
	`)
	Bot.SortReplies()
}

func TestIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(IndexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("IndexHandler returned wrong status code: expected %v, got %v",
			http.StatusOK, status,
		)
	}
	_ = req
}

func TestUsernameError(t *testing.T) {
	res := post(t, Request{
		Message: "Hello bot",
	})
	assertError(t, res, "username is required")
}

func TestSimple(t *testing.T) {
	res := post(t, Request{
		Username: "alice",
		Message:  "Hello bot",
	})
	assert(t, res, "Hello human.")
}

func TestAliceVars(t *testing.T) {
	// The request sends all vars for the user. Assert that existing
	// vars were changed and default ones (topic) added.
	res := post(t, Request{
		Username: "alice",
		Message:  "my name is Alice",
		Vars: map[string]string{
			"name": "Bob",
			"age":  "10",
		},
	})
	assert(t, res, "Nice to meet you, Alice.")
	assertVars(t, res, map[string]string{
		"topic": "random",
		"name":  "Alice",
		"age":   "10",
	})

	// This request doesn't send the vars, but the server remembers the user.
	// As long as the server is running it caches user vars.
	res = post(t, Request{
		Username: "alice",
		Message:  "What is my name?",
	})
	assert(t, res, "Your name is Alice.")
	assertVars(t, res, map[string]string{
		"name":  "Alice",
		"topic": "random",
		"age":   "10",
	})
}

func TestBobVars(t *testing.T) {
	// This user will not send any vars initially, and we'll slowly build
	// them up over time.
	expect := map[string]string{}

	// Reusable function to send a message, expect a reply, and assert
	// a new variable was added.
	testReply := func(message, reply string) {
		res := post(t, Request{
			Username: "bob",
			Message:  message,
		})
		assert(t, res, reply)
		assertVars(t, res, expect)
	}

	// The first request should only set the default topic.
	expect["topic"] = "random"
	testReply("Hello bot.", "Hello human.")

	// Test default (missing) variables.
	testReply("What is my name?", "Your name is undefined.")
	testReply("How old am I?", "You are undefined years old.")

	// Now we tell it our name.
	expect["name"] = "Bob"
	testReply("My name is Bob.", "Nice to meet you, Bob.")

	// And age.
	expect["age"] = "10"
	testReply("I am 10 years old", "I will remember you are 10 years old.")
}

// post handles the common logic for POSTing to the /reply endpoint.
func post(t *testing.T, params Request) Response {
	payload, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	// Make an HTTP Request for the handler.
	req, err := http.NewRequest("POST", "/reply", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Call the handler.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ReplyHandler)
	handler.ServeHTTP(rr, req)

	// Read the response body.
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Parse it.
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatal(err)
	}

	return response
}

// assert verifies that the request was successful and the reply was given.
func assert(t *testing.T, res Response, expect string) {
	if res.Status != "ok" {
		t.Errorf("bad response status: expected 'ok', got '%v'", res.Status)
	}

	if res.Reply != expect {
		t.Errorf(
			"didn't get expected reply from bot\n"+
				"expected: '%s'\n"+
				"     got: '%s'",
			expect,
			res.Reply,
		)
	}
}

// assertVars makes sure user vars are set.
func assertVars(t *testing.T, res Response, expect map[string]string) {
	if !reflect.DeepEqual(res.Vars, expect) {
		t.Errorf(
			"user vars are not what I expected\n"+
				"expected: %v\n"+
				"     got: %v",
			expect,
			res.Vars,
		)
	}
}

// assertError verifies that an error was given.
func assertError(t *testing.T, res Response, expect string) {
	if res.Status != "error" {
		t.Errorf("bad response status: expected 'error', got '%v'", res.Status)
	}

	if res.Error != expect {
		t.Errorf(
			"didn't get expected error message\n"+
				"expected: '%s'\n"+
				"     got: '%s'",
			expect,
			res.Error,
		)
	}
}
