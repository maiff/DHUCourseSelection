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
    errChan chan string
    SchoolList []string = []string{"DHU"}
    CourseMap  map[string](map[string]*mCourseInfo) = map[string](map[string]*mCourseInfo){}
    SchoolSession *mgo.Session
    SchoolStructs map[string]Monitor = map[string]Monitor{}
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
func InitSchoolStructs() {
    SchoolStructs["DHU"] = NewDHUStruct()
}
func getSchoolStruct(school string) Monitor{
    return SchoolStructs[school]
}

func InitCourseMap() {
    session := GetSession()
    defer session.Close()
    for _,school := range SchoolList{
        var err error
        var allCourse []CourseContent
        cMap := map[string]*mCourseInfo{}
        cTable := session.DB(school).C("CourseTable")
        err = cTable.Find(nil).All(&allCourse)
        if err != nil{
            panic(err)
        }
        for _,course := range allCourse{
            cMap[course.CourseID] = &mCourseInfo{}
        }
        CourseMap[school] = cMap
        // for key,value := range CourseMap{
        //     fmt.Println(key)
        //     fmt.Println(value)
        // }
    }
}
func getErrChan() chan string{
    if errChan == nil{
        errChan = make(chan string,50)
    }
    return errChan
}
