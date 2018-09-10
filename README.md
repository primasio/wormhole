# Wormhole

[![Build Status](https://travis-ci.com/primasio/wormhole.svg?branch=master)](https://travis-ci.com/primasio/wormhole)
[![codecov](https://codecov.io/gh/primasio/wormhole/branch/master/graph/badge.svg)](https://codecov.io/gh/primasio/wormhole)
[![GolangCI](https://golangci.com/badges/github.com/primasio/wormhole.svg)](https://golangci.com)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/primasio/wormhole)

The wormhole, or Einstein–Rosen bridge, connects two different 3D locations from another dimension. It provides a way
to travel through the space in no time.

The Wormhole project, connects the Internet to the Blockchain, or more precisely, the
[Primas Decentralized Network](https://primas.io).

Wormhole is the common part of all the traditional applications that integrate with Primas. It is a showcase of how Primas
API can be used to implement various of use cases. It can be used as the boilerplate for new applications that want to
connect to Primas.

### Centralized Account System

Wormhole is a centralized platform built upon Primas, which encapsulates the cryptographic account to the
outside world and provides a user-friendly traditional account system. Users could sign up using username and password.
They can even sign up with Twitter or Facebook account.

### Integration with Primas SDK

Wormhole connects to [Primas API](https://github.com/primasio/primas-api-doc)
using [Primas Go SDK](https://github.com/primasio/primas-api-sdk-go). It provides the same collection of APIs that
Primas offers in an access token based authentication model where access token is granted by providing the username and
password, and access token can be used to authenticate upcoming API requests. As a working example, Wormhole hosted by
Primas development team is used by Primas browser extension and several other side projects.

Wormhole implements offline signing mechanism to protect its own private key. Other than that, no signature is needed in
any cases for Wormhole users or applications connecting to Wormhole.

A working instance of Wormhole can be found at:

[https://api.connect2.cc](https://api.connect2.cc)

Which is used by the browser extension project called [Connect](https://www.connect2.cc)

Check the [github wiki page](https://github.com/primasio/wormhole/wiki) for API documentation.

### Independent Economic Incentives Model

Wormhole isolates Primas Token, or PST, from its users. Users of Wormhole won't need to know anything about PST.
Instead, an independent token, or point, WORM, is used in the system. How WORMs are used in the system, how WORMs are
distributed to Wormhole users, can WORMs be traded on an exchange, all depend on Wormhole's decision and **CAN BE**
changed at any time.

Wormhole itself, however, needs to hold some amount of PSTs in its root account to use Primas API.

### Development

#### Deploy your own instance

Wormhole is deployed using Docker. You can find the latest released docker image in Docker Hub.

```bash
$ sudo docker pull primasio/wormhole:latest
```

There's also a docker-compose.yml file provided in the release package to show how the docker image is used to start
a Wormhole API server. In most cases, you can simply start the server using only one docker-compose command:

```bash
$ cd wormhole/release/package
$ sudo docker-compose up
```

#### Build the release package

After modification to the source, you need to rebuild the binary and then update the Docker image:

**Rebuild and Test**

There's a Makefile for convenient build and test of Womrhole project. To test the project, simply run:

```bash
$ make test
```

To build the release package run:

```bash
$ make
```

By default the static Linux x64 binary is built.

After rebuild you can find the release package under the folder named "dist".

**Update Docker image**

You can find the required Dockerfile in the release package.
To generate a new docker image:

```bash
$ cd wormhole/project/folder/dist
$ sudo docker build -t primasio/wormhole .
```

You may need to change the target name according to your environment.

### Contribution

Pull request is always welcome.

Join Primas developer community for discussions and instant help:

[https://slack.primas.io](https://slack.primas.io)

