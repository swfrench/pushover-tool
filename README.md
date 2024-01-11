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
