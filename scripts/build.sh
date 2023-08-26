#!/usr/bin/env bash

custom_package=$1
if [[ -z "$custom_package" ]]; then
  echo "Usage: $0 diffr [dir1/file1 dir2/file2]]"
  exit 1
fi

custom_tag="0.1.0"

custom_platforms=(
 "darwin/amd64"
 "darwin/arm64"
 "linux/386"
 "linux/amd64"
 "linux/arm"
 "linux/arm64"
 "windows/386"
 "windows/amd64"
 "windows/arm"
)

rm -rf custom-platforms-$custom_tag/*
for custom_platform in "${custom_platforms[@]}"
do
    platform_split=(${custom_platform//\// })
    CUSTOM_GOOS=${platform_split[0]}
    CUSTOM_GOARCH=${platform_split[1]}
    echo 'Building' $CUSTOM_GOOS-$CUSTOM_GOARCH
    custom_output_name='diffr-'$CUSTOM_GOOS-$CUSTOM_GOARCH

    env GOOS=$CUSTOM_GOOS GOARCH=$CUSTOM_GOARCH VERSION=$custom_tag go build -ldflags "-X main.Version=$custom_tag" -v -o custom-platforms-$custom_tag/$custom_output_name $custom_package

    if [ $? -ne 0 ]; then
        echo 'An error occurred! Aborting the script execution...'
        exit 1
    fi

    custom_bin_name='diffr'
    if [ "$CUSTOM_GOOS" == 'windows' ]; then
      custom_bin_name='diffr.exe'
    fi

    cd custom-platforms-$custom_tag
    mv $custom_output_name $custom_bin_name
    tar -czvf $custom_output_name-$custom_tag.tar.gz $custom_bin_name
    rm -rf $custom_bin_name
    cd ..
done
