package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func maybeFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func setupDb(u string) {
	req, err := http.NewRequest("PUT", u, nil)
	maybeFatal(err)
	res, err := http.DefaultClient.Do(req)
	maybeFatal(err)
	res.Body.Close()
}

func sendOne(u, k string, body []byte) {
	resp, err := http.DefaultClient.Post(u+"?ts="+k,
		"application/json", bytes.NewReader(body))
	maybeFatal(err)
	defer resp.Body.Close()
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		log.Fatalf("HTTP Error on %v: %v", k, err)
	}
}

func main() {
	u := os.Args[1]
	setupDb(u)

	t := time.Tick(5 * time.Second)
	i := 0

	d := json.NewDecoder(os.Stdin)
	for {
		kv := map[string]*json.RawMessage{}

		err := d.Decode(&kv)
		if err == io.EOF {
			log.Printf("Done!")
			break
		}
		maybeFatal(err)

		for k, v := range kv {
			body := []byte(*v)
			sendOne(u, k, body)
		}

		i++
		select {
		case <-t:
			var k string
			for k = range kv {
			}
			log.Printf("Processed %v items, latest was %v", i, k)
		default:
		}
	}
}
