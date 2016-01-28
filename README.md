# getopt
a traditional getopt library for golang.

# example
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
  
# functions
    func Getopt(tokens []string, shortopts string, longopts []string) (opts [][]string, args []string, err error)

## tokens
usually os.Args[1:], developer can also specify a different string array.

## shortopts
a traditional shortopt definition. a singal letter means this is a flag, as shown in the example '-h'. a single letter followed by a ':' means this is a option, user must specify option value while execute program, as shown in the example '-c'.

## longopts
a array of string. a option defined as --_option_name_[=]. '--' prefix is required. with '=' as suffix means this is a option, without '=' followed means this is a flag.


  
