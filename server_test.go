package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestJSON(t *testing.T) {
	header := http.Header{}
	headerKey := "Content-Type"
	headerValue := "application/json; charset=utf-8"
	header.Add(headerKey, headerValue)

	testCases := []struct {
		in     interface{}
		header http.Header
		out    string
	}{
		{map[string]string{"hello": "world"}, header, `{"hello":"world"}`},
		{map[string]string{"go": "lang"}, header, `{"go":"lang"}`},
		{map[string]string{"key": "value"}, header, `{"key":"value"}`},
		{make(chan bool), header, `{"error":"json: unsupported type: chan bool"}`},
	}
	for _, test := range testCases {

		recorder := httptest.NewRecorder()

		JSON(recorder, test.in)

		response := recorder.Result()
		defer response.Body.Close()

		got, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %s\n", err)
		}

		if string(got) != test.out {
			t.Errorf("Got %s, expected %s\n", string(got), test.out)
		}

		if contentType := response.Header.Get(headerKey); contentType != headerValue {
			t.Errorf("Got %s, expected %s\n", contentType, headerValue)
		}
	}
}

func TestGet(t *testing.T) {
	makeStorage(t)
	defer cleanupStorage(t)

	kvStore := map[string]string{
		"hello": "world",
		"key1":  "value1",
		"go":    "lang",
	}

	encodedStore := map[string]string{}
	for k, v := range kvStore {
		encodedKey := base64.URLEncoding.EncodeToString([]byte(k))
		encodedValue := base64.URLEncoding.EncodeToString([]byte(v))
		encodedStore[encodedKey] = encodedValue
	}

	fileContents, _ := json.Marshal(encodedStore)
	os.WriteFile(StoragePath+"/data.json", fileContents, 0644)

	testCases := []struct {
		in  string
		out string
		err error
	}{
		{"hello", "world", nil},
		{"key1", "value1", nil},
		{"test", "", nil},
	}

	for _, test := range testCases {
		got, err := Get(context.Background(), test.in)
		if err != test.err {
			t.Errorf("Got: %s, expected: %s\n", err, test.err)
		}

		if got != test.out {
			t.Errorf("Got %s, expected %ss\n", got, test.out)
		}
	}
}

func makeStorage(t *testing.T) {
	err := os.Mkdir("testdata", 0755)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Couldn't create directory testdata: %s\n", err)
	}

	StoragePath = "testdata"
}

func cleanupStorage(t *testing.T) {
	if err := os.RemoveAll(StoragePath); err != nil {
		t.Errorf("Failed to delete storage path: %s\n", StoragePath)
	}
	StoragePath = "/tmp/kv"
}
