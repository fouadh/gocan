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

