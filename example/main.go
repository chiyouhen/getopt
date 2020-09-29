package main

import (
    "fmt"
    "os"
    "github.com/chiyouhen/getopt"
)

func Usage() {
    fmt.Println("usage: xxx")
}

func main() {
    var configPath = "/user/local/etc/hello.conf"
    var opts, args, err = getopt.Getopt(os.Args[1:], "hc:", []string{"--help", "--config="})
    if err != nil {
        fmt.Printf("error while getopt: %v\n", err)
        Usage()
        os.Exit(1)
    }

    for _, opt := range opts {
        var k, v = opt[0], opt[1]
        switch k {
            case "-h", "--help":
                Usage()
                os.Exit(0)
            case "-c", "--config":
                configPath = v
        }
    }
    fmt.Printf("configPath: %s\n", configPath)
    fmt.Println(opts, args)
}