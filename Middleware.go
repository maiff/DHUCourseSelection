package CourseSelection
import (
    "os"
    "log"
    "net/http"
    "gopkg.in/mgo.v2"
    "github.com/urfave/negroni"
    "github.com/unrolled/render"
    "github.com/gorilla/context"
    "github.com/gorilla/sessions"
)
func RunFunc() {
    r := render.New()
    mgoSession,err := mgo.Dial("localhost:27017")
    if err != nil {
        panic(err)
    }
    defer mgoSession.Close()
    mux,store := InitServerMux(r)
    n   := negroni.Classic()
    n.Use(negroni.HandlerFunc(CookieMiddleware(store)))
    n.Use(negroni.HandlerFunc(SetAllowOrigin()))
    n.UseHandler(context.ClearHandler(mux))
    //All the middlerware we used
    l := log.New(os.Stdout, "[negroni] ", 0)
    l.Printf("listening on :8080")
    server := http.Server{Addr:":8080",Handler:n}
    server.SetKeepAlivesEnabled(true)
    l.Fatal(server.ListenAndServe())
}

type MiddlewareFunc func(rw http.ResponseWriter,r *http.Request,next http.HandlerFunc)
func CookieMiddleware(c *sessions.CookieStore) MiddlewareFunc{
    //That Middleware is used to detect the session
    UrlNotDetectSessionList := map[string]string{"/":"","/index":"","/login":"","/errMessage":"","/feedback":""}
    return func(rw http.ResponseWriter,r *http.Request,next http.HandlerFunc){
        urlpath := r.URL.Path
        if _,ok := UrlNotDetectSessionList[urlpath];!ok{
            session, _ := c.Get(r, "sessionid")
            if session.IsNew || session.Values["stuid"] == ""{
                next(rw,r)
                http.Redirect(rw,r,"/index",http.StatusMovedPermanently)
            }
        }
        next(rw,r)
    }
}
func SetAllowOrigin() MiddlewareFunc{
    return func(rw http.ResponseWriter,r *http.Request,next http.HandlerFunc){
        rw.Header().Set("Access-Control-Allow-Origin", "*")
        next(rw,r)
    }
}
