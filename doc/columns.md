# Columns configuration options

## Name

`name` is **mandatory**

The name of the column

```code
  "columns": [
     {
       "name": "name",
     }
   ]
```


## Generator

`generator` is **optional**

Control which data will be inserted in the target table, all values are documented in [Generators](generator.md).

If `generator` is not defined the value in the source table will be used

```code
  "columns": [
     {
       "name": "name",
       "generator": "null"
     }
   ]
```

## Min

`min` is **optional**

When using a generator with parameter, define the minimal value, like `randomIntMinMax`

```code
  "columns": [
     {
       "name": "name",
       "generator": "randomIntMinMax",
       "min": 0,
       "max": 42

     }
   ]
```

## Max

`max` is **optional**

When using a generator with parameter, define the maximal value, like `randomIntMinMax`

```code
  "columns": [
     {
       "name": "name",
       "generator": "randomIntMinMax",
       "min": 0,
       "max": 42
     }
   ]
```

## SQLFunction

`function` is **optional**

When using the generator `sql`, define the SQL function to use. The
function must return a value; for example you may want to use the
function `now()` to fill a column that contains the last update of the
row.

```code
  "columns": [
     {
       "name": "name",
       "generator": "sql",
       "function": "now()"
     }
   ]
```
