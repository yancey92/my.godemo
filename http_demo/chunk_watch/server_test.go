package main2

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func Test_test1(t *testing.T) {
	url := "http://localhost:8080/register?watch=1&s=2"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')

		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				break
			}
			fmt.Printf("Error: %v\n", err)
			break
		}

		fmt.Print(string(line))
	}

}
