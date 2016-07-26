package CourseSelection
import (
    "strconv"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/unrolled/render"
    "github.com/gorilla/sessions"
    "github.com/gorilla/securecookie"
)
type StudentRigisterCourse struct{
    StudentID       string
    CourseList      []RigisteredCourse
}
type RigisteredCourse struct{
    CourseID        string          `json:"courseID"`
    CourseNo        string          `json:"courseNo"`
    CourseName      string          `json:"courseName"`
    TeacherName     string          `json:"teacherName"`
    CourseState     int             `json:"courseState"`
    // 0 means the Course in the queue
    // 1 means the Course finished
    // -1 means the Course deleted
    QueueNumber     int             `json:"queueNumber"`
}
type SelectLists struct{
    SelectList []RigisterCourse     `json:"SelectList"`
}
type RigisterCourse struct{
    CourseID        string          `json:"courseID"`
    CourseNo        string          `json:"courseNo"`
}
type RigisterCourseInfo struct{
    StudentID   string
    CourseNo    string
    QueueNum    int
}
type StudentInfo struct{
    StudentID       string
    StudentPW       string
    PWNotEffective  bool
}
type ServePrimeFunc func (w http.ResponseWriter, req *http.Request)
var (
    courseInQueue                    = 0
    courseDelete                     = -1
    courseFinish                     = 1
    homeStatusList          []string = []string{"0"}
    homeSelectCourseType    []string = []string{"0","1","2","3","4","5","6"}
)
//If you want to know what's the global variables meaning,you should read the README.md
//in the github
var (
    NotifyCourse chan [2]string = make(chan [2]string,100)
)
func InitServerMux(r *render.Render) (*http.ServeMux,*sessions.CookieStore){
    mux := http.NewServeMux()
    store := sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(32)))
    mux.HandleFunc("/",RootFunc())
    mux.HandleFunc("/index",IndexFunc(r))
    mux.HandleFunc("/login",LoginFunc(r,store))
    mux.HandleFunc("/home",HomeFunc(r,store))
    mux.HandleFunc("/home/select",HomeSelectFunc(r,store))
    mux.HandleFunc("/home/delete",HomeDeleteFunc(store))
    mux.HandleFunc("/commonselect",CommonSelectFunc(r))
    mux.HandleFunc("/home/register",HomeRegisterFunc(store))
    mux.HandleFunc("/feedback",FeedBackFunc(store))
    mux.HandleFunc("/errMessage",ErrorMessageFunc(r))
    return mux,store
}
func RootFunc() ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        http.Redirect(w,req,"/index",http.StatusMovedPermanently)
    }
}
func IndexFunc(r *render.Render) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        r.Text(w,http.StatusOK,"Hello World!This is the index of the website")
    }
}
func HomeDeleteFunc(store *sessions.CookieStore) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        courseNo := req.PostFormValue("courseNo")
        stuid,school,ok := validateSession(req,store)
        if ok{
            DBsession := GetSession()
            defer DBsession.Close()
            cRigister := DBsession.DB(school).C("RigisterInfo")

            //TODO test here
            err := cRigister.Update(bson.M{"studentid":stuid,"courselist.courseno":courseNo,"courselist.coursestate":courseInQueue},bson.M{"$set":bson.M{"courselist.$.coursestate":courseDelete}})
            if err == nil{
                http.Redirect(w,req,"/home",http.StatusMovedPermanently)
                return
            }
        }
            http.Redirect(w,req,"/errMessage",http.StatusMovedPermanently)
    }
}
func LoginFunc(r *render.Render,store *sessions.CookieStore) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        id := req.PostFormValue("UserID")
        pw := req.PostFormValue("UserPassword")
        //TODO It should be a form value that come from the user request
        school := "DHU"
        // DBsession := GetSession()
        // defer DBsession.Close()
        // cLogin := DBsession.DB(school).C("StudentInfo")
        if validateLogin(id,pw,nil){
            session,_ := store.Get(req,"sessionid")
            session.Values["stuid"] = id
            session.Values["school"] = school
            session.Save(req, w)
            http.Redirect(w,req,"/home",http.StatusMovedPermanently)
        }else{
            http.Redirect(w,req,"/errMessage",http.StatusMovedPermanently)
        }
    }
}
func validateLogin(id,pw string,cLogin *mgo.Collection) (flag bool){
    var err error
    err = cLogin.Find(bson.M{"studentid":id,"studentpw":pw})
    if err != nil{
        _,err = strconv.Atoi(id)
        if err != nil{
            return
        }else{
            //TODO request to the school website
            flag = true
        }
    }else{
        flag = true
    }
    return
}
func HomeFunc(r *render.Render,store *sessions.CookieStore) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        status := req.PostFormValue("userstatus")
        if validateHomeStatus(status){
            stuid,school,ok := validateSession(req,store)
            if ok{
                if status == ""{
                    hello := "Hello " + stuid + "!This is the home page of the website"
                    r.Text(w,http.StatusOK,hello)
                }else{
                    var err error
                    var done bool
                    var courselist []RigisteredCourse
                    DBsession := GetSession()
                    defer DBsession.Close()
                    cRigister := DBsession.DB(school).C("RigisterInfo")
                    for i := 0; i < 3;i++ {
                        courselist,err = getrigisteredFunc(stuid,cRigister)
                        if err == nil{
                            done = true
                            break
                        }
                    }
                    if done{
                        r.JSON(w,http.StatusOK,map[string]([]RigisteredCourse){"RigisterCourse":courselist})
                    }
                }
            }
        }
        http.Redirect(w,req,"/errMessage",http.StatusMovedPermanently)
    }
}
func validateHomeStatus(status string) (flag bool){
    if status == ""{
        flag = true
        return
    }
    for _,value := range homeStatusList{
        if value == status {
            flag = true
            break
        }
    }
    return
}
func validateSession(req *http.Request,store *sessions.CookieStore) (sessionid,schoolid string,flag bool){
    session,_ := store.Get(req,"sessionid")
    id := session.Values["stuid"]
    school := session.Values["shool"]
    schooln,ok := school.(string)
    if ok{
        _,ok := SchoolDB[schooln]
        if !ok {
            return
        }
    }else{
        return
    }
    stringid,ok := id.(string)
    if ok && stringid != ""{
        sessionid = stringid
        schoolid = schooln
        flag = true
    }
    return
}
func getrigisteredFunc(stuid string,cRigister *mgo.Collection) ([]RigisteredCourse,error){
    //TODO The collection will be nil here,so we will finish the database collection
    //and take it to the function
    // var stuRigister StudentRigisterCourse
    // err := cRigister.Find(bson.M{"studentid":stuid}).One(&stuRigister)
    // if err != nil && err != mgo.ErrNotFound{
    //     return []RigisteredCourse{},err
    // }else{
    //     return stuRigister.CourseList,nil
    // }
    return []RigisteredCourse{
                RigisteredCourse{"131441","专业英语","李悦",0,1},
                RigisteredCourse{"130153","计算机网络","朱明",0,1}},nil
}
func HomeSelectFunc(r *render.Render,store *sessions.CookieStore) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        coursetype := req.PostFormValue("coursetype")
        if validateCourseType(coursetype){
            sessionid,school,ok := validateSession(req,store)
            if ok{
                if coursetype == ""{
                    hello := "Hello " + sessionid + "!This is the home/select page of the website"
                    r.Text(w,http.StatusOK,hello)
                    return
                }else{
                    var done bool
                    var err  error
                    var teachSchemas []TeachSchema
                    DBsession := GetSession()
                    defer DBsession.Close()
                    for i := 0; i < 3; i++ {
                        cTable := DBsession.DB(school).C("CourseTable")
                        cIndex := DBsession.DB(school).C("CourseIndex")
                        teachSchemas,err = APIHomeSelect(cTable,cIndex,sessionid[:6],coursetype)
                        if err == nil{
                            done = true
                            break
                        }
                    }
                    if done{
                        r.JSON(w,http.StatusOK,map[string]([]TeachSchema){"TeachSchema":teachSchemas})
                        return
                    }
                }
            }
        }
        http.Redirect(w,req,"/errMessage",http.StatusMovedPermanently)
    }
}
func validateCourseType(courseType string) (flag bool){
    if courseType == ""{
        flag = true
        return
    }
    for _,value := range homeSelectCourseType{
        if value == courseType{
            flag = true
            break
        }
    }
    return
}
func CommonSelectFunc(r *render.Render) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        id := req.PostFormValue("courseID")
        if id == ""{
            r.Text(w,http.StatusOK,"Hello World!This is the commonselection page of the website")
        }else{
            var done bool
            var err error
            var course CourseContent
            DBsession := GetSession()
            defer DBsession.Close()
            cTable := DBsession.DB(school).C("CourseTable")
            for i := 0; i < 3; i++ {
                course,err = APICommonselect(cTable,id)
                if err == nil{
                    done = true
                    break
                }
            }
            if done{
                r.JSON(w,http.StatusOK,course)
            }else{
                http.Redirect(w,req,"/errMessage",http.StatusMovedPermanently)
            }
        }
    }
}
func HomeRegisterFunc(store *sessions.CookieStore) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        var slist SelectLists
        stuid,school,ok := validateSession(req,store)
        if ok{
            result,err := ioutil.ReadAll(req.Body)
            if err == nil{
                 err = json.Unmarshal([]byte(result), &slist)
                 if err == nil{
                     DBsession := GetSession()
                     defer DBsession.Close()
                     cTable := DBsession.DB(school).C("CourseTable")
                     if validateCourse(slist,cTable){
                         saveAndRegister(stuid,school,slist,DBsession.DB(school))
                         http.Redirect(w,req,"/home",http.StatusMovedPermanently)
                     }else{
                         http.Redirect(w,req,"/errMessage",http.StatusMovedPermanently)
                     }
                     return
                 }
            }
        }
        http.Redirect(w,req,"/errMessage",http.StatusMovedPermanently)
    }
}
func saveAndRegister(stuid,school string,slist SelectLists,DBsession *mgo.Database) {
    var err error
    cTable := DBsession.C("CourseTable")
    cRigister := DBsession.C("RigisterInfo")
    cCourseSelector := DBsession.C("CourseSelector")
    rigisterlist := map[string]RigisterCourseInfo{}
    for _,coursevalue := range slist.SelectList{
        var teachername string
        var courseContent CourseContent
        err = cTable.Find(bson.M{"courseid":coursevalue.CourseID}).One(&courseContent)
        if err == nil{
            for _,courselist := range courseContent.CourseList{
                if coursevalue.CourseNo == courselist.CourseNo{
                    teachername = courselist.TeacherName
                    break
                }
            }
            err = cRigister.Find(bson.M{"studentid":stuid,"courselist.courseno":coursevalue.CourseNo}).One(nil)
            if err == nil{
                err = cRigister.Find(bson.M{"studentid":stuid,"courselist.courseno":coursevalue.CourseNo,"courselist.coursestate":courseInQueue}).One(nil)
                if err == nil{
                    continue
                }else{
                    cRigister.Update(bson.M{"studentid":stuid,"courselist.courseno":coursevalue.CourseNo},bson.M{"$set":bson.M{"courselist.$.coursestate":courseInQueue}})
                }
            }else{
                CourseMap[school][coursevalue.CourseID].mutexDB.RLock()
                Squeuenum := getQueueNum(coursevalue.CourseNo,cCourseSelector) + 1
                rInfo ：= RigisteredCourse{coursevalue.CourseID,coursevalue.CourseNo,courseContent.CourseName,teachername,courseInQueue,Squeuenum}
                for{
                    err := cRigister.Insert(bson.M{"studentid":stuid},bson.M{"$push":bson.M{"courselist":rInfo}})
                    if err != nil{
                        break
                    }
                }
                CourseMap[school][coursevalue.CourseID].mutexDB.RUnlock()
            }
            rigisterlist[coursevalue.CourseID] = RigisterCourseInfo{stuid,coursevalue.CourseNo,Squeuenum}
        }
    }
    // take a goroutine to avoid blocking
    //TODO CourseMap here
    go func() {
        for key,value := range rigisterlist{
            CourseMap[school][key].mutexDB.Lock()
            for{
                err := cCourseSelector.Insert(value)
                if err != nil{
                    break
                }
            }
            CourseMap[school][key].mutexDB.Unlock()
        }
        for _,value := range slist.SelectList{
            NotifyCourse <- [2]string{school,value.CourseID}
        }
    }()
    return
}
func getQueueNum(courseNo string,cCourseSelector *mgo.Collection) (num int){
    var err error
    for{
        num,err = cCourseSelector.Find(bson.M{"CourseNo":courseNo}).Count()
        if err == nil{
            break
        }
    }
    return
}
func validateCourse(slist SelectLists,cTable *mgo.Collection) (flag bool){
    var err error
    var course CourseContent
    for _,courseInfo := range slist.SelectList{
        var ok bool
        err = cTable.Find(bson.M{"courseid":courseInfo.CourseID}).One(&course)
        if err == nil{
            for _,coursevalue := range course.CourseList{
                if coursevalue.CourseNo == courseInfo.CourseNo{
                    ok = true
                    break
                }
            }
            if !ok{
                return
            }
        }else{
            return
        }
    }
    flag = true
    return
}
func ErrorMessageFunc(r *render.Render) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        r.Text(w,http.StatusOK,"Hello World!This is the error page of the website")
    }
}
func FeedBackFunc(store *sessions.CookieStore) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        stuid,school,ok := validateSession(req,store)
        message := req.PostFormValue("Message")
        if message == ""{
            http.Redirect(w,req,"/errMessage",http.StatusMovedPermanently)
        }else{
            go storeFeedBackMessage(stuid,school,message)
            http.Redirect(w,req,"/index",http.StatusMovedPermanently)
        }
    }
}
func storeFeedBackMessage(stuid,school,message string){
    // DBsession := GetSession()
    // defer DBsession.Close()
    //TODO That we will take the message into the database and notify the developer
    return
}
