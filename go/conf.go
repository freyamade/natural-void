package naturalvoid

import (
    "os"
    "sync"
    "github.com/gorilla/sessions"
)

// Struct of configuration stuff
type Conf struct {
    SessionStore *sessions.Store
}

var instance *Conf
var once sync.Once

func GetConf() *Conf {
    once.Do(func() {
        instance = &Conf{}
        err := instance.new()
        if err != nil {
            panic(err)
        }
    })
    return instance
}

func (conf *Conf) new() error {
    // Initialize a new Conf struct
    // Create a gorilla session store (use Env vars)
    conf.SessionStore = sessions.NewCookieStore("replacethiswithanactualsecretkey")
}
