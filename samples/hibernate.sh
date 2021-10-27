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
tmp_dir=$(mktemp -d -t gocan-XXXXXX)
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
      --directory $tmp_dir/orm/ \
      --spaces 4

echo "Looking for the modus operandi"
$GOCAN modus-operandi orm --scene hibernate | head -n 10

echo "Looking for the relationships between authors and entities"
$GOCAN revisions-authors orm --scene hibernate | head -n 10

echo "Measure temporal coupling for AbstractEntityPersister"
$GOCAN coupling orm --scene hibernate --min-revisions-average 20 | grep AbstractEntityPersister

echo "Identify main developer of AbstractEntityPersister"
$GOCAN main-devs orm --scene hibernate | grep AbstractEntityPersister

echo "Identify developers of the modules coupled to AbstractEntityPersister"
$GOCAN main-devs orm --scene hibernate | grep CustomPersister
$GOCAN main-devs orm --scene hibernate | grep entity/EntityPersister
$GOCAN main-devs orm --scene hibernate | grep GoofyPersisterClassProvider

echo "Calculate individual contributions to EntityPersister"
$GOCAN entity-efforts orm --scene hibernate | grep entity/EntityPersister

echo "Rename author"
$GOCAN rename-dev --app orm --scene hibernate --current edalquist --new "Eric Dalquist"

echo "Verify the impact on the individual contributions"
$GOCAN entity-efforts orm --scene hibernate | grep entity/EntityPersister

echo "Verify the main developer of EntityPersister"
$GOCAN main-devs orm --scene hibernate | grep entity/EntityPersister

echo "Calculate individual contributions to AbstractEntityPersister"
$GOCAN entity-efforts orm --scene hibernate | grep AbstractEntityPersister