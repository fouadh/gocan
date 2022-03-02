#!/usr/bin/env bash

set -ex

if [[ -z "${GOCAN}" ]]; then
  GOCAN="gocan"
fi

"$GOCAN" --version

echo "Cleaning up before starting"
"$GOCAN" delete-scene craft

echo "Creating the scene & the app"
"$GOCAN" create-scene craft
"$GOCAN" create-app craft --scene craft

echo "Importing the history"
tmp_dir=$(mktemp -d -t gocan-XXXXXX)
echo "Cloning craft in $tmp_dir/craft"
git clone https://github.com/SirCmpwn/Craft.Net.git $tmp_dir/craft

cd $tmp_dir/craft
"$GOCAN" import-history craft --scene craft --before 2014-08-08 --interval-between-analyses 30

echo "Sum of coupling before 2014-08-08"
"$GOCAN" soc craft --scene craft | head -n 10

echo "Analyze coupling of MinecraftServer before 2013-01-01"
"$GOCAN" coupling craft --scene craft --before 2013-01-01 --min-degree 40 | grep MinecraftServer.cs

echo "Analyze coupling of MinecraftServer between 2013-01-01 and 2014-08-08"
"$GOCAN" coupling craft --scene craft --after 2013-01-01 --min-degree 40 | grep MinecraftServer.cs

echo "Identify modules coupled with MinecraftServer.cs"
"$GOCAN" coupling craft --scene craft --min-revisions-average 15 | grep MinecraftServer.cs