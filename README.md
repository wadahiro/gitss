# GitSS

[![wercker status](https://app.wercker.com/status/8bb1c52941262ac810ef6219e02937cf/s/develop "wercker status")](https://app.wercker.com/project/byKey/8bb1c52941262ac810ef6219e02937cf)

- [About](#about)
- [How to use](#how-to-use)
- [Development](#development)
 - [Requirements](#requirements)
 - [Setup](#setup)
 - [Run with development mode](#run-with-development-mode)
 - [Release Build](#release-build)
- [License](#license)

## About

**GitSS** is full text source search app for git repositories.
This repository is heavily under development and unstable now.

## How to use

### Requirements

* [Git](https://git-scm.com/)

GitSS use `git` command internaly. So you need to set `PATH` environment variable to use `git` command.

### Install Binary

Download binary file for your environment in the [Release Page](https://github.com/wadahiro/gitss/releases). Currently, you can download binary files for Linux(64bit), Darwin(64bit) and Windows(64bit).
Then put `gitss` (or `gitss.exe` for Windows) file in the archive into the directory where you'd like to install.

### Add Setting for a git repository

To add sync setting, run `gitss add` command as follows.

 ```bash
./gitss add yourOrgName yourProjectName http://your-git-site/your-git-repo.git
 ```

This will create `data/conf/yourOrgName.json` file as follows. Using this setting file, GitSS do syncing with remote git repository and indexing the contents automatically.

 ```json
{
  "name": "yourOrgName",
  "projects": [
    {
      "name": "yourProjectName",
      "repositories": [
        {
          "url": "http://your-git-site/your-git-repo.git",
          "sizeLimit": 1048576,
          "includeBranches": ".*",
          "includeTags": ".*"
        }
      ]
    }
  ]
}
 ```

Also there are more options for `gitss add`. Please check `gitss add --help`.


### Add Setting for a Bitbucket server

If you use Bitbucket server, you can sync and index the all repositories easily.
Run `gitss bitbucket add` command as follows.

 ```bash
./gitss bitbucket add yourOrgName http://your-bitbucket-server-site/bitbucket --user=yourId --password=yourPassword
 ```

This will create `data/conf/yourOrgName.json` file as follows. Using this setting file, GitSS fetch all repositories from Bitbucket server, then do syncing git repositories and indexing.

 ```json
{
  "name": "yourOrgName",
  "scm": {
    "excludeProjects": "",
    "excludeRepositories": "",
    "includeProjects": ".*",
    "includeRepositories": ".*",
    "password": "yourPassword",
    "type": "bitbucket",
    "url": "http://your-bitbucket-server-site/bitbucket",
    "user": "yourId"
  },
  "sizeLimit": 1048576,
  "includeBranches": ".*",
  "includeTags": ".*"
}
 ```

Also there are more options for `gitss bitbucket add`. Please check `gitss bitbucket add --help`.


### Manual syncing & indexing

After adding setting file, run `gitss sync` command with `--all` option as follows. GitSS read all setting files and sync git repository and index the contents.

 ```bash
./gitss sync --all
 ```

 If you'd like to sync one repository only, you can use `gitss sync` command as follows.

 ```bash
./gitss sync yourOrgName yourProjectName your-git-repo
 ```

### Run server

Run `gitss server` command as follows. Then open http://your-server:3000 in your browser. In addition, the sync scheduler is started when starting GitSS server. The scheduler do syncing git repository and indexing the contents automatically.

 ```bash
./gitss server
 ```

If you'd like to change the port (Default: 3000), you can use `--port` option.

 ```bash
./gitss server --port=5000
 ```

If you'd like to change the sync schedule (Default: 0 */10 * * * *), you can use `--schedule` option as follows.

 ```bash
./gitss server --schedule="0 */30 * * * *"
 ```


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

2. Generate vendor.js for development mode.

 ```
npm run build:client-dll
 ```

3. Start webpack and gin with watch mode.

 ```bash
npm start & fresh
 ```
 
4. Open http://localhost:9000

### Release Build

Run webpack with production mode, go-bindata and go build in turn. All you have to do is run `npm run build`. The artifact is created under `./dist` directory.

 ```bash
npm run build
 ```

## License

Licensed under the [MIT](/LICENSE.txt) license.
