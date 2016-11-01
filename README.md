# GITSS

- [About](#about)
- [How to use](#how-to-use)
- [Development](#development)
 - [Requirements](#requirements)
 - [Setup](#setup)
 - [Run with development mode](#run-with-development-mode)
 - [Release Build](#release-build)
- [License](#license)

## About

`GITSS` is full text source search app for git repositories.
This repository is heavily under development and unstable now.

## How to use

1. Download binary file for your environment in the [Release Page](https://github.com/wadahiro/gitss/releases). Currently, you can download binary files for Linux(64bit) and Windows(64bit). 

2. To import git repository, run `gitss` command with `import` option as follows.

 ```bash
./gitss-0.1.0-linux-amd64 import yourOrgName yourProjectName http://your-git-site/your-git-repo.git
 ```

3. Then you just start with server mode. Run `gitss` command with `server` option as follows.

 ```bash
./gitss-0.1.0-linux-amd64 server
 ```

4. Open http://your-server:3000

## Development

### Requirements 

* [Golang](http://golang.org/)
* [Node.js](https://nodejs.org/)
* [Yarn](https://yarnpkg.com/)

### Setup

1. Install [gom](https://github.com/mattn/gom)(old version), [fresh](https://github.com/pilu/fresh).

 ```bash
go get github.com/mattn/gom
go get github.com/pilu/fresh

# install old version gom
cd $GOPATH/src/github.com/mattn/gom
git checkout 393e714d663c35e121a47fec32964c44a630219b
go install
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

1. Generate bindata.go for development mode.

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
