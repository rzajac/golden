package golden

import (
	"bufio"
	"net/http"
	"net/textproto"
	"strings"
)

// parseHeaders parses HTTP header lines from golden file.
func parseHeaders(t T, lines ...string) http.Header {
	rdr := strings.NewReader(strings.Join(lines, "\n") + "\r\n\r\n")
	tp := textproto.NewReader(bufio.NewReader(rdr))
	h, err := tp.ReadMIMEHeader()
	if err != nil {
		t.Fatal(err)
	}
	return http.Header(h)
}
