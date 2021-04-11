# WordPress XMLRPC Exploit

Brute-force a WordPress website credentials

## Usage

```md
Usage of ./main:
  -dict string
    A text dictionary containing passwords to try (default "dict.txt")

  -url string
    URL to WordPress website to target (without trailing /)
    (default "https://www.wordpress-site.com")

  -user string
    Username to brute-force (default "admin")
```

### Where do I get the Website Usernames?

You can use the excellent [WPScan](https://github.com/wpscanteam/wpscan) to Enumerate username for any vulnerable WordPress website.
