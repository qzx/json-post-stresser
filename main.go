package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type config struct {
	url         string
	bearerToken string
	method      string
	payload     string
	connections int
	timeout     int
	isHTTPS     bool
}

func main() {
	count := 0
	errn := 0

	start := time.Now()
	if len(os.Args) < 5 {
		fmt.Println("Not enough args: .\\client.exe bearer url body-file.json CONNECTION_COUNT(int)")
		return
	}

	c := configFromArgs(os.Args)

	var wg sync.WaitGroup

	for i := 1; i <= c.connections; i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			for {
				err := makeCall(&c)
				if err != nil {
					errn++
				} else {
					count++
				}
				break
			}
		}(&wg)
	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("requests took: %s, had %v errors, and %v successes", elapsed, errn, count)
}

func configFromArgs(a []string) config {
	var c config
	reqUrl, err := url.Parse(os.Args[2])
	if err != nil {
		panic("Couldn't parse url")
	}
	c.url = os.Args[2]
	if reqUrl.Scheme == "https" {
		c.isHTTPS = true
	} else {
		c.isHTTPS = false
	}

	c.bearerToken = fmt.Sprintf("Bearer %s", a[1])
	c.method = "POST"
	payloadData, err := ioutil.ReadFile(a[3])
	if err != nil {
		panic("Couldn't open json payload file")
	}
	payload := string(payloadData)
	c.payload = payload
	connections, err := strconv.Atoi(os.Args[4])
	if err != nil {
		panic("connections should be a integer")
	}
	c.connections = connections

	return c
}

func makeCall(c *config) error {
	payload := strings.NewReader(c.payload)

	client := &http.Client{}
	if c.isHTTPS {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client.Transport = customTransport
	}
	req, err := http.NewRequest(c.method, c.url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.bearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
