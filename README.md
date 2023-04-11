# Email Indexer

This project is a full-stack application that indexes emails for visualization. While the emails can be any type of email, the application was made with the [Enron emails database](https://en.wikipedia.org/wiki/Enron_Corpus) in mind.

## Tech Stack

### Front-end

- [Vue 3](https://v3.vuejs.org/)
- [Vite](https://vitejs.dev/)
- [Tailwind CSS](https://tailwindcss.com/)

### Back-end

- [Go](https://golang.org/)
- [Zinc Search](https://github.com/zincsearch/zincsearch)

## Requirements (local development)

The application is configured to run locally using Docker via [Docker Compose](https://docs.docker.com/compose/compose-file/). It's recommended to use Docker Engine, although you may use Docker Desktop.

- [Docker Engine](https://docs.docker.com/engine/install/)
- [Docker Desktop](https://docs.docker.com/get-docker/)

Check that the docker daemon is running and the version with:

```sh
docker version
```

The `docker-compose.yml` file specification requires a version `19.03.0+`.

## Running the application locally

Get the code:

```sh
git clone https://github.com/amoralesc/email-indexer.git
cd email-indexer
```

Copy the `*.env.example` files into their respective `*.env` files and update the values as needed (see [Configuration](#configuration)):

```sh
cp .env.example .env
cp .indexer.env.example .indexer.env
```

Get the Enron emails database (or any other emails, just place them in the `emails` directory)

```sh
# you may need to: chmod +x get-enron-emails.sh
./get-enron-emails.sh
```

Build the Docker images:

```sh
docker compose build
```

Run the application:

```sh
docker compose up -d
```

The application should be running at [http://localhost:8080](http://localhost:8080).

Check out the logs with:

```sh
docker compose logs -f
```

## Configuration

### Email files

The application expects the emails to be in the `emails` directory. The emails should be in the syntax RFC 5322 / RFC 6532. The application will recursively search for any files in the `emails` directory. Any valid email file will be indexed.

The `emails` directory is directly mounted into the `indexer` container as a volume. The `indexer` container will then parse and upload theses emails to the Zinc server (see [Indexing](#indexing)).

The environment variable `EMAILS_DIR` can be used to change the directory where the emails are stored. However, this may break the application if configured incorrectly.

### Indexing

The indexing process is done by the `indexer` container. The `indexer` container will parse the emails and upload them to the Zinc server. This process uses goroutines to speed up the indexing process.

Some environment variables can be used to configure the indexing performance process. Fine tuning these variables may improve the indexing performance.

| Variable | Description | Default |
| --- | --- | --- |
| `NUM_PARSER_WORKERS` | Number of goroutines spawned to parse email files into JSON | `128` |
| `NUM_UPLOADER_WORKERS` | Number of goroutines spawned to upload JSON emails from the indexer to Zinc | `32` |
| `BULK_UPLOAD_SIZE` | Number of emails sent in a single bulk upload operation to Zinc | `5000` |

Other environment variables control the behavior of the `indexer` container, specially when the container is restarted.

| Variable | Description | Default |
| --- | --- | --- |
| `REMOVE_INDEX_IF_EXISTS` | If `true`, the `indexer` container will remove the index from Zinc if it already exists | `false` |
| `SKIP_UPLOAD_IF_INDEX_EXISTS` | If `true`, the `indexer` container will skip uploading emails to Zinc if the index already exists. This is useful for preventing re-upload of emails when the attached directory hasn't changed | `true` |

The `REMOVE_INDEX_IF_EXISTS` and `SKIP_UPLOAD_IF_INDEX_EXISTS` variables are meant to be overriden in the `.indexer.env` file.

### Profiling

The application can be configured to enable a profiling server. This is useful for debugging performance issues. The profiler is the default Go profiler, which is based on the [pprof](
https://golang.org/pkg/net/http/pprof/) package.

Both back-end containers (the `api` and `indexer` containers) can be configured to enable the profiler. The profiler is disabled by default.

| Variable | Description | Default |
| --- | --- | --- |
| `ENABLE_PROFILING` | If `true`, the containers enable profiling | `false` |
| `INDEXER_ENABLE_PROFILING` | If `true`, the `indexer` container enables profiling | `false` |
| `API_ENABLE_PROFILING` | If `true`, the `api` container enables profiling | `false` |
| `PROFILING_PORT` | The port that the profiler is exposed on | `6060` |
| `INDEXER_PROFILING_PORT` | The port that the profiler is exposed on for the `indexer` container | `6060` |
| `API_PROFILING_PORT` | The port that the profiler is exposed on for the `api` container | `6061` |
| `SLEEP_TIME_AFTER_INDEXING` | The seconds to sleep after the `indexer` container finishes indexing | `0` |

The `ENABLE_PROFILING` variable is meant to be overriden by `INDEXER_ENABLE_PROFILING` and `API_ENABLE_PROFILING`. The `PROFILING_PORT` variable is meant to be overriden by `INDEXER_PROFILING_PORT` and `API_PROFILING_PORT`. This behavior is done automatically by the `docker-compose.yml` file.

The `SLEEP_TIME_AFTER_INDEXING` variable is meant to be overriden in the `.indexer.env` file. This variable is useful for debugging the `indexer` container. It allows the `indexer` container to sleep after it finishes indexing. This allows the CPU profiling tool to keep running after the indexing finishes. If this variable is set to `0`, the `indexer` container will exit immediately after it finishes indexing and the CPU profiling tool might not have finished running, which will cause the profiling tool to fail.

#### Running the profiling tool

The profiling tool requires `go` to be installed. Refer to the [Go installation instructions](https://go.dev/doc/install) for more information.

For example, if the `INDEXER_ENABLE_PROFILING` is set to `true` and its exposed port is set to `6060`, the CPU can be profiled for 60 seconds with:

```sh
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=60
```

which will activate a pprof interactive shell after 60 seconds. We then can produce a PNG graph of the CPU profile with:

```sh
(pprof) png
```

or a text report with:

```sh
(pprof) top
```

More info about the Go profiler can be found in the [Go's blog about pprof](https://go.dev/blog/pprof).

### Environment variables

As shown above, application defines multiple environment variables. The following tables describe each variable and its default value.

#### `.env` file

| Variable | Description | Default |
| --- | --- | --- |
| `WEB_PORT` | The port that the web app container is exposed on | `8080` |
| `ZINC_ADMIN_USER` | The username for the Zinc server admin user | `admin` |
| `ZINC_ADMIN_PASSWORD` | The password for the Zinc server admin user | `Complexpass#123` |
| `ZINC_HOST` | The host where the other containers find Zinc. WARNING: not supposed to be changed | `zinc` |
| `ZINC_PORT` | The port that the Zinc server is exposed on | `4080` |
| `ZINC_RETRY_INTERVAL` | The containers' entrypoint use it to retry connecting to the Zinc server (in seconds) | `5` |
| `EMAILS_DIR` | The directory where the emails are stored. WARNING: not supposed to be changed, this may break the app | `emails` |
| `REMOVE_INDEX_IF_EXISTS` | If `true`, the `indexer` container will remove the index from Zinc if it already exists | `false` |
| `SKIP_UPLOAD_IF_INDEX_EXISTS` | If `true`, the `indexer` container will skip uploading emails to Zinc if the index already exists | `true` |
| `ENABLE_PROFILING` | If `true`, the containers enable profiling | `false` |
| `INDEXER_ENABLE_PROFILING` | If `true`, the `indexer` container enables profiling | `false` |
| `API_ENABLE_PROFILING` | If `true`, the `api` container enables profiling | `false` |
| `PROFILING_PORT` | The port that the profiler is exposed on | `6060` |
| `INDEXER_PROFILING_PORT` | The port that the profiler is exposed on for the `indexer` container | `6060` |
| `API_PROFILING_PORT` | The port that the profiler is exposed on for the `api` container | `6061` |
| `API_PORT` | The port that the API container is exposed on | `3000` |
| `NUM_PARSER_WORKERS` | Number of goroutines spawned to parse email files into JSON | `128` |
| `NUM_UPLOADER_WORKERS` | Number of goroutines spawned to upload JSON emails from the indexer to Zinc | `32` |
| `BULK_UPLOAD_SIZE` | Number of emails sent in a single bulk upload operation to Zinc | `5000` |

#### `.indexer.env` file

Overrides the default values of the `.env` file for the `indexer` container.

| Variable | Description | Default |
| --- | --- | --- |
| `REMOVE_INDEX_IF_EXISTS` | If `true`, the `indexer` container will remove the index from Zinc if it already exists | `false` |
| `SKIP_UPLOAD_IF_INDEX_EXISTS` | If `true`, the `indexer` container will skip uploading emails to Zinc if the index already exists | `true` |
| `SLEEP_TIME_AFTER_INDEXING` | The seconds to sleep after the `indexer` container finishes indexing | `0` |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
