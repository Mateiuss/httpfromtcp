package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	numRequestLineElements = 3
	numHttpVersionElements = 2

	expectedHttpVersion = "1.1"
)

func parseRequestLine(requestLine string) (*RequestLine, error) {
	elements := strings.Split(requestLine, " ")
	if len(elements) != numRequestLineElements {
		return nil, errors.New("invalid number of request elements")
	}

	if elements[0] != strings.ToUpper(elements[0]) {
		return nil, errors.New("invalid method")
	}

	httpVersion := strings.Split(elements[2], "/")
	if len(httpVersion) != numHttpVersionElements {
		return nil, errors.New("invalid http version")
	}

	if httpVersion[1] != expectedHttpVersion {
		return nil, errors.New("expected version " + expectedHttpVersion + ", not " + httpVersion[1])
	}

	return &RequestLine{
		HttpVersion:   httpVersion[1],
		RequestTarget: elements[1],
		Method:        elements[0],
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("could not get data from reader")
	}

	lines := strings.Split(string(data), "\r\n")
	if len(lines) == 0 {
		return nil, errors.New("empty request")
	}

	requestLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, errors.New("could not parse request line: " + err.Error())
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}