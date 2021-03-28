package cache

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/silenceper/log"
	"ginblog/setup"
	"github.com/kms9/publicyc/pkg/store/oredis"
	"github.com/kms9/publicyc/pkg/util/ostring"
)

var redisInstance *htRedis

type htRedis struct {
	*oredis.Redis
	LifeSpan int64
}

func Redis() *htRedis {
	return &htRedis{setup.Cache, 864e2}
}

func (htr *htRedis) key(k string) string {
	return fmt.Sprintf("%s:custom:%s", htr.GetConfig().Name, k)
}

func (htr *htRedis) Do(command string, args ...interface{}) (interface{}, error) {
	return htr.Redis.Do(command, args...)
}

func (htr *htRedis) TryLockWithTimeout(key, value string, timeout int64) (ok bool, result interface{}, err error) {
	_, err = redis.String(htr.Redis.Do("SET", key, value, "EX", timeout, "NX"))
	if err == redis.ErrNil {
		str, err := redis.String(htr.Redis.Do("GET", key))
		if err != nil {
			return false, nil, err
		}
		return false, str, nil
	}
	if err != nil {
		return false, nil, err
	}
	return true, nil, nil
}

func (htr *htRedis) DelLock(key, value string) (err error) {
	_, err = htr.Redis.Do("EVAL", `if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end`, 1, key, value)
	return err
}

func (htr *htRedis) LRUGet(k string) string {
	count, err := redis.Int64(htr.Redis.Do("ZCARD", k))
	if err != nil || count == 0 {
		return ""
	}
	index := rand.Int63n(count)
	result, err := redis.Strings(htr.Redis.Do("ZRANGE", k, index, index))
	if err != nil {
		log.Warnf("LRUGet Redis ZRANGE k:%s err: %s", k, err.Error())
	}
	if len(result) == 0 {
		return ""
	}
	return result[0]
}

func (htr *htRedis) LRUSet(k string, v interface{}, limitAmount int64, expire int64) error {
	_, err := htr.Redis.Do("ZADD", k, time.Now().Unix(), v)
	if err != nil {
		log.Warnf("LRUSet Redis ZADD k:%s err: %s", k, err.Error())
		return err
	}
	_, err = htr.Redis.Do("ZREMRANGEBYRANK", k, 0, -1*(limitAmount+1))
	if err != nil {
		log.Warnf("LRUSet Set Redis ZREMRANGEBYRANK k:%s err: %s", k, err.Error())
	}
	if expire != 0 {
		_, err = htr.Redis.Do("EXPIRE", k, expire) // 延长expire存储
		if err != nil {
			log.Warnf("LRUSet Set Redis EXPIRE k:%s err: %s", k, err.Error())
		}
	}
	return nil
}

// ZrangeByScore 有序集合按照分数排序
func (htr *htRedis) ZrangeByScore(k string, n, m, point interface{}) error {
	values, err := redis.Values(htr.Redis.Do("ZRANGEBYSCORE", k, n, m, "WITHSCORES"))
	if err != nil {
		return err
	}
	return redis.ScanSlice(values, point)
}

// Zrange 向有序集合获取第n位到第m位排序
//     n：其实位置；m：截止位置，-1，则为全部；point：指针，
//     type rangeInfo struct {
//     	   UserID string
//     	   Score  int
//     }
//     point: &[]*rangeInfo{}
func (htr *htRedis) Zrange(k string, n, m int64, point interface{}) error {
	values, err := redis.Values(htr.Redis.Do("ZRANGE", k, n, m, "WITHSCORES"))
	if err != nil {
		return err
	}
	return redis.ScanSlice(values, point)
}

// PipelineParams ..
type PipelineParams struct {
	Command string
	Args    []interface{}
}

// Pipeline 批量调用
func (htr *htRedis) Pipeline(params []*PipelineParams) (interface{}, error) {
	conn, closeFunc := htr.Redis.GetConnect()
	defer func() {
		_ = closeFunc()
	}()

	_ = conn.Send("MULTI")
	for _, param := range params {
		_ = conn.Send(param.Command, param.Args...)
	}
	return conn.Do("EXEC")
}

// LoadDetail 加载映射详情
func (t *htRedis) LoadDetail(k string, v interface{}) interface{} {
	var err error
	var value []byte
	if value, err = redis.Bytes(t.Do("GET", t.key(k))); err != nil {
		return nil
	}

	if err := json.Unmarshal(value, &v); err != nil {
		return nil
	}
	return v
}

// Put 重新设置缓存
func (t *htRedis) Put(k string, o interface{}, liftSpan int64) error {
	value := string(ostring.MustJson(o, false))
	if ostring.InterfaceIsNil(o) {
		return fmt.Errorf(" value is nil")
	}
	if liftSpan == 0 {
		liftSpan = t.LifeSpan
	}
	if _, err := t.Do("SET", t.key(k), value, "EX", liftSpan); err != nil {
		return err
	}
	return nil
}

// Del 移除缓存
func (t *htRedis) Del(k string) error {
	_, err := t.Do("DEL", t.key(k))
	return err
}

// Exists 是否存在
func (t *htRedis) Exists(k string) bool {
	exists, err := t.Do("EXISTS", t.key(k))
	if err != nil {
		return false
	} else {
		return exists.(int64) == 1
	}
}
