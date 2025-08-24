package controller

import (
	"strconv"
	"web-app/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// --- 跟社区相关的 ---


func CommunityHandler(c *gin.Context){
	// 查询到所有的社区 （community_id,community_name）以列表的形式返回给前端	
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)   				// 不轻易把服务端报错暴露给外部
		return
	}
	ResponseSuccess(c, data)
}


// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(c *gin.Context){
	// 1. 获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)          // 转换成 十进制
	if err != nil{
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)   				// 不轻易把服务端报错暴露给外部
		return
	}
	ResponseSuccess(c, data)
}