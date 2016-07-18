package CourseSelection
import (
    // "fmt"
    "errors"
    "strconv"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)
type TeachSchema struct{
    CourseType      string          `json:"courseType"`
    CourseContent   []CourseContent `json:"courseContent"`
}
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
var (
    StatusTypeList = []string{"1","2","3"}
)
var (
    ErrHomeSelectTypeNotFound = errors.New("The course status type not found")
    ErrHomeSelectStatusIllegal = errors.New("Illegal parameter of statustype")
    ErrHomeSelectCourseNotFound = errors.New("The course not found in database table CourseTable")
    ErrCommentSelectParameterIllegal = errors.New("Illegal parameter of course ID")
)
func APICommonselect(cTable *mgo.Collection,courseid string) (CourseContent,error){
    var courselist CourseContent
    var err error
    _,err = strconv.Atoi(courseid)
    if err != nil{
        return courselist,ErrCommentSelectParameterIllegal
    }else{
        return apiCommonselect(cTable,courseid)
    }
}
func apiCommonselect(cTable *mgo.Collection,courseid string) (CourseContent,error){
    var courselist CourseContent
    err := cTable.Find(bson.M{"courseid":courseid}).One(&courselist)
    return courselist,err
}
func APIHomeSelect(cTable,cIndex *mgo.Collection,gradeMajor,statusType string) ([]TeachSchema,error){
    var err error
    var errflag bool
    var teachSchemas   []TeachSchema
    var coursecontent  []CourseContent
    // var coursecontents []CourseContent
    _,err = strconv.Atoi(statusType)
    if err != nil{
        return teachSchemas,ErrHomeSelectStatusIllegal
    }
    _,err = strconv.Atoi(gradeMajor)
    if err != nil{
        return teachSchemas,ErrHomeSelectStatusIllegal
    }
    if statusType == "0"{
        for _,index := range StatusTypeList{
            var teachSchema    TeachSchema
            coursecontent,err = apiHomeSelect(cTable,cIndex,gradeMajor,index)
            switch err {
            case nil:
                teachSchema.CourseType = index
                teachSchema.CourseContent = coursecontent
                teachSchemas = append(teachSchemas,teachSchema)
            case ErrHomeSelectTypeNotFound:
                continue
            case ErrHomeSelectCourseNotFound:
                errflag = true
                teachSchema.CourseType = index
                teachSchema.CourseContent = coursecontent
                teachSchemas = append(teachSchemas,teachSchema)
            default:
                return teachSchemas,err
            }
        }
        if errflag{
            return teachSchemas,ErrHomeSelectCourseNotFound
        }else{
            // fmt.Println(err)
            return teachSchemas,err
        }
    }else{
        var teachSchema    TeachSchema
        coursecontent,err = apiHomeSelect(cTable,cIndex,gradeMajor,statusType)
        teachSchema.CourseType = statusType
        teachSchema.CourseContent = coursecontent
        teachSchemas = append(teachSchemas,teachSchema)
        return teachSchemas,err
    }
}
func apiHomeSelect(cTable,cIndex *mgo.Collection,gradeMajor,statusType string) ([]CourseContent,error){
    var err error
    var errflag bool
    var courseList CourseTypeList
    var courseContent  CourseContent
    var courseContents []CourseContent
    MajorStruct := CourseIndex{}
    err = cIndex.Find(bson.M{"grademajor":gradeMajor}).One(&MajorStruct)
    if err != nil{
        return courseContents,err
    }
    for _,TestForCourse := range MajorStruct.CoureseTypeList{
        if TestForCourse.CourseType == statusType{
            courseList = TestForCourse
            break
        }
    }
    if courseList.CourseType == ""{
        return courseContents,ErrHomeSelectTypeNotFound
    }
    for _,CourseID := range courseList.CourseList{
        if CourseID == ""{
            continue
        }
        err = cTable.Find(bson.M{"courseid":CourseID}).One(&courseContent)
        if err != nil{
            // fmt.Println(err)
            // fmt.Println(CourseID)
            errflag = true
        }else{
            courseContents = append(courseContents,courseContent)
        }
    }
    if errflag{
        return courseContents,ErrHomeSelectCourseNotFound
    }else{
        return courseContents,nil
    }
}
