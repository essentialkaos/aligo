<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/aligo"><img src="https://kaos.sh/r/aligo.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/y/aligo"><img src="https://kaos.sh/y/be732041f34d4e92a12a28a386b3558a.svg" alt="Codacy badge" /></a>
  <a href="https://kaos.sh/w/aligo/ci"><img src="https://kaos.sh/w/aligo/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/aligo/codeql"><img src="https://kaos.sh/w/aligo/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#screenshots">Screenshots</a> • <a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#faq">FAQ</a> • <a href="#usage">Usage</a><br/><a href="#ci-status">CI Status</a> • <a href="#contributing">Contributing</a> • <a href="#thanks">Thanks</a> • <a href="#license">License</a></p>

<br/>

𝑎𝑙𝑖𝑔𝑜 is a utility for checking and viewing Golang [struct alignment info](https://medium.com/@codewithkushal/understanding-struct-padding-in-go-in-depth-guide-ed70c0432c63).

### Introduction

Struct alignment is the extra space added between fields in a struct to align them in memory according to the CPU's word size.

By understanding and managing struct padding (e.g., reordering fields), you can improve program performance, reduce memory usage, and ensure data integrity.

This tool aims to provide a visual way, but also a GitHub Action to report possible improvements in your code.

You can refer [@codingwithkushal][github-codingwithkushal]'s [Go struct alignment issues][medium-codewithkushal-struct-padding] article, if you want to understand further.

[github-codingwithkushal]: https://github.com/codingwithkushal
[medium-codewithkushal-struct-padding]: https://medium.com/@codewithkushal/understanding-struct-padding-in-go-in-depth-guide-ed70c0432c63

### Screenshots

<p align="center">
  <img src=".github/images/screenshot1.png" alt="aligo preview">
  <img src=".github/images/screenshot2.png" alt="aligo preview">
</p>

### Installation

#### From source

To build the _aligo_ from scratch, make sure you have a working [Go 1.22+](https://github.com/essentialkaos/.github/blob/master/GO-VERSION-SUPPORT.md) workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/aligo/v2@latest
```

#### Using with Github Actions

For using _aligo_ with GitHub Actions use this workflow file or add job `Aligo` to your workflow:

```yml
name: Aligo

on:
  push:
    branches: [master, develop]
  pull_request:
    branches: [master]

jobs:
  Aligo:
    name: Aligo
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: '1.22.x'

      - name: Check Golang sources with Aligo
        uses: essentialkaos/aligo-action@v2
        with:
          files: ./...
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

You can generate man page for _aligo_ using next command:

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

<img src=".github/images/usage.svg" />

### CI Status

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
