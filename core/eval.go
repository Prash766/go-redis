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
				exprMs, err := strconv.ParseInt(args[3], 10, 64)
				if err != nil {
					return errors.New("(error) ERR value is not an integer or out of range")
				}
				expiryTime = time.Now().UnixMilli() + exprMs
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
	default:
		return evalPing(cmd.Args, c)
	}
}
