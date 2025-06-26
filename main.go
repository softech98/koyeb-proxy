package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "time"
)

func main() {
    target := os.Getenv("TARGET_URL")
    if target == "" {
        log.Fatal("TARGET_URL not set")
    }

    remote, err := url.Parse(target)
    if err != nil {
        log.Fatal("Invalid TARGET_URL: ", err)
    }

    proxy := httputil.NewSingleHostReverseProxy(remote)
    
    transport := &http.Transport{
    Proxy: http.ProxyFromEnvironment,
    MaxIdleConns: 100,
    IdleConnTimeout: 90 * time.Second,
    TLSHandshakeTimeout: 10 * time.Second,
}

proxy.Transport = transport


    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        r.Host = remote.Host
        proxy.ServeHTTP(w, r)
    })

    log.Println("Proxy running on :8080 to", target)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
