package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"moul.io/http2curl"
)

// HistoryItem holds one HTTP request data
type HistoryItem struct {
	*http.Request
	Date         time.Time
	bodyOriginal string
	bodyMock     string
	headersOrder []string
	CurlCommand  string
}

// History holds a set of requests
type History struct {
	values []HistoryItem
}

// AddItem adds one item to history
func (h *History) AddItem(item *HistoryItem) {
	item.headersOrder = make([]string, 0, len(item.Header))
	for s, _ := range item.Header {
		item.headersOrder = append(item.headersOrder, s)
	}
	curlCommand, _ := http2curl.GetCurlCommand(item.Request)
	item.CurlCommand = curlCommand.String()

	h.values = append(h.values, *item)
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
	for _, header := range hi.headersOrder {
		fmt.Fprintf(buff, "<p> %s: %s</p>", header, strings.Join(hi.Header[header], ""))
	}

	fmt.Fprintf(buff, "<pre style=\"max-width:770px;\">reqest body:<code class=\"language-json\">%s</code></pre>", hi.bodyOriginal)
	if hi.bodyMock != "" {
		fmt.Fprintf(buff, "<pre style=\"max-width:770px;\">mock response:<code class=\"language-json\">%s</code></pre>", hi.bodyMock)

	}
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
	// headers
	for _, header := range hi.headersOrder {
		w.Header().Add(header, strings.Join(hi.Header[header], ""))
		fmt.Printf("\n %s: %s", header, strings.Join(hi.Header[header], ""))
	}
	hi.printBody(w)
	// form
	for s, header := range hi.Form {
		w.Header().Add(s, strings.Join(header, ""))
		fmt.Printf("\n %s: %s", s, strings.Join(header, ""))
	}
	fmt.Printf("\n")
	w.Header().Add("referrer", "http-echo-server")
}

func (hi *HistoryItem) printBody(w io.Writer) {
	fmt.Printf("\nBody:")
	fmt.Printf("\n %s", hi.bodyOriginal)
	if hi.bodyMock != "" {
		fmt.Printf("\nMock response:")
		fmt.Printf("\n %s", hi.bodyMock)
		_, _ = w.Write([]byte(hi.bodyMock))
	} else {
		fmt.Printf("\nResponse is the same")
		_, _ = w.Write([]byte(hi.bodyOriginal))
	}
}

func (hi *HistoryItem) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{
		"Header":     hi.Header,
		"Body":       hi.bodyOriginal,
		"Method":     hi.Method,
		"URL":        hi.URL,
		"RemoteAddr": hi.RemoteAddr,
		"RequestURI": hi.RequestURI,
	}

	return json.Marshal(result)
}
