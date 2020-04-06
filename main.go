package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"
)

var blogPosts = make(map[string]*Post)
var dataLock = sync.RWMutex{}

type Post struct {
	Title string
	Content string
	Time time.Time
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", newRouter()))
}

func newRouter() *mux.Router{
	router := mux.NewRouter()

	// Register pprof handlers
	registerProfilingHandlers(router)

	// Register our business logic handler
	router.Methods(http.MethodPut).Path("/blog/{title}").HandlerFunc(CreateBlogPost)
	router.Methods(http.MethodGet).Path("/blog/{title}").HandlerFunc(GetBlogPost)

	return router
}

func registerProfilingHandlers(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
}

func CreateBlogPost(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	params := mux.Vars(r)

	title, found := params["title"]
	if !found {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	dataLock.Lock()
	defer dataLock.Unlock()

	post := string(body)
	blogPosts[title] = &Post{
		Title:   title,
		Content: post,
		Time:    time.Now(),
	}
	log.Print("Created blogpost")
	rw.WriteHeader(http.StatusCreated)
}

func GetBlogPost(rw http.ResponseWriter, r *http.Request) {
	dataLock.RLock()
	defer dataLock.RUnlock()

	params := mux.Vars(r)
	title, found := params["title"]
	if !found {
		rw.WriteHeader(http.StatusNotFound)
		log.Print("Blogpost not found")
		return
	}

	post, found := blogPosts[title]
	if !found {
		rw.WriteHeader(http.StatusNotFound)
		log.Print("Blogpost not found")
		return
	}

	_,err := rw.Write([]byte(post.Content))
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print("Fetched blogpost")

}