1. get readonly commands
    ```go
    package main
    
    import (
        "context"
        "fmt"
        "log"

        "github.com/go-redis/redis/v8"
    )
    
    func main() {
        rdb := redis.NewClient(&redis.Options{
            Addr: "localhost:6379",
        })
    
        ctx := context.Background()
        cmdsRes := rdb.Command(ctx)
    
        var readonly []string
    
        cc, err := cmdsRes.Result()
        if err != nil {
            log.Fatal(err)
        }
        for _, r := range cc {
            if !r.ReadOnly {
                continue
            }
    
            readonly = append(readonly, r.Name)
        }
    
        fmt.Println(readonly)
    }
    ```
2. remove all lines by regexp `^(?!(?:del|set...)\().*`
3. replace by pattern

    Find:
    `^([^(]+)(\((?:((?:, )?\w+)(?: (?:[^,)]+))?)?(?:((?:, )?\w+)(?: (?:[^,)]+))?)?(?:((?:, )?\w+)(?: (?:[^,)]+))?)?(?:((?:, )?\w+)(?: (?:[^,)]+))?)?(?:((?:, )?\w+)(?: (?:[^,)]+))?)?\)) (\*?)(.+)`
    
    Replace:
    ```regexp
    func(p *WritePipeliner) $1$2 $8redis.$9 {\n\treturn p.read.$1($3$4$5$6$7)\n}\n
    ```
4. Set redis for arguments
   
    Find: `(\w)\((.*\w )\*(\w+.*)\)`
    replace: `$1($2*redis.$3)`
5. Set ... for arguments

    Find: `(keys|members|pos|fields)\)`
    replace: `$1...)`