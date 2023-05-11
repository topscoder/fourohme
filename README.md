# FourOhMe

FourOhMe is a tool for testing HTTP headers on a website in order to try to bypass 40* HTTP codes. 
It tries to bypass the 40* status code by manipulating the HTTP request via HTTP headers and URL payloads.

## Installation

Install Golang, then run:

`go install -v github.com/topscoder/fourohme@latest`

## Usage

FourOhMe can be used in three ways:

1. Reading URLs from a file:

   ```
   $ fourohme -file urls.txt
   ```

   where `urls.txt` is a file containing a list of URLs separated by newlines.

2. Reading URLs from standard input (`STDIN`):

   ```
   $ cat urls.txt | fourohme
   ```

   where `urls.txt` is a file containing a list of URLs separated by newlines.
   
   Or attach in your favourite chain commands:
   
   ```
   $ cat domains.txt | subfinder | httpx -mc 401,402,403,404,405 -silent | fourohme 
   ```

3. Reading single URL:

    ```
    $ fourohme -url https://foo.bar
    ```

## Options

- `-file`: Path to a file containing URLs.
- `-force`: Force the scanner to scan all URL's regardless of the initial HTTP status code.
- `-silent`: Don't print shizzle. Only what matters.
- `-threads`: The amount of threads to be used to execute the HTTP requests. Be gentle or get blocked. (default 4)
- `-url`: URL to make requests to

## Output

FourOhMe prints the responses of each request in a colored format. A successful response (HTTP status code 200) is printed in green, while an unsuccessful response is printed in red.

## Contributing

Contributions are welcome! If you find a bug or want to suggest a new feature, please open an issue or submit a pull request.

## TODO

[ ] Add a WAF check
[ ] Add content check filter false positives ("The requested URL was rejected")

## License

FourOhMe is released under the [MIT License](https://github.com/topscoder/fourohme/blob/main/LICENSE).
