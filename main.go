package main

import (
    "fmt"
    "net/http"
    "net/url"
    "sync"
    "time"
)

const endpoint = "https://nyanpass.com/add.php"
const maxThreads = 200

func main() {
    var reqCount int
    fmt.Println("Auto Nyanpasser")
    fmt.Print("にゃんぱすーする回数を指定(0なら無限)：")
    fmt.Scan(&reqCount)

    if reqCount == 0 {
        sendUnlimitedRequests()
    } else {
        sendLimitedRequests(reqCount)
    }
}

func sendUnlimitedRequests() {
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, maxThreads)

    for i := 0; ; i++ {
        wg.Add(1)
        semaphore <- struct{}{}
        go func(i int) {
            defer wg.Done()
            sendRequest(i)
            <-semaphore
        }(i)
    }
    wg.Wait()
}

func sendLimitedRequests(reqCount int) {
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, maxThreads)

    for i := 0; i < reqCount; i++ {
        wg.Add(1)
        semaphore <- struct{}{}
        go func(i int) {
            defer wg.Done()
            sendRequest(i)
            <-semaphore
        }(i)
    }
    wg.Wait()
}

func sendRequest(i int) {
    data := url.Values{}
    data.Set("nyan", "pass")

    resp, err := http.PostForm(endpoint, data)
    if err != nil {
        fmt.Printf("%d回目のにゃんぱすーに失敗\n", i+1)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        fmt.Printf("%d回目のにゃんぱすーに成功\n", i+1)
    } else {
        fmt.Printf("%d回目のにゃんぱすーに失敗\n", i+1)
    }
    time.Sleep(time.Millisecond * 10)
}
