package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

// 全局 node 实例，用于生成唯一ID
var node *sf.Node

// 初始化雪花算法节点
// startTime: 起始时间（格式如 "2020-07-01"），machineID: 机器ID
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime) // 解析起始时间
	if err != nil {
		return
	}

	// 设置雪花算法的起始时间戳（毫秒）
	sf.Epoch = st.UnixNano() / 1000000
	// 创建节点实例，machineID 用于分布式唯一性
	node, err = sf.NewNode(machineID)
	return
}

// 生成唯一ID
func GenID() int64 {
	return node.Generate().Int64()
}
