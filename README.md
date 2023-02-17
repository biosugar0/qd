# qd

Take Obsidian Daily Notes Easily


## Usage

```
NAME:
   qd - Quick Daily Notes

USAGE:
   qd [section title]

VERSION:
   0.0.5

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## config

write `~/.config/qd/config.toml`

 `dailynotedir` is required.

```toml
dailynotedir = "~/memo/daily"
editor = "vim"
```


## example

you run below command,

```
$qd
```

Then, today's daily note will be created and the file will open.


```md
[[daily notes]]

[[2022-07-09]] | [[2022-07-11]]

### 07:18
```

You can write freely.

A short time later, you have breakfast.

```
$qd "have a breakfast"
```

Then the following statement will be added to the daily note.

```md
[[daily notes]]

[[2022-07-09]] | [[2022-07-11]]

### 07:18 Good Morning

I'm sleepy...

### 08:13 have a breakfast
```
