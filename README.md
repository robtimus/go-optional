# go-optional
[![Build Status](https://github.com/robtimus/go-optional/actions/workflows/build.yml/badge.svg)](https://github.com/robtimus/go-optional/actions/workflows/build.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=robtimus%3Ago-optional&metric=alert_status)](https://sonarcloud.io/summary/overall?id=robtimus%3Ago-optional)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=robtimus%3Ago-optional&metric=coverage)](https://sonarcloud.io/summary/overall?id=robtimus%3Ago-optional)
[![Known Vulnerabilities](https://snyk.io/test/github/robtimus/go-optional/badge.svg)](https://snyk.io/test/github/robtimus/go-optional)

A simple implementation of optionals in Go, based on Java's [Optional](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/Optional.html) class.

Main differences:

* In Go, only pointers can be `nil`. The function provided to the `Map` operation must return a value of the same type, and therefore will only result in an empty `Optional` if the receiver was already empty. The `MapNillable` operation is added that takes a function that returns a pointer to the `Optional`'s generic type.
* Go does not support method generic types. That means that the `Optional` returned by `Map`, `MapNillable` and `FlatMap` operations cannot have a different generic type. To overcome this functions with the same name are provided that take the `Optional` as first argument. That does mean that method chaining is not always possible and needs to be replaced with call chaining:
    ```go
    // o1 is some Optional
    // f1 and f3 has different input and out types
    // f2 has the same input and output types
    o2 := optional.Map(optional.Map(o1, f1).Map(f2), f3)
    ```
* Go does not support method overloading. Java's `orElseThrow` is implemented in three ways:
    * `OrElsePanic` panics if called on an empty `Optional`.
    * `OrElseError` returns a default error if called on an empty `Optional`.
    * `OrElseSupplyError` returns an error provided by a function if called on an empty `Optional`.
* Go does not have the concept of streams the way that Java does. Java's `stream` operation has therefore been replaced by `Slice` that returns a slice with 0 or 1 elements, depending on the `Optional`.
