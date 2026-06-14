package core

import (
	"fmt"
	"log"
)

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty data")
	}
	value, _, err := DecodeOne(data)
	if err != nil {
		fmt.Errorf("Error while parsing the data")
	}

	return value, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("empty data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInteger(data)
	case '$':
		return readBulkString(data)
	case '*':
		val, len, _ := readArray(data)
		fmt.Println("READ ARRAY RETUEN D VALUE ", val, len)
		return readArray(data)
	default:
		return nil, 0, fmt.Errorf("invalid data")
	}

}

func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for ; data[pos] != '\r'; pos++ {
	}
	return string(data[1:pos]), pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readInteger(data []byte) (int64, int, error) {
	pos := 1
	var value int64 = 0
	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}
	return value, pos + 2, nil
}

func readBulkString(data []byte) (string, int, error) {
	pos := 1
	length, delta := readLength(data[pos:])
	fullPos := pos + delta
	return string(data[fullPos : fullPos+length]), fullPos + length + 2, nil
}

func readLength(data []byte) (int, int) {
	pos, length := 0, 0
	for pos = range data {
		if !(data[pos] >= '0' && data[pos] <= '9') {
			fmt.Println("length of readlenght inside the func", length)
			return length, pos + 2
		}
		length = length*10 + int(data[pos]-'0')
	}

	return 0, 0
}

func readArray(data []byte) (interface{}, int, error) {
	pos := 1
	length, delta := readLength(data[pos:])
	fmt.Println("read LENGHT VALUE", length, delta)
	pos += delta
	value := make([]interface{}, length)
	for i := range length {
		log.Println("position ", pos)
		val, len, err := DecodeOne(data[pos:])
		fmt.Println("Decode one values", val, len)
		if err != nil {
			fmt.Errorf("Unable to read the data")
		}
		value[i] = val
		pos += len
	}
	return value, pos, nil
}

func DecodeArrayString(data []byte) ([]string, error) {
	value, err := Decode(data)
	fmt.Println("avalue ", value)
	if err != nil {
		return nil, err
	}
	ts := value.([]interface{})
	fmt.Println("ts ", ts)
	result := make([]string, len(ts))
	for i, v := range ts {
		result[i] = v.(string)
	}
	return result, nil
}

func Encode(value interface{}, isSimple bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", v))
		} else {
			return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
		}

	}
	return []byte{}
}
