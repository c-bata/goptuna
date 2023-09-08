# Dashboard

## Running dashboard

Pre-built JavaScript files are embedded using `go:embed`.
So it's ok you just run `goptuna dashboard` command like this:

```
$ go build cmd/main.go -o goptuna
$ ./goptuna dashboard --storage sqlite:///example.db --host 127.0.0.1 --port 8000
```

<details>

<summary>more command line options</summary>

```
$ ./bin/goptuna --help
A command line interface for Goptuna

Usage:
  goptuna [command]

Available Commands:
  create-study Create a study in your relational database storage.
  dashboard    Launch web dashboard
  delete-study Delete a study in your relational database storage.
  help         Help about any command

Flags:
  -h, --help      help for goptuna
      --version   version for goptuna

Use "goptuna [command] --help" for more information about a command.
```

</details>


## How to compile TypeScript files

### Compiling TypeScript files manually

Node.js v20 is required to compile TypeScript files.

```
$ npm install
$ npm run build:dev
```

<details>
<summary>Watch for files changes</summary>

```
$ npm run watch
```

</details>

<details>
<summary>Production builds</summary>

```
$ npm run build:prd
```

</details>

## Running Dashboard server for reloading

Please pass `-tags=develop` the custom build tag to return JS files inside local directory (not embedded files).
This is useful while devleopment.

```
$ go run -tags=develop cmd/main.go dashboard --storage sqlite:///db.sqlite3
Started to serve at http://127.0.0.1:8000
```
