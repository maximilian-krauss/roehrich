#! /bin/sh

# this is a work in progress install script which "works on my machine".
# feel free to make this more robust and reliable.

# dependencies:
# - curl
# - jq
set -eo pipefail

UNAME=$(echo uname)
ARCH=$(arch)
NAME=roehrich
LATEST_VERSION=$(curl -s "https://api.github.com/repos/maximilian-krauss/$NAME/releases/latest" | jq -r '.tag_name')
DOWNLOAD_URL="https://github.com/maximilian-krauss/$NAME/releases/download/$LATEST_VERSION/${NAME}_$(UNAME)_${ARCH}.tar.gz"
DOWNLOAD_DESTINATION="/tmp/${NAME}_$LATEST_VERSION.tar.gz"

curl -L -s -o "$DOWNLOAD_DESTINATION" "$DOWNLOAD_URL"
tar -xzvf "$DOWNLOAD_DESTINATION" -C /tmp

chmod +x "/tmp/$NAME"
sudo mv -f "/tmp/$NAME" "/usr/local/bin"
rm "$DOWNLOAD_DESTINATION"
