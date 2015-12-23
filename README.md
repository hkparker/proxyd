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

`TTPD_CONFIG` is the environment variable containing a JSON encoded string describing services and configuration options.  It is an array of TTPD service config objects, which contain two required fields and two optional fields.

Inside the application, `TTPD_CONFIG` is unmarshaled into a TTPDServicePack.
```
type TLSConfig map[string]string

type TTPDServiceConfig struct {
	Front		string
	Back		string
	FrontConfig	TLSConfig
	BackConfig	TLSConfig
}

type TTPDServicePack []TTPDServiceConfig
```

The required fields "Front" and "Back" define the listen address and the address to proxy connections back to.  These are in the form proto://host:port, such as `tls://0.0.0.0:443`.

Each address has a corresponding optional field, called "FrontConfig" and "BackConfig".  These fields descirbe TLS configuration options and are only required if the corresponding protocol is `tls`.

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
This will bind to 0.0.0.0:443 and present the certificate/key in the `TTPD_DEMO_CERT` and `TTPD_DEMO_KEY` environment variables.  Incomming connections will be proxied to 127.0.0.1:8080.

**Valid TLSConfig options**

`CERT`

Names an environment variable containing PEM encoded cert data to present.

`KEY`

Names an environment variable containing PEM encoded key data for the certificate in `CERT`.

`CIPHERSUITES` (*devel*)

Names an environment variable containing a comma separated list of enabled TLS cipher suites.

`CURVEPREFERENCES` (*devel*)

Names an environment variable containing a comma separated list of ECC curves to preference.

`SERVERNAME` (*devel*)

`ROOT_CAS` (*devel*)

`CLIENT_AUTH_POLICY` (*devel*)

`CLIENT_CAS` (*devel*)

Why?
----

[STUD](https://github.com/bumptech/stud) wasn't maintined, Go is memory safe.

Tests
-----

```
Running Suite: TTPD Suite
=========================
Random Seed: 1450833210
Will run 4 of 4 specs

••••
Ran 4 of 4 Specs in 0.887 seconds
SUCCESS! -- 4 Passed | 0 Failed | 0 Pending | 0 Skipped PASS
coverage: 19.3% of statements
ok      github.com/hkparker/TTPD        0.906s
```

Performance
-----------

[Vegeta](https://github.com/tsenart/vegeta) was used to load test a simple node server on amazon t2.small instance.

**Node directly over TCP**
```
Requests      [total, rate]            6000, 100.02
Duration      [total, attack, wait]    1m0.154433404s, 59.989999949s, 164.433455ms
Latencies     [mean, 50, 95, 99, max]  166.443566ms, 165.695926ms, 168.08417ms, 178.272255ms, 491.472446ms
Bytes In      [total, mean]            1158000, 193.00
Bytes Out     [total, mean]            0, 0.00
Success       [ratio]                  100.00%
Status Codes  [code:count]             200:6000
```

**TTPD with RSA 4096**
```
Requests      [total, rate]            6000, 100.02
Duration      [total, attack, wait]    1m0.073718057s, 59.98999995s, 83.718107ms
Latencies     [mean, 50, 95, 99, max]  91.028404ms, 83.601273ms, 85.360568ms, 447.419508ms, 610.772342ms
Bytes In      [total, mean]            1158000, 193.00
Bytes Out     [total, mean]            0, 0.00
Success       [ratio]                  100.00%
Status Codes  [code:count]             200:6000
```

**TTPD with ECC P521**
```

```

Logging
-------

All logging is done on standard output in JSON format to enable analysis in elasticsearch.

```
{"Now":"2015-12-21 17:05:38.177097504 -0800 PST","Hostname":"Haydens-MacBook-Pro.local","Event":"services_started","Environment":"unknown"}
{"Now":"2015-12-21 17:05:53.044758231 -0800 PST","Service":"tcp://localhost:5000","Hostname":"Haydens-MacBook-Pro.local","Environment":"unknown","Event":"service_down","Error":{"Op":"dial","Net":"tcp","Source":null,"Addr":{"IP":"::1","Port":5000,"Zone":""},"Err":{"Syscall":"getsockopt","Err":61}}}
{"Now":"2015-12-21 17:06:25.881547328 -0800 PST","Service":"tcp://localhost:5000","Hostname":"Haydens-MacBook-Pro.local","Environment":"unknown","Event":"service_recovered"}
```

Alerting
--------

If a slack webhook is specified in the environment variable `TTPD_SLACK_ENDPOINT`, alerts will be sent to the channel specified in the webhook.

![slack example](http://i.imgur.com/gBNHJWf.png)

License
-------

This project is licensed under the MIT license, see LICENSE for more information.
