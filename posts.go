package bellows

import (
	//"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	//"net/http"
	//"github.com/SlyMarbo/spdy"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Post struct {
	Id        int64  `json:"id"`
	ParentId  int64  `json:"parent_id"`
	Pseudonym string `json:"pseudonym"`
	Body      string `json:"body"`
}

func (Post) Storable() {}

//////// REST handlers ////////
func getPosts(p martini.Params, rend render.Render, l *log.Logger, c *Context) {
	id, err := strconv.ParseInt(p["pid"], 10, 64)
	if err != nil {
		rend.Error(http.StatusNotAcceptable)
		return
	}

	l.Printf("getting post %v", id)
	post := &Post{Id: id}
	if err := c.Storage.Get(post); err != nil {
		log.Printf("error getting post: ", err)
		rend.Error(http.StatusNotFound)
		return
	}

	rend.JSON(http.StatusOK, post)
}

func postPosts(rend render.Render, res http.ResponseWriter, req *http.Request, c *Context, l *log.Logger) {
	l.Println("got post req")

	//dec := json.NewDecoder(req.Body)
	//defer req.Body.Close()
	post := new(Post)
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		l.Printf("error on read: %v", err)
		rend.Error(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf, post); err != nil {
		l.Println("error during decode: ", err)
		rend.Error(http.StatusBadRequest)
		return
	}

	// overwrite any user input here
	post.Id = 0

	if len(post.Pseudonym) == 0 || len(post.Body) == 0 {
		rend.Error(http.StatusBadRequest)
		return
	}

	if err := c.Storage.Store(post); err != nil || post.Id == 0 {
		rend.Error(http.StatusInternalServerError)
		return
	}

	res.Header().Add("Location", fmt.Sprintf("/api/v1/posts/%d", post.Id))
	rend.JSON(http.StatusCreated, post)
}
