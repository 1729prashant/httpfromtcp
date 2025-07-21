package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Request represents an HTTP request with a RequestLine and parser state
type Request struct {
	RequestLine RequestLine
	state       parserState // internal state to track parsing progress
}

// RequestLine holds components of the HTTP start-line
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// Constants for line termination
const crlf = "\r\n"

// parserState is an enum to track parsing stages
type parserState int

const (
	stateInitialized parserState = iota // parsing just started
	stateDone                           // parsing is complete
)

// RequestFromReader reads and parses an HTTP request from a stream incrementally
func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{state: stateInitialized}
	buf := make([]byte, 0, 8) // initial buffer size
	temp := make([]byte, 8)   // temporary read buffer

	for r.state != stateDone {
		// Read from reader
		n, err := reader.Read(temp)
		if err != nil && err != io.EOF {
			return nil, err
		}
		buf = append(buf, temp[:n]...) // append new data to buffer

		// Attempt to parse from buffer
		consumed, parseErr := r.parse(buf)
		if parseErr != nil {
			return nil, parseErr
		}

		// Remove consumed bytes from buffer
		buf = buf[consumed:]

		if err == io.EOF {
			break // end of stream
		}
	}

	if r.state != stateDone {
		return nil, fmt.Errorf("incomplete request")
	}

	return r, nil
}

// parse attempts to parse the next section of data into the Request struct
func (r *Request) parse(data []byte) (int, error) {
	if r.state == stateInitialized {
		consumed, requestLine, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if consumed == 0 {
			return 0, nil // not enough data to parse
		}
		r.RequestLine = *requestLine
		r.state = stateDone
		return consumed, nil
	}
	return 0, nil // already done
}

// parseRequestLine attempts to extract and parse the start-line from input
func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil, nil // need more data, not an error
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(strings.TrimSpace(requestLineText))
	if err != nil {
		return 0, nil, err
	}
	return idx + len(crlf), requestLine, nil
}

// requestLineFromString validates and parses the method, target, and HTTP version
func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Fields(str)
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]
	if requestTarget == "" {
		return nil, fmt.Errorf("empty request target")
	}

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}
