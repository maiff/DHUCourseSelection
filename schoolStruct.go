package CourseSelection
import (
    "sync"
    "time"
    "net/url"
    "net/http"
    "net/http/cookiejar"
)
const (
    updateTime = 
    //
)
var (
    errChan chan string
)
func getErrChan() chan string{
    if errChan == nil{
        errChan = make(chan string,50)
    }
    return errChan
}
type SchoolStruct struct{
    //    SchoolName string
    ErrChan     chan string
    Client      *http.Client
    mutexClient *sync.RWMutex
}
func (s *SchoolStruct) SetErrorMessage(message string){
    s.ErrChan <- message
    return
}
func NewClient() *http.Client{
    jar,_ := cookiejar.New(nil)
    return &http.Client{
        Jar:jar,
    }
}
func MakeParameters(para map[string]string) url.Values{
    data := make(url.Values)
    for key,value := range para{
        data.Set(key,value)
    }
    return data
}
