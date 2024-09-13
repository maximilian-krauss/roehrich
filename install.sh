#! /bin/sh

# this is a work in progress install script which "works on my machine".
# feel free to make this more robust and reliable.

# dependencies:
# - curl
# - jq
set -eo pipefail

UNAME=$(echo uname)
ARCH=$(arch)
LATEST_VERSION=$(curl -s https://api.github.com/repos/maximilian-krauss/roehrich/releases/latest | jq -r '.tag_name')
DOWNLOAD_URL="https://github.com/maximilian-krauss/roehrich/releases/download/$LATEST_VERSION/roehrich_$(UNAME)_${ARCH}.tar.gz"
DOWNLOAD_DESTINATION="/tmp/roehrich_$LATEST_VERSION.tar.gz"

curl -L -s -o "$DOWNLOAD_DESTINATION" "$DOWNLOAD_URL"
tar -xzvf "$DOWNLOAD_DESTINATION" -C /tmp

chmod +x /tmp/roehrich
sudo mv -f "/tmp/roehrich" "/usr/local/bin"
rm "$DOWNLOAD_DESTINATION"
