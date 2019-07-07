Timetracker
===========

This is a simple time tracker tool based on wakatime. It does not send
the data to a remote service - it simply writes it into a local log
file (by default ~/.timetracker.log but can be overridden by setting
the environment string VELOTRACKER_LOG).

To use it simply install the wakatime IDE plugin, and copy this binary
to your path named as /usr/bin/wakatime. This binary has the same
command line flags as the wakatime binary and so it is compatible with
the IDE plugins.
