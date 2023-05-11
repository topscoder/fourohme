package fourohme

import (
	"fmt"
	"net/http"
	"sync"
)

func TalkHttpBaby(ch chan Request, wg *sync.WaitGroup) {
	defer wg.Done() // Schedule the wg.Done() function call to be executed when the function returns

	request := <-ch

	statusCode := executeHttpRequest(request)

	printOutput(statusCode, request.Verb, request.Url, request.Headers)
}

func executeHttpRequest(request Request) int {
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

	resp, err := http.DefaultClient.Do(req)
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
