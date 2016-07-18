package main
import (
    "fmt"
    "time"
    "strconv"
    "mahonia"
    "net/url"
    "net/http"
    "net/http/cookiejar"
    "gopkg.in/mgo.v2"
    // "gopkg.in/mgo.v2/bson"
    "github.com/PuerkitoBio/goquery"
)
type CourseContent struct{
    CourseID        string          `json:"courseID"`
    CourseName      string          `json:"courseName"`
    CourseList      []CourseList    `json:"courseList"`
}
type CourseList struct{
    CourseNo        string          `json:"courseNo"`
    TeacherName     string          `json:"teacherName"`
    CourseInfo      []CourseInfo    `json:"courseInfo"`
}
type CourseInfo struct{
    CourseWeek      string          `json:"courseWeek"`
    CourseTime      string          `json:"courseTime"`
}
const (
    DHUHostUrl          =   "http://jw.dhu.edu.cn/dhu"
    DHULoginUrl         =   "/login_wz.jsp"
    DHUCommonqueryUrl   =   "/commonquery/selectcoursetermcourses.jsp?pageSize=2200"
    DHUCourseTableUrl   =   "/commonquery/coursetimetableinfo.jsp?courseId="
    DHUTeachSchemaUrl   =   "/commonquery/teachschemasquery.jsp"
    //Parameters: "gradeYear"  "majorId"
)
func main() {
    var key   string
    var value string
    session,err := mgo.Dial("localhost:27017")
    if err != nil {
        panic(err)
    }
    defer session.Close()
    c := session.DB("DHU").C("CourseTable")
    content := []CourseContent{}
    client  := login()
    if client != nil{
        res,err := client.Get(DHUHostUrl + DHUCommonqueryUrl)
        if err != nil{
            fmt.Println(err)
        }else{
            CourseList := GetAllCourse(res)
            for key,value = range CourseList{
                // fmt.Println(key)
                // fmt.Println(value)
                res,err := client.Get(DHUHostUrl + DHUCourseTableUrl + key)
                if err != nil{
                    fmt.Println(err)
                    break
                }else{
                    newlist := GetAllLessons(res)
                    newContent := CourseContent{key,value,newlist}
                    content = append(content,newContent)
                }
                time.Sleep(1 * time.Second)
            }
            for _,values := range content{
                c.Insert(values)
            }
        }
    }
}
func GetAllLessons(res *http.Response) []CourseList{
    doc,err := goquery.NewDocumentFromResponse(res)
    if err != nil{
        fmt.Println(err)
        return nil
    }else{
        newlist := []CourseList{}
        dec := mahonia.NewDecoder("GB18030")
        doc.Find("tr").Each(func (i int,s *goquery.Selection){
            selectid := s.Find("td").Eq(0).Text()
            _,err := strconv.Atoi(selectid)
            if err == nil{
                _,teachername,_ := dec.Translate([]byte(s.Find("td").Eq(6).Text()),true)
                newinfo := []CourseInfo{}
                s.Find("td").Eq(7).Find("tr").Each(func (i int,s *goquery.Selection){
                    _,weektime,_ := dec.Translate([]byte(s.Find("td").Eq(0).Text()),true)
                    _,daytime,_ := dec.Translate([]byte(s.Find("td").Eq(1).Text()),true)
                    info := CourseInfo{string(weektime),string(daytime)}
                    newinfo = append(newinfo,info)
                })
                courselist := CourseList{selectid,string(teachername),newinfo}
                newlist = append(newlist,courselist)
            }else{
                return
            }
        })
        return newlist
    }
}
func GetAllCourse(res *http.Response) map[string]string{
    doc,err := goquery.NewDocumentFromResponse(res)
    if err != nil{
        fmt.Println(err)
        return nil
    }else{
        CourseList := map[string]string{}
        dec := mahonia.NewDecoder("GB18030")
        doc.Find("tr").Each(func (i int,s *goquery.Selection){
            courseid := s.Find("td").Eq(1).Text()
            _,err := strconv.Atoi(courseid)
            if err == nil{
                coursename := s.Find("td").Eq(0).Text()
                _,data,_ := dec.Translate([]byte(coursename),true)
                CourseList[courseid] = string(data)
            }else{
                return
            }
        })
        return CourseList
    }
}





func login() *http.Client{
    client := client_with_cookiejar()
    value := map[string]string{"userName":"141320131","userPwd":"130681199507125816"}
    urlvalue := url_value(value)
    _,err := client.PostForm(DHUHostUrl + DHULoginUrl,urlvalue)
    if err != nil{
        fmt.Println(err)
        fmt.Println("Fuck it!Something wrong in the login function!")
        return nil
    }else{
        return client
    }
}
func url_value(para map[string]string)  url.Values{
    data := make(url.Values)
    for key,value := range para{
        data.Set(key,value)
    }
    return data
}
//Return the http client with cookiejar so it can keep the cookie
func client_with_cookiejar() *http.Client {
    jar,_ := cookiejar.New(nil)
    client := &http.Client{
        Jar:jar,
    }
    return client
}
