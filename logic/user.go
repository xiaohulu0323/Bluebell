package logic

import (
	"web-app/dao/mysql"
	"web-app/models"
	"web-app/pkg/jwt"
	"web-app/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *models.ParamsSignUp) (err error) {
	// 1.判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		// 用户已存在
		return err
	}

	// 2.生成UID
	userID := snowflake.GenID()

	// 构造一个User 实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	// 3.保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamsLogin) (token string, err error) {
	// 1.根据用户名去数据库查询用户信息
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	// 传递的是指针，就能拿到user.UserID
	if err := mysql.Login(user); err != nil {
		return "", err
	}
	
	return jwt.GenToken(user.UserID, user.Username)
	
}
