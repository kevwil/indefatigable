# indefatigable – http/2 proxy written in Go

__Indefatigable__ is a very basic http/2 proxy and TLS terminator. I created
this to replace stud/stunnel/pound/etc. and offer http/2 in hopes of having
a performance benefit over standard HTTP/1.1 + TLS.

__*The backend target server will be called using http (unencrypted).*__

SPDY style server push will hopefully be an option, sometime later.

## Options

 * __*bindAddr*__ IP address to bind to for incoming traffic (default "0.0.0.0")
 * __*bindPort*__ Port to listen on (default "443")
 * __*targetAddr*__ Address of target server to proxy traffic for (default "127.0.0.1")
 * __*targetPort*__ Port of target server (default "8080")
 * __*host*__ Host domain name matching TLS cert (default "localhost")
 * __*key*__ TLS key file (default "secure.key")
 * __*cert*__ TLS certificate pem file (default "secure.pem")

###### License

The MIT License

Copyright 2015 Kevin D. Williams
