package tdog

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

/**
*  时间戳占41 bit 机器ID（机房ID和机器ID）6 bit  毫秒内产生的ID数  16 bit
 * 1 23456781234567812345678123456781234567812  345681                          2345567812345678
  * 雪花算法
*/
const (
	//初始化时间值 2^41 - 1   差不多可以用69年
	timeInitValue = 1585644268888
	//机器ID
	machineIdBits = 4
	//机房ID
	computerRoomIdBits = 2
	//每毫秒产生的ID数
	sequenceBits = 6
)

var (
	LastTimestamp = -1
	Lock          sync.Mutex
)

type Snowflake struct {
	//机器ID
	WorkerId int
	//机房ID
	DatacenterId int
	//代表一毫秒内生成的多个id的最新序号 最多12位 共 4096 -1 = 4095 个id
	Sequence int
}

// 初始化雪花算法模块
func NewSnowflake(workId, dataId, seq int) *Snowflake {
	return &Snowflake{
		WorkerId:     workId,
		DatacenterId: dataId,
		Sequence:     seq,
	}
}

func getCurrentTime() int {
	return int(time.Now().Unix())
}

func tilNextMillis(lastStamp int) int {
	for {
		timeStamp := getCurrentTime()
		if timeStamp > lastStamp {
			return timeStamp
		}
	}
}

func (sf *Snowflake) Get() (int64, error) {
	Lock.Lock()
	defer Lock.Unlock()
	workerIdShift := sequenceBits
	datacenterIdShift := sequenceBits + machineIdBits
	timestampLeftShift := sequenceBits + machineIdBits + computerRoomIdBits
	sequenceMask := (1 << sequenceBits) - 1
	timeStamp := getCurrentTime()
	if timeStamp < LastTimestamp {
		return 0, errors.New("系统时钟发生倒退，生成ID异常，请仔细检查。")
	}
	//如果两个时间相同，则自增一
	if LastTimestamp == timeStamp {
		sf.Sequence = (sf.Sequence + 1) & sequenceMask
		//当某一毫秒的时间，产生的id数 超过4095，系统会进入等待，直到下一毫秒，系统继续产生ID
		if sf.Sequence == 0 {
			timeStamp = tilNextMillis(LastTimestamp)
		}
	} else {
		sf.Sequence = 0
	}
	//记录最后一次产生的时间戳
	LastTimestamp = timeStamp
	num := ((timeStamp - timeInitValue) << timestampLeftShift) | (sf.DatacenterId << datacenterIdShift) | (sf.WorkerId << workerIdShift) | sf.Sequence
	id, err := strconv.ParseInt(fmt.Sprintf("%d", &num), 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
