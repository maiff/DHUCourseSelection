package CourseSelection
import (
    "sync"
    "gopkg.in/mgo.v2"
)
type mCourseInfo struct{
    running       bool
    mutexDB       *sync.RWMutex
    mutexBool     *sync.RWMutex
}
var (
    SchoolList []string = []string{"DHU"}
    CourseMap  map[string](map[string]mCourseInfo) = map[string](map[string]mCourseInfo){}
    SchoolSession *mgo.Session
    // SchoolDB   map[string](*mgo.Database) = map[string](*mgo.Database){}
)
func GetSession() *mgo.Session{
    return SchoolSession.Clone()
}
func InitSchoolDB() {
    mgoSession,err := mgo.Dial("localhost:27017")
    if err != nil {
        panic(err)
    }
    mgoSession.SetMode(mgo.Strong,true)
    SchoolSession = mgoSession
}
//TODO Rebuild here
func InitCourseMap() {
    for school,DB := range SchoolDB{
        var err error
        var allCourse []CourseContent
        cTable := DB.C("CourseTable")
        err = cTable.Find(nil).All(&allCourse)
        if err != nil{
            panic(err)
        }
        for _,course := range allCourse{
            CourseMap[school][course.CourseID] = mCourseInfo{}
        }
        // for key,value := range CourseMap{
        //     fmt.Println(key)
        //     fmt.Println(value)
        // }
    }
}
