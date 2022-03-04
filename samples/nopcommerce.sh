#!/usr/bin/env bash

set -ex

if [[ -z "${GOCAN}" ]]; then
  GOCAN="gocan"
fi

"$GOCAN" --version

echo "Cleaning up before starting"
"$GOCAN" delete-scene nopcommerce

echo "Creating the scene & the app"
"$GOCAN" create-scene nopcommerce
"$GOCAN" create-app nopcommerce --scene nopcommerce

echo "Importing the history"
tmp_dir=$(mktemp -d -t gocan-XXXXXX)
echo "Cloning nopcommerce in $tmp_dir/nopcommerce"
git clone https://github.com/nopSolutions/nopCommerce $tmp_dir/nopcommerce

cd $tmp_dir/nopcommerce
"$GOCAN" import-history nopcommerce --scene nopcommerce --before 2014-09-25 --after 2014-01-01

echo "Create software architecture boundaries"
"$GOCAN" create-boundary architecture --scene nopcommerce --app nopcommerce \
        --module "Admin Models:src/Presentation/Nop.Web/Administration/Models" \
        --module "Admin Views:src/Presentation/Nop.Web/Administration/Views" \
        --module "Admin Controllers:src/Presentation/Nop.Web/Administration/Controllers" \
        --module "Services:src/Libraries/Nop.Services" \
        --module "Core:src/Libraries/Nop.Core" \
        --module "Data Access:src/Libraries/Nop.Data" \
        --module "Business Access Layer:src/Libraries/Nop.Services" \
        --module "Models:src/Presentation/Nop.Web/Models" \
        --module "Views:src/Presentation/Nop.Web/Views" \
        --module "Controllers:src/Presentation/Nop.Web/Controllers"

echo "Analyze coupling between layers"
"$GOCAN" coupling nopcommerce --scene nopcommerce --boundary architecture

echo "Analyze revisions between layers"
"$GOCAN" revisions nopcommerce --scene nopcommerce --boundary architecture