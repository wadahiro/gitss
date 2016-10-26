# GITSS

- [About](#about)
- [Development](#development)
 - [Requirements](#requirements)
 - [Setup](#setup)
 - [Run with development mode](#run-with-development-mode)
 - [Release Build](#release-build)
- [License](#license)

## About

`gitss` is full text search server for git repositories.

## Development

### Requirements 

* [Golang](http://golang.org/)
* [Node.js](https://nodejs.org/)
* [Yarn](https://yarnpkg.com/)

### Setup

1. Install [gom](https://github.com/mattn/gom), [fresh](https://github.com/pilu/fresh).

```bash
go get github.com/mattn/gom
go get github.com/pilu/fresh
```

2. Install Golang dependencies.

```bash
gom install
```

3. Install JavaScript dependencies.

```bash
yarn
```

### Run with development mode

1. Generate bindata.go.

 ```bash
npm run bindata
 ```

2. Start webpack and gin with watch mode.

 ```bash
npm run devserver & fresh
 ```
 
3. Open http://localhost:9000

### Release Build

Run webpack with production mode, go-bindata and go build in turn. All you have to do is run `npm run build`. The artifact is created under `./dist` directory.

```bash
npm run build
```

## License

Licensed under the [MIT](/LICENSE.txt) license.
