package server

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	endpoints := `[
			{
				"path":"/graphQl",
				"request":"{\"foo1\":\"bar\"}",
				"mock":"{\"foo_resp1\":\"bar_resp1\"}"
			},
			{
				"path":"/graphQl",
				"request":"{\"foo2\":\"bar\"}",
				"mock":"{\"foo_resp2\":\"bar_resp2\"}"
			},
			{
				"path":"/auth",
				"request":"",
				"mock":"{\"key\":\"auth_key\"}"
			}
		]`
	s := NewServer(":9999", nil, endpoints)
	go func() {
		s.Start()
	}()
	time.Sleep(100 * time.Millisecond)

	r, err := http.Post("http://localhost:9999/graphQl", "", bytes.NewBufferString(`{"foo1":"bar"}`))
	if err != nil {
		t.Fatalf("error creating request: %+v", err)
	}

	if !checkResponseEqual(r, `{"foo_resp1":"bar_resp1"}`) {
		t.Fatal("response is not equal")
	}

	r, err = http.Post("http://localhost:9999/graphQl", "", bytes.NewBufferString(`{"foo2":"bar"}`))
	if err != nil {
		t.Fatal("error creating request")
	}

	if !checkResponseEqual(r, `{"foo_resp2":"bar_resp2"}`) {
		t.Fatal("response is not equal")
	}

	r, err = http.Post("http://localhost:9999/auth", "", bytes.NewBufferString(`{"login":"user"}`))
	if err != nil {
		t.Fatal("error creating request")
	}

	if !checkResponseEqual(r, `{"key":"auth_key"}`) {
		t.Fatal("response is not equal")
	}

	s.Stop(context.Background())
	time.Sleep(100 * time.Millisecond)
}

func checkResponseEqual(r *http.Response, sample string) bool {
	if r.Body == nil && sample != "" {
		return false
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("error reading bodyOriginal:%+v", err)
	}

	return string(b) == sample
}
