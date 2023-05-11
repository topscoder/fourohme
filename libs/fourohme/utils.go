package fourohme

import (
	"flag"
	"fmt"
	"net/url"
)

func ParseCommandLineFlags() (*string, *string, *bool, *int, *bool) {
	urlPtr := flag.String("url", "", "URL to make requests to")
	filePtr := flag.String("file", "", "Path to a file containing URLs")
	silentPtr := flag.Bool("silent", false, "Don't print shizzle. Only what matters.")
	threadsPtr := flag.Int("threads", 4, "The amount of threads to be used to execute the HTTP requests. Be gentle or get blocked.")
	forcePtr := flag.Bool("force", false, "Force the scanner to scan all URL's regardless of the initial HTTP status code.")
	flag.Parse()

	return urlPtr, filePtr, silentPtr, threadsPtr, forcePtr
}

func GetHostAndPath(parsedURL *url.URL) (string, string) {
	sUrl := parsedURL.Scheme + "://" + parsedURL.Host
	sPath := parsedURL.Path
	if sPath == "" {
		sPath = "/"
	}

	return sUrl, sPath
}

func printOutput(statusCode int, verb string, url string, headers []Header) {
	// Print in green if it's 200
	if statusCode == 200 {
		fmt.Printf("\033[32m%d => HTTP %s %s %v\033[0m\n", statusCode, verb, url, headers)
	} else {
		fmt.Printf("\033[31m%d => HTTP %s %s %v\033[0m\n", statusCode, verb, url, headers)
	}
}

func ShowBanner() {
	const banner = `


	███████╗ ██████╗ ██╗   ██╗██████╗      ██████╗ ██╗  ██╗    ███╗   ███╗███████╗
	██╔════╝██╔═══██╗██║   ██║██╔══██╗    ██╔═══██╗██║  ██║    ████╗ ████║██╔════╝
	█████╗  ██║   ██║██║   ██║██████╔╝    ██║   ██║███████║    ██╔████╔██║█████╗  
	██╔══╝  ██║   ██║██║   ██║██╔══██╗    ██║   ██║██╔══██║    ██║╚██╔╝██║██╔══╝  
	██║     ╚██████╔╝╚██████╔╝██║  ██║    ╚██████╔╝██║  ██║    ██║ ╚═╝ ██║███████╗
	╚═╝      ╚═════╝  ╚═════╝ ╚═╝  ╚═╝     ╚═════╝ ╚═╝  ╚═╝    ╚═╝     ╚═╝╚══════╝

	by @topscoder

	`

	fmt.Println(banner)
}
