#!/bin/sh

export COMMIT_HASH=$(git rev-parse HEAD)
export VERSION=`node -pe 'require("./package.json").version'`
export NODE_ENV=production

echo "Build $VERSION ($COMMIT_HASH) for $NODE_ENV"

# generate bindata for client assets
vendor/bin/go-bindata -o ./server/bindata.go assets/... ./server/templates/... 

# build server
vendor/bin/gox \
 -ldflags "-X main.CommitHash=$COMMIT_HASH -X main.Version=$VERSION -X main.BuildTarget=production" \
 -osarch="$1/$2" \
 -output=build/gitss-${VERSION}-{{.OS}}-{{.Arch}} \
 ./server/...

mkdir dist
pushd build
tar cv gitss-${VERSION}-$1-$2* | gzip > ../dist/gitss-${VERSION}-$1-$2.tar.gz
popd