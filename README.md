Sit
===

SSL is terminated

What?
-----

This Go application listens for incoming TLS connections and proxies them back to another service over TCP.

Why?
----

STUD wasn't maintined, Go is memory safe.

Tests
-----

```
```

Coverage is not accurate (too low), as much of the testing is done on the compiled binary and coverage checks are not used when building with gexec.

Performance
-----------

License
-------

This project is licensed under the MIT license, see LICENSE for more information.
