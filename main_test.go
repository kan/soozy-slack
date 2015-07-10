package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func testHttp(t *testing.T, resp *http.Response, err error) interface{} {
	if err != nil {
		t.Error(err)
		return nil
	}

	if resp.StatusCode != 200 {
		t.Error("Status code error")
		return nil
	}

	defer resp.Body.Close()
	var data interface{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data)
	return data
}

func testHttpGet(t *testing.T, uri string) interface{} {
	resp, err := http.Get(uri)
	return testHttp(t, resp, err)
}

func testHttpPost(t *testing.T, uri string, values url.Values) interface{} {
	resp, err := http.PostForm(uri, values)
	return testHttp(t, resp, err)
}

func TestMain(t *testing.T) {
	config := Config{"soozy.slack.com", "8080", "test-token"}
	ts := httptest.NewServer(http.HandlerFunc(InviteFunc(config)))
	defer ts.Close()

	// invite not allow GET request
	result := testHttpGet(t, ts.URL)
	if result != nil {
		err := result.(map[string]interface{})["error"]
		if err != "not allow GET reuqest" {
			t.Errorf("%v", result)
		}
	}

	// fail when empty "email" field
	result = testHttpPost(t, ts.URL, url.Values{})
	if result != nil {
		err := result.(map[string]interface{})["error"]
		if err != "empty email" {
			t.Errorf("%v", result)
		}
	}

	// success case
	values := url.Values{}
	values.Add("email", "example@example.com")
	result = testHttpPost(t, ts.URL, values)
	if result != nil {
		status := result.(map[string]interface{})["success"].(float64)
		if status != 1 {
			t.Errorf("%v", result)
		}
	}
}
