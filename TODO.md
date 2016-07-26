TODO List
====================
- WEB
    - -----分割线
    - 权限控制
    - 登陆验证行为(向教务处发起请求并确认账号密码正确)
    - 账号密码数据库表加盐
    - 将Session信息完善(session clone,session mode,session safe)
- 监视器
    - School Struct
    - 查询课程的goroutine是否存在需要加读锁,退出goroutine的监视器需要加入写锁然后复查数据库是否为空
    - 成功选课后同步数据库
    
