package main

import (
    "math/rand"
    "flag"
    "fmt"
    "net/http"
    "sync"
    "os"
    "bufio"
    "time"
)

var hostList = []string{}
var hostRoot string
var wg sync.WaitGroup
var timeout_secs = 10 * time.Second

func doRequest(host string) {
        url := "http://" + host
        newHost := host + "." + hostRoot
        // define headers
        headers := map[string]string{
            "Host": newHost,
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) ",
        }
        // create request
        client := &http.Client{
            Timeout: timeout_secs,
        }
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            fmt.Fprintln(os.Stderr, "new request error: ", err)
            wg.Done()
            return
        }
        // attach headers
        for key, value := range headers {
            req.Header.Add(key, value)
        }
        // send request
        resp, err := client.Do(req)
        if err != nil {
            fmt.Fprintln(os.Stderr, "get error: ", err)
            wg.Done()
            return
        }
        defer resp.Body.Close()
        fmt.Print(resp.StatusCode, " ", host, "\n")
        wg.Done()
}


func main() {
    var inputFilename string
    var maxJobs int
    flag.StringVar(&inputFilename, "f", "", "input file name")
    flag.StringVar(&hostRoot, "h", "yoururl.com", "collaborator or interactsh host")
    flag.IntVar(&maxJobs, "j", 8, "max jobs")
    flag.Parse()
    // check input file
    if inputFilename == "" {
        fmt.Fprintln(os.Stderr, "input file name is required")
        // print usage
        flag.Usage()
        os.Exit(1)
    }
    // make sure job count is not negative
    if maxJobs < 0 {
        fmt.Fprintln(os.Stderr, "max jobs must be greater than 0")
        flag.Usage()
        os.Exit(1)
    }
    // read input file
    openFile, err := os.Open(inputFilename)
    if err != nil {
        fmt.Fprintln(os.Stderr, "open file error: ", err)
        flag.Usage()
        os.Exit(1)
    }
    scanner := bufio.NewScanner(openFile)
    for scanner.Scan() {
        hostList = append(hostList, scanner.Text())
    }
    openFile.Close()
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "scanner error: ", err)
        os.Exit(1)
    }
    currentJobs := 0
    rand.Shuffle(len(hostList), func(i, j int) { hostList[i], hostList[j] = hostList[j], hostList[i] })
    for _, host := range hostList {
        if currentJobs >= maxJobs {
            wg.Wait()
            currentJobs = 0
        }
        wg.Add(1)
        currentJobs += 1
        go doRequest(host)
    }
    wg.Wait()
}

