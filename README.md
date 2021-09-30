# Introduction

Pet project heavily inspired from the book [Your Code as a Crime Scene]() from Adam Tornhill.

This cli tool will allow you to build some of the charts described in that book in order to start conversations regarding the design or organisation of some application.

It has very similar interface to [code-maat](), the tool created by Adam Tornhill. The differences are mainly:

- written in golang
- only git support
- usage of a database to store the data
- ui section to get the charts

# Using the app

## Setup the database

### To setup the embedded database

```
gocan setup-db
```

## Start the embedded database

```
gocan start-db
```

## Run the UI

```
gocan ui
```

## Create a forensics scene

```
gocan create-scene my-scene
```

## Add an application to a scene

```
gocan create-app my-app -s my-scene
```

## Stop the embedded database

```
gocan stop-db
```

# A Few Minutes Tutorial

Let's use one of the examples in the book.

```
gocan create-scene hibernate
gocan create-app orm -s hibernate
git clone https://github.com/hibernate/hibernate-orm.git
gocan import-history orm -s hibernate --after 2011-12-31 --before 2013-09-05 --path ./hibernate-orm

```


# Building the app

## Requirements

* golang 1.16
* nodejs
* yarn

## Build

```
make build
```

# Running the tests

```
make test
```

