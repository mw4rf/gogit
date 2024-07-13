# gogit

gogit is a CLI tool, written in Go lang, to batch manage multiple git repositories.

- [Configuration](#configuration)
- [Usage](#usage)

## Installation

## Configuration

All the configuration is done in a `config.toml` file.

The location of this file depends on your OS:

- Linux: `~/.config/gogit/config.toml`
- macOS: `~/Library/Application Support/gogit/config.toml`
- Windows: `%APPDATA%\gogit\config.toml`

```toml
[github]
local = "/home/user/git/"

[github.repo1]
remote = "git@github.com:user/repo1.git"

[github.repo2]
remote = "git@github.com:user/repo2.git"
```

- The `local` field specifies the local _absolute_ root path where repositories are located.
- Each repository under a root is specified by its name and remote URL.

An example configuration file is provided: [config.example.toml](config.example.toml).

## Usage

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
