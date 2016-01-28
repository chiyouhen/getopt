package getopt

import (
    "fmt"
    "bytes"
    "strings"
)

type GetoptError struct {
    message string
}

func (err *GetoptError) Error() string {
    return err.message
}

func ParseShortOptions(shortopts string) (shortoptions map[string]bool, err error) {
    var buf = bytes.NewBufferString(shortopts)
    var curr string
    var b byte
    var ok bool
    shortoptions = make(map[string]bool)
    for buf.Len() > 0 {
        b, _ = buf.ReadByte()
        if b == ':' {
            if curr == "" {
                err = &GetoptError{"invalid short options"}
                return
            }
            shortoptions[curr] = false
        } else if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
            curr = fmt.Sprintf("-%s", string(b))
            _, ok = shortoptions[curr]
            if ok {
                err = &GetoptError{"duplicate definition of short options"}
                return
            }
            shortoptions[curr] = true
        } else {
            err = &GetoptError{"invalid short options"}
        }
    }
    return
}

func ParseLongOptions(longopts []string) (longoptions map[string]bool, err error) {
    longoptions = make(map[string]bool)
    for len(longopts) > 0 {
        var c = longopts[0]
        longopts = longopts[1:]
        var flag = false
        if ! strings.HasPrefix(c, "--") {
            err = &GetoptError{"invalid long options"}
            return
        }
        if strings.HasSuffix(c, "=") {
            c = c[:len(c) - 1]
        } else {
            flag = true
        }
        if len(c) == len("--") {
            err = &GetoptError{"invalid long options"}
            return
        }
        longoptions[c] = flag
    }
    return
}

func DoLongs(longoptions map[string]bool, cmd string, tokens *[]string, opts *[][]string, args *[]string) (err error) {
    var value string
    var i = strings.Index(cmd, "=")
    if i > 0 {
        value = cmd[i + 1:]
        cmd = cmd[:i]
    }
    var flag, ok = longoptions[cmd]
    if ! ok {
        err = &GetoptError{fmt.Sprintf("invalid argument: %s", cmd)}
        return
    }

    if flag {
        if i > 0 {
            err = &GetoptError{fmt.Sprintf("flag %s defined, but option specified", cmd)}
            return
        }
        *opts = append(*opts, []string{cmd, ""})
    } else {
        if i < 0 {
            if len(*tokens) < 1 {
                err = &GetoptError{fmt.Sprintf("option %s defined, but no value specified", cmd)}
                return
            } else {
                value = (*tokens)[0]
                *tokens = (*tokens)[1:]
            }
        }
        *opts = append(*opts, []string{cmd, value})
    }
    return
}

func DoShorts(shortoptions map[string]bool, cmd string, tokens *[]string, opts *[][]string, args *[]string) (err error) {
    cmd = strings.TrimPrefix(cmd, "-")
    var buf = bytes.NewBufferString(cmd)
    for buf.Len() > 0 {
        var b, _ = buf.ReadByte()
        var c = fmt.Sprintf("-%s", string(b))
        var flag, ok = shortoptions[c]
        var value string
        if ! ok {
            err = &GetoptError{fmt.Sprintf("invalid argument: %s", c)}
            return
        }
        if flag {
            *opts = append(*opts, []string{c, ""})
        } else {
            if buf.Len() > 0 {
                value = buf.String()
                *opts = append(*opts, []string{c, value})
            } else {
                if len(*tokens) < 1 {
                    err = &GetoptError{fmt.Sprintf("option %s defined, but no value specified", c)}
                    return
                } else {
                    value = (*tokens)[0]
                    *tokens = (*tokens)[1:]
                    *opts = append(*opts, []string{c, value})
                }
            }
            break
        }
    }
    return
}

func Getopt(tokens []string, shortopts string, longopts []string) (opts [][]string, args []string, err error) {
    var longoptions map[string]bool
    var shortoptions map[string]bool
    longoptions, err = ParseLongOptions(longopts)
    if err != nil {
        return
    }
    shortoptions, err = ParseShortOptions(shortopts)
    if err != nil {
        return
    }

    opts = make([][]string, 0, len(tokens))
    args = make([]string, 0, len(tokens))

    for len(tokens) > 0 {
        var c = string(tokens[0])
        tokens = tokens[1:]
        if c == "--" {
            args = append(args, tokens...)
            break
        }
        if strings.HasPrefix(c, "--") {
            err = DoLongs(longoptions, c, &tokens, &opts, &args)
            if err != nil {
                return
            }
        } else if strings.HasPrefix(c, "-") {
            err = DoShorts(shortoptions, c, &tokens, &opts, &args)
            if err != nil {
                return
            }
        } else {
            args = append(args, c)
        }
    }
    return
}
