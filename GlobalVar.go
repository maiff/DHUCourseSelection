package CourseSelection
import (
    "sync"
    "gopkg.in/mgo.v2"
)
type mCourseInfo struct{
    running     bool
    mutex       *sync.RWMutex
}
var (
    SchoolList []string = []string{"DHU"}
    CourseMap  map[string]mCourseInfo = map[mCourseID]mCourseInfo{}
    SchoolDB   map[string](*mgo.Database) = map[string](*mgo.Database){}
)
func InitSchoolDB(mgoSession *mgo.Session) {
    var err error
    for _,school := range SchoolList{
        SchoolDB[school] = mgoSession.DB(school)
    }
}
//TODO Rebuild here
func InitCourseMap() {
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
