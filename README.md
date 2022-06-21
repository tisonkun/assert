# Assert

## Description

This package is heavily inspired from [stretchr/testify](https://github.com/stretchr/testify), and you can regard it as a fork of the upstream repository.

This package extracts all assertions and republish them with new ones including:

* `(*Assertion).ErrorRegexp`

The `Assert` package servers as a supplement of Golang's `testing` for convenient assertions. And thus I don't want to implement anything like `suite` or `mock`.

* `suite` can be simply implemented leveraging [Golang's Subtests](https://go.dev/blog/subtests).
* `mock` is not a good practice as it's hard to sync logics between the mock and the real object.

## Usage

```shell
go get github.com/tisonkun/assert
```

## Copyright & License

The bundle itself is licensed under the [Apache License](LICENSE).

Copyright 2022 tison wander4096@gmail.com.

You can see all transitive licenses and notices under the [LICENSES](LICENSES) folder.
