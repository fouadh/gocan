#!/usr/bin/env bash

set -ex

if [[ -z "${GOCAN}" ]]; then
  GOCAN="gocan"
fi

$GOCAN --version

echo "Cleaning up before starting"
$GOCAN delete-scene scala

echo "Creating the scene & the app"
$GOCAN create-scene scala
$GOCAN create-app scala --scene scala

echo "Importing the history"
tmp_dir=$(mktemp -d -t gocan-)
echo "Cloning scala in $tmp_dir/scala"
git clone https://github.com/nopSolutions/nopCommerce $tmp_dir/scala

cd $tmp_dir/scala
$GOCAN import-history scala --scene scala --before 2013-12-31 --after 2011-12-31

echo "Display the main developers"
$GOCAN main-devs scala --scene scala | head -n 10
