## goimage - A simple image uploader/server
[![Travis branch](https://img.shields.io/travis/omar-h/goimage/master.svg?style=flat-square)](https://travis-ci.org/omar-h/goimage)
[![GitHub tag](https://img.shields.io/github/release/omar-h/goimage.svg?style=flat-square)](https://github.com/omar-h/goimage/releases)
[![Report Card](https://img.shields.io/badge/report%20card-a%2B-c0392b.svg?style=flat-square)](https://goreportcard.com/report/github.com/omar-h/goimage)
![Powered By](https://img.shields.io/badge/powered%20by-go-blue.svg?style=flat-square)
[![License](https://img.shields.io/badge/license-MIT%20License-1abc9c.svg?style=flat-square)](https://github.com/omar-h/goimage/blob/master/LICENSE.txt)

Goimage is a simple, fully functional, go server for handling image uploads. It is fully standalone, but can be put behind a reverse-proxy.

[Releases](https://github.com/omar-h/goimage/releases) | [Docs](https://godoc.org/github.com/omar-h/goimage)

## Menu
* [Features](#features)
* [Getting Started](#getting-started)
* [Benchmarks](#benchmarks)
* [Running from Source](#running-from-source)
* [Contributing](#contributing)

## Features
* Simple and easy configuration
* CSRF Protection
* Performant and efficient
* More coming soon..

## Getting Started
**Requires Go 1.9**

Download the latest [release](https://github.com/omar-h/goimage/releases) for your platform, and extract the files.

Create a config.json file in the root directory:
```JSON
{
    "port": 8080,
    "image-name-length": 4,
    "max-file-size": 1,
    "allowed_mime_types": [
        "image/x-icon",
        "image/jpeg",
        "image/pjpeg",
        "image/png",
        "image/tiff",
        "image/x-tiff",
        "image/webp",
        "image/gif"
    ],
    "allowed_extensions": [
        "png",
        "jpeg",
        "jpg",
        "jiff",
        "ico",
        "gif",
        "tif",
        "webp"
    ],
    "directories": {
        "image": "public/",
        "template": "templates/",
        "public": "public/"
    },
    "csrf": {
        "enabled": true,
        "32_byte_auth_key": "62caed6a7842b5470c2e89693f92c9ba",
        "secure": false,
        "httpOnly": true
    }
}
```

Then run the program.
For Linux, in a terminal window, you would run:
```BASH
$ ./goimage
```
On a Windows machine, you would open a command prompt window, and run:
```BASH
> goimage.exe
```

## Benchmarks
Coming soon..

## Running from Source
You can do this two ways; either by using go get, or git clone.

**Not recommended for production use.**

Note: Requires Go installed on the system.

**Git Clone**

```BASH
git clone https://github.com/Omar-H/goimage.git
cd goimage/cmd/goimage
go build .
./goimage
```

**Go Get**

Note: You need to have GOPATH set.
```BASH
go get github.com/Omar-H/goimage
cd $GOPATH/src/github.com/Omar-H/goimage/cmd/goimage
go build .
./goimage
```

## Contributing
Please check out our [wiki](https://github.com/omar-h/goimage/wiki) for more information about contributing.

You can contact the author on Discord: Omar H.#6299 or via email: [contact@omarh.net](mailto:contact@omarh.net).

If you feel something is missing, or you find a bug, you can feel free to open an issue or a pull request.
We will check your issue or pull request as soon as possible.
