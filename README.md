# bsky stats api

An API written in Go (1.23) to get various statistics about Bluesky. Powered by [jaz.bsky.social](https://bsky.app/profile/jaz.bsky.social)'s [bsky stats](https://bsky.jazco.dev/stats).

Derived from [espeon/bcounter-backend](https://github.com/espeon/bcounter-backend).

Ths API esentially just wraps the parent API, adds some growth rate calculations, and caches the results for continuous fetching.

![GitHub Sponsors](https://img.shields.io/github/sponsors/dickeyy)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/dickeyy/bsky-stats-api)

## Usage

This API is free to use, but it is also extremely easy to self-host. If you can, please self-host this API as it will save on resources for my server.

**Endpoint:** `GET https://bsky-stats.kyle.so`
Returns:

```json
{
	"total_users": 123456,
	"total_posts": 123456,
	"total_follows": 123456,
	"total_likes": 123456,
	"users_growth_rate_per_second": 0.123456,
	"last_update_time": "2023-01-01T00:00:00Z",
	"next_update_time": "2023-01-01T00:00:00Z"
}
```

The API will attempt to serve cahced data if it is available. Cached data is invalidated every 60 seconds, so hitting the API over and over will not return new data unless the update time has passed.

### Self-hosting

You can self-host this API via Docker (see [Dockerfile](./Dockerfile)) or build it from source.

The Docker file does not provide any reverse proxying, so you will need to do that on your own.

See [.env.example](./.env.example) for a list of environment variables you can set.

#### Docker

```bash
docker build -t bsky-stats-api .
docker run -p 8080:8080 bsky-stats-api
```

#### Build from source

```bash
go build -o bsky-stats-api .
./bsky-stats-api
```

## Development

If you want to run this API locally, you will need to set the following environment variables:

```bash
export ENV=dev
export PORT=8080
```

Then run the API:

```bash
go run main.go
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Credits

-   [espeon/bcounter-backend](https://github.com/espeon/bcounter-backend) - The original API that this is derived from, and growth rate calculations.
-   [jaz.bsky.social](https://bsky.app/profile/jaz.bsky.social) - The parent API that this API wraps.

## License

Licensed under the MIT License. See [LICENSE](./LICENSE) for details.
