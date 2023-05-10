# Four-Oh-Me

Four-Oh-Me is a tool for testing HTTP headers on a website in order to try to bypass 40* HTTP codes. It makes requests to a given URL with different headers and prints the responses.

## Usage

Four-Oh-Me can be used in two ways:

1. Reading URLs from a file:

   ```
   $ four-oh-me --file urls.txt
   ```

   where `urls.txt` is a file containing a list of URLs separated by newlines.

2. Reading URLs from standard input:

   ```
   $ cat urls.txt | four-oh-me
   ```

   where `urls.txt` is a file containing a list of URLs separated by newlines.

## Options

- `--url`: URL to make requests to
- `--file`: Path to a file containing URLs

## Headers

By default, Four-Oh-Me makes requests with the following headers:

- `X-Forwarded-For: 127.0.0.1:80`
- `X-Custom-IP-Authorization: 127.0.0.1`
- `X-Original-URL: %URL%`
- `X-Original-URL: %PATH%`
- `HTTP: OPTIONS`

You can add or modify these headers by editing the `headersList` variable in the script.

## Output

Four-Oh-Me prints the responses of each request in a colored format. A successful response (HTTP status code 200) is printed in green, while an unsuccessful response is printed in red.

## Contributing

Contributions are welcome! If you find a bug or want to suggest a new feature, please open an issue or submit a pull request.

## License

Four-Oh-Me is released under the [MIT License](https://github.com/yourname/four-oh-me/blob/main/LICENSE).