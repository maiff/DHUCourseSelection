东华大学选课文档
===============

#### 欢迎页面
- URL : /index
- 概括 : 使用帮助以及登陆
- 按钮 :
    - 登陆
    - 留言
- 可选 : 链接->东华大学教务处

#### 登陆
- URL : /login
- 概括 : 登陆(弹窗形式 or 直接将form写在index.html中)
- 参数(Post) :
    - UserID       -> 用户账号
    - UserPassword -> 用户密码
    - ActionURL    -> /login

#### 用户个人主页
- URL : /home
- 概括 : 用户个人主页
- 内容 : 用户已申请排队选取课程
- 链接 :
    - 用户选课界面
    - 按照选课序号选课
- JSON :
    - 目的 : 从数据库中读出该用户在网站上申请的课程并展示在页面上.
    - APIURL : /home?userstatus=0   
        - userstatus = 0 则返回所有已选课程
        - 其他值留着以后加其他功能
    - 格式 :
```
    {
        "RegisterCourse" : [
            {"courseID": ,"courseName": ,"teacherName": ,"courseState": ,"queueNumber": ,},
            {"courseID": ,"courseName": ,"teacherName": ,"courseState": ,"queueNumber": ,},
            ...]
    }
```
    - 参数说明 :
        - courseID : 课程序号
        - courseName : 课程名称
        - teacherName : 教师姓名
        - courseState : 课程状态(排队中 or 已选上 or 被取消)
        - queueNumber : 排队序号

#### 用户选课界面
- URL : /home/select
- 概括 : 用户选课界面
- 内容 : 用户教学计划查询中本学期开设的课程(包括公共课程以及专业课程)
- JSON :
    - 目的 : 渲染可选课程
    - APIURL :/home/select?coursetype=0  
        - coursetype = 0 则返回所有的课程
        - coursetype = 1 选修和必修
        - coursetype = 2 政治法律
        - coursetype = 3 自然科学
        - coursetype = 4 文化素质
        - coursetype = 5 体育
        - coursetype = 6 外语

    - 格式 :
```
    {
        "TeachSchema":[
            {"courseType":,
                "courseContent":[
                {"courseID":,"courseName":,"courseList":[
                    {"courseNo":,"teacherName":,"courseWeek":,"courseTime":,}
                    {"courseNo":,"teacherName":,"courseWeek":,"courseTime":,}
                    ...]}
                {"courseID":,"courseName":,"courseList":[
                    {"courseNo":,"teacherName":,"courseWeek":,"courseTime":,}
                    {"courseNo":,"teacherName":,"courseWeek":,"courseTime":,}
                    ...]}
                ...]
            },
            {"courseType":,
                "courseContent":[
                {"courseID":,"courseName":,"courseList":[
                    {"courseNo":,"teacherName":,"courseWeek":,"courseTime":,}
                    {"courseNo":,"teacherName":,"courseWeek":,"courseTime":,}
                    ...]}
                {"courseID":,"courseName":,"courseList":[
                    {"courseNo":,"teacherName":,"courseInfo":
                        [{"courseWeek":,"courseTime":},{"courseWeek":,"courseTime":}...]}
                    {"courseNo":,"teacherName":,"courseInfo":
                        [{"courseWeek":,"courseTime":},{"courseWeek":,"courseTime":}...]}
                    ...]}
                ...]
            },
            ...
        ]
    }
```

    - 参数说明 :
        - courseType : 课程类型(必修,选修,政治法律,自然科学,文化素质,体育,外语)
        - courseID : 课程代码
        - courseName : 课程名称
        - courseNo : 选课序号
        - teacherName : 教师姓名
        - courseWeek : 上课周次
        - courseTime : 上课时间
- 参数(JSON) :
    - TargetURL -> /home/register
    - 格式       
```
    {
        "SelectList":[
            {"courseID": ,"courseNo":},
            {"courseID": ,"courseNo":},
            ...
        ]
    }
```
    - 参数说明
        - courseID  -> 课程编码
        - courseNo  -> 选课序号
- 备注 :
    - 数据做成按格式可收缩，层次和JSON相同,课程类型(如选课对照表里的政治法律那一栏)->课程名称->可选课程列表.
      [参考链接](http://zhidao.baidu.com/link?url=08Zuu4QEF_VI1yO4ck0qWfRzRGENZeyEodd_UYCbxm8JgocuxFBu9Ji3YdO4R8U6j5tFs9D5E36gI-WUNu8GE_)
    - 将选择做成复选框的形式，点击发送后用js处理成JSON格式
    - 传送数据也用JSON

#### 按照序号选课
- URL : /commonselect
- 概括 : 通过courseID查询课程并选课
- 内容 : 查询框，发送请求后在查询框下面渲染出JSON格式的课程列表.
- JSON :
    - 目的 : 渲染用户请求的课程
    - APIURL(Post) : /commonselect
        - 参数: courseName
    - 格式 :
```
        {"courseID":,"courseName":,"courseList":[
            {"courseNo":,"teacherName":,"courseWeek":,"courseTime":,}
            {"courseNo":,"teacherName":,"courseWeek":,"courseTime":,}
            ...]}
```
- 参数 :
    - courseName  -> 课程名称
    - 参数说明 :
        - 点击查询后用ajax的方式将input里面的内容发送出去
- 参数 :
    - 与用户选课界面(/home/select)发送参数的格式相同
    - 参数说明 :
        - 点击提交后将数据以JSON的形式发送出去
        - TargetURL -> /register
- 备注 :
    - 在无参数请求这个页面的状态下,页面的主体部分只有一个搜索框
    - 在点击搜索后，后端接受到courseName后返回JSON,将JSON渲染到页面上(页面其他部分不变),并产生发送按钮
    - 在用户选好课程以后点击发送将数据以JSON的格式发送到指定URL

#### 留言(联系我们)
- URL : /feedback
- 概括 : 获取用户意见
- 参数(Post) :
    - Message -> 反馈信息，前端限制不要超过60个汉字
    - TargetURL -> /feedback
- 备注 :
    - 做成弹窗形式

#### 错误信息
- URL : /errMessage
- 概括 : 发生错误后重定向到错误信息
- 备注 :
    - 可以留下联系方式，邮箱XXX
    - 其他的随便写吧...
