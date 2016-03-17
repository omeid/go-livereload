# Livereload
[![GoDoc](https://godoc.org/github.com/omeid/livereload?status.svg)](https://godoc.org/github.com/omeid/livereload)

Go implementation of [livereload](http://livereload.com/) and a server.


## Server Implementation
 See [GoDoc](https://godoc.org/github.com/omeid/livereload).

## livereload (cmd/livereload)


### Install

```
go install github.com/omeid/go-livereload/cmd/livereload
```
### Usage

```sh 
$ livereload --help
  Usage of livereload:
  -livereload=":35729": liverloead servera addr.
  -serve="": static files folders.
  -server=":8082": static server addr. Requires -serve 
  -strip="": path to strip from static files.
```



### LICENSE
  [MIT](LICENSE).

[slurp-contrib/livereload](https://github.com/slurp-contrib/livereload/) for Slurp bindings.
