package gitlab

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"
)

// XXX: move this func out if necessary
// content read a file, return its content as string
func content(f string) string {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c := string(b)
	return c
}

// TODO: better test posting empty string
func TestValidate(t *testing.T) {
	mux, server, client := setup()
	defer teardown(server)

	var res string

	valid := `{
			"status": "valid",
			"errors": []
		}`

	invalid := `{
			"status": "invalid",
			"errors": [
				"variables config should be a hash of key value pairs"
			]
		}`

	mux.HandleFunc("/ci/lint", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintf(w, res)
	})

	testFunc := func(t *testing.T, got, want *LintResult) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Validate returned \ngot:\n%v\nwant:\n%v", Stringify(got), Stringify(want))
		}
	}

	t.Run("valid yaml", func(t *testing.T) {
		f := "testdata/validate_valid.yml"
		c := content(f)
		res = valid

		got, _, err := client.Validate.Lint(c)

		if err != nil {
			t.Errorf("Validate returned error: %v", err)
		}

		e := make([]string, 0)

		want := &LintResult{
			Status: "valid",
			Errors: e,
		}

		testFunc(t, got, want)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		f := "testdata/validate_invalid.yml"
		c := content(f)
		res = invalid

		got, _, err := client.Validate.Lint(c)

		if err != nil {
			t.Errorf("Validate returned error: %v", err)
		}

		e := make([]string, 1)
		e[0] = "variables config should be a hash of key value pairs"

		want := &LintResult{
			Status: "invalid",
			Errors: e,
		}

		testFunc(t, got, want)
	})
}
