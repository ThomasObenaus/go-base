# Changelog

## v0.0.5 (2020-02-19)

- With [#9](https://github.com/ThomasObenaus/go-base/issues/9) it is now possible to register a shutdown.Stopable even after the shutdown handler was already set uo.

## v0.0.4 (2020-01-15)

- With [#8](https://github.com/ThomasObenaus/go-base/pull/7) it is possible to wait until the sigterm or sigkill signal was issued to the process and until all stoppables were stopped.

## v0.0.3 (2020-01-14)

- With [#7](https://github.com/ThomasObenaus/go-base/pull/7) it is possible to register multiple health checks at once.

## v0.0.2 (2020-01-07)

- Moved from BSD to MIT license

## v0.0.1 (2020-01-06)

- First release contains:
  - build information endpoint
  - configuration support (CLI, ENV and file)
  - health check support
  - logging
  - shutdown handling
