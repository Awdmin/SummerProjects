package main

import(
	"fmt"
	"log"
	"io"
	"bytes"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string{
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		curr_line := ""

		for {
			data := make([]byte, 8) 
			n, err := f.Read(data)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal("error", err)
			}

			data = data[:n]
			
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				curr_line += string(data[:i])
				out <- curr_line
				curr_line = string(data[i+1:])
			} else {
				curr_line += string(data[:n])
			}

		}
		if len(curr_line) != 0 {
			out <- curr_line
		}
	}()

	

	return out
}


func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}


	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("error", err)
		}

		for line := range getLinesChannel(conn){
			fmt.Printf("read: %s\n", line)
		}
	}

	

}
