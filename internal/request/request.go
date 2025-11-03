package request

import (
	"errors"
	"io"
	"strings"
)

const (
	stateInitialized = iota
	stateDone
)

const (
	numRequestLineElements = 3
	numHttpVersionElements = 2
	expectedHttpVersion = "1.1"
	
	bufferSize = 8
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	State       int
}

func parseRequestLine(request string) (int, *RequestLine, error) {
	idx := strings.Index(request, "\r\n")
	if idx == -1 {
		return 0, nil, nil
	}

	requestLine := request[:idx]

	elements := strings.Split(requestLine, " ")
	if len(elements) != numRequestLineElements {
		return 0, nil, errors.New("invalid number of request elements")
	}

	if elements[0] != strings.ToUpper(elements[0]) {
		return 0, nil, errors.New("invalid method")
	}

	httpVersion := strings.Split(elements[2], "/")
	if len(httpVersion) != numHttpVersionElements {
		return 0, nil, errors.New("invalid http version")
	}

	if httpVersion[1] != expectedHttpVersion {
		return 0, nil, errors.New("expected version " + expectedHttpVersion + ", not " + httpVersion[1])
	}

	return idx + len("\r\n"), &RequestLine{
		HttpVersion:   httpVersion[1],
		RequestTarget: elements[1],
		Method:        elements[0],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case stateInitialized:
		n, requestLine, err := parseRequestLine(string(data))

		if err != nil {
			return 0, errors.New("error while parsing the request line: " + err.Error())
		} else if n == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.State = stateDone

		return n, nil
	case stateDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	r := Request{
		State: stateInitialized,
	}

	for r.State != stateDone {
		if readToIndex >= len(buf) {
			tmp := make([]byte, cap(buf) * 2)
			copy(tmp, buf)
			buf = tmp
		}

		n, err := reader.Read(buf[readToIndex:])

		if err == io.EOF {
			r.State = stateDone
			break
		} else if err != nil {
			return nil, errors.New("error while reading data: " + err.Error())
		}

		readToIndex += n
		
		n, err = r.parse(buf)
		if err != nil {
			return nil, errors.New("error while parsing data: " + err.Error())
		}

		readToIndex -= n
	}

	return &r, nil
}