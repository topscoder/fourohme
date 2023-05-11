package fourohme

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
)

func TalkHttpBaby(ch chan Request, wg *sync.WaitGroup) {
	defer wg.Done() // Schedule the wg.Done() function call to be executed when the function returns

	request := <-ch

	statusCode := ExecuteHttpRequest(request)

	printOutput(statusCode, request.Verb, request.Url, request.Headers)
}

func ExecuteHttpRequest(request Request) int {
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
		fmt.Println(err)
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
