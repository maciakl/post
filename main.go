package main

import (
    "os"
    "io"
    "fmt"
    "flag"
    "bytes"
    "net/url"
    "net/http"
    "mime/multipart"
    "path/filepath"
)

const version = "0.1.1"
const useragent = "Go-http-client/1.1/ (luke at maciak.net) post/" + version

func main() {

    var urlFlag string
    flag.StringVar(&urlFlag, "url", "https://0x0.st", "URL to upload file to")
    flag.StringVar(&urlFlag, "u", "https://0x0.st", "URL to upload file to")

    var proxy string
    flag.StringVar(&proxy, "proxy", "", "HTTP proxy to use")
    flag.StringVar(&proxy, "p", "", "HTTP proxy to use")

    var ver bool
    flag.BoolVar(&ver, "version", false, "show version")
    flag.BoolVar(&ver, "v", false, "show version")

    flag.Usage = Usage
    flag.Parse()

    if ver {
        Version()
    }
    
    // get the non-flag arguments
    args := flag.Args()

    if len(args) > 0 {
        switch args[0] {

        case "-h", "-help":
            Usage()

        default:
            Post(urlFlag, args[0], proxy)
        } 
    } else {
        Usage()
    }
    
}

func Version() {
    fmt.Println(filepath.Base(os.Args[0]), "version", version)
}

func Usage() {
    fmt.Println("Usage:", filepath.Base(os.Args[0]), "[options] <file>")
    fmt.Println("Options:")
    fmt.Println("  -v, --version    Print version information")
    fmt.Println("  -h, --help       Print this message and exit")
    fmt.Println("  -u, --url        URL to upload file to (default: https://0x0.st)")
    fmt.Println("  -p, --proxy      HTTP proxy to use (eg. http://user:pass@host:port)")
    os.Exit(0)
}

func Post(urlPost string, file string, proxy string) {

    fmt.Println("\nURL:", urlPost)
    fmt.Println("File:", file)
    if proxy != "" {
        fmt.Println("Proxy:", proxy)
    }

    // open file
    file_reader, err := os.Open(file)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error opening file:", err)
        os.Exit(1)
    }
    defer file_reader.Close()
    
    // create buffer
    var buff bytes.Buffer
    writer := multipart.NewWriter(&buff)

    // create form field named file and attach the file to it
    form, err := writer.CreateFormFile("file", filepath.Base(file))
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error creating a form:", err)
        os.Exit(1)
    }

    // copy the file into form
    written_bytes, err := io.Copy(form, file_reader)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error copying the file:", err)
        os.Exit(1)
    }
    fmt.Println("Bytes:", written_bytes)
    fmt.Println("\n--- Response ---")

    // have to explicitly close the writer instead of deferring
    // because we need to get the boundary string
    writer.Close()

    // create a new request
    req, err := http.NewRequest("POST", urlPost, &buff)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error creating HTTP request:", err)
        os.Exit(1)
    }
    defer req.Body.Close()

    // set the content type
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Set("User-Agent", useragent)

    // create a new http client
    client := &http.Client{}

    // if a proxy is specified, create a new transport and set the proxy
    if proxy != "" {
        proxyUrl, err := url.Parse(proxy)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Error parsing proxy URL:", err)
            os.Exit(1)
        }
        transport := &http.Transport{
            Proxy: http.ProxyURL(proxyUrl),
        }
        client.Transport = transport
    }

    resp, err := client.Do(req)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error sending the HTTP request:", err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    response, _ := io.ReadAll(resp.Body)
    msg := string(response)

    fmt.Println("Status:", resp.Status)
    for k, v := range resp.Header {
        fmt.Printf("%s: %s\n", k, v)
    }
    fmt.Println("Response:", msg)

}
