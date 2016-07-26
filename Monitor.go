package CourseSelection
import (
    "fmt"
    "time"
    "errors"
    "net/url"
    "net/http"
    "net/http/cookiejar"
)
type Monitor interface{
    MonitorParameter
    MonitorAction
}
type MonitorParameter interface{
    LoginPara(para ...string)     map[string]string
    DeletePara(para ...string)    map[string]string
    monitorPara(para ...string)   map[string]string
    RigisterPara(para ...string)  map[string]string
}
type MonitorAction interface{
    Login(map[string]string) *http.Client
    // DefaultLogin() *http.Client
    UpdateClient() *http.Client
    Delete(map[string]string,*http.Client) error
    Register(map[string]string,*http.Client) error
    Monitor(map[string]string,string) ([]string,error)
    RollBack(string,string)
    SetErrorMessage(string)
    // GetSchoolName() string
    ValidateStuCourseConflict(string,*http.Client) (string,bool)
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
    DBsession := GetSession()
    defer DBsession.Close()
    cCourseSelector := DBsession.DB(school).C("CourseSelector")
    cLogin := DBsession.DB(school).C("StudentInfo")
    for{
        courselist,err = m.Monitor(courseid)
        if err == nil{
            for _,courseno := range courselist{
                wg.Add(1)
                go monitorFunc(m,school,courseno,&wg)
            }
            wg.Wait()
        }
        if detectDatabase(school,courseid){
            return
        }
        time.Sleep(SleepTime * time.Second)
    }
}
func monitorFunc(m Monitor,school,courseno string,wg *sync.WaitGroup) {
    var err error
    var stuInfo    StudentInfo
    var courseInfo RigisterCourseInfo
    defer wg.Done()
    err = cCourseSelector.Find(bson.M{"courseno":courseNo,"queuenum":1}).One(&courseInfo)
    if err == nil{
        err = cLogin.Find(bson.M{"studentid":courseInfo.StudentID}).One(&stuInfo)
        if err == nil{
            if stuInfo.PWEffective{
                return
            }
            client,err := m.Login(m.LoginPara(stuInfo.StudentID,stuInfo.StudentPW))
            if checkLoginErr(err){
                return
            }
            err = m.Register(m.RigisterPara(courseNo),client)
            errflag,done := checkRegisterErr(err)
            if errflag{
                return
            }
            if !done{
                conflictCourse,isConflict := ValidateStuCourseConflict(courseid,client)
                if isConflict{
                    err = m.Delete(m.DeletePara(conflictCourse),client)
                    errflag,done = checkDeleteErr(err)
                    if errflag || !done{
                        return
                    }
                }
                err = m.Register(m.RigisterPara(courseNo),client)
                errflag,done = checkRegisterErr(err)
                if errflag || !done {
                    m.SetErrorMessage(err.String())
                    m.RollBack(courseno,conflictCourse)
                }
            }
            if m.ValidateStuCourseSelected(courseid,courseNo,client){
                alterDatabase(courseid,courseNo,school)
            }else{
                m.SetErrorMessage(selectErrStr + courseno)
            }
        }
    }
}
func alterDatabase(courseid,courseno,school string) {

}
func synchronizeDatabase() {
    
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
func checkLoginErr(err error) (flag bool){
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
func detectDatabase(school,courseid string) (Empty bool){

}
func setcourseMap(school,courseid string,flag bool) {
    CourseMap[school][courseid].mutexBool.Lock()
    defer CourseMap[school][courseid].mutexBool.Unlock()
    CourseMap[school][courseid].running = flag
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
