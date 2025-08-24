package logic

import (
	"web-app/dao/mysql"
	"web-app/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查询数据库 查找到所有的community 并返回

	return mysql.GetCommunityList()

}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	// 查询数据库
	return mysql.GetCommunityDetailByID(id)
}
