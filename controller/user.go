package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"web-app/dao/mysql"
	"web-app/logic"
	"web-app/models"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// SignUpHandler 处理注册请求的函数
// @Summary      用户注册
// @Description  注册新用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        body  body      models.ParamsSignUp  true  "注册参数"
// @Success      200   {object}  ResponseData
// @Failure      200   {object}  ResponseData
// @Router       /signup [post]
func SignUpHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	p := new(models.ParamsSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断 err 是不是validator
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// // 手动对请求参数进行参数校验
	// if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0|| p.Password != p.RePassword {
	// 	zap.L().Error("SignUp with invalid param")
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"msg": "请求参数有误",
	// 	})
	// 	return
	// }

	fmt.Println(p)
	// 2.业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("SignUp failed", zap.Error(err))

		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
		}

		ResponseError(c, CodeServerBusy)
		return
	}

	// 3.返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 用户登录
// @Summary      用户登录
// @Description  获取登录 Token
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        body  body      models.ParamsLogin  true  "登录参数"
// @Success      200   {object}  ResponseData
// @Failure      200   {object}  ResponseData
// @Router       /login [post]
func LoginHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamsLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断 err 是不是 validator.ValidationErrors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2. 业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}

		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"user_id":  strconv.FormatInt(user.UserID, 10), // id 值大于 1<<53-1 （JSON        int64类型的最大值 1<<63-1
		"username": user.Username,
		"token":    user.Token,
	})
}

// 去除提示信息中的结构体名称
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
