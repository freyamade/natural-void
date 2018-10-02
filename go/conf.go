package naturalvoid

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
)

// Struct of configuration stuff
type Conf struct {
	SessionStore *sessions.CookieStore
	CSRF         func(http.Handler) http.Handler
	// Database stuff
	DBHost   string
	DBPort   int
	DBUser   string
	DBPass   string
	DBName   string
	DBSecure string
	// LDAP Stuff
	LDAPHost   string
	LDAPPort   int
	LDAPUser   string
	LDAPPass   string
	LDAPSecure bool
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
	key := []byte(getEnv("SECRET_KEY", "replacethiswithanactualsecretkey"))
	production := getEnv("PRODUCTION", "false") == "true"
	conf.SessionStore = sessions.NewCookieStore(key)
	conf.CSRF = csrf.Protect(key, csrf.Secure(production))

	// Store other env variables in the conf
	// DB stuff
	conf.DBHost = getEnv("POSTGRES_HOST", "localhost")
	dbPort, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		fmt.Println("Invalid LDAP Port value, not an integer.")
		panic(err)
	}
	conf.DBPort = dbPort
	conf.DBUser = getEnv("POSTGRES_USER", "postgres")
	conf.DBPass = getEnv("POSTGRES_PASS", "postgres")
	conf.DBName = getEnv("POSTGRES_NAME", "postgres")
	dbSecure := getEnv("POSTGRES_SECURE", "false") == "true"
	if dbSecure {
		conf.DBSecure = "enable"
	} else {
		conf.DBSecure = "disable"
	}

	// LDAP Stuff
	conf.LDAPHost = getEnv("LDAP_HOST", "localhost")
	ldapPort, err := strconv.Atoi(getEnv("LDAP_PORT", "389"))
	if err != nil {
		fmt.Println("Invalid LDAP Port value, not an integer.")
		panic(err)
	}
	conf.LDAPPort = ldapPort
	conf.LDAPUser = getEnv("LDAP_USER", "cn=admin,dc=naturalvoid,dc=com")
	conf.LDAPPass = getEnv("LDAP_PASS", "nv")
	conf.LDAPSecure = getEnv("LDAP_SECURE", "false") == "true"
}

// Helper to provide a default for getting environment variables
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
