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

### Integration with Primas SDK

Wormhole is a centralized platform built upon Primas, connecting to [Primas API](https://github.com/primasio/primas-api-doc)
using [Primas Go SDK](https://github.com/primasio/primas-api-sdk-go), which encapsulates the cryptographic account to the
outside world and provides a user-friendly traditional account system. Users could sign up using username and password.
They can even sign up with Twitter or Facebook account.

Wormhole implements offline signing mechanism to protect its own private key. Other than that, no signature is needed in
any cases for Wormhole users or applications connecting to Wormhole.

### Centralized Account System

Wormhole provides the same collection of APIs that Primas offers in an access token based authentication model where
access token is granted by providing the username and password, and access token can be used to authenticate upcoming
API requests. As a working example, Wormhole hosted by Primas development team is used by Primas browser extension and
several other side projects.

### Independent Economic Incentives Model

Wormhole isolates Primas Token, or PST, from its users. Users of Wormhole won't need to know anything about PST.
Instead, an independent token, or point, WORM, is used in the system. How WORMs are used in the system, how WORMs are
distributed to Wormhole users, Can WORMs be traded on an exchange, all depend on Wormhole's decision and **CAN BE**
changed at any time.

Wormhole itself, however, needs to hold some amount of PSTs in its root account to use Primas API.
