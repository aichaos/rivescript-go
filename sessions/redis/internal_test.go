package redis

// This tests the internal interface of just the Redis specific bits,
// independent of RiveScript. At time of writing, this test file gets
// us to 93.4% coverage!

import (
	"fmt"
	"os"
	"testing"

	"github.com/aichaos/rivescript-go/sessions"
)

// newTest creates a new test environment with a new Redis prefix.
// The generated prefix takes the form: `redis_test:<PID>:<name>`
func newTest(name string) *Session {
	s := New(&Config{
		Prefix: fmt.Sprintf("rivescript:%d:%s/", os.Getpid(), name),
	})
	s.ClearAll()
	return s
}

// tearDown deletes a Redis prefix after the test is done.
func tearDown(s *Session) {
	s.ClearAll()
}

func TestRedis(t *testing.T) {
	s := newTest("main")
	defer tearDown(s)

	// There should be no user data yet.
	s.expectCount(t, 0)

	// Create the first user.
	{
		username := "alice"

		s.Set(username, map[string]string{
			"name": "Alice",
		})

		// Sanity check that the default topic was implied with that.
		s.checkVariable(t, username, "topic", "random", false)

		// Check the variable we just set, and one we didn't.
		s.checkVariable(t, username, "name", "Alice", false)
		s.checkVariable(t, username, "age", "5", true)

		// See if we have as many variables as we expect.
		vars, _ := s.GetAny(username)
		if len(vars.Variables) != 2 {
			t.Errorf(
				"expected to have 2 variables, but had %d: %v",
				len(vars.Variables),
				vars.Variables,
			)
		}

		// She should have an empty history.
		history, _ := s.GetHistory(username)
		for i := 0; i < sessions.HistorySize; i++ {
			if history.Input[i] != "undefined" {
				t.Errorf(
					"expected to have a blank history, but input[%d] = %s",
					i,
					history.Input[i],
				)
			}
			if history.Reply[i] != "undefined" {
				t.Errorf(
					"expected to have a blank history, but reply[%d] = %s",
					i,
					history.Reply[i],
				)
			}
		}

		// Add some history.
		s.AddHistory(username, "hello bot", "hello human")
		history, _ = s.GetHistory(username)
		if history.Input[0] != "hello bot" {
			t.Errorf(
				"got unexpected input history: expected 'hello bot', got %s",
				history.Input[0],
			)
		}
		if history.Reply[0] != "hello human" {
			t.Errorf(
				"got unexpected reply history: expected 'hello human', got %s",
				history.Reply[0],
			)
		}

		// LastMatch.
		lastMatch, _ := s.GetLastMatch(username)
		if lastMatch != "" {
			t.Errorf(
				"didn't expect to have a LastMatch, but had: %s",
				lastMatch,
			)
		}
		s.SetLastMatch(username, "hello bot")
		lastMatch, _ = s.GetLastMatch(username)
		if lastMatch != "hello bot" {
			t.Errorf(
				"LastMatch wasn't '%s' like I expected, but was: %s",
				"hello bot",
				lastMatch,
			)
		}

		// Verify we only have one user so far.
		s.expectCount(t, 1)
	}

	// Create the second user.
	{
		username := "bob"

		s.Init(username)

		// Verify we now have two users.
		s.expectCount(t, 2)

		// Delete this user.
		s.Clear(username)

		// We should be back to one.
		s.expectCount(t, 1)
	}

	// Create the new second user.
	{
		username := "barry"

		// Set some variables.
		s.Set(username, map[string]string{
			"name":   "Barry",
			"age":    "20",
			"gender": "male",
		})

		// Freeze his variables.
		s.Freeze(username)

		// Happy birthday!
		birthday := map[string]string{
			"age": "21",
		}
		s.Set(username, birthday)

		// Thaw the variables and make sure it was restored.
		s.Thaw(username, sessions.Thaw)
		s.checkVariable(t, username, "age", "20", false)

		// Make sure trying to thaw again gives an error because the frozen
		// copy isn't there.
		err := s.Thaw(username, sessions.Thaw)
		expectError(t, "thawing again after sessions.Thaw", err)

		// Freeze it again and repeat.
		s.Freeze(username)
		s.Set(username, birthday)

		// Thaw with the 'keep' option this time.
		s.Thaw(username, sessions.Keep)
		s.checkVariable(t, username, "age", "20", false)

		// One more time. The frozen copy is still there so just update
		// the user var and try the last thaw option.
		s.Set(username, birthday)

		// Discard should just delete the frozen copy and not restore it.
		s.Thaw(username, sessions.Discard)
		s.checkVariable(t, username, "age", "20", false)

		// One last call to Thaw should error out now.
		err = s.Thaw(username, sessions.Thaw)
		expectError(t, "thawing again after discard", err)
	}

	// Create the third and fourth users.
	{
		s.Init("charlie")
		s.Init("dave")
		s.expectCount(t, 4)

		// Clear all data and expect to have nothing left.
		s.ClearAll()
		s.expectCount(t, 0)
	}

	// Test all the error cases.
	{
		var err error

		_, err = s.Get("nobody", "name")
		expectError(t, "get variable from missing user", err)

		_, err = s.GetAny("nobody")
		expectError(t, "get any variables for missing user", err)

		_, err = s.GetLastMatch("nobody")
		expectError(t, "get a LastMatch for missing user", err)

		_, err = s.GetHistory("nobody")
		expectError(t, "get history for missing user", err)

		err = s.Freeze("nobody")
		expectError(t, "freeze missing user", err)

		s.Init("nobody")
		s.Freeze("nobody")
		err = s.Thaw("nobody", 42)
		expectError(t, "invalid thaw action", err)
	}
}

// checkVariable handles tests on user variables.
func (s *Session) checkVariable(t *testing.T, username, name, expected string, expectError bool) {
	value, err := s.Get(username, name)

	// Got an error when we aren't expecting one?
	if err != nil {
		if !expectError {
			t.Errorf(
				"got an unexpected error when getting variable '%s' for %s: %s",
				name,
				username,
				err,
			)
		}
		return
	}

	// Didn't get an error when we expected to?
	if err == nil {
		if expectError {
			t.Errorf(
				"was expecting an error when getting variable '%s' for %s, but did not get one",
				name,
				username,
			)
		}
		return
	}

	// Was it what we expected?
	if value != expected {
		t.Errorf(
			"got unexpected user variable '%s' for %s:\n"+
				"expected: %s\n"+
				"     got: %s",
			name,
			username,
			expected,
			value,
		)
	}
}

// expectCount expects the user count from GetAll() to be a certain number.
func (s *Session) expectCount(t *testing.T, expect int) {
	users := s.GetAll()
	if len(users) != expect {
		t.Errorf(
			"expected to have %d users, but had: %d",
			expect,
			len(users),
		)
	}
}

func expectError(t *testing.T, name string, err error) {
	if err == nil {
		t.Errorf(
			"expected to get an error from '%s', but did not get one",
			name,
		)
	}
}
