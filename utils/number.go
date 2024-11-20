package utils

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/sony/sonyflake"
)

// Random 得到一个范围的随机数
func Random(num ...int) int {
	minInt := 0
	maxInt := 0

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(num) == 0 {
		return r.Int()
	}
	if len(num) == 1 {
		maxInt = num[0]
	} else if len(num) > 1 {
		minInt = num[0]
		maxInt = num[1]
	}

	if minInt > maxInt {
		minInt, maxInt = maxInt, minInt
	}
	if maxInt == 0 || maxInt == minInt {
		return maxInt
	}

	return r.Intn(maxInt-minInt) + minInt
}

var nodeMap sync.Map
var newSony *sonyflake.Sonyflake

// GetID 获取唯一ID，base为转换进制
func GetID(base int) (string, error) {
	if base < 2 || base > 36 {
		return "", fmt.Errorf("base from 2 to 36")
	}
	retInt, err := GetInt64ID()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(retInt, base), nil
}

// GetInt64ID 获取int64ID
func GetInt64ID() (int64, error) {
	snowflake.Epoch = 1609430400000 // 2021-01-01 00:00:00
	snowflake.NodeBits = 10
	snowflake.StepBits = 12
	nodeNum := Random(1, 20)

	var node *snowflake.Node
	var err error
	if nodeIn, ok := nodeMap.Load(nodeNum); ok {
		nodeTemp, ok := nodeIn.(*snowflake.Node)
		if ok {
			node = nodeTemp
		}
	}

	if node == nil {
		node, err = snowflake.NewNode(int64(nodeNum))
		if err == nil && node != nil {
			nodeMap.Store(nodeNum, node)
		}
	}

	if node != nil {
		// Generate a snowflake ID.
		id := node.Generate()
		return id.Int64(), nil
	}

	// sonyFlake 不用保存参数
	if newSony == nil {
		newSony = sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime:      time.Time{},
			MachineID:      nil,
			CheckMachineID: nil,
		})
	}

	if newSony == nil {
		// 这个老报这个错，调试困难
		return 0, fmt.Errorf("init NewSonyflake error")
	}

	intID, err := newSony.NextID()
	if err != nil {
		return 0, err
	}
	return int64(intID), nil
}

// Checksum 分表等后缀方法使用，用于获取某一个特定的数字
func Checksum(key string, maxSize int) int {
	if maxSize <= 0 {
		return 0
	}
	return int(crc32.ChecksumIEEE([]byte(key)) % uint32(maxSize))
}
