TTPD
====

TLS and TCP Proxy Daemon

What?
-----

This Go application listens for incoming TLS/TCP connections and proxies them back to another service over TLS/TCP.  The primary use case is to terminate incoming TLS connections and proxy the traffic back over TCP to another service listening on localhost.  TTPD could also be used to listen for TCP connections and send the traffic to another host over TLS, where TLS is terminated and traffic is again sent over TCP to the service, creating an encrypted pipe over an untrusted network.

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

TTPD configuration is done with environment variables.

`TTPD_CONFIG` contains a JSON encoded string describing services and congifuration options.  It is an array of TTPD service config objects, which contain two required fields and two optional fields.

The required fields "Front" and "Back" define the listen address and the address to proxy connections back to.  These are in the form proto://host:port, such as `tls://0.0.0.0:443`.

Each address has a corresponding optional field, called "FrontConfig" and "BackConfig".  These fields descirbe TLS configuration options and are only required with the 

**Example:**

Terminating TLS on 0.0.0.0:443 and proxying back to 127.0.0.1:8080
```
[
	{
		"Front":	"tls://0.0.0.0:443",
		"Back":		"tcp://127.0.0.1:8080",
		"FrontConfig":	{
			"CERT":	"TTPD_DEMO_CERT",
			"KEY":	"TTPD_DEMO_KEY"
		}
	}
]
```
As an environment variable:
```
TTPD_CONFIG="[{"Front":"tls://0.0.0.0:443","Back":"tcp://127.0.0.1:8080","FrontConfig": {}}]"
```

**All TLS configuration options:**

####CERT
If the `CERT` option is an a configuration object it names an environment variable that will contain the PEM encoded cert data to present.

####KEY

**TLS options in development:**
####ROOT_CAS
####SERVER_NAME
####CLIENT_AUTH_POLICY
####CLIENT_CAS
####CIPHER_SUITES
####CURVE_PREFERENCES

Why?
----

[STUD](https://github.com/bumptech/stud) wasn't maintined, Go is memory safe.

Tests
-----

```
```

Performance
-----------

[Vegeta](https://github.com/tsenart/vegeta) was used to load test a simple node server on  

Logging
-------

All logging is done on standard output in JSON to enable analysis in elasticsearch.

Alerting
--------

If a slack webhook is specified in the environment variable `TTPD_SLACK_ENDPOINT`, alerts will be sent to the channel specified in the webhook.



License
-------

This project is licensed under the MIT license, see LICENSE for more information.
