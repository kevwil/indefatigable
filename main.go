// Copyright (c) 2015 Kevin D. Williams
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
    "flag"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
    "crypto/tls"
    "io/ioutil"
    "log"
    "time"
    "github.com/bradfitz/http2"
)


var (
    bindAddr string
    bindPort string
    targetAddr string
    targetPort string
    host string
    keyFile string
    certFile string
)

func parseOptions() {
    flag.StringVar(&bindAddr, "bindAddr", "0.0.0.0", "IP address to bind to for incoming traffic")
    flag.StringVar(&bindPort, "bindPort", "443", "Port to listen on")
    flag.StringVar(&targetAddr, "targetAddr", "127.0.0.1", "Address of target server to proxy traffic for")
    flag.StringVar(&targetPort, "targetPort", "8080", "Port of target server")
    flag.StringVar(&host, "host", "localhost", "host domain name matching TLS cert")
    flag.StringVar(&keyFile, "key", "secure.key", "TLS key file")
    flag.StringVar(&certFile, "cert", "secure.pem", "TLS certificate pem file")
    flag.BoolVar(&http2.VerboseLogs, "verbose", false, "Verbose HTTP/2 debugging.")
    flag.Parse()
}

func serveTls() error {
    keyPem, err := ioutil.ReadFile(keyFile)
    if err != nil {
        return err
    }
    certPem, err := ioutil.ReadFile(certFile)
    if err != nil {
        return err
    }
    cert, err := tls.X509KeyPair(certPem, keyPem)
    if err != nil {
        return err
    }
    u, err := url.Parse("http://"+targetAddr+":"+targetPort)
    if err != nil {
        return err
    }
    proxy := httputil.NewSingleHostReverseProxy(u)
    srv := &http.Server{
        Handler: proxy,
        TLSConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
            MinVersion: tls.VersionTLS11,
        },
    }
    http2.ConfigureServer(srv, &http2.Server{})
    ln, err := net.Listen("tcp", bindAddr + ":" + bindPort)
	if err != nil {
		return err
	}
	return srv.Serve(tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, srv.TLSConfig))
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func serveHttp2() error {
    errc := make(chan error, 2)
    go func() { errc <- serveTls() }()
    return <- errc
}

func main() {
    parseOptions()

    log.Fatal(serveHttp2())

    select {}
}
