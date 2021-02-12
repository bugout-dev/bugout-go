# bugout-go

This repository contains the Bugout Go client library. It is also the home of the `bugout` command
line tool.

## Installation

### Pre-built binaries

You can get the latest pre-built release of the `bugout` command line tool on the
[Releases](https://github.com/bugout-dev/bugout-go/releases) page.

### go get

If you are familiar with golang and have it installed, you can also get bugout using:

```bash
go get github.com/bugout-dev/bugout-go
```

This will install the Bugout client library in `$GOPATH/src/github.com/bugout-dev/bugout-go`.

It will also put the `bugout` command line interface in your `$GOPATH/bin` directory.

## Using the bugout command line tool

### Access tokens and the BUGOUT_ACCESS_TOKEN environment variable

Many `bugout` commands require you to pass a Bugout token as the `-t`/`--token` argument. You can
generate an access token by logging into https://bugout.dev/account/tokens.

Once you have generated an access token, if you would like `bugout` to use it automatically without
having to explicitly pass it using `-t`/`--token` every time, you can set it as the
`BUGOUT_ACCESS_TOKEN` environment variable.

If `BUGOUT_ACCESS_TOKEN` is set and you pass a `-t`/`--token` argument, the `-t`/`--token` value
takes precedence.

On a Mac or on Linux:

```bash
export BUGOUT_ACCESS_TOKEN="<access token from https://bugout.dev/account/tokens>"
```

On Windows:

```powershell
setx BUGOUT_ACCESS_TOKEN "<access token from https://bugout.dev/account/tokens>"
```

### Journal IDs and the BUGOUT_JOURNAL_ID environment variable

Some `bugout` commands require you to pass a journal ID using the `-j`/`--journal` argument. If you
find yourself using these commands often, you can set the `BUGOUT_JOURNAL_ID` environment variable
and omit `-j`/`--journal`.

If `BUGOUT_JOURNAL_ID` is set and you pass a `-j`/`--journal` argument, the `-j`/`--journal` value
takes precedence.

On a Mac or on Linux:

```bash
export BUGOUT_JOURNAL_ID="<uuid of bugout journal>"
```

On Windows:

```powershell
setx BUGOUT_JOURNAL_ID "<uuid of bugout journal>"
```
