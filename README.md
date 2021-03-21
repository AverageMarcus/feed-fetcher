![feed-fetcher](logo.png)

Returns the RSS feed associated with the given URL

Available at https://feed-fetcher.cluster.fun/

## Usage

```sh
GET https://feed-fetcher.cluster.fun/?url=${URL_TO_CHECK}
```

Example:

```sh
curl -v http://localhost:8000/\?url\=https://marcusnoble.co.uk/
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8000 (#0)
> GET /?url=https://marcusnoble.co.uk/ HTTP/1.1
> Host: localhost:8000
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 307 Temporary Redirect
< Date: Wed, 17 Mar 2021 09:44:32 GMT
< Content-Type: text/plain; charset=utf-8
< Content-Length: 18
< Location: https://marcusnoble.co.uk/feed.xml
<
* Connection #0 to host localhost left intact
Temporary Redirect* Closing connection 0
```

### Possible status code responses

* **300** - Multiple possible feeds found on page (the first is returned on the `Location` header)
* **301** - URL provided was already a valid feed URL
* **307** - Feed URL found on provided page
* **400** - No URL provided
* **404** - No feed URL found on provided webpage
* **500** - Server error while trying to fetch feed

## Building from source

With Docker:

```sh
make docker-build
```

Standalone:

```sh
make build
```

## Contributing

If you find a bug or have an idea for a new feature please [raise an issue](issues/new) to discuss it.

Pull requests are welcomed but please try and follow similar code style as the rest of the project and ensure all tests and code checkers are passing.

Thank you 💛

## License

See [LICENSE](LICENSE)
