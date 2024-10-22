package main

import (
    "os"
    "io"
    "fmt"
    "flag"
    "bytes"
    "net/http"
    "mime/multipart"
    "path/filepath"
)

const version = "0.1.0"
const useragent = "Go-http-client/1.1/ (luke at maciak.net) post/" + version

func main() {

    var url string
    flag.StringVar(&url, "url", "https://0x0.st", "URL to upload file to")
    flag.StringVar(&url, "u", "https://0x0.st", "URL to upload file to")

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
            Post(url, args[0])
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
    os.Exit(0)
}

func Post(url string, file string) {

    fmt.Println("\nURL:", url)
    fmt.Println("File:", file)

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
    req, err := http.NewRequest("POST", url, &buff)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error creating HTTP request:", err)
        os.Exit(1)
    }
    defer req.Body.Close()

    // set the content type
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Set("User-Agent", useragent)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error sending the HTTP request:", err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    response, _ := io.ReadAll(resp.Body)
    msg := string(response)

    fmt.Println("Status:", resp.Status)
    fmt.Println("Response:", msg)

}
