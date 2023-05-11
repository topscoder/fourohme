package fourohme

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
)

func TalkHttpBaby(ch chan Request, wg *sync.WaitGroup, silent bool) {
	defer wg.Done() // Schedule the wg.Done() function call to be executed when the function returns

	request := <-ch

	statusCode := ExecuteHttpRequest(request)

	if !silent || (statusCode >= 200 && statusCode <= 303) {
		printOutput(statusCode, request.Verb, request.Url, request.Headers)
	}
}

func ExecuteHttpRequest(request Request) int {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from panic:", err)
		}
	}()

	verb := request.Verb
	url := request.Url
	headers := request.Headers

	req := createRequest(verb, url)

	if req == nil {
		return -1
	}

	for _, header := range headers {
		req.Header.Add(header.Key, header.Value)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")

	// Create a transport with insecure skip verify
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create a client with the custom transport
	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println(err)
		resp.Body.Close()
		return -1
	}

	resp.Body.Close()
	return resp.StatusCode
}

func createRequest(verb string, pUrl string) *http.Request {
	req, err := http.NewRequest(verb, pUrl, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return req
}
