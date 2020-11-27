package tiddlywiki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Field is a property of a tiddler
type Field struct {
	Name  string
	Value interface{}
}

// Tiddler is the smallest piece of information in Tiddlywiki
type Tiddler struct {
	Title     string            `json:"title"`
	Created   string            `json:"created"`
	Creator   string            `json:"creator"`
	Modified  string            `json:"modified"`
	Modifier  string            `json:"modifier"`
	Tags      string            `json:"tags"`
	Type      string            `json:"type"`
	Text      string            `json:"text"`
	Fields    []Field           `json:"-"`
	RawFields map[string]string `json:"fields"`
}

// NewTiddler returns a new tiddler
func NewTiddler(title string) *Tiddler {
	return &Tiddler{
		Title: title,
	}
}

// TW describes a running Tiddlywiki instance
type TW struct {
	BaseURL    *url.URL
	httpClient *http.Client
	tiddlers   []*Tiddler
}

// AddField adds a new field to a Tiddler
func (t *Tiddler) AddField(name, value string) {
	// Don't accept empty values
	if name == "" {
		return
	}
	f := Field{Name: name, Value: value}
	t.Fields = append(t.Fields, f)
}

// JSON will marshal current tiddler to JSON format
func (t *Tiddler) JSON() string {
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}

// NewTW returns a new Tiddlywiki instance
func NewTW(host string) *TW {
	u, err := url.Parse(host)
	if err != nil {
		log.Fatal(err)
	}
	return &TW{
		BaseURL:    u,
		httpClient: &http.Client{},
	}
}

// GetAll retrieves all available tiddlers from Tiddlywiki
func (t *TW) GetAll() ([]*Tiddler, error) {
	var decoded []map[string]interface{}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/recipes/default/tiddlers.json", t.BaseURL.String()),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	// Send request
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read everything
	body, _ := ioutil.ReadAll(resp.Body)

	// Unmarshal JSON results
	err = json.Unmarshal([]byte(body), &decoded)
	if err != nil {
		fmt.Printf("Error unmarshalling: %s\n\n", err)
		return nil, err
	}

	// Convert JSON list to Tiddlers
	return t.Convert2Tiddlers(decoded), nil
}

// Convert2Tiddlers will convert a list of map[string]interface{}
// to a list of Tiddlers
func (t *TW) Convert2Tiddlers(data []map[string]interface{}) []*Tiddler {
	var tiddlers []*Tiddler

	// values contains JSON data
	for _, values := range data {

		// First create json string from map[string]interface{}
		tid := Tiddler{}
		jsonStr, err := json.Marshal(values)
		if err != nil {
			fmt.Errorf("Error on convert")
		}

		// Then unmarshall json string to object
		err = json.Unmarshal([]byte(jsonStr), &tid)
		if err != nil {
			fmt.Errorf("Error creating tiddler")
		}

		tiddlers = append(tiddlers, &tid)
	}
	return tiddlers
}

// Get gets a Tiddler from Tiddlywiki
func (t *TW) Get(name string) (*Tiddler, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/recipes/default/tiddlers/%s", t.BaseURL.String(), name),
		nil,
	)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return nil, err
	}

	// Send request
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read everything
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Convert to Tiddler
	tiddler := Tiddler{}
	err = json.Unmarshal(body, &tiddler)
	if err != nil {
		return nil, err
	}

	// Create fields from raw fields
	for k, v := range tiddler.RawFields {
		tiddler.AddField(k, v)
	}
	// Free space
	tiddler.RawFields = nil

	return &tiddler, nil
}

// Put adds a new tiddler
func (t *TW) Put(tid *Tiddler) error {
	// Marshal tiddler to JSON
	b, err := json.Marshal(tid)
	if err != nil {
		return err
	}

	// Create PUT request with JSON data
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/recipes/default/tiddlers/%s", t.BaseURL.String(), tid.Title),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}

	// Send request
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		fmt.Printf("Couldn't put tiddler: %s", tid.Title)
	}

	return nil
}

// Append appends some text to an existing tiddler
func (t *TW) Append(tiddlerName string, text string) error {
	// Get tiddler
	tid, err := t.Get(tiddlerName)
	if err != nil {
		return err
	}

	// Append text
	// TODO: Have new line here or somewhere else?
	tid.Text += fmt.Sprintf("\n\n%s", text)

	// Put back tiddler
	err = t.Put(tid)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a tiddler specified by name
func (t *TW) Delete(tiddlerName string) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/bags/default/tiddlers/%s", t.BaseURL.String(), tiddlerName),
		nil,
	)

	// Send request
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read Response Body
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		log.Printf("Couldn't delete tiddler: %s\n", tiddlerName)
	}
	return err
}
