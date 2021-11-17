package main

import (
  "fmt"
  "net/http"
  "io"
  "log"
  "time"
)

const PORT = "3000"
const DOMAIN = "https://pkg.go.dev" 
// should use net/url.Parse method to check urls are valid

func main() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    proxyUrl := fmt.Sprintf("%s%s", DOMAIN, r.URL.Path)

    if r.URL.RawQuery != "" {
      proxyUrl = fmt.Sprintf("%s?%s", proxyUrl, r.URL.RawQuery)
    } 

    if r.URL.RawFragment != "" {
      proxyUrl = fmt.Sprintf("%s?%s", proxyUrl, r.URL.RawFragment)
    } 

    resp, err := http.Get(proxyUrl)
    if err != nil {
      log.Fatal(err)
    }

    for name, values := range resp.Header {
      for _, value := range values  {
        w.Header().Add(name, value)
      }
    }

    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)

    w.Write(body)
  })

  s := &http.Server{
    Addr: fmt.Sprintf(":%s", PORT),
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  fmt.Println("Listening on port", PORT )
  log.Fatal(s.ListenAndServe())
}
