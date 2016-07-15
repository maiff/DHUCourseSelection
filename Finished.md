Finished Job
=============
- Session实现
- 数据库









#### Session实现
Session实现为登陆时通过Set-Cookie:session-name=XXXXX.  
其中不验证的页面有
- /
- /index
- /feedback
- /login
- /errMessage  

其余页面都做了验证机制，如果判断session失效则重定向到/index  
如果密码错误则跳转到/errMessage(希望可以用通讯的方式将密码错误信息显示在登陆框上)  
登陆信息为:UserID : (数字就可以) ; UserPassword : (随便填,必须不为空)
PS:其实只是简单的能用了...

#### 数据库
- 所有专业的数据库已经收录,(不包括留学生以及文化素质,体育,外语类课程)
- 登陆信息更新了，请看Session实现
- 判断返回课程的方式为截取登陆账号的前六位,比如我的学号141320131,那么会截取141320。
- 备注:
    - 体育外语等课程我暂时没找到链接可以拿到目录
    - 名字叫SH(School Helper)怎么样?我觉得咱们尽量不要局限于东华大学，也不要局限于选课这一个功能
