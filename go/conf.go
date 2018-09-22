package naturalvoid

import (
	"net/http"
	// "os"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"sync"
)

// Struct of configuration stuff
type Conf struct {
	SessionStore *sessions.CookieStore
	CSRF         func(http.Handler) http.Handler
}

var confInstance *Conf
var confOnce sync.Once

func GetConf() *Conf {
	confOnce.Do(func() {
		confInstance = &Conf{}
		confInstance.new()
	})
	return confInstance
}

func (conf *Conf) new() {
	// Initialize a new Conf struct
	// Create a gorilla session store (use Env vars)
	key := []byte("replacethiswithanactualsecretkey") // os.GetEnv
	production := false                               // os.GetEnv
	conf.SessionStore = sessions.NewCookieStore(key)
	conf.CSRF = csrf.Protect(key, csrf.Secure(production))
}
