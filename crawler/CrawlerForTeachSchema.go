package main
import (
    "fmt"
    // "time"
    "strconv"
    "mahonia"
    "net/url"
    "net/http"
    "net/http/cookiejar"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
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
type CourseIndex struct{
    GradeMajor      string           `json:"gradeMajor"`
    CoureseTypeList []CourseTypeList `json:"courseTypeList"`
}
type CourseTypeList struct{
    CourseType      string          `json:"courseType"`
    CourseList      []string        `json:"courseList"`
}
const (
    DHUHostUrl          =   "http://jw.dhu.edu.cn/dhu"
    DHULoginUrl         =   "/login_wz.jsp"
    DHUCommonqueryUrl   =   "/commonquery/selectcoursetermcourses.jsp?pageSize=2200"
    DHUCourseTableUrl   =   "/commonquery/coursetimetableinfo.jsp?courseId="
    DHUTeachSchemaUrl   =   "/commonquery/teachschemasquery.jsp"
    DHUGetMajorURL      =   "/commonquery/selectgradeyearmajor.jsp"
    //Parameters: "gradeYear"  "majorId"
    BreakTitle          =   "实践教学"
)
func main() {
    gradelist   := []string{"2012a","2013a","2014a","2015a","2016a"}
    client := login()
    getmajorlist := GetMajorList(client)
    GetTeachSchema(client,gradelist,getmajorlist)
}
func GetTeachSchema(client *http.Client,gradelist,majorlist []string) {
    courseTypeList := map[string]string{"政治法律":"2","自然科学":"3","语言文字":"3","other":"1"}
    dec := mahonia.NewDecoder("GB18030")
    session,err := mgo.Dial("localhost:27017")
    if err != nil {
        panic(err)
    }
    defer session.Close()
    courseTable := session.DB("DHU").C("CourseTable")
    courseindexDB := session.DB("DHU").C("CourseIndex")
    for _,grade := range gradelist{
        for _,major := range majorlist{
            var ok bool
            var titletype string
            courseindexs := CourseIndex{grade[2:4] + major[2:],[]CourseTypeList{}}
            Other := CourseTypeList{"1",make([]string,1)}
            PolicyAndLaw := CourseTypeList{"2",make([]string,1)}
            NatureAndLanguage := CourseTypeList{"3",make([]string,1)}
            res,err := client.PostForm(DHUHostUrl + DHUTeachSchemaUrl,url_value(map[string]string{"gradeYear":grade,"majorId":major}))
            if err != nil{
                fmt.Println(err)
                return
            }
            doc,_ := goquery.NewDocumentFromResponse(res)
            doc.Find("tr").EachWithBreak(func (i int,s *goquery.Selection) bool{
                _,titleForBreak,_ := dec.Translate([]byte(s.Find("td").Eq(0).Text()),true)
                if string(titleForBreak) == BreakTitle{
                    return false
                }else{
                    element := s.Find("td").Eq(1).Text()
                    _,err := strconv.Atoi(element)
                    if err != nil{
                        _,title,_:= dec.Translate([]byte(element),true)
                        titletype,ok = courseTypeList[string(title)]
                        if !ok{
                            titletype = courseTypeList["other"]
                        }
                    }else{
                        // fmt.Println(element)
                        err := courseTable.Find(bson.M{"courseid":element}).One(nil)
                        if err == nil{
                            // for index := 9;index < 17; index ++ {
                            //     if s.Find("td").Eq(index).Text() != " "{
                            //         if index % 2 == 0{
                            //             return true
                            //         }else{
                            //             break
                            //         }
                            //     }else{
                            //         continue
                            //     }
                            // }
                            switch titletype {
                            case "1":
                                Other.CourseList = append(Other.CourseList,element)
                            case "2":
                                PolicyAndLaw.CourseList = append(PolicyAndLaw.CourseList,element)
                            case "3":
                                NatureAndLanguage.CourseList = append(NatureAndLanguage.CourseList,element)
                            default :
                                fmt.Println("Something Wrong In Switch")
                            }
                        }else{
                            return true
                        }
                    }
                }
                return true
            })
            courseindexs.CoureseTypeList = append(courseindexs.CoureseTypeList,Other)
            courseindexs.CoureseTypeList = append(courseindexs.CoureseTypeList,PolicyAndLaw)
            courseindexs.CoureseTypeList = append(courseindexs.CoureseTypeList,NatureAndLanguage)
            courseindexDB.Insert(courseindexs)
            // fmt.Println(grade)
            // fmt.Println(major)
            // fmt.Println(Other)
            // fmt.Println(PolicyAndLaw)
            // fmt.Println(NatureAndLanguage)
        }
    }
}
func GetMajorList(client *http.Client) []string{
    banlist      := map[string]string{"110431":"","110432":"","110511":"","110521":"","113110":"","113111":""}
    getmajorlist := []string{}
    res,err := client.Get(DHUHostUrl + DHUGetMajorURL)
    if err != nil{
        fmt.Println(err)
        fmt.Println("Something Wrong In GetMajorList")
    }else{
        doc,_ := goquery.NewDocumentFromResponse(res)
        doc.Find("select").Eq(1).Find("option").Each(func (i int,s *goquery.Selection){
            sss,ok := s.Attr("value")
            if ok{
                _,testmap := banlist[sss]
                if testmap{
                    return
                }else{
                    getmajorlist = append(getmajorlist,sss)
                }
            }else{
                fmt.Println("Something Wrong")
            }
        })
        return getmajorlist
    }
    return nil
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
