package CourseSelection
import (
    "time"
    "errors"
    "strconv"
    "net/http"
    "github.com/PuerkitoBio/goquery"
)
const (
    DHUHostUrl          =   "http://jw.dhu.edu.cn/dhu"
    DHULoginUrl         =   "/login_wz.jsp"
    DHUSelectedUrl      =   "/student/selectcourse/seeselectedcourse.jsp"
    DHUMonitorUrl       =   "/commonquery/coursetimetableinfo.jsp"
    DHURegisterUrl      =   "/servlet/com.collegesoft.eduadmin.tables.selectcourse.SelectCourseController"
)
const (
    DHULoginNameKey     =   "userName"
    DHULoginPwdKey      =   "userPwd"
    DHUMonitorKey       =   "courseId"
    DHURegisterKey      =   "courseNo"
    DHUDefaultID        =   "141320131"
    DHUDefaultPW        =   "130681199507125816"
)
const (
    DHULoginSuccessURL  =   "http://jw.dhu.edu.cn/dhu/student/"
)
type DHUStruct struct {
    SchoolStruct
    SchoolName string
}
func NewDHUStruct() *DHUStruct{
    dhuStruct := &DHUStruct{SchoolName:"DHU",}
    dhuStruct.SchoolStruct = SchoolStruct{ErrChan : getErrChan(),Client : dhuStruct.defaultLogin()}
    return dhuStruct
}
func (d *DHUStruct) LoginPara(para ...string) map[string]string{
    loginPara := make(map[string]string)
    if len(para) != 2{
        d.SetErrorMessage("Parameter error in Login Action")
    }else{
        loginPara["Action"] = "Post"
        loginPara[DHULoginNameKey] = para[0]
        loginPara[DHULoginPwdKey] = para[1]
    }
    return loginPara
}
// func (d *DHUStruct) DeletePara(para ...string) map[string]string{
//
// }
func (d *DHUStruct) MonitorPara(para ...string) map[string]string{
    monitorPara := make(map[string]string)
    if len(para) != 1{
        d.SetErrorMessage("Parameter error in Monitor Action")
    }else{
        monitorPara["Action"] = "Get"
        monitorPara[DHUMonitorKey] = para[0]
    }
    return monitorPara
}
func (d *DHUStruct) RigisterPara(para ...string) map[string]string{
    registerPara := make(map[string]string)
    if len(para) != 1{
        d.SetErrorMessage("Parameter error in Register Action")
    }else{
        registerPara["Action"] = "Get"
        registerPara["doWhat"] = "selectcourse"
        registerPara["courseName"] = ""
        registerPara[DHURegisterKey] = para[0]
    }
    return registerPara
}
func (d *DHUStruct) Login(m map[string]string) (*http.Client,error){
    client := NewClient()
    res,err := sendRequest(DHUHostUrl + DHULoginUrl,m,client)
    if res != nil{
        loc,_ := res.Location()
        if loc.String() != DHULoginSuccessURL{
            err = passwordErr
        }
        res.Body.Close()
    }
    return client,err
}

func (d *DHUStruct) defaultLogin() *http.Client{
    var err error
    var client *http.Client
    for{
        client,err = d.Login(d.LoginPara(DHUDefaultID,DHUDefaultPW))
        if err == nil{
            break
        }
    }
    return client
}
// func (d *DHUStruct) Delete(map[string]string,*http.Client) error{
//
// }
func (d *DHUStruct) Register(m map[string]string,client *http.Client) error{
    res,err := sendRequest(DHUHostUrl + DHURegisterUrl,m,client)
    res.Body.Close()
    return err
}
func (d *DHUStruct) Monitor(m map[string]string) ([]string,error){
    var lessonlist []string
    d.mutexClient.RLock()
    res,err := sendRequest(DHUHostUrl + DHUMonitorUrl,m,d.Client)
    d.mutexClient.RUnlock()
    if err == nil{
        doc,_ := goquery.NewDocumentFromResponse(res)
        doc.Find("tr").Each(func (i int,s *goquery.Selection)  {
            lessonid := s.Find("td").Eq(0).Text()
            _,err := strconv.Atoi(lessonid)
            if err == nil {
                max,_ := strconv.Atoi(s.Find("td").Eq(2).Text())
                now,_ := strconv.Atoi(s.Find("td").Eq(4).Text())
                if max > now{
                    lessonlist = append(lessonlist,s.Find("td").Eq(0).Text())
                }
            }
        })
    }
    return lessonlist,err
}
func sendRequest(url string,m map[string]string,client *http.Client) (*http.Response,error){
    var err error
    var res *http.Response
    isPost,err := checkPost(m)
    delete(m,"Action")
    if err == nil{
        for i := 0;i < 3;i ++ {
            if isPost{
                res,err = client.PostForm(url,MakeParameters(m))
            }else{
                res,err = client.Get(url + getParaString(m))
            }
            if err == nil{
                break
            }
            time.Sleep(time.Duration(i * 2 + 1) * time.Second)
        }
        if err != nil{
            err = networkErr
        }
        return res,err
    }else{
        return nil,errors.New("Error in Login")
    }
}
func (d *DHUStruct) ValidateStuCourseSelected(courseid,courseNo string,client *http.Client) bool{
    var err error
    var res *http.Response
    var Done bool
    for i := 0; i < 3; i++ {
        res,err = client.Get(DHUHostUrl + DHUSelectedUrl)
        if err == nil{
            break
        }else{
            time.Sleep(time.Duration(i * 2 + 1) * time.Second)
        }
    }
    if err != nil{
        return false
    }
    doc,_ := goquery.NewDocumentFromResponse(res)
    doc.Find("tr").Each(func (i int,s *goquery.Selection){
        lessonid := s.Find("td").Eq(0).Text()
        _,err := strconv.Atoi(lessonid)
        if err == nil {
            if lessonid == courseid{
                Done = true
            }
        }
    })
    return Done
}
func (d *DHUStruct) UpdateClient() func() {
    var timeFroClient time.Time
    tickerForClient := time.NewTicker(time.Hour * updateTime)
    return func() {
        select{
            case timeFroClient = <- tickerForClient.C:
                newClient := d.defaultLogin()
                d.mutexClient.Lock()
                defer d.mutexClient.Unlock()
                d.Client = newClient
                return
            default :
                return
        }
    }
}
func checkPost(m map[string]string) (bool,error){
    action,ok := m["Action"]
    if ok{
        var flag bool
        if action == "Post"{
            flag = true
        }else if action == "Get"{
            flag = false
        }
        return flag,nil
    }
    return false,errors.New("Action Wrong")
}
func getParaString(m map[string]string) string{
    para := "?"
    for key,value := range m{
        if len(para) == 1{
            para = para + key + "=" + value
        }else{
            para = para + "&" + key + "=" + value
        }
    }
    return para
}
