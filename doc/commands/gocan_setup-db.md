## gocan setup-db

Configure the database options. Caution: if you are changing the embedded db properties, you will lose all the existing data.

```
gocan setup-db [flags]
```

### Options

```
  -n, --database string    Database name (default "postgres")
  -d, --directory string   Directory where the data will be stored. Only valid for the embedded database. (default "/Users/fhamdi/.gocan/data")
  -e, --external-db        Set this flag if you prefer to use an external db rather than the embedded one
  -h, --help               help for setup-db
      --host string        Database host (default "localhost")
  -p, --password string    Database password (default "postgres")
      --port int           Port to which the database listens. (default 5432)
  -u, --user string        Database user (default "postgres")
      --verbose            display the log information
```

### SEE ALSO

* [gocan](gocan.md)	 - 

###### Auto generated by spf13/cobra on 19-Dec-2022
