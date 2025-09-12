package request

import(
	"io"
	"fmt"
	"bytes"
)

type parserState string
const(
	StateInit parserState = "init"
	StateDone parserState = "done"
	StateError parserState = "error"
)

type RequestLine struct {
	HttpVersion   	string
	RequestTarget 	string
	Method        	string
}

type Request struct {
	RequestLine		RequestLine
	State			parserState
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.State {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := ParseRequestLine(data[read:])
			if err != nil {
				r.State = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.State = StateDone

		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.State == StateDone || r.State == StateError
}


var ERROR_BAD_REQUEST_LINE = fmt.Errorf("bad request line")
var ERROR_UNSUPORTED_HTTP_VERSION = fmt.Errorf("unsuported HTTP version")
var ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")
var ERROR_INCOMPLETE_REQUEST_LINE = fmt.Errorf("incomplete request line")
var SEPERATOR = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func ParseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, ERROR_INCOMPLETE_REQUEST_LINE
	}

	startLine := b[:idx]
	read := idx+len(SEPERATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_BAD_REQUEST_LINE
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ERROR_BAD_REQUEST_LINE
	}

	rl := &RequestLine {
		HttpVersion: 	string(httpParts[1]),
		RequestTarget: 	string(parts[1]),
		Method: 		string(parts[0]),
	}

	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error){
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil

}


