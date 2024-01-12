# pushover-tool

A simple tool for interacting with the [Pushover API](https://pushover.net/api).

## Status

This tool is highly incomplete in terms of Pushover API endpoints (see below),
and is likely to remain so for the foreseeable future.

## Supported API endpoints

At the moment, only the following are supported:

* Message API: See the `message` subcommand.
* Receipt API: See the `receipt` subcommand.

These alone have been sufficient to meet the vast majority of my intended use
cases (e.g., completion / failure notification for long-running or unattended
operations).

## Token

The tool expects the Pushover API token to be provided via a JSON format file,
containing a single object with a `"token"` field, e.g.

```json
{
    "token": "<YOUR TOKEN>"
}
```

The path to the token file is specified with the `token_path` flag.

## Examples

Send a normal-priority message to `$USER`:

```shell
pushover-tool -token_path ~/.config/pushover/token.json message \
    -message "This is fine" -user $USER
```

Send an emergency-priority message to `$USER` and print the message receipt ID
to stdout:

```shell
pushover-tool -token_path ~/.config/pushover/token.json message \
    -message "My pants are on fire" -user $USER -emergency
```

Wait for the message associated with `$RECEIPT` to be acknowledged, for up to
2m (polling every 10s):

```shell
pushover-tool -token_path ~/.config/pushover/token.json -timeout 2m receipt \
    -receipt $RECEIPT -interval 10s
```

Note that if the `-message` or `-receipt` flags are elided from these commands,
their values will instead be read from stdin (e.g., you can pipe an emergency
`message` to `receipt` in order to await acknowledgement).
