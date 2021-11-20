package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HistoryItem holds one HTTP request data
type HistoryItem struct {
	*http.Request
	Date time.Time
	body string
}

// History holds a set of requests
type History struct {
	values []HistoryItem
}

// AddItem adds one item to history
func (h *History) AddItem(item HistoryItem) {
	h.values = append(h.values, item)
}


//GetHistory returns requests history in reverse order
func (h *History) GetHistory(reverse bool) []HistoryItem {
	if !reverse {
		return h.values
	}

	reversed := append([]HistoryItem{}, h.values...)
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	return reversed
}

// PrintString returns unescaped html string
func (hi *HistoryItem) PrintString() template.HTML {
	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "<p>RequestURI: %s</p>", hi.RequestURI)
	for s, header := range hi.Header {
		fmt.Fprintf(buff, "<p> %s: %s</p>", s, strings.Join(header, ""))
	}

	fmt.Fprintf(buff, "<pre>%s</pre>", hi.body)
	for s, header := range hi.Form {
		fmt.Fprintf(buff, "<p> %s: %s<p>", s, strings.Join(header, ""))
	}

	return template.HTML(buff.String())
}

// PrintConsole writes a content of the item to the console
func (hi *HistoryItem) PrintConsole(w http.ResponseWriter) {
	fmt.Printf("\n%s", hi.Method)
	fmt.Printf("\nRemoteAddr: %s", hi.RemoteAddr)
	fmt.Printf("\nRequestURI: %s", hi.RequestURI)
	printHeaders(hi.Header, w)
	data := printBody(hi.Body, w)
	hi.body = data
	printHeaders((map[string][]string(hi.Form)), nil)
	fmt.Printf("\n")
	w.Header().Add("referrer", "http-echo-server")
}

func printBody(body io.ReadCloser, w io.Writer) string {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ""
	}
	fmt.Printf("\nBody:")
	fmt.Printf("\n %s", data)
	_, _ = w.Write(data)

	return string(data)
}

func printHeaders(headers http.Header, w http.ResponseWriter) {
	for s, header := range headers {
		w.Header().Add(s, strings.Join(header, ""))
		fmt.Printf("\n %s: %s", s, strings.Join(header, ""))
	}
}
