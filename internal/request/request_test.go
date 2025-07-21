package request

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good POST Request line
	r, err = RequestFromReader(strings.NewReader("POST /submit HTTP/1.1\r\nHost: localhost\r\n\r\n"))
	require.NoError(t, err)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/submit", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line (missing method)
	_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost\r\n\r\n"))
	require.Error(t, err)

	// Test: Invalid method with lowercase letters
	_, err = RequestFromReader(strings.NewReader("get / HTTP/1.1\r\nHost: localhost\r\n\r\n"))
	require.Error(t, err)

	// Test: Invalid method with special characters
	_, err = RequestFromReader(strings.NewReader("G3T / HTTP/1.1\r\nHost: localhost\r\n\r\n"))
	require.Error(t, err)

	// Test: Invalid method with mixed case
	_, err = RequestFromReader(strings.NewReader("GeT / HTTP/1.1\r\nHost: localhost\r\n\r\n"))
	require.Error(t, err)

	// Test: Invalid HTTP version
	_, err = RequestFromReader(strings.NewReader("GET / HTTP/1.0\r\nHost: localhost\r\n\r\n"))
	require.Error(t, err)

	// Test: Valid method but no path
	_, err = RequestFromReader(strings.NewReader("GET  HTTP/1.1\r\nHost: localhost\r\n\r\n")) // double space
	require.Error(t, err)

	// Test: Extra space in request line (4 parts)
	_, err = RequestFromReader(strings.NewReader("GET / something HTTP/1.1\r\nHost: localhost\r\n\r\n"))
	require.Error(t, err)

	// Test: Trailing whitespace after request line
	r, err = RequestFromReader(strings.NewReader("GET / HTTP/1.1   \r\nHost: localhost\r\n\r\n"))
	require.NoError(t, err) // still valid: we trim line before parsing
	assert.Equal(t, "GET", r.RequestLine.Method)

	// Test: Empty input
	_, err = RequestFromReader(strings.NewReader(""))
	require.Error(t, err)

	// Test: Only CRLF (no request line)
	_, err = RequestFromReader(strings.NewReader("\r\nHost: localhost\r\n\r\n"))
	require.Error(t, err)

	// Test: Good GET Request line
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}

	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	reader = &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

}

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}
