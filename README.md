# go-optional
[![Build Status](https://github.com/robtimus/go-optional/actions/workflows/build.yml/badge.svg)](https://github.com/robtimus/go-optional/actions/workflows/build.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=robtimus%3Ago-optional&metric=alert_status)](https://sonarcloud.io/summary/overall?id=robtimus%3Ago-optional)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=robtimus%3Ago-optional&metric=coverage)](https://sonarcloud.io/summary/overall?id=robtimus%3Ago-optional)

A simple implementation of optionals in Go, based on Java's [Optional](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/util/Optional.html) class.

Main differences:

* In Go, only pointers can be `nil`. The function provided to the `Map` operation must return a value of the same type, and therefore will only result in an empty `Optional` if the receiver was already empty. The `MapNillable` operation is added that takes a function that returns a pointer to the `Optional`'s generic type.
* Go does not support method overloading. Java's `orElseThrow` is implemented in three ways:
    * `OrElsePanic` panics if called on an empty `Optional`.
    * `OrElseError` returns a default error if called on an empty `Optional`.
    * `OrElseSupplyError` returns an error provided by a function if called on an empty `Optional`.
* Go does not have the concept of streams the way that Java does. Java's `stream` operation has therefore been replaced by `Slice` that returns a slice with 0 or 1 elements, depending on the `Optional`.
