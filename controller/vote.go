package controller

import (
	"web-app/logic"
	"web-app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 投票  本来是和帖子相关的 为了清晰一点，重新放一个文件

type VoteData struct {
	// UserID  谁发请求就是谁投票 这个可以不写  从请求中获取当前用户
	PostID    int64 `json:"post_id,string" binding:"required"` // 帖子id
	Direction int8  `json:"direction,string" `                 // 赞成票 1 反对票 -1 取消投票 0

}

func PostVoteController(c *gin.Context) {
	// @Summary      帖子投票
	// @Description  对帖子进行投票
	// @Tags         投票
	// @Accept       json
	// @Produce      json
	// @Security     ApiKeyAuth
	// @Param        body  body      models.ParamsVote  true  "投票参数"
	// @Success      200   {object}  ResponseData
	// @Router       /vote [post]
	// 参数校验
	p := new(models.ParamsVote)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return

	}
	// 获取当前请求用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
