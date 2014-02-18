package bellows

import (
	//"fmt"
	//"github.com/SlyMarbo/spdy"
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/pabbott0/martini-spdy"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Bellows struct {
	Context *Context
	martini *martini.Martini
}

type Context struct {
	Conf    *Config
	Storage *StorageEngine
	Events  *EventEngine
}

// Generally what you want to call
// to get everything set up
func Server(confPath string) *Bellows {
	// slurp in the conf file
	rawConf, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(err)
	}

	conf := new(Config)
	jserr := json.Unmarshal(rawConf, conf)
	if err != nil {
		panic(jserr)
	}

	// set the env var before initializing Martini
	if len(conf.Env) == 0 {
		log.Println("no env defined in conf - using default")
		conf.Env = "development"
	}
	log.Printf("env: %s", conf.Env)
	os.Setenv("MARTINI_ENV", conf.Env)

	b := NewBellows(conf)
	return b
}

func NewBellows(conf *Config) *Bellows {
	// set up Martini
	m := martini.New()

	m.Handlers(
		martini.Recovery,
		martini.Logger,
		render.Renderer(),
		mspdy.All(),
	)

	r := martini.NewRouter()
	defRoutes(r)
	m.Action(r.Handle)

	eve := NewEventEngine(conf)
	ste := NewStorageEngine(conf, eve)

	// set up context
	c := &Context{conf, ste, eve}
	m.Map(c)

	return &Bellows{c, m}
}

func (b *Bellows) Handler() http.Handler {
	return b.martini
}

func defRoutes(r martini.Router) {
	r.Get("/", func() string {
		return "this is /"
	})
	r.Get("/conf", func(c *Context, r render.Render) {
		//return c.conf.Env
		r.JSON(200, c.Conf)
	})
	r.Post("/test", func(rend render.Render, res http.ResponseWriter, req *http.Request, c *Context, l *log.Logger) string {
		l.Println("in test")
		buf, _ := ioutil.ReadAll(req.Body)
		return string(buf)
	})

	r.Get("/api/v1/posts/:pid", getPosts)
	r.Post("/api/v1/posts", postPosts)

}
