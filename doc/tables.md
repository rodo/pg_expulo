# Table configuration options

## Name

`name` is **mandatory**

The name of the table, is mandatory. If the table is not present in the configuration file it will not be used by pg_expulo.

```code
  "tables": [
     {
       "name": "name",
     }
   ]
```


## Schema

`schema` is **optional**

Contains the schema name where the table is in the postgresql database.

```code
  "tables": [
     {
       "name": "boat",
       "schema": "public"
     }
   ]
```

## Clean

`clean` is **optional**, default value `TRUNCATE`

Define how the table will be purged on target. `clean` option take the following values :

* **truncate**, the `TRUNCATE` sql command will be used

* **delete**, the `DELETE` sql command will be used

* **append**, the table will be not purged, data will be added run after run


```code
  "tables": [
     {
       "name": "boat",
       "schema": "public",
       "clean": "truncate" || "delete" || "append"
     }
   ]
```

## Deletion filter

`deletion_filter` is **optional**

Define the filter to delete data on target, this option works only with `"clean": "delete"`. The string will be added after the `WHERE` keywords in the sql statement.

```code
  "tables": [
     {
       "name": "boat",
       "schema": "public",
       "clean": "delete",
       "deletion_filter": "area != 'South'"
     }
   ]
```