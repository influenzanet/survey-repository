#!/bin/sh
export VERSION=$(cut -d '|' -f 1  pkg/version/version.txt)
echo "Building rpm for $VERSION"
nfpm package -f build/nfpm.yaml -p rpm --target survey-repository-${VERSION}-linux-x86_64.rpm