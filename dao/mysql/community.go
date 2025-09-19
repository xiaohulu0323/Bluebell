package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"web-app/models"

	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	// 查询数据库 查找到所有的community 并返回
	sqlStr := `select community_id,community_name from community`
	// 读操作使用读数据库
	readDB := GetReadDB()
	if err := readDB.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

// GetCommunityDetailByID 根据ID查询分类详情
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	// 查询数据库
	community = new(models.CommunityDetail)
	sqlStr := `select community_id,community_name, introduction,create_time 
				from community
				where community_id = ?`
	// 读操作使用读数据库
	readDB := GetReadDB()
	if err := readDB.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = ErrorInvalidID
		}
	}
	return community, err
}

// BatchGetCommunitiesByIDs 批量根据社区ID列表获取社区详情信息
// 解决N+1查询问题的核心函数
func BatchGetCommunitiesByIDs(communityIDs []int64) (communityMap map[int64]*models.CommunityDetail, err error) {
	if len(communityIDs) == 0 {
		return make(map[int64]*models.CommunityDetail), nil
	}

	// 防止IN查询参数过多导致性能问题
	const maxBatchSize = 1000
	if len(communityIDs) > maxBatchSize {
		return nil, fmt.Errorf("too many community IDs: %d, max allowed: %d", len(communityIDs), maxBatchSize)
	}

	// 去重社区ID
	uniqueIDs := make([]int64, 0, len(communityIDs))
	idSet := make(map[int64]bool)
	for _, id := range communityIDs {
		if !idSet[id] {
			uniqueIDs = append(uniqueIDs, id)
			idSet[id] = true
		}
	}

	// 构建IN查询
	sqlStr := `select community_id, community_name, introduction, create_time 
			   from community 
			   where community_id in (?` + strings.Repeat(`,?`, len(uniqueIDs)-1) + `)`

	// 准备参数
	args := make([]interface{}, len(uniqueIDs))
	for i, id := range uniqueIDs {
		args[i] = id
	}

	// 执行查询
	var communities []*models.CommunityDetail
	// 读操作使用读数据库
	readDB := GetReadDB()
	err = readDB.Select(&communities, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("batch get communities failed: %w", err)
	}

	// 转换为map便于查找
	communityMap = make(map[int64]*models.CommunityDetail, len(communities))
	for _, community := range communities {
		communityMap[community.ID] = community
	}

	// 检查是否所有社区都查询到了
	if len(communityMap) != len(uniqueIDs) {
		missingIDs := make([]int64, 0)
		for _, id := range uniqueIDs {
			if _, exists := communityMap[id]; !exists {
				missingIDs = append(missingIDs, id)
			}
		}
		// 只记录警告，不返回错误，因为某些社区可能被删除了
		zap.L().Warn("Some communities not found in batch query",
			zap.Int64s("missing_community_ids", missingIDs),
			zap.Int("requested_count", len(uniqueIDs)),
			zap.Int("found_count", len(communityMap)))
	}

	return communityMap, nil
}
