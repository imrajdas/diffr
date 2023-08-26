# Diffr - Compare Directory Content Differences with Ease


Diffr is an open-source web-based tool designed to make comparing content differences between two directories a simple and intuitive process. Whether you're a developer comparing source code, a designer comparing image assets, or anyone dealing with files, Diffr provides a user-friendly interface to quickly identify changes and similarities between directories.

Visit the project on GitHub: [https://github.com/imrajdas/diffr](https://github.com/imrajdas/diffr)

<img src="./static/images/demo.png">

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Flags](#flags)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Installation

To use Diffr, you need to have [Go](https://golang.org/) installed on your system. Once you have Go set up, you can install Diffr using the following command:

```bash
go get -u github.com/imrajdas/diffr
```

## Usage

Diffr simplifies the process of comparing content differences between two directories. The basic usage is as follows:

```bash
diffr [dir1] [dir2] [flags]
```

You can also use the command to access specific features:

```bash
diffr [command]
```

## Commands

Diffr supports the following commands:

- `help`: Displays help information about any command.
- `version`: Displays the version of Diffr.

## Flags

Diffr provides the following flags to customize its behavior:

- `-a, --address string`: Set the address for the web server to listen on. The default is `http://localhost`.
- `-h, --help`: Display help information about Diffr.
- `-p, --port int`: Set the port for the web server to listen on. The default is `8080`.

## Examples

Here are some examples of how to use Diffr:

```bash
# Compare contents of two directories
diffr /path/to/dir1 /path/to/dir2

# Compare contents with custom server address and port
diffr /path/to/dir1 /path/to/dir2 -a http://127.0.0.1 -p 9000
```

## Contributing

Contributions to Diffr are welcomed and encouraged! If you find a bug or have a feature request, please open an issue on the [GitHub repository](https://github.com/imrajdas/diffr). If you'd like to contribute code, feel free to fork the repository and submit a pull request.

## License

Diffr is released under the [Apache](LICENSE). You are free to use, modify, and distribute this software in accordance with the terms of the license.

---

Diffr makes directory content comparison hassle-free, allowing you to focus on identifying differences rather than dealing with complex tools. Give it a try, and make directory comparison a breeze!