proxyd
======

proxyd proxies data between TCP, TLS, and unix sockets

```
TLS termination:     Connecting to a remote application's unix socket:
   +---------+         +----------+        +----------+
   |         |         |          |        |          |
->TLS -> TCP |         | -> unix -+------>TLS -> unix |
   |         |         |          |        |          |
   +---------+         +----------+        +----------+
                                                                     _
Proxying through an encrypted pipe:                                 | |
                                       _ __  _ __ _____  ___   _  __| |
   +---------+      +---------+       | '_ \| '__/ _ \ \/ / | | |/ _` |
   |         |      |         |       | |_) | | | (_) >  <| |_| | (_| |
->TCP -------+---> TLS -> TCP |       | .__/|_|  \___/_/\_\\__, |\__,_|
   |         |      |         |       | |                   __/ |
   +---------+      +---------+       |_|                  |___/
```

Usage
-----

```
$ proxyd --help
Usage of proxyd:
  -config string
        name of configuration file (default "proxyd_config.json")
```

proxyd takes a file with a configuration object that describes listeners and proxies.  The object is in the form of a [`ServicePack`](https://github.com/hkparker/proxyd/blob/master/services.go#L17).

Examples
--------

Listen on a port and forward to another host

```json
[
	{
		"Front": "tcp://0.0.0.0:80",
		"Back": "tcp://192.168.1.2:8080",
	}
]
```

Accept TCP connections and forward to a local TLS server without verification

```json
[
	{
		"Front": "tcp://0.0.0.0:80",
		"Back": "tls://127.0.0.1:443",
		"BackConfig": {
			"InsecureSkipVerify": "true"
		}
	}
]
```

Terminate TLS

```json
[
	{
		"Front": "tls://0.0.0.0:443",
		"Back": "tcp://127.0.0.1:80",
		"FrontConfig": {
			"CERT": "-----BEGIN CERTIFICATE-----\n...",
			"KEY": "-----BEGIN RSA PRIVATE KEY-----\n..."
		}
	}
]
```

Terminate TLS with client certificate verification

```json
```

Serve the same unix socket on 80 and 443 with TLS

```json
[
	{
		"Front": "tcp://0.0.0.0:80",
		"Back": "unix:///tmp/ipc.sock"
	},
	{
		"Front": "tls://0.0.0.0:443",
		"Back": "unix:///tmp/ipc.sock",
		"FrontConfig": {
			"CERT": "-----BEGIN CERTIFICATE-----\n...",
			"KEY": "-----BEGIN RSA PRIVATE KEY-----\n..."
		}
	}
]
```

Supported TLS features
----------------------



Tests
-----

```
$ go test -race -cover .
ok      github.com/hkparker/proxyd      1.016s  coverage: 59.6% of statements
```
