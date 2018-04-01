package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
)

var (
	fileCache sync.Map
)

func getFileAsBytes(p string) []byte {
	if f, exists := fileCache.Load(p); exists {
		return f.([]byte)
	}

	f := fmt.Sprintf("%s/index.html", p)
	fp := filepath.Join(configurationServer.GetValue("frontend_location"), f)
	fileBytes, _ := ioutil.ReadFile(fp)
	fileCache.Store(p, fileBytes)
	return fileBytes
}

func getURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Printf("%s request for %s Origin: %s", r.Method, r.URL, r.RemoteAddr)
	w.Header().Set("Content-Type", "text/html")
	index := getFileAsBytes(r.URL.String())
	w.Write(index)
}

func addBodyToHTML(htmlTags string) string {
	htmlTemplate := getFileAsBytes("/templates")
	lines := strings.Split(string(htmlTemplate), "\n")
	lines = append(lines, htmlTags)
	lines = append(lines, "</body>\n</html>")
	return strings.Join(lines, "\n")
}

// SuccessResponse returns the successful response
func SuccessResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	html := addBodyToHTML(fmt.Sprintf("<b>%s</b>", "Command Successful"))
	w.Write([]byte(html))
}

// ErrorResponse creates a page which tells you what your error is
func ErrorResponse(w http.ResponseWriter, ErrorString string) {
	w.Header().Set("Content-Type", "text/html")
	html := addBodyToHTML(fmt.Sprintf("<b>%s</b>", ErrorString))
	w.Write([]byte(html))
}
