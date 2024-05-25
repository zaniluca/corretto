# Corretto

Corretto (Italian for "Free from errors")[^1] is a simple but powerful schema validation package for go structs. It is designed with code readability in mind.

The library is inspired by popular JavaScript libraries like [Yup](https://github.com/jquense/yup) and [Zod](https://github.com/colinhacks/zod), and it aims to provide a similar experience in Go.

## Getting Started

```bash
go get github.com/zaniluca/corretto
```

To create a schema, use the `Schema` constructor, and define the shape of the data you expect. Each field in the schema can have a number of methods attached to it, which define how the field should be parsed and validated.

```go
package main

import (
    c "github.com/zaniluca/corretto"
)

type User struct {
    FirstName string
    Age       int
    BirthDate string
    Email     string
}

func main() {
    user := &User{
      FirstName: "John",
      Age:       17,
      BirthDate: "2007-01-01",
      Email:     "john@doe.com",
    }

    schema := c.Schema{
      "FirstName": c.Field("Name").Min(3),
      "Age":       c.Field().Min(18),
      "BirthDate": c.Field().Required(),
      "Email":     c.Field().Email(),
    }

    // ❌ ValidationError{Message: "Age must be at least 18"}
    err := schema.Parse(user)
}

```

**Table of Contents**

- [The Schema](#the-schema)
  - [Validation](#validation)
  - [Custom validations](#custom-validations)
    - [Customizing errors](#customizing-errors)
  - [Composition and Reuse](#composition-and-reuse)
  - [Primitive Validators](#primitive-validators)
  - [Nested Schemas](#nested-schemas)
- [License](#license)

## The Schema

Corretto's `Schema` is nothing more than a map of fields to their respective validation rules. Each Field in the schema, which must be explicitly declared with the `Field()` func can have a number of methods attached to it, which define how the field should be parsed and validated.

### Validation

The core of a validation schema is to check that a given value conforms to a set of rules. This is done by calling the `Parse` method on the schema, which will return an error if the value does not conform to the schema.

```go
err := schema.Parse(user)
if err != nil {
    log.Println(err)
}
```

And if you want you can work with **Json** data too

```go
json := byte[]`{"FirstName":"John","Age":17,"BirthDate":"2007-01-01", ...}`
u := User{}
// Unmarshal the json data into `u` and validate it
// If both the json data and the struct are valid, `u` will be filled with the data
// Otherwise, `err` will contain the error
err := schema.Unmarshal(json, u)
```

> There are also `MustParse` and `MustUnmarshal` methods that will panic if the value does not conform to the schema.

### Custom validations

TODO

#### Customizing errors

You can customize the field name in the error message by passing it as an argument to the `Field()` func.

```go
// The error message will be "Name must be at least 3 characters long"
s := c.Schema{
    "FirstName": c.Field("Name").Min(3),
    // ...
}
```

If you want to customize the entire error message, you can pass a second argument to most validation methods.

```go
// The error message will be "Name not long enough"
s := c.Schema{
    "FirstName": c.Field("Name")
                    .Min(3, "%v not long enough (min %v)"),
  // ...
}
```

> As you can see `Min` accepts passing a string with placeholders like you do in the `fmt` package. The first placeholder will be replaced with the field name, and the second with the value of the `Min(3)` method (in this case, 3), if the method has more than one argument or none it will have an according number of placeholders.

### Composition and Reuse

Schemas can be composed and reused in a number of ways. The most common is to use the `Field()` func to define a field and its validation rules, and then reuse that field in multiple schemas.

```go
nameValidator := c.Field("Name").Min(3)
```

Another way to reuse schemas is to use the `Schema` constructor to define a schema and then use `Concat()` to combine it with other schemas.

```go
nameSchema := c.Schema{
    "FirstName": nameValidator,
}

userSchema := nameSchema.Concat(c.Schema{
    "Age": c.Field().Min(18),
})
```

> Note: if you're concatenating two schemas that have the same field, the field in the second schema will override the field in the first schema.

### Primitive Validators

Not all validations are designed to be used with all types, for example, the `Email` validation should only be applied to strings.

To enforce this corretto offers a set of **Primitive Validators** that can be used to restrict the types of values that can be validated.

```go
schema := c.Schema{
    "Email": c.Field().String().Email(), // Will error if the value is not a string (and also if it's not a valid email)
}
```

When applying a primitive validator to a field, the field will be restricted to the type of the primitive validator, and only the methods that are valid for that type will be available. For example, if you apply the `Email()` validator to a field, you will only be able to use string methods (and base ones like `Required()`) on that field.

This means that you'll get a compile-time error if you try to use a method that is not valid for the type of the field. _(and also methods suggestions from your IDE)_

Currently we support `String()` and `Array()` as primitive validators.

### Nested Schemas

Schemas can be used to validate nested structs. Let's say you have a `User` struct that contains an `Address` struct.

```go
type Address struct {
    Street string
    City   string
}

type User struct {
    FirstName string
    LastName  string
    Address   Address
}
```

You can define a schema for the `Address` struct and then use that schema in the schema for the `User` struct.

```go
addressSchema := c.Schema{
    "Street": c.Field().String().Required(),
    "City":   c.Field().String().Required(),
}

userSchema := c.Schema{
    "FirstName": c.Field().String().Required(),
    "LastName":  c.Field().String().Required(),
    "Address":   c.Field().Schema(addressSchema),
}
```

> Note: in this case `Address` was an **exported** field, if it was unexported the validator would not be able to access it and will panic.

## Package

### `corretto`

#### `Schema`

This simply an alias for `map[string]*Validator`

#### `Field(fieldName ...string) *Validator`

Returns a new `*Validator` with the given field name. If no field name is provided, the field name will be the same as the struct field name.

#### `(s Schema) Parse(value any) error`

Checks that the given value conforms to the schema. If the value does not conform to the schema, an error is returned.

If the struct contains a field that is not present in the schema it will **panic**.

> Note: This method will not modify the value, it will only check that it conforms to the schema.

#### `(s Schema) MustParse(value any)`

Like `Parse`, but panics if the value does not conform to the schema.

#### `(s Schema) Concat(other Schema)`

Returns a new `Schema` that is the combination of the two schemas. If a field is present in both schemas, the field from the second schema will be used.

#### `(s Schema) Unmarshal(data []byte, value any) error`

Unmarshals the given JSON data into the given value and checks that it conforms to the schema. If the value does not conform to the schema, an error is returned.

### Mixed

A list of all the methods that can be called on a `*Validator` with no specific type.

> Mixed methods are available for most of the basic types like `int`, `string`, `float`, etc. If a type is not supported they will **panic**

#### `(*Validator) Required(msg ...string) *Validator`

Checks that the value is not `nil` or the zero value for the type.

#### `(*Validator) Min(min int, msg ...string) *Validator`

Checks that the value is greater than or equal to the given minimum value.

#### `(*Validator) Max(max int, msg ...string) *Validator`

Checks that the value is less than or equal to the given maximum value.

### Strings

Methods that can be called on a `*Validator` applied to a `string`

#### `(v *Validator) Matches(regex string, msg ...string) *Validator`

Matches checks if the field matches the provided regex pattern

if the string is empty, it will return true, use `Required()` to check for empty strings

Corretto also provides some predefined regex validations:

#### `(v *Validator) Email(msg ...string) *Validator`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

[^1]: "Corretto" is also used as a term for a type of coffee in Italy, one that is "corrected" with a shot of liquor.
