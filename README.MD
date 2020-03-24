GOREDIS

# Get Started
### Import
```go
import redis "github.com/ac5tin/goredis"
```
### create client
```go
func main(){
    rc := NewRedisClient(os.Getenv("REDIS_ADDR"), 0, os.Getenv("REDIS_PW"))
    rconn := rc.Get() // rc.Get() returns redis.Conn
}
```
### pubsub
```go
rconn := rc.Get()
rsc := make(chan map[string]string)
ready := make(chan bool)
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
go func() {
    err := rc.Subscribe(ctx, func(channel string, message []byte) error {
        if channel == channelName {
            data := map[string]interface{}
            json.Unmarshal(message, &data)
            if data["Type"] == "abc" {
                d := returnme // map[string]string
                rsc <- d
            }
        }
        return nil
    }, &ready, channelName)
    if err != nil {
        // failed to subscribe
        log.Println(err.Error())
        cancel()
        return
    }
}()

// publish
<-ready
res := map[string]interface{}{"Type":"abc"}
b, err := json.Marshal(&res)
log.Printf("publishing : %s\n", "abc")
_, err = rc.Publish(channelName, string(b))
if err != nil {
    // failed to publish
    log.Println(err.Error())
    cancel()
    return
}
```