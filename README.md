A cmd app for generating combination between a content (it can be a password) with totp token. The implementation can be seen in [example](example/) directory.

**Use with your own risk !!!**

# Installing

You can download the (compressed) static binary in [release page](https://github.com/iomarmochtar/content-plus-totp/releases)

## Compiling

- Clone this repository.

- Make sure you have `make` and `go` configured properly.

- Compile it, the static binary output will be created in `dist` directory.

```bash
make compile
```

# Configuration

Run following command for generating json configuration file, the contents are protected using `AES` encryption make sure to remember the key otherwise you cannot decrypt it.

```bash
content-plus-totp -g
```

**note:** You can add `-b` for copy the output directly to clipboard.


# Executing

Generate the combination by run the cmd directly with specified configuration.

```bash
content-plus-totp -c config.json
```

**note:** You can add `-b` for copy the output directly to clipboard.


