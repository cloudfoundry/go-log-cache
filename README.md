# Golang Log Cache Client

> **NOTE**: This client library is only compatible with Log Cache versions prior to `2.0.0` and will no longer receive new features. Log Cache `2.0.0` and above will be supported by the [`client` package found within `log-cache` itself][new-client]. That client is also backwards compatible with the `1.4.x` releases of Log Cache.

[![GoDoc][go-doc-badge]][go-doc] [![travis][travis-badge]][travis] [![slack.cloudfoundry.org][slack-badge]][log-cache-slack]

This is a golang client library for [Log-Cache][log-cache].

## Usage

This repository should be imported as:

`import logcache "code.cloudfoundry.org/go-log-cache"`

[slack-badge]:              https://slack.cloudfoundry.org/badge.svg
[log-cache-slack]:          https://cloudfoundry.slack.com/archives/log-cache
[log-cache]:                https://code.cloudfoundry.org/log-cache
[go-doc-badge]:             https://godoc.org/code.cloudfoundry.org/go-log-cache?status.svg
[go-doc]:                   https://godoc.org/code.cloudfoundry.org/go-log-cache
[travis-badge]:             https://travis-ci.org/cloudfoundry/go-log-cache.svg?branch=master
[travis]:                   https://travis-ci.org/cloudfoundry/go-log-cache?branch=master
[new-client]:               https://github.com/cloudfoundry/log-cache/tree/master/pkg/client
