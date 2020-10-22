# Dashboard

## Development

```
$ go run -tags=develop cmd/main.go dashboard --storage sqlite:///db.sqlite3
Started to serve at http://127.0.0.1:8000
```

```
$ npm run watch
```

## Production build

```
$ npm run build:prd
$ statik -src=./public -include=bundle.js,bundle.js.LICENSE.txt
```
