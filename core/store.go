package core

import (
	"time"
)

var store map[string]*RedisObj

type RedisObj struct {
	Value     interface{}
	ExpiredAt int64
}

func Init() {
	store = make(map[string]*RedisObj)
}

func NewObj(value interface{}, durationMs int) *RedisObj {
	var expiresAt int64 = -1
	if durationMs > 0 {
		expiresAt = time.Now().UnixMilli() + int64(durationMs)
	}
	return &RedisObj{
		Value:     value,
		ExpiredAt: int64(expiresAt),
	}
}

func Put(key string, value string, expiresAt int64) {
	store[key] = &RedisObj{Value: value, ExpiredAt: expiresAt}
}

func Get(key string) *RedisObj {
	return store[key]
}

func Delete(key string) {
	delete(store, key)
}
