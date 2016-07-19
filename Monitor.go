package CourseSelection
import (
    "fmt"
    "sync"
    "gopkg.in/mgo.v2"
    // "gopkg.in/mgo.v2/bson"
)
type mCourseInfo struct{
    running     bool
    mutex       *sync.RWMutex
}
var (
    CourseMap map[string]mCourseInfo = map[mCourseID]mCourseInfo{}
)
func InitCourseMap(DBsession *mgo.Database) {
    var err error
    var allCourse []CourseContent
    cTable := DBsession.C("CourseTable")
    err = cTable.Find(nil).All(&allCourse)
    if err != nil{
        panic(err)
    }
    for _,course := range allCourse{
        CourseMap[course.CourseID] = mCourseInfo{}
    }
    for key,value := range CourseMap{
        fmt.Println(key)
        fmt.Println(value)
    }
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
