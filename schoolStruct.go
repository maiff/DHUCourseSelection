package CourseSelection
import (
    "sync"
    "time"
    "net/url"
    "net/http"
    "net/http/cookiejar"
)
const (
    updateTime = 12
    //
)
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
        Timeout:time.Duration(10 * time.Second),
    }
}
func MakeParameters(para map[string]string) url.Values{
    data := make(url.Values)
    for key,value := range para{
        data.Set(key,value)
    }
    return data
}
