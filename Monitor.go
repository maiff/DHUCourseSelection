package CourseSelection
import (
    "fmt"
    "net/url"
    "net/http"
    "net/http/cookiejar"
)
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
