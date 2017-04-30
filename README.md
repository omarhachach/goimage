## goimage - A simple image uploader/server

[![Build Status](https://travis-ci.org/Omar-H/goimage.svg?branch=master)](https://travis-ci.org/Omar-H/goimage)

Goimage is a simple, fully functional, go server for handling image uploads. It is fully standalone, but can be put behind a reverse-proxy.

[Releases](https://github.com/Omar-H/goimage/releases) | [Docs](https://github.com/Omar-H/goimage/wiki)

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
Download the latest [release](https://github.com/Omar-H/goimage/releases) for your platform, and extract the files.

Create a config.json file in the root directory:
```JSON
{
    "port": 8080,
    "secure": false,
    "32-byte-auth-key": "62caed6a7842b5470c2e89693f92c9bab01219f8ebc0c9c0785b97cfd7a68187",
    "allowed-mime-types": [
        "image/x-icon",
        "image/jpeg",
        "image/pjpeg",
        "image/png",
        "image/tiff",
        "image/x-tiff",
        "image/webp",
        "image/gif"
    ],
    "allowed-extensions": [
        ".png",
        ".jpeg",
        ".jpg",
        ".jiff",
        ".png",
        ".ico",
        ".gif",
        ".tif",
        ".webp"
    ],
    "max-file-size": 30000000,
    "image-directory": "public/i/",
    "template-directory": "templates/",
    "public-directory": "public/",
    "image-url": "i/",
    "csrf": false
}
```

Then run the program.
For Linux, in a terminal window, you would run:
```BASH
$ ./goimage
```
On a Windows machine, you would open a command prompt window, and run:
```BASH
C:\goimage-folder> goimage.exe
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
cd goimage/
go run main.go
```

**Go Get**

Note: You need to have GOPATH set.
```BASH
go get github.com/Omar-H/goimage
cd $GOPATH/src/github.com/Omar-H/goimage
go run main.go
```

## Contributing
Please check out our [wiki](https://github.com/Omar-H/goimage/wiki) for more information about contributing.

You can contact the author on Discord: Omar H#6299 or via email: [contact@omarh.net](mailto:contact@omarh.net).

If you feel something is missing, or you find a bug, you can feel free to open an issue or a pull request.
We will check your issue or pull request as soon as possible.
