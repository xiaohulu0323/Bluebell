package logic

import (
	"web-app/dao/mysql"
	"web-app/models"
	"web-app/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *models.SignUpParams) {
	// 1.判断用户存不存在
	mysql.QueryUserByUsername()

	// 2.生成UID
	snowflake.GenID()

	// 3.保存进数据库
	mysql.InsertUser()
}