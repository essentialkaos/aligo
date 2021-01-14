<p align="center"><a href="#readme"><img src="https://gh.kaos.st/aligo.svg"/></a></p>

<p align="center">
  <a href="https://github.com/essentialkaos/aligo/actions"><img src="https://github.com/essentialkaos/aligo/workflows/CI/badge.svg" alt="GitHub Actions Status" /></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-aligo-master"><img alt="codebeat badge" src="https://codebeat.co/badges/18a359f5-50dd-4bfc-95b2-07dee23d018a" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/aligo"><img src="https://goreportcard.com/badge/github.com/essentialkaos/aligo" alt="GoReportCard" /></a>
  <a href="https://github.com/essentialkaos/aligo/actions?query=workflow%3ACodeQL"><img src="https://github.com/essentialkaos/aligo/workflows/CodeQL/badge.svg" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#screenshots">Screenshots</a> • <a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#thanks">Thanks</a> • <a href="#license">License</a></p>

<br/>

`aligo` is a utility for checking and viewing Golang struct alignment info.

### Screenshots

<p align="center">
  <img src="https://gh.kaos.st/aligo-1.png" alt="aligo preview">
  <img src="https://gh.kaos.st/aligo-2.png" alt="aligo preview">
</p>

### Installation

#### From source

To build the `aligo` from scratch, make sure you have a working Go 1.12+ workspace (_[instructions](https://golang.org/doc/install)_), then:

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

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo aligo --completion=bash 1> /etc/bash_completion.d/aligo
```


ZSH:
```bash
sudo aligo --completion=zsh 1> /usr/share/zsh/site-functions/aligo
```


Fish:
```bash
sudo aligo --completion=fish 1> /usr/share/fish/vendor_completions.d/aligo.fish
```

### Man documentation

You can generate man page for aligo using next command:

```bash
aligo --generate-man | sudo gzip > /usr/share/man/man1/aligo.1.gz
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
| `master` | ![CI](https://github.com/essentialkaos/aligo/workflows/CI/badge.svg?branch=master) |
| `develop` | ![CI](https://github.com/essentialkaos/aligo/workflows/CI/badge.svg?branch=develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### Thanks

We would like to thank:

- [@mdempsky](https://github.com/mdempsky) for [maligned](https://github.com/mdempsky/maligned) utility;
- [@tyranron](https://github.com/tyranron) for [golang-sizeof.tips](http://golang-sizeof.tips/) website.

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
