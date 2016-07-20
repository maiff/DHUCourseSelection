- CourseTable  
    数据表中存储了本学期开设的所有课程的所有信息  
    数据格式为:  
    - type CourseContent struct{
        CourseID        string          `json:"courseID"`
        CourseName      string          `json:"courseName"`
        CourseList      []CourseList    `json:"courseList"`
    }
    - type CourseList struct{
        CourseNo        string          `json:"courseNo"`
        TeacherName     string          `json:"teacherName"`
        CourseInfo      []CourseInfo    `json:"courseInfo"`
    }
    - type CourseInfo struct{
        CourseWeek      string          `json:"courseWeek"`
        CourseTime      string          `json:"courseTime"`
    }
- CourseIndex  
    数据表中存储了不同年级不同专业课程信息作为索引,在选课渲染的时候可以通过CourseID来得到具体的CourseNo以及TeacherName
    - type CourseIndex struct{
        GradeMajor      string           `json:"gradeMajor"`
        CoureseTypeList []CourseTypeList `json:"courseTypeList"`
    }
    - type CourseTypeList struct{
        CourseType      string          `json:"courseType"`
        CourseList      []string        `json:"courseList"`
    }
- StudentInfo  
    数据表中存储了所有学生的账号密码  
    - type StudentInfo struct{
        StudentID       string
        StudentPW       string
    }
- RigisterInfo  
    数据表中存储了学生注册的所有课程  
    - type StudentRigisterCourse struct{
        StudentID       string
        CourseList      []RigisteredCourse
    }
    - type RigisteredCourse struct{
        CourseID        string          `json:"courseID"`
        CourseNo        string          `json:"courseNo"`
        CourseName      string          `json:"courseName"`
        TeacherName     string          `json:"teacherName"`
        CourseState     int          `json:"courseState"`
        // 0 means the Course in the queue
        //1 means the Course finished
        QueueNumber     int          `json:"queueNumber"`
    }
- CourseSelector  
    数据表中存放了监视器需要的选课信息
    - type RigisterCourseInfo struct{
        StudentID   string
        CourseNo    string
        QueueNum    int
    }
