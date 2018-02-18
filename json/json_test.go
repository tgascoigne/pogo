package json_test

import (
	"os"
	"testing"

	"github.com/tgascoigne/pogo/json"
)

var testFiles = []string{
	"test_data/glossary.json",
	"test_data/menu.json",
	"test_data/webapp.json",
}

func TestJson(t *testing.T) {
	doTest := func(path string) {
		t.Run(path, func(t *testing.T) {
			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			// Not a proper test; good enough for an example!
			json.Parse(file)
		})
	}

	for _, f := range testFiles {
		doTest(f)
	}
}
