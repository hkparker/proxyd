proxyd
======

proxyd proxies data between TCP, TLS, and unix sockets.

A common use case would be to terminate incoming TLS connections and proxy the traffic back to another service listening on a TCP port on localhost.  proxyd could also be used to listen for TCP connections and send the traffic to another host over TLS, where TLS is terminated and traffic is again sent over TCP to the service, creating an encrypted pipe over an untrusted network.

```
TLS termination:
   +---------+
   |         |
->TLS -> TCP |
   |         |
   +---------+

Proxying over an encrypted pipe:
   +---------+      +---------+
   |         |      |         |
->TCP -> TLS-+---> TLS -> TCP |
   |         |      |         |
   +---------+      +---------+
```

Usage
-----

```
$ proxyd --help
Usage of proxyd:
  -config string
        name of configuration file (default "proxyd_config.json")
```

Example
-------

```json
```

Tests
-----

```
```

Logging
-------

All logging is done on standard output in JSON format.

```
```

License
-------

This project is licensed under the MIT license, see LICENSE for more information.
