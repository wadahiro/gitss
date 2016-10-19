#!/bin/sh

export COMMIT_HASH=$(git rev-parse HEAD)
export VERSION=`node -pe 'require("./package.json").version'`
export NODE_ENV=production

echo "Build $VERSION ($COMMIT_HASH) for $NODE_ENV"

# build client
node_modules/.bin/node-sass client/css/main.scss assets/css/style.css
node_modules/.bin/webpack -p --config ./client/webpack/webpack.config.js

# generate bindata for client assets
vendor/bin/go-bindata -o ./server/bindata.go assets/... ./server/templates/... 

# build server
vendor/bin/gox \
 -ldflags "-X main.CommitHash=$COMMIT_HASH -X main.Version=$VERSION -X main.BuildTarget=production" \
 -osarch="$1/$2" \
 -output=build/gitS-${VERSION}-{{.OS}}-{{.Arch}} \
 ./server/...

mkdir dist
pushd dist
tar cv ../build/gitS-${VERSION}-$1-$2* | gzip > gitS-${VERSION}-$1-$2.tar.gz
popd