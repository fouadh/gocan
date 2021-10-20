#!/usr/bin/env bash

set -ex

if [[ -z "${GOCAN}" ]]; then
  GOCAN="gocan"
fi

echo "Cleaning up before starting"
$GOCAN delete-scene hibernate

echo "Creating the scene & the app"
$GOCAN create-scene hibernate
$GOCAN create-app orm --scene hibernate

echo "Importing the history"
tmp_dir=$(mktemp -d -t gocan-)
echo "Cloning hibernate in $tmp_dir/orm"
git clone https://github.com/hibernate/hibernate-orm.git $tmp_dir/orm

cd $tmp_dir/orm
$GOCAN import-history orm --scene hibernate --after 2012-01-01 --before 2013-09-05

echo "Looking at app summary"
$GOCAN app-summary orm --scene hibernate

echo "Looking at first revisions"
$GOCAN revisions orm --scene hibernate | head -n 10

echo "Analyze complexity trends of the Configuration class"
$GOCAN create-complexity-analysis configuration-analysis \
      --app orm \
      --scene hibernate \
      --filename hibernate-core/src/main/java/org/hibernate/cfg/Configuration.java \
      --directory $tmp_dir/orm/

