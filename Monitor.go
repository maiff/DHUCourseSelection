package CourseSelection
import (
    "fmt"
    "time"
    "errors"
    "net/http"
)
type Monitor interface{
    MonitorParameter
    MonitorAction
}
type MonitorParameter interface{
    LoginPara(para ...string)     map[string]string
    // DeletePara(para ...string)    map[string]string
    MonitorPara(para ...string)   map[string]string
    RigisterPara(para ...string)  map[string]string
}
type MonitorAction interface{
    Login(map[string]string) (*http.Client,error)
    // DefaultLogin() *http.Client
    UpdateClient() func () *http.Client
    // Delete(map[string]string,*http.Client) error
    Register(map[string]string,*http.Client) error
    Monitor(map[string]string) ([]string,error)
    // RollBack(string,string)
    SetErrorMessage(string)
    // GetSchoolName() string
    // ValidateStuCourseConflict(string,*http.Client) (string,bool)
    ValidateStuCourseSelected(string,string,*http.Client) bool
}
var (
    SleepTime int = 1
    selectErrStr string = "Selected error "
    networkErr error = errors.New("Three times time out in the request")
    passwordErr error = errors.New("User's password can not pass authentication")
)
func MainMonitor() {
    var school string
    var courseid string
    var courseIsRun bool
    var newMessage [2]string
    for{
        newMessage <- NotifyCourse
        school = newMessage[0]
        courseid = newMessage[1]
        CourseMap[school][courseid].mutexBool.RLock()
        courseIsRun = CourseMap[school][courseid].running
        CourseMap[school][courseid].mutexBool.RUnlock()
        if !courseIsRun{
            setcourseMap(school,courseid,true)
            go MonitorFunc(getSchoolStruct(school),school,courseid)
        }
    }
}
func MonitorFunc(m Monitor,school,courseid string) {
    var err error
    var courselist []string
    var wg sync.WaitGroup
    detectDatabase := DetectDatabase(school,courseid)
    for{
        courselist,err = m.Monitor(m.MonitorPara(courseid))
        if err == nil{
            for _,courseno := range courselist{
                wg.Add(1)
                go registerFunc(m,school,courseno,&wg)
            }
            wg.Wait()
        }
        if detectDatabase(){
            setcourseMap(school,courseid,false)
            return
        }
        time.Sleep(SleepTime * time.Second)
    }
}
func registerFunc(m Monitor,school,courseno string,wg *sync.WaitGroup) {
    defer wg.Done()
    var err error
    var stuInfo    StudentInfo
    var courseInfo RigisterCourseInfo
    DBsession := GetSession()
    defer DBsession.Close()
    cCourseSelector := DBsession.DB(school).C("CourseSelector")
    cLogin := DBsession.DB(school).C("StudentInfo")
    err = cCourseSelector.Find(bson.M{"courseno":courseNo,"queuenum":1}).One(&courseInfo)
    if err == nil{
        err = cLogin.Find(bson.M{"studentid":courseInfo.StudentID}).One(&stuInfo)
        if err == nil{
            if stuInfo.PWEffective{
                return
            }
            client,err := m.Login(m.LoginPara(stuInfo.StudentID,stuInfo.StudentPW))
            if checkLoginErr(err,cLogin){
                return
            }
            err = m.Register(m.RigisterPara(courseNo),client)
            errflag,done := checkRegisterErr(err)
            if errflag{
                return
            }
            if !done{
                return
            }
            //I think we should not delete course automatically
            //So we should notify the course is conflicted and let users delete by the themselves
            // if !done{
            //     conflictCourse,isConflict := ValidateStuCourseConflict(courseid,client)
            //     if isConflict{
            //         err = m.Delete(m.DeletePara(conflictCourse),client)
            //         errflag,done = checkDeleteErr(err)
            //         if errflag || !done{
            //             return
            //         }
            //     }
            //     err = m.Register(m.RigisterPara(courseNo),client)
            //     errflag,done = checkRegisterErr(err)
            //     if errflag || !done {
            //         m.SetErrorMessage(err.String())
            //         m.RollBack(courseno,conflictCourse)
            //     }
            // }
            if m.ValidateStuCourseSelected(courseid,courseNo,client){
                alterDatabase(school,courseNo,courseInfo.StudentID,cCourseSelector)
            }else{
                m.SetErrorMessage(selectErrStr + courseno)
            }
        }
    }
}
func DetectDatabase(school,courseid string) func () bool{
    DBsession := GetSession()
    courselist,notFound := getCourselist(school,courseid,DBsession)
    cSelector := DBsession.DB(school).C("CourseSelector")
    if notFound{
        return func() bool{
            return true
        }
    }else{
        return func() bool{
            var err error
            for _,value := range courselist{
                err = cSelector.Find(bson.M{"courseno":value}).One(nil)
                if err == nil{
                    return false
                }
            }
            DBsession.Close()
            return true
        }
    }
}
func getCourselist(school,courseid string,s *mgo.Session) (courselist []string,notFound bool){
    cTable := s.DB(school).C("CourseTable")
    err := cTable.Find(bson.M{"courseid":courseid}).One(&course)
    if err != nil{
        notFound = true
    }else{
        for _,value := range course.CourseList{
            courselist.append(courselist,value.CourseNo)
        }
    }
    return
}
func alterDatabase(school,courseno,stuid string,cCourseSelector *mgo.Collection) {
    err := cCourseSelector.Remove(bson.M{"courseno":courseNo,"queuenum":1})
    if err == nil{
        cCourseSelector.UpdateAll(bson.M{"courseno":courseNo},bson.M{"$inc":bson.M{"queuenum":-1}})
        synchronizeDatabase(school,courseno,stuid)
    }
}
func synchronizeDatabase(school,courseno,stuid string) {
    DBsession := GetSession()
    defer DBsession.Close()
    cRigister := DBsession.DB(school).C("RigisterInfo")
    cRigister.UpdateAll(bson.M{"courselist.courseno":courseNo},bson.M{"$inc":bson.M{"courselist.queuenumber":-1}})
}
func checkDeleteErr(err error) (flag,done bool) {
    //The same as checkRegisterErr
    //And the deleteErr will make the done be false
    return checkRegisterErr(err)
}
func checkRegisterErr(err error) (flag,done bool){
    switch err {
    case networkErr:
        flag = true
    case nil:
        done = true
    // case conflictErr:
    }
    return
}
func checkLoginErr(err error,cLogin *mgo.Collection) (flag bool){
    switch err {
    case networkErr:
        flag = true
    case passwordErr:
        for{
            err = cLogin.Update(bson.M{"studentid":courseInfo.StudentID},bson.M{"$set":bson.M{"pwnoteffective":true}})
            if err == nil{
                flag = true
                break
            }
        }
    }
    return
}
func setcourseMap(school,courseid string,flag bool) {
    CourseMap[school][courseid].mutexBool.Lock()
    defer CourseMap[school][courseid].mutexBool.Unlock()
    CourseMap[school][courseid].running = flag
}
