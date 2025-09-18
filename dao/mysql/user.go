package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"web-app/models"

	"go.uber.org/zap"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

const secret = "xp"

// CheckUserExist 检查用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}

	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)

	// 执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?, ?, ?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))

	return hex.EncodeToString(h.Sum([]byte(oPassword))) // 转换成16进制字符串
}

func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户登录的密码

	sqlStr := `select user_id, username, password from user where username = ?`
	err = db.Get(user, sqlStr, user.Username) // Get()将数据库查询结果填充到指针指向的结构体
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}

	if err != nil {
		// 查询用户信息失败
		return err
	}
	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return

}

// GetUserByID 根据id 获取用户信息
func GetUserByID(userID int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, userID)

	return
}

// BatchGetUsersByIDs 批量根据用户ID列表获取用户信息
// 解决N+1查询问题的核心函数
func BatchGetUsersByIDs(userIDs []int64) (userMap map[int64]*models.User, err error) {
	if len(userIDs) == 0 {
		return make(map[int64]*models.User), nil
	}

	// 防止IN查询参数过多导致性能问题
	const maxBatchSize = 1000
	if len(userIDs) > maxBatchSize {
		return nil, fmt.Errorf("too many user IDs: %d, max allowed: %d", len(userIDs), maxBatchSize)
	}

	// 去重用户ID
	uniqueIDs := make([]int64, 0, len(userIDs))
	idSet := make(map[int64]bool)
	for _, id := range userIDs {
		if !idSet[id] {
			uniqueIDs = append(uniqueIDs, id)
			idSet[id] = true
		}
	}

	// 构建IN查询
	sqlStr := `select user_id, username from user where user_id in (?` + strings.Repeat(`,?`, len(uniqueIDs)-1) + `)`

	// 准备参数
	args := make([]interface{}, len(uniqueIDs))
	for i, id := range uniqueIDs {
		args[i] = id
	}

	// 执行查询
	var users []*models.User
	err = db.Select(&users, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("batch get users failed: %w", err)
	}

	// 转换为map便于查找
	userMap = make(map[int64]*models.User, len(users))
	for _, user := range users {
		userMap[user.UserID] = user
	}

	// 检查是否所有用户都查询到了
	if len(userMap) != len(uniqueIDs) {
		missingIDs := make([]int64, 0)
		for _, id := range uniqueIDs {
			if _, exists := userMap[id]; !exists {
				missingIDs = append(missingIDs, id)
			}
		}
		// 只记录警告，不返回错误，因为某些用户可能被删除了
		zap.L().Warn("Some users not found in batch query",
			zap.Int64s("missing_user_ids", missingIDs),
			zap.Int("requested_count", len(uniqueIDs)),
			zap.Int("found_count", len(userMap)))
	}

	return userMap, nil
}
