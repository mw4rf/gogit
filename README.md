# gogit

gogit is a CLI tool, written in Go lang, to batch manage multiple git repositories.

- [Configuration](#configuration)
- [Usage](#usage)

## Configuration

All the configuration is done in a `repos.json` file.

The location of this file depends on your OS:

- Linux: `~/.config/gogit/repos.json`
- macOS: `~/Library/Application Support/gogit/repos.json`
- Windows: `%APPDATA%\gogit\repos.json`

The configuration files uses the very simple JSON syntax to declare the repositories.

```json
{
    "name": "Ventanas",
    "local": "/home/bill/worlddomination/git/ventanas",
    "remote": "git@gitpuertas.com:bill/ventanas.git"
},
{
    "name": "AdjectiveAnimal"
    "local": "/home/bill/worldemancipation/git/adjectiveanimal",
    "remote": "git@freeforall.org:bill/adjectiveanimal.git"
}
```

- The `local` field specifies the local _absolute_ root path where repository is located.
- The `remote` field specifies the URL to the remote git repository.

If you already have a folder, let's say `~/git`, with a bunch of cloned repos, you can generate the `repos.json` file with the `genrepos` command.

``` sh
# Print JSON content to screen
gogit genrepos ~/git

# Generate repos.json
gogit genrepos ~/git > ~/.config/gogit/repos.json
```

## Usage

``` sh
Usage: gogit <command> [arguments]
Commands:
  list                  List the repositories in a simple and compact format
  list full             List the repositories in a detailed format
  genrepos [root]       Generate and print a JSON string with the details of all git repositories in a given root folder
  clone                 Check all repositories and clone the ones that are missing
  help                  Print this help message
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Todo

- [x] Clone absent repos.
- [ ] Fetch/Pull/Push all repos.
- [ ] Fetch/Pull/Push a single repo (by name).
- [ ] Show repos behind/ahead of remote, or with changes to be committed.
- [x] Write help output.
- [ ] Make a log-compliant output to run gogit as a cronjob.
