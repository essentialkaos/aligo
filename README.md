<p align="center"><a href="#readme"><img src="https://gh.kaos.st/aligo.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/aligo/ci"><img src="https://kaos.sh/w/aligo/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/r/aligo"><img src="https://kaos.sh/r/aligo.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/aligo"><img src="https://codebeat.co/badges/18a359f5-50dd-4bfc-95b2-07dee23d018a" alt="codebeat badge" /></a>
  <a href="https://kaos.sh/w/aligo/codeql"><img src="https://kaos.sh/w/aligo/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#screenshots">Screenshots</a> • <a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#faq">FAQ</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#thanks">Thanks</a> • <a href="#license">License</a></p>

<br/>

`aligo` is a utility for checking and viewing Golang struct alignment info.

### Screenshots

<p align="center">
  <img src="https://gh.kaos.st/aligo-1.png" alt="aligo preview">
  <img src="https://gh.kaos.st/aligo-2.png" alt="aligo preview">
</p>

### Installation

#### From source

To build the `aligo` from scratch, make sure you have a working Go 1.17+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go install github.com/essentialkaos/aligo
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

### FAQ

**Q:** I think my struct is well aligned. How can I disable check for it?

**A:** You could add a comment with text `aligo:ignore` for this struct, and _aligo_ will ignore all problems with it. Example:

```go
// This is my supa-dupa struct
// aligo:ignore
type MyStruct struct {
  A bool
  B int
}
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
| `master` | [![CI](https://kaos.sh/w/aligo/ci.svg?branch=master)](https://kaos.sh/w/aligo/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/aligo/ci.svg?branch=develop)](https://kaos.sh/w/aligo/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### Thanks

We would like to thank:

- [@mdempsky](https://github.com/mdempsky) for [maligned](https://github.com/mdempsky/maligned) utility;
- [@tyranron](https://github.com/tyranron) for `golang-sizeof.tips` website.

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
