# Crate [![Build Status][travis_img]][travis_href]

**File archival and meta data synchronization tool (experimental)**

[![Crate by Steve][crate.jpg]][crate.jpg]


## About

This project is intended mostly to show off a complete Go project that I wrote myself (my other Go projects are closed source, unfortunately) as part of my larger profile. It should demonstrate a systems program that interacts with the file system, uses a database, and performs network traffic. It is also useful to me, since I'm using it to extract and backup all the images my family has to a cloud service.

I also wanted to have some record of what a _good_ Go project looks like. It is routine for me to setup a quality Python library and deploy it on PyPi - but it's less clear about how to do that for Go. After researching other Go project libraries, this is what I came up with:

- A main.go that interacted with a self-contained library and used [CLI][cli]
- A test suite that uses [Ginkgo][ginkgo] and [Gomega](gomega)
- A Makefile that fetches dependencies and builds locally rather than on the $GOPATH

Together I hope this demonstrates to you and future me how to best set up a Go project, especially as I may be moving away from Go programming back to other languages at least in the immediate future.

### Contributing

Crate is open source, and I would be happy to have you contribute! You can do so in the following ways:

1. Create a Pull Request in Github: [https://github.com/bbengfort/crate](https://github.com/bbengfort/crate)
2. Add issues or bugs to the bug tracker: [https://github.com/bbengfort/crate/issues](https://github.com/bbengfort/crate/issues)
3. Work on a card on the dev board: [https://waffle.io/bbengfort/crate](https://waffle.io/bbengfort/crate)

You can connect with me on Twitter for other discussions: [@bbengfort](https://twitter.com/bbengfort)

### Name Origin

crate<br />
krāt<br />
<small><em>noun</em></small>

1. a slatted wooden case used for transporting or storing goods.

    "a crate of bananas"

2. an old and dilapidated vehicle.

<small><em>verb</em></small>

1. pack (something) in a crate for transportation.

Ok, so I'm a bit of a nerd - but I was thinking about other projects like Box, Dropbox, S3, etc. Crate seemed like an appropriate name. Unfortunately there is already a [crate.io](https://crate.io/) which similarly is an elastic data store. My thing is more about storage and transportation, but hey - it's close! 

### Attribution

Crate does a lot of work on Images, and so it's only fair to specify the attribution of those images. For the test images in the fixtures directory, see the [README.md](fixtures/README.md) there. The header image in this README is attributed as follows:

[Astro Crate](https://flic.kr/p/4LvZAE) by [Steve](https://www.flickr.com/photos/3dking/) is licensed under [CC-BY-NC](https://creativecommons.org/licenses/by-nc/2.0/)

<!-- Link References -->

[travis_img]: https://travis-ci.org/bbengfort/crate.svg
[travis_href]: https://travis-ci.org/bbengfort/crate
[crate.jpg]: fixtures/crate.jpg
[cli]: https://github.com/codegangsta/cli
[ginkgo]: https://github.com/onsi/ginkgo
[gomega]: https://github.com/onsi/gomgea
