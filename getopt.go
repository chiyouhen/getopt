package getopt

import (
    "os"
    "io"
    "fmt"
    "strings"
)

type Option struct {
    LongOpt string,
    ShortOpt string,
    EnvOpt string,
    WithLong bool,
    WithShort bool,
    WithEnv bool,
    IsFlag bool,
    Description bool,
    DefaultValue string,
}

type Config struct {
    Description string,
    Options []*Option,
}

type Definition struct {
    Conf *Config,
    Tokens []string,
    ShortMap map[string]*Option,
    LongMap map[string]*Option,
    EnvMap map[string]*Option,
    Opts [][]string,
    Args []string,
}

type GetoptError struct {
    msg string
}

func (err *GetoptError) Error() string {
    return err.msg
}

func (cf *Config) ParseCommandLine() ([][]string, []string, error) {
    var def = &Definition{
        Conf: cf,
    }
    return def.ParseCommandLine()
}

func (def *Definition) CreateOptMaps() (err error) {
    err = nil
    def.ShortMap = make(map[string]*Option)
    def.LongMap = make(map[string]*Option)
    def.EnvMap = make(map[string]*Option)

    for _, o := range def.Conf.Options {
        if o.WithLong {
            def.LongMap[o.LongOpt] = o
        }
        if o.WithShort {
            def.ShortMap[o.ShortOpt] = o
        }
        if o.WithEnv {
            def.EnvMap[o.EnvOpt] = o
        }
    }
    return
}

func (def *Definition) ReadToken() (token string, ok bool) {
    if len(def.Tokens) > 0 {
        token = def.Tokens[0]
        def.Tokens = def.Tokens[1:]
        ok = true
    } else {
        ok = false
    }
    return
}

func (def *Definition) DoLongs(cmd string) (err error) {
    var i = strings.Index(cmd, "=")
    var value string
    var ok bool
    if i > 0 {
        value = cmd[i + 1:]
        cmd = cmd[:i]
    }
    var o, ok = def.LongMap[cmd]
    if ! ok {
        return &GetoptError{fmt.Sprintf("invalid argument: %s", cmd)}
    }
    if i > 0 && o.IsFlag {
        return &GetoptError{fmt.Sprintf("argument %s defined as flag, but key-value format found", cmd)}
    }
    if o.IsFlag {
        def.Opts = append(def.Opts, []string{cmd, ""})
    } else {
        if i < 0 {
            value, ok = def.ReadToken()
            if ! ok {
                return &GetoptError{fmt.Sprintf("argument %s defined as option, but no value specified", cmd)}
            }
        }
        def.Opts = append(def.Opts, []string{cmd, value})
    }
    return
}

func (def *Definition) DoShorts(cmd string) (err error) {
    for len(cmd) > 0 {
        var value string
        var c = string(cmd[0])
        cmd = string(cmd[1:])
        var o, ok = def.ShortMap[c]
        if ! ok {
            return &GetoptError{fmt.Sprintf("invalid argument %s", c)}
        }
        if o.IsFlag {
            def.Opts = append(def.Opts, []string{c, ""})
        } else {
            if len(cmd) > 0 {
                def.Opts = append(def.Opts, []string{c, cmd})
                break
            } else {
                value, ok = def.ReadToken()
                if ! ok {
                    return &GetoptError{fmt.Sprintf("argument %s defined as option, but no value specified", c)}
                }
                def.Opts = append(def.Opts, []string{c, value})
            }
        }
    }
    return
}

func (def *Definition) ParseTokens() (err error) {
    err = nil
    var cmd string
    var ok bool
    for cmd, ok = def.ReadToken(); ok {
        if cmd == "--" {
            def.Args = append(def.Args, def.Tokens...)
            break
        }
        if strings.HasPrefix(cmd, "--") {
            cmd = strings.TrimPrefix(cmd, "--")
            err = def.DoLongs(cmd)
            if err != nil {
                return
            }
        } else if strings.HasPrefix(cmd, "-") {
            cmd = strings.TrimPrefix(cmd, "-")
            err = def.DoShorts(cmd)
            if err != nil {
                return
            }
        } else {
            def.Args = append(def.Args, cmd)
        }
    }
}

func (def *Definition) ParseCommandLine() (opts [][]string, args []string, err error) {
    err = def.CreateOptMaps()
    def.Tokens = make([]string, len(os.Args))
    copy(def.Tokens, os.Args)
    err = def.ParseTokens()
    opts = def.Opts
    args = def.Args
    return
}
