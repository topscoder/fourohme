package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	urlPtr, filePtr, silentPtr := parseCommandLineFlags()

	if *silentPtr == false {
		showBanner()
	}

	headersList := []map[string]string{
		{"X-Forwarded-For": "127.0.0.1:80"},
		{"X-Custom-IP-Authorization": "127.0.0.1"},
		{"X-Host": "127.0.0.1"},
		{"X-Original-URL": "%URL%"},
		{"X-Original-URL": "%PATH%"},
		{"X-rewrite-url": "%PATH%"},
		{"Content-Length": "0", "HTTP": "POST"},
		{"HTTP": "OPTIONS"},
		{"HTTP": "TRACE"},
	}

	urls := readUrlsFromInput(urlPtr, filePtr)
	for i, pUrl := range urls {
		parsedURL, err := url.Parse(pUrl)
		if err != nil {
			panic(err)
		}

		sUrl, sPath := getHostAndPath(parsedURL)

		for _, headers := range headersList {
			verb := getVerb(headers)
			headers = replacePlaceholders(headers, sUrl, sPath)
			req := createRequest(verb, pUrl, headers)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err, i)
				return
			}

			printOutput(resp.StatusCode, verb, sUrl, sPath, headers)

			resp.Body.Close()
		}

		fmt.Println("")
	}
}

func parseCommandLineFlags() (*string, *string, *bool) {
	urlPtr := flag.String("url", "", "URL to make requests to")
	filePtr := flag.String("file", "", "Path to a file containing URLs")
	silentPtr := flag.Bool("silent", false, "Don't print shizzle. Only what matters.")
	flag.Parse()

	return urlPtr, filePtr, silentPtr
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

func getVerb(headers map[string]string) string {
	verb := "GET"
	for key, val := range headers {
		if key == "HTTP" {
			verb = val
			break
		}
	}
	return verb
}

func replacePlaceholders(headers map[string]string, sUrl, sPath string) map[string]string {
	for key, val := range headers {
		if key != "HTTP" {
			val = strings.ReplaceAll(val, "%URL%", sUrl)
			val = strings.ReplaceAll(val, "%PATH%", sPath)
			headers[key] = val
		}
	}
	return headers
}

func createRequest(verb string, pUrl string, headers map[string]string) *http.Request {
	req, err := http.NewRequest(verb, pUrl, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return req
}

func printOutput(statusCode int, verb string, url string, path string, headers map[string]string) {
	// Print in green if it's 200
	if statusCode == 200 {
		fmt.Printf("\033[32m%d => HTTP %s %s%s %v\033[0m\n", statusCode, verb, url, path, headers)
	} else {
		fmt.Printf("\033[31m%d => HTTP %s %s%s %v\033[0m\n", statusCode, verb, url, path, headers)
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
