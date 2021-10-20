#!/usr/bin/env bash

set -ex

if [[ -z "${GOCAN}" ]]; then
  GOCAN="gocan"
fi

$GOCAN --version

echo "Cleaning up before starting"
$GOCAN delete-scene maat

echo "Creating the scene & the app"
$GOCAN create-scene maat
$GOCAN create-app maat --scene maat

echo "Importing the history"
tmp_dir=$(mktemp -d -t gocan-)
echo "Cloning maat in $tmp_dir/maat"
git clone https://github.com/adamtornhill/code-maat.git $tmp_dir/maat

cd $tmp_dir/maat
$GOCAN import-history maat --scene maat --before 2013-11-01

echo "Get app summary"
$GOCAN app maat --scene maat

echo "Analyze change frequency"
$GOCAN revisions maat --scene maat | head -n 10

echo "Sum of coupling"
$GOCAN soc maat --scene maat | head -n 10

echo "Measuring coupling"
$GOCAN coupling maat --scene maat
