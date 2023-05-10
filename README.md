# FourOhMe

FourOhMe is a tool for testing HTTP headers on a website in order to try to bypass 40* HTTP codes. It makes requests to a given URL with different headers and prints the responses.

## Installation

Install Golang, then run:

`go install -v github.com/topscoder/fourohme`

## Usage

FourOhMe can be used in three ways:

1. Reading URLs from a file:

   ```
   $ fourohme -file urls.txt
   ```

   where `urls.txt` is a file containing a list of URLs separated by newlines.

2. Reading URLs from standard input:

   ```
   $ cat urls.txt | fourohme
   ```

   where `urls.txt` is a file containing a list of URLs separated by newlines.

3. Reading single URL:

    ```
    $ fourohme -url https://foo.bar
    ```

## Options

- `-file`: Path to a file containing URLs
- `-url`: URL to make requests to
- `-silent`: Don't print shizzle

## Headers

By default, FourOhMe makes requests with the following headers:

- `X-Forwarded-For: 127.0.0.1:80`
- `X-Custom-IP-Authorization: 127.0.0.1`
- `X-Original-URL: %URL%`
- `X-Original-URL: %PATH%`
- `HTTP: OPTIONS`

You can add or modify these headers by editing the `headersList` variable in the script.

## Output

FourOhMe prints the responses of each request in a colored format. A successful response (HTTP status code 200) is printed in green, while an unsuccessful response is printed in red.

## Contributing

Contributions are welcome! If you find a bug or want to suggest a new feature, please open an issue or submit a pull request.

## License

FourOhMe is released under the [MIT License](https://github.com/topscoder/fourohme/blob/main/LICENSE).