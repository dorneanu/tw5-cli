package tiddlywiki

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestTiddler(t *testing.T) {
	want := "tiddler 1"
	got := NewTiddler("tiddler 1")
	if got.Title != want {
		t.Errorf("got %s want %s", got.Title, want)
	}
}

func TestTiddler_fields(t *testing.T) {
	field := Field{
		Name:  "field_name",
		Value: "field_value",
	}
	t1 := NewTiddler("t1")
	t2 := NewTiddler("t2")

	t.Run("check for empty fields", func(t *testing.T) {
		t1.AddField("", "")
		fields := t1.Fields
		if len(fields) != 0 {
			t.Errorf("len(t1.Fields) should be 0")
		}

	})

	t.Run("check if fields equal", func(t *testing.T) {
		t2.AddField("field_name", "field_value")
		fields := t2.Fields

		// Check number of fields
		if len(fields) != 1 {
			t.Errorf("len(t2.Fields) should be 1")
		}

		// Check if same field
		if !reflect.DeepEqual(fields[0], field) {
			t.Errorf("got %v want %v", fields[0], field)
		}
	})
}

func TestTiddlywikiGetAll(t *testing.T) {
	// Sample json
	jsonResponse := `
	[
		{
			"created": "20191229203445271",
			"type": "text/vnd.tiddlywiki",
			"title": "Tiddler 1",
			"tags": "Tag1 Tag2",
			"modified": "20200104204233898",
			"revision": 0
		  },
		  {
			"title": "Tiddler 2",
			"created": "20180709135924116",
			"modified": "20191128165532546",
			"modifier": "boru",
			"tags": "Tag1 Tag2 Tag3",
			"type": "text/vnd.tiddlywiki",
			"revision": 0
		  }
	]
	`

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, jsonResponse)

	}))
	defer ts.Close()

	// Create new TW instance
	tw := NewTW(ts.URL)

	t.Run("check get all tiddlers", func(t *testing.T) {
		tiddlers, err := tw.GetAll()
		if err != nil {
			t.Errorf("Couldn't get all tiddlers: %s", err)
		}

		if len(tiddlers) != 2 {
			t.Errorf("len of tiddlers is %d", len(tiddlers))
		}
	})
}

func TestTiddlywikiGet(t *testing.T) {
	// Sample json
	jsonResponse := `
	{
		"created": "20191229203445271",
		"type": "text/vnd.tiddlywiki",
		"title": "Tiddler 1",
		"tags": "Tag1 Tag2",
		"modified": "20200104204233898",
		"revision": 0
	}
	`

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, jsonResponse)

	}))
	defer ts.Close()

	// Create new TW instance
	tw := NewTW(ts.URL)

	t.Run("check single tiddlers", func(t *testing.T) {
		tiddler, err := tw.Get("Tiddler 1")
		want := "Tiddler 1"

		if err != nil {
			t.Errorf("Couldn't get tiddler: %s", err)
		}

		if tiddler.Title != want {
			t.Errorf("got %s want %s", tiddler.Title, want)
		}
	})
}
