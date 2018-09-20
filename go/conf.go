package naturalvoid

import (
//    "os"
    "sync"
    "github.com/gorilla/sessions"
)

// Struct of configuration stuff
type Conf struct {
    SessionStore *sessions.CookieStore
}

var confInstance *Conf
var confOnce sync.Once

func GetConf() *Conf {
    confOnce.Do(func() {
        confInstance = &Conf{}
        err := confInstance.new()
        if err != nil {
            panic(err)
        }
    })
    return confInstance
}

func (conf *Conf) new() error {
    // Initialize a new Conf struct
    // Create a gorilla session store (use Env vars)
    conf.SessionStore = sessions.NewCookieStore([]byte("replacethiswithanactualsecretkey"))
    return nil
}
