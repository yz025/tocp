# tocp
**TOML Copy** is a CLI tool for managing file/directory synchronization using TOML config files.

tocp allows you to define copy pairs in a simple TOML config and easily sync them with two commands, `push` and `pull`.

Requires Go 1.26 or higher.

## Installation
```
go install github.com/yz025/tocp@latest
```

## Usage
Create a `tocp.toml` file:

```
log = false # default value is true
remove_before_copy = true # default value is true

[[pairs]]
src = "$HOME/.config/nvim" # you can use environment variable, but not tilde(~).
dst = "$DOTFILES/nvim"
# since this pair lack of remove_before_copy, it uses the global value

[[pairs]]
src = "$HOME/.bashrc"
dst = "$DOTFILES/.bashrc"
# for a single file, the old file will be replaced even when remove_before_copy = false.

[[pairs]]
src = "$HOME/Videos/Screencasts"
dst = "$CLIPS"
remove_before_copy = false
```
Then run:
```
# Copy src to dst
tocp push

# Copy dst to src
tocp pull
```
## Configuration
tocp searches for `tocp.toml` in the following order and first found is used.
1. Explicit path via `--path` or `-p` flag `tocp push -p path/to/tocp.toml`
2. Current working directory `./tocp.toml`
3. Custom config directory `$TOCP_CONFIG_HOME/tocp.toml`
4. Global config directory
    - Linux: `~/.config/tocp.toml`
    - Windows: `%APPDATA%\tocp.toml`
## Behaviour
- It will shutdown with an error if it fails to find a config file.
- It will continue to process other pairs even if it fails with some pairs.
- It will log completely nothing if you set `log = false`. 
## Notes
- Not tested on Windows yet.
- Check more examples on [my dotfiles](https://github.com/yz025/dotfiles).
## Todo
- [ ] Improve `--help`.
