# Generator

The different values supported by `generator` are :


* null

The column will be set to **null**

* ignore

The column will not be present in the `INSERT` statement. As is it permits to use the default value in target table

* mask

The column will be filled with a 8 long char '********'

* randomInt

  A random integer

* md5

The md5 sum will be compute on the data source

* sql

  A SQL function that return value

* FakeEmail

  Generate a fake email address, based on [Faker](https://pkg.go.dev/github.com/go-faker/faker/v4)

* FakeFirstName

  Generate a fake email address, based on [Faker](https://pkg.go.dev/github.com/go-faker/faker/v4)

* FakeName

  Generate a fake email address, based on [Faker](https://pkg.go.dev/github.com/go-faker/faker/v4)