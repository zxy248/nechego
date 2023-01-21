package danbooru

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	data := []byte(`{"file_url": "asdf", "file_size": 123, "rating": "s"}`)
	want := file{URL: "asdf", Size: 123, Rating: Sensitive}

	handler := func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, bytes.NewReader(data))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	f, err := get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if *f != want {
		t.Errorf("%v != %v", *f, want)
	}
}
