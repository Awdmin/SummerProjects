package request

import(
	"io"
	"fmt"
	"strings"
	"errors"
)

type RequestLine struct {
	HttpVersion   	string
	RequestTarget 	string
	Method        	string
}

type Request struct {
	RequestLine		RequestLine
}

var ERROR_BAD_REQUEST_LINE = fmt.Errorf("bad request line")
var ERROR_UNSUPORTED_HTTP_VERSION = fmt.Errorf("unsuported HTTP version")
var ERROR_INCOMPLETE_REQUEST_LINE = fmt.Errorf("incomplete request line")
var SEPERATOR = "\r\n"

func ParseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, b, ERROR_INCOMPLETE_REQUEST_LINE
	}

	startLine := b[:idx]
	restOfMsg := b[idx+len(SEPERATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfMsg, ERROR_BAD_REQUEST_LINE
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMsg, ERROR_BAD_REQUEST_LINE
	}

	rl := &RequestLine {
		HttpVersion: 	httpParts[1],
		RequestTarget: 	parts[1],
		Method: 		parts[0],
	}

	return rl, restOfMsg, nil
}

func RequestFromReader(reader io.Reader) (*Request, error){
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("unable to read all"), 
			err,
		)
	}

	str := string(data)
	rl, _, err := ParseRequestLine(str)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: 	*rl,
	}, err

}


