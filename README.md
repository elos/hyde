hyde [![GoDoc](https://godoc.org/github.com/elos/hyde?status.svg)](https://godoc.org/github.com/elos/hyde)
----

Package hyde provides a engine for generating static documentation sites.

#### Architecture
The architecture of hyde is quite simple. At the core there are `hyde.Engine`s, which
server as an abstraction over a filesystem server, providing some hyde-specific
information like their root directory. They recursively inspect all non-hidden
files.

One step up from a `hyde.Engine` is a `hyde.Pod`. Each pod has an engine, and manages its rendered
templates and serves http responses. At the top level there is a hull, which is a colleciton
of documentation pods. An example would be the elos documentation, which is itself a hull, but
has a colleciton of pods of documentation distributed throughout the elos repositories. Elos builds
an elos documentaiton hull and then attaches pods of api, app, server, and culture documentation
to the hull.
