package gitlab

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"
)

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
	validContent := content("testdata/validate_valid.yml")
	invalidContent := content("testdata/validate_invalid.yml")

	validRes := `{
			"status": "valid",
			"errors": []
		}`

	invalidRes := `{
			"status": "invalid",
			"errors": [
				"variables config should be a hash of key value pairs"
			]
		}`

	testFunc := func(t *testing.T, got, want *LintResult) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Validate returned \ngot:\n%v\nwant:\n%v", Stringify(got), Stringify(want))
		}
	}

	wantValid := &LintResult{
		Status: "valid",
		Errors: make([]string, 0),
	}

	e := make([]string, 1)
	e[0] = "variables config should be a hash of key value pairs"

	wantInvalid := &LintResult{
		Status: "invalid",
		Errors: e,
	}

	testCases := []struct {
		desc     string
		contents string
		res      string
		want     *LintResult
	}{
		{"valid case", validContent, validRes, wantValid},
		{"invalid case", invalidContent, invalidRes, wantInvalid},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			mux, server, client := setup()
			defer teardown(server)

			mux.HandleFunc("/ci/lint", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				fmt.Fprintf(w, tc.res)
			})

			got, _, err := client.Validate.Lint(tc.contents)

			if err != nil {
				t.Errorf("Validate returned error: %v", err)
			}

			want := tc.want
			testFunc(t, got, want)
		})
	}
}
