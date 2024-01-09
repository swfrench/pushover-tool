# pushover-tool

A simple tool for interacting with the [Pushover API](https://pushover.net/api).

## Supported API endpoints

At the moment, only the following are supported:

* Message API: See the `message` subcommand.

## Token

The tool expects the Pushover API token to be provided via a JSON format file,
containing a single object with a `"token"` field.

The path to the token file is specified with the `token_path` flag.
