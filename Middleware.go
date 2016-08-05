package CourseSelection
import (
    "os"
    "log"
    "time"
    "net/http"
    "github.com/urfave/negroni"
    "github.com/unrolled/render"
    "github.com/gorilla/context"
    "github.com/gorilla/sessions"
)
const (
    timeoutErrMessage   =   ""
)
func RunFunc() {
    r := render.New()
    mux,store := InitServerMux(r)
    n   := negroni.Classic()
    n.Use(negroni.HandlerFunc(CookieMiddleware(store)))
    n.Use(negroni.HandlerFunc(SetAllowOrigin()))
    timeourHandler := http.TimeoutHandler(context.ClearHandler(mux),time.Duration(3 * time.Second),timeoutErrMessage)
    n.UseHandler(timeourHandler)
    //All the middlerware we used
    l := log.New(os.Stdout, "[negroni] ", 0)
    l.Printf("listening on :8080")
    server := http.Server{Addr:":8080",Handler:n}
    server.SetKeepAlivesEnabled(true)
    MainMonitor()
    InitCourseMap()
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
