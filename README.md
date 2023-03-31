# Email Indexer

This project is a full-stack application that indexes emails for visualization. While the emails can be any type of email, the application was made with the [Enron emails database](http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz) in mind.

## Tech Stack

### Front-end

- [Vue 3](https://v3.vuejs.org/)
- [Vite](https://vitejs.dev/)
- [Tailwind CSS](https://tailwindcss.com/)

### Back-end

- [Go](https://golang.org/)
- [Zinc Search](https://github.com/zincsearch/zincsearch)

## Building the application

Get the code:

```sh
git clone https://github.com/amoralesc/email-indexer.git
cd email-indexer
```

### Zinc Search

Zinc Search is a search engine server that is used to index and search the emails. To run the application, you need to have Zinc Search installed and running.

Download the appropiate binary for your OS from the [Zinc Search releases page](https://github.com/zincsearch/zincsearch/releases) and run it with:

```sh
mkdir data
ZINC_FIRST_ADMIN_USER=admin ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123 ./zincsearch
```

Admin User and Password are defined by you.

## Running the application

### Back-end

After the Zinc Search server is running, you can run the back-end server:

```sh
cd indexer
go run .
```

The application defines the flags:

- `-dir`: Directory to look for email files recursively. Default: `none`, will not index emails.
- `-r`: Remove the existing index. Default: `false`, used to re-index the emails (creating a new index with the `-dir` flag).

It also defines environment variables:

- 'PORT': Port to run the server on. Default: `8080`.
- 'ZINC_PORT': Port where the Zinc Search server is running. Default: `4080`.
- 'NUM_UPLOADER_WORKERS': Number of concurrent workers to use for uploading emails to Zinc Search. Default: `4`. This parameter can be fine-tuned to improve performance.
- 'NUM_PARSER_WORKERS': Number of concurrent workers to use for parsing emails. Default: `8`. This parameter can be fine-tuned to improve performance.
- 'BULK_UPLOAD_SIZE': Number of emails to upload to Zinc Search in a single request. Default: `5000`. This parameter can be fine-tuned to improve performance.

### Front-end

After the back-end server is running, you can run the front-end server:

```sh
cd emails-app
```

Install dependencies:

```sh
npm install
```

Run the server:

```sh
npm run dev
```

Check it out at: http://localhost:5173/