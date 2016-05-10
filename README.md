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

Configuration
-------------


**Example:**

```json
[
    {
        "Front":    "tls://0.0.0.0:443",
        "Back":     "unix:///path/to/file",
        "FrontConfig":  {
            "CERT": "----BEGIN-RSA-.....",
            "KEY":  "----BEING-RSA-......"
        }
    },
    {
        "Front":    "tcp://0.0.0.0:443",
        "Back":     "tcp://127.0.0.1:8080"
    },
    {
        "Front":    "tls://0.0.0.0:443",
        "Back":     "tls://",
        "FrontConfig":  {
            "CERT": "DEMO_CERT",
            "KEY":  "DEMO_KEY"
        },
        "FrontConfig":  {
            "CERT": "DEMO_CERT",
            "KEY":  "DEMO_KEY"
        }
    }
]
```

Tests
-----

```
```

Performance
-----------


Logging
-------

All logging is done on standard output in JSON format.

```
```

License
-------

This project is licensed under the MIT license, see LICENSE for more information.
