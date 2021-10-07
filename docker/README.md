[gocan](https://github.com/fouadh/gocan) is a code analyzer that can be used to generate some of the charts described in the book [Your Code as a Crime Scene](https://pragprog.com/titles/atcrime/your-code-as-a-crime-scene/) written by Adam Tornhill.

# Usage Example

## Starting gocan

```
docker run --publish 8888:80 --name gocan --volume /code:/code -d fouadhamdi/gocan:latest
```

## Creating a scene and an app

```
git clone https://github.com/hibernate/hibernate-orm.git /code/hibernate
docker exec gocan gocan create-scene hibernate
docker exec gocan gocan create-app orm --scene hibernate
docker exec gocan gocan import-history orm --scene hibernate --after 2012-01-01 --before 2013-09-04 --directory /code/hibernate
```

## Visualizing the charts

Open your browser at the address `http://localhost:8888`

