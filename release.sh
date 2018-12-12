#!/usr/bin/env bash

set -e

# TODO: detect versions from main.go and dockerfiles and make this more dynamic
FORCE=0
VERSION=""

for i in "$@"
do
    case ${i} in
        -f|--force)
        FORCE=1
        shift
        ;;
        *)
        ;;
    esac
done

function get_version() {
    V=$(grep 'const version' cmd/fdb-clusterfile/main.go | awk '{print $4}')
    VERSION="${V%\"}"
    VERSION="${VERSION#\"}"

    if ! [[ ${VERSION} =~ [0-9]+\.[0-9]+\.[0-9] ]]; then
     echo "Invalid version from cmd/fdb-clusterfile/main.go"
     exit
    fi
}

function already_built() {
    OS=$1
    if [[ -f ./releases/${OS}/${VERSION}/fdb-clusterfile ]]; then
        echo "1"
    else
        echo "0"
    fi
}


function build() {
    OS=$1

    built=$(already_built ${OS})

    if [[ "$built" == "1" ]] && [[ ${FORCE} -eq 0 ]]; then
        echo "$OS: binary already in ./releases/${OS}/${VERSION}/fdb-clusterfile. Use -f to rebuild"
        return
    fi

    mkdir -p ./releases/${OS}/${VERSION}
    docker build -t fdb-clusterfile-${OS}:${VERSION} -f ${OS}.dockerfile .
    docker run -d --name fdb-clusterfile-${OS} fdb-clusterfile-${OS}:${VERSION}
    docker cp fdb-clusterfile-${OS}:/build/cmd/fdb-clusterfile/fdb-clusterfile ./releases/${OS}/${VERSION}/
    docker stop fdb-clusterfile-${OS}
    docker rm fdb-clusterfile-${OS}
    echo "$OS: binary built and placed in ./releases/${OS}/${VERSION}/fdb-clusterfile"
}

get_version
build alpine
build debian
build ubuntu