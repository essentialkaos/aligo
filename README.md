<p align="center"><a href="#readme"><img src="https://gh.kaos.st/aligo.svg"/></a></p>

<p align="center">
  <a href="https://travis-ci.com/essentialkaos/aligo"><img src="https://travis-ci.com/essentialkaos/aligo.svg?branch=master" alt="TravisCI" /></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-aligo-master"><img alt="codebeat badge" src="https://codebeat.co/badges/18a359f5-50dd-4bfc-95b2-07dee23d018a" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/aligo"><img src="https://goreportcard.com/badge/github.com/essentialkaos/aligo" alt="GoReportCard" /></a>
  <a href="https://github.com/essentialkaos//actions?query=workflow%3ACodeQL"><img src="https://github.com/essentialkaos//workflows/CodeQL/badge.svg" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg" alt="License" /></a>
</p>

<p align="center"><a href="#screenshots">Screenshots</a> • <a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#thanks">Thanks</a> • <a href="#license">License</a></p>

<br/>

`aligo` is a utility for checking and viewing Golang struct alignment info.

### Screenshots

<p align="center">
  <img src="https://gh.kaos.st/aligo-1.png" alt="aligo preview">
  <img src="https://gh.kaos.st/aligo-2.png" alt="aligo preview">
</p>

### Installation

#### From source

Before the initial install, allow git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the `aligo` from scratch, make sure you have a working Go 1.10+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get github.com/essentialkaos/aligo
```

If you want to update `aligo` to latest stable release, do:

```
go get -u github.com/essentialkaos/aligo
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.st/aligo/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) aligo
```

### Usage

```
Usage: aligo {options} {command} package…

Commands

  check    Check package for alignment problems
  view     Print alignment info for all structs

Options

  --arch, -a name      Architecture name
  --struct, -s name    Print info only about struct with given name
  --detailed, -d       Print detailed alignment info (useful with check command)
  --no-color, -nc      Disable colors in output
  --help, -h           Show this help message
  --version, -v        Show version

Examples

  aligo view .
  Show info about all structs in current package

  aligo check .
  Check current package for alignment problems

  aligo -s PostMessageParameters view .
  Show info about PostMessageParameters struct

```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![Build Status](https://travis-ci.com/essentialkaos/aligo.svg?branch=master)](https://travis-ci.com/essentialkaos/aligo) |
| `develop` | [![Build Status](https://travis-ci.com/essentialkaos/aligo.svg?branch=develop)](https://travis-ci.com/essentialkaos/aligo) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### Thanks

We would like to thank:

- [@mdempsky](https://github.com/mdempsky) for [maligned](https://github.com/mdempsky/maligned) utility;
- [@tyranron](https://github.com/tyranron) for [golang-sizeof.tips](http://golang-sizeof.tips/) website.

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
