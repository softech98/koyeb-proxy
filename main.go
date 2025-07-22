package main

import (
    "fmt"
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
        log.Fatal("Invalid TARGET_URL:", err)
    }

    proxy := httputil.NewSingleHostReverseProxy(remote)

    // Transport setup (optional tuning)
    proxy.Transport = &http.Transport{
        Proxy:               http.ProxyFromEnvironment,
        MaxIdleConns:        100,
        IdleConnTimeout:     90 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
    }

    // === Handler utama proxy ===
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        region := os.Getenv("KOYEB_REGION")
        if region == "" {
            region = "unknown"
        }

        // Tambahkan informasi region ke response header
        w.Header().Set("X-Region", region)

        // Pastikan host tujuan diset dengan benar
        r.Host = remote.Host

        proxy.ServeHTTP(w, r)
    })

    // === Endpoint untuk cek region secara langsung ===
    http.HandleFunc("/region", func(w http.ResponseWriter, r *http.Request) {
        region := os.Getenv("KOYEB_REGION")
        if region == "" {
            region = "unknown"
        }
        clientIP := r.Header.Get("X-Forwarded-For")
    if clientIP == "" {
        clientIP = r.RemoteAddr
    }

    fmt.Fprintf(w, "Client IP: %s\nActive region: %s\n", clientIP, region)
    })

    log.Println("✅ Proxy running on :80 to", target)
    log.Println("🌍 Region:", os.Getenv("KOYEB_REGION"))
    log.Fatal(http.ListenAndServe(":80", nil))
}
