package core

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

func evalPing(args []string, c io.ReadWriter) error {
	var b []byte
	if len(args) >= 2 {
		return errors.New("(error) ERR wrong number of arguments for 'ping' command")
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	_, err := c.Write(b)
	return err
}

func evalGet(args []string, c io.ReadWriter) error {
	if len(args) < 1 {
		return errors.New("(err) Syntax error")
	}
	key := args[0]
	obj := Get(key)
	fmt.Println("objOBJ", obj)
	if obj == nil {
		c.Write([]byte(":-2\r\n"))
		return nil
	}
	strValue, ok := obj.Value.(string)
	if !ok {
		return fmt.Errorf("value is not a string")
	}
	if obj.ExpiredAt == -1 {
		c.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(strValue), obj.Value)))
	}
	fmt.Println("whats the time now", time.Now().UnixMilli())
	fmt.Println("whats the diff of the time now - expired at", obj.ExpiredAt-time.Now().UnixMilli())
	if expiry := time.Now().UnixMilli(); obj.ExpiredAt-expiry < 0 {
		c.Write([]byte(":-2\r\n"))
		return nil
	}
	str := fmt.Sprintf("$%d\r\n%s\r\n", len(obj.Value.(string)), obj.Value)
	fmt.Println("str", str)
	c.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(obj.Value.(string)), obj.Value)))
	return nil
}

func evalSet(args []string, c io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("(err) Wrong number of inputs for 'set' command")
	}
	key, value := args[0], args[1]
	var expiryTime int64 = -1

	if len(args) > 2 {
		for i := 2; i < len(args); i++ {
			val := args[i]
			switch val {
			case "EX", "ex":
				i++
				if i == len(args) {
					return errors.New("(err) ERR syntax error")
				}
				exprSec, err := strconv.ParseInt(args[3], 10, 64)
				if err != nil {
					return errors.New("(error) ERR value is not an integer or out of range")
				}
				expiryTime = time.Now().UnixMilli() + exprSec*1000
			default:
				return errors.New("err unsupported type")
			}
		}

	}
	Put(key, value, int64(expiryTime))
	c.Write([]byte("+OK\r\n"))
	return nil
}

func evalTTL(args []string, c io.ReadWriter) error {

	if len(args) > 1 {
		return errors.New("(error) ERR wrong number of arguments for 'ttl' command")
	}
	key := args[0]
	if Get(key) == nil {
		c.Write([]byte(":-2\r\n"))
		return nil
	}
	obj := Get(key)
	if obj.ExpiredAt-time.Now().UnixMilli() < 0 {
		c.Write([]byte(":-2\r\n"))
		return nil
	}
	c.Write([]byte(fmt.Sprintf(":%d\r\n", int(obj.ExpiredAt-time.Now().UnixMilli())/1000)))
	return nil

}

func evalExpire(args []string, c io.ReadWriter) error {
	if len(args) > 2 || len(args) < 2 {
		return errors.New("(error) ERR wrong number of arguments for 'expire' command")
	}
	key := args[0]
	expireSec := args[1]
	if Get(key) == nil {
		c.Write([]byte(":0\r\n"))
		return nil
	}
	obj := Get(key)
	expireSecInt, err := strconv.ParseInt(expireSec, 10, 64)
	if err != nil {
		return errors.New("(error) ERR value is not an integer or out of range")
	}
	obj.ExpiredAt = time.Now().UnixMilli() + expireSecInt*1000
	c.Write([]byte(":1\r\n"))

	return nil
}

func evalDEL(args []string, c io.ReadWriter) error {
	var count int64 = 0
	for _, key := range args {
		if Get(key) == nil {
			continue
		}
		Delete(key)
		count++
	}
	c.Write([]byte(fmt.Sprintf(":%d\r\n", count)))
	return nil
}

func EvalAndRespond(cmd *RedisCmd, c io.ReadWriter) error {
	switch cmd.Cmd {
	case "PING":
		return evalPing(cmd.Args, c)
	case "GET":
		return evalGet(cmd.Args, c)
	case "SET":
		return evalSet(cmd.Args, c)
	case "TTL":
		return evalTTL(cmd.Args, c)
	case "EXPIRE":
		return evalExpire(cmd.Args, c)
	case "DEL":
		return evalDEL(cmd.Args, c)
	default:
		return evalPing(cmd.Args, c)
	}
}
