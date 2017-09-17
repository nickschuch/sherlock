Sherlock
========

**Maintainer**: Nick Schuch

When a Pod is murdered, Sherlock isn't far away to solve the mystery.

## Components

* Watson - Daemon for storing data
* Sherlock - Command line for developers

## How it works

When a Pod exits, Watson will store data in S3:

* Object
* Logs
* Events

This data can then be loaded by `sherlock`.

**Listing**

```bash
$ sherlock list
INCIDENT ID             	TIMESTAMP                    	NAMESPACE  	POD                            	CONTAINER
eNVFSnCJdXEHYyQJsNjTSKDN	2017-09-17 06:41:52 +0000 UTC	kube-system	project1-3908849101-nmxpx	exporter 
qzTvepsIRXQXkEMfbJNeJDDI	2017-09-17 06:42:27 +0000 UTC	kube-system	project1-3908849101-nmxpx	exporter
```

**Inspecting**

```bash
$ sherlock inspect eNVFSnCJdXEHYyQJsNjTSKDN

###########################################################################
events.log
###########################################################################

......................

###########################################################################
object.yaml
###########################################################################

......................

###########################################################################
output.log
###########################################################################

......................
```

## Resources

* [Dave Cheney - Reproducible Builds](https://www.youtube.com/watch?v=c3dW80eO88I)

## Development

### Principles

* Code lives in the `workspace` directory

### Tools

* **Dependency management** - https://getgb.io
* **Build** - https://github.com/mitchellh/gox
* **Linting** - https://github.com/golang/lint

### Workflow

(While in the `workspace` directory)

**Installing a new dependency**

```bash
gb vendor fetch github.com/foo/bar
```

**Running quality checks**

```bash
make lint test
```

**Building binaries**

```bash
make build
```
