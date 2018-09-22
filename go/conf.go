package naturalvoid

import (
	// "os"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"net/http"
	"sync"
)

// Struct of configuration stuff
type Conf struct {
	SessionStore *sessions.CookieStore
	CSRF func(http.Handler) http.Handler
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
	key := []byte("replacethiswithanactualsecretkey")  // os.GetEnv
	production := false  // os.GetEnv
	conf.SessionStore = sessions.NewCookieStore(key)
	conf.CSRF = csrf.Protect(key, csrf.Secure(production))
	return nil
}
