package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	idx := bytes.Index(rawBytes, []byte(crlf))
	if idx == -1 {
		return nil, errors.New("missing CRLF after request line")
	}

	requestLineBytes := rawBytes[:idx]
	requestLineText := string(requestLineBytes)

	requestLine, err := parseRequestLine(requestLineText)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	// Use Fields to be resilient against multiple spaces or leading/trailing whitespace
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed request line: expected 3 parts but got %d", len(parts))
	}

	method := parts[0]
	target := parts[1]
	versionStr := parts[2]

	if !isAllUpperAlpha(method) {
		return nil, fmt.Errorf("invalid method: must be all uppercase letters, got '%s'", method)
	}

	versionParts := strings.Split(versionStr, "/")
	if len(versionParts) != 2 || versionParts[0] != "HTTP" {
		return nil, fmt.Errorf("malformed HTTP version: '%s'", versionStr)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unsupported HTTP version: '%s'", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version, // normalized to "1.1"
	}, nil
}

func isAllUpperAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}
