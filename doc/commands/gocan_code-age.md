## gocan code-age

Retrieve the age of the entities in months

```
gocan code-age [flags]
```

### Examples

```

gocan age myapp --scene myscene
gocan age myapp --scene myscene --after 2022-01-01 --before 2022-06-30
gocan age myapp --scene myscene --initial-date 2021-01-01

```

### Options

```
      --after string          Calculate the code age after this day
      --before string         Calculate the code age before this day
      --csv                   get the results in csv format
  -h, --help                  help for code-age
      --initial-date string   From when to calculate the age (default "2022-12-19")
  -s, --scene string          Scene name
```

### SEE ALSO

* [gocan](gocan.md)	 - 

###### Auto generated by spf13/cobra on 19-Dec-2022
