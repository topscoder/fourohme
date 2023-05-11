package fourohme

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadUrlsFromInput(urlPtr, filePtr *string) []string {
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
