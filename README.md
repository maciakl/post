# post

Post a file to [0x0.st](https://0x0.st) and get the URL.

## Usage

Simply run `post filename.ext`.

There are couple of options available:

    Usage: post.exe [options] <file>
    Options:
      -v, --version    Print version information
      -h, --help       Print this message and exit
      -u, --url        URL to upload file to (default: https://0x0.st)

Please make sure you specify the options before the filename.

## Installing

Install via go:
 
    go install github.com/maciakl/post@latest

On Windows, this tool is distributed via `scoop` (see [scoop.sh](https://scoop.sh)).

First, you need to add my bucket:

    scoop bucket add maciak https://github.com/maciakl/bucket
    scoop update

Next simply run:
 
    scoop install post

If you don't want to use `scoop` you can simply download the executable from the release page and extract it somewhere in your path.


