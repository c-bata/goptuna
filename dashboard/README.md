# Dashboard

## Running dashboard

### Compiling TypeScript files

Node.js v14.14.0 is required to compile TypeScript files.

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

## Running Dashboard server with live-reloading

Please pass `-tags=develop` the custom build tag to return JS files inside local directory (not embedded files).
This is useful while development.

```
$ go run -tags=develop cmd/main.go dashboard --storage sqlite:///db.sqlite3
Started to serve at http://127.0.0.1:8000
```
