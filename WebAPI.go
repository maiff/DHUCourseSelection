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
    //1 means the Course finished
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
    NotifyCourseID chan string = make(chan string,100)
)
func InitServerMux(r *render.Render,mgoSession *mgo.Session) (*http.ServeMux,*sessions.CookieStore){
    mux := http.NewServeMux()
    DHUDB := mgoSession.DB("DHU")
    // cTable := session.DB("DHU").C("CourseTable")
    // cIndex := session.DB("DHU").C("CourseIndex")
    // cLogin := session.DB("DHU").C("StudentInfo")
    // cRigister := session.DB("DHU").C("RigisterInfo")
    store := sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(32)))
    mux.HandleFunc("/",RootFunc())
    mux.HandleFunc("/index",IndexFunc(r))
    mux.HandleFunc("/login",LoginFunc(r,DHUDB,store))
    mux.HandleFunc("/home",HomeFunc(r,DHUDB,store))
    mux.HandleFunc("/home/select",HomeSelectFunc(r,DHUDB,store))
    mux.HandleFunc("/home/delete",HomeDeleteFunc(DHUDB,store))
    mux.HandleFunc("/commonselect",CommonSelectFunc(r,DHUDB))
    mux.HandleFunc("/home/register",HomeRegisterFunc(DHUDB,store))
    mux.HandleFunc("/feedback",FeedBackFunc(DHUDB))
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
func HomeDeleteFunc(DBsession *mgo.Database,store *sessions.CookieStore) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        courseNo := req.PostFormValue("courseNo")
        stuid,ok := validateSession(req,store)
        if ok{
            cRigister := DBsession.C("RigisterInfo")
            err := cRigister.Find(bson.M{"studentid":stuid}).Update(bson.M{"$set":bson.M{"courselist.$.coursestate":courseDelete}})
            if err == nil{
                http.Redirect(w,req,"/home",http.StatusMovedPermanently)
                return
            }
        }
            http.Redirect(w,req,"/errMessage",http.StatusMovedPermanently)
    }
}
func LoginFunc(r *render.Render,DBsession *mgo.Database,store *sessions.CookieStore) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        id := req.PostFormValue("UserID")
        pw := req.PostFormValue("UserPassword")
        // cLogin := DBsession.C("StudentInfo")
        if validateLogin(id,pw,nil){
            session,_ := store.Get(req,"sessionid")
            session.Values["stuid"] = id
            session.Save(req, w)
            http.Redirect(w,req,"/home",http.StatusMovedPermanently)
        }else{
            http.Redirect(w,req,"/errMessage",http.StatusMovedPermanently)
        }
    }
}
func validateLogin(id,pw string,cLogin *mgo.Collection) (flag bool){
    //TODO Test it from database or request to the school website
    var err error
    err = cLogin.Find(bson.M{"studentid":id,"studentpw":pw})
    if err != nil{
        _,err = strconv.Atoi(id)
        if err != nil{
            return
        }

    }else{
        flag = true
    }
    return
}
func HomeFunc(r *render.Render,DBsession *mgo.Database,store *sessions.CookieStore) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        status := req.PostFormValue("userstatus")
        if validateHomeStatus(status){
            stuid,ok := validateSession(req,store)
            if ok{
                if status == ""{
                    hello := "Hello " + stuid + "!This is the home page of the website"
                    r.Text(w,http.StatusOK,hello)
                }else{
                    var err error
                    var done bool
                    var courselist []RigisteredCourse
                    for i := 0; i < 3;i++ {
                        cRigister := DBsession.C("RigisterInfo")
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
func validateSession(req *http.Request,store *sessions.CookieStore) (sessionid string,flag bool){
    session,_ := store.Get(req,"sessionid")
    id := session.Values["stuid"]
    stringid,ok := id.(string)
    if ok && stringid != ""{
        sessionid = stringid
        flag = true
    }
    return
}
func getrigisteredFunc(stuid string,cRigister *mgo.Collection) ([]RigisteredCourse,error){
    //TODO The collection will be nil here,so we will finish the database collection
    //and take it to the function
    //mgo.ErrNotFound
    return []RigisteredCourse{
                RigisteredCourse{"131441","专业英语","李悦",0,1},
                RigisteredCourse{"130153","计算机网络","朱明",0,1}},nil
}
func HomeSelectFunc(r *render.Render,DBsession *mgo.Database,store *sessions.CookieStore) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        coursetype := req.PostFormValue("coursetype")
        if validateCourseType(coursetype){
            sessionid,ok := validateSession(req,store)
            if ok{
                if coursetype == ""{
                    hello := "Hello " + sessionid + "!This is the home/select page of the website"
                    r.Text(w,http.StatusOK,hello)
                    return
                }else{
                    var done bool
                    var err  error
                    var teachSchemas []TeachSchema
                    for i := 0; i < 3; i++ {
                        cTable := DBsession.C("CourseTable")
                        cIndex := DBsession.C("CourseIndex")
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
func CommonSelectFunc(r *render.Render,DBsession *mgo.Database) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        id := req.PostFormValue("courseID")
        if id == ""{
            r.Text(w,http.StatusOK,"Hello World!This is the commonselection page of the website")
        }else{
            var done bool
            var err error
            var course CourseContent
            cTable := DBsession.C("CourseTable")
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
func HomeRegisterFunc(DBsession *mgo.Database,store *sessions.CookieStore) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        var slist SelectLists
        stuid,ok := validateSession(req,store)
        if ok{
            result,err := ioutil.ReadAll(req.Body)
            if err == nil{
                 err = json.Unmarshal([]byte(result), &slist)
                 if err == nil{
                     cTable := DBsession.C("CourseTable")
                     if validateCourse(slist,cTable){
                         //TODO saveAndRegister use the incorrect date struct
                         saveAndRegister(stuid,slist,DBsession)
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
func saveAndRegister(stuid string,slist SelectLists,DBsession *mgo.Database) {
    var err error
    var oldInfo StudentRigisterCourse
    var rigisterlist map[string]RigisterCourseInfo
    rInfo := StudentRigisterCourse{StudentID:stuid}
    cTable := DBsession.C("CourseTable")
    cRigister := DBsession.C("RigisterInfo")
    cCourseSelector := DBsession.C("CourseSelector")
    for _,coursevalue := range slist.SelectList{
        var teachername string
        var courseContent CourseContent
        err = cTable.Find(bson.M{"courseid":coursevalue.CourseID,"courselist.courseno":coursevalue.CourseNo}).One(&courseContent)
        //TODO roll back
        if err == nil{
            for _,courselist := range courseContent.CourseList{
                if coursevalue.CourseNo == courselist.CourseNo{
                    teachername = courselist.TeacherName
                    break
                }
            }
            Squeuenum := getQueueNum(coursevalue.CourseID,coursevalue.CourseNo,cCourseSelector) + 1
            rInfo.CourseList = append(rInfo.CourseList,RigisteredCourse{coursevalue.CourseID,courseContent.CourseName,teachername,courseInQueue,Squeuenum})
            rigisterlist[coursevalue.CourseID] = RigisterCourseInfo{stuid,coursevalue.CourseNo,Squeuenum}
        }
    }
    err = cRigister.Find(bson.M{"studentid":stuid}).One(&oldInfo)
    if err == nil{
        rInfo.CourseList = append(rInfo.CourseList,oldInfo.CourseList...)
        cRigister.Update(bson.M{"studentid":stuid},rInfo)
    }else{
        cRigister.Insert(rInfo)
    }
    // take a goroutine to avoid blocking
    go func() {
        for key,value := range rigisterlist{
            CourseMap[key].mutex.Lock()
            cCourseSelector.Insert(value)
            CourseMap[key].mutex.Unlock()
        }
        for _,value := range slist.SelectList{
            NotifyCourseID <- value.CourseID
        }
    }()
    return
}
func getQueueNum(courseID,courseNo string,cCourseSelector *mgo.Collection) (num int){
    var err error
    CourseMap[courseID].mutex.RLock()
    defer CourseMap[courseID].mutex.RUnlock()
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
func FeedBackFunc(DBsession *mgo.Database) ServePrimeFunc{
    return func (w http.ResponseWriter, req *http.Request){
        message := req.PostFormValue("Message")
        if message == ""{
            http.Redirect(w,req,"/errMessage",http.StatusMovedPermanently)
        }else{
            go storeFeedBackMessage(message,DBsession)
            http.Redirect(w,req,"/index",http.StatusMovedPermanently)
        }
    }
}
func storeFeedBackMessage(message string,DBsession *mgo.Database){
    //TODO That we will take the message into the database and notify the developer
    return
}
