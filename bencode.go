package bencode

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"sort"
	"strconv"
)

const (
	colon        = ":"
	numberLeader = "i"
	listLeader   = "l"
	dictLeader   = "d"
	commonEnd    = "e"

	colonByte        = ':'
	numberLeaderByte = 'i'
	listLeaderByte   = 'l'
	dictLeaderByte   = 'd'
	commonEndByte    = 'e'
)

// 1.  字符串表示方法   4:ssss      长度:string
// 2.  整形数据表示方法 i3e = 3 , i-3e  =  -3   , i0e  代表0   但是  i03e 是错误的.  i0xe 是错误的表示方法.
// 3.  list 列表：  l4:aaaa5:xxxxxe  ->  ["aaaa","xxxxx"]   以e结尾
// 4.  字典:  d3:cow3:moo4:spam4:eggse  以e结尾  ->  {"cow":"moo","spam":"eggs"}

type marshal struct{}

func (marshal) marshalString(s string) string {
	i := len(s)
	if i == 0 {
		return ""
	}
	leader := strconv.Itoa(i)
	return leader + colon + s
}

func (marshal) marshalNumber(i int64) string {
	return numberLeader + strconv.FormatInt(i, 10) + commonEnd
}

// list have follow types:
// [string, number,list,dict]
func (m marshal) marshalList(s []interface{}) string {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString(listLeader)
	for _, i := range s {
		switch v := i.(type) {
		case string:
			buf.WriteString(m.marshalString(v))
		case int:
			buf.WriteString(m.marshalNumber(int64(v)))
		case int8:
			buf.WriteString(m.marshalNumber(int64(v)))
		case int16:
			buf.WriteString(m.marshalNumber(int64(v)))
		case int32:
			buf.WriteString(m.marshalNumber(int64(v)))
		case int64:
			buf.WriteString(m.marshalNumber(v))
		case uint:
			buf.WriteString(m.marshalNumber(int64(v)))
		case uint8:
			buf.WriteString(m.marshalNumber(int64(v)))
		case uint16:
			buf.WriteString(m.marshalNumber(int64(v)))
		case uint32:
			buf.WriteString(m.marshalNumber(int64(v)))
		case uint64:
			buf.WriteString(m.marshalNumber(int64(v)))
		case map[string]interface{}:
			buf.WriteString(m.marshalMap(v))
		case []interface{}:
			buf.WriteString(m.marshalList(v))
		}
	}
	buf.WriteString(commonEnd)
	return buf.String()
}

func (m marshal) marshal(val interface{}) string {
	switch v := val.(type) {
	case string:
		return m.marshalString(v)
	case int:
		return m.marshalNumber(int64(v))
	case int8:
		return m.marshalNumber(int64(v))
	case int16:
		return m.marshalNumber(int64(v))
	case int32:
		return m.marshalNumber(int64(v))
	case int64:
		return m.marshalNumber(int64(v))
	case uint:
		return m.marshalNumber(int64(v))
	case uint8:
		return m.marshalNumber(int64(v))
	case uint16:
		return m.marshalNumber(int64(v))
	case uint32:
		return m.marshalNumber(int64(v))
	case uint64:
		return m.marshalNumber(int64(v))
	case map[string]interface{}:
		return m.marshalMap(v)
	case []interface{}:
		return m.marshalList(v)
	default:
		panic("un support type")
	}
}

func (m marshal) marshalMap(data map[string]interface{}) string {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString(dictLeader)
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	fn := func(val interface{}) {
		switch v := val.(type) {
		case string:
			buf.WriteString(m.marshalString(v))
		case int:
			buf.WriteString(m.marshalNumber(int64(v)))
		case int8:
			buf.WriteString(m.marshalNumber(int64(v)))
		case int16:
			buf.WriteString(m.marshalNumber(int64(v)))
		case int32:
			buf.WriteString(m.marshalNumber(int64(v)))
		case int64:
			buf.WriteString(m.marshalNumber(v))
		case uint:
			buf.WriteString(m.marshalNumber(int64(v)))
		case uint8:
			buf.WriteString(m.marshalNumber(int64(v)))
		case uint16:
			buf.WriteString(m.marshalNumber(int64(v)))
		case uint32:
			buf.WriteString(m.marshalNumber(int64(v)))
		case uint64:
			buf.WriteString(m.marshalNumber(int64(v)))
		case map[string]interface{}:
			buf.WriteString(m.marshalMap(v))
		case []interface{}:
			buf.WriteString(m.marshalList(v))
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		fn(key)
		fn(data[key])
	}
	buf.WriteString(commonEnd)
	return buf.String()
}

type unmarshal struct {
	buf *bufio.Reader
	r   *bytes.Buffer
}

func (m *unmarshal) unMarshal() (interface{}, error) {
	for {
		b, err := m.next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch b {
		case numberLeaderByte:
			return m.unmarshalNumber()
		case listLeaderByte:
			return m.unmarshalList()
		case dictLeaderByte:
			return m.unmarshalMap()
		default:
			if b <= '9' && b >= '0' {
				return m.unmarshalString()
			}
			return nil, errors.New("invalid structure: " + string(b))
		}
	}
	return nil, io.EOF
}

func (m *unmarshal) next() (byte, error) {
	b, err := m.buf.Peek(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (m *unmarshal) unmarshalNumber() (int64, error) {
	_, err := m.buf.Discard(1)
	if err != nil {
		return 0, err
	}
	data, err := m.buf.ReadBytes(commonEndByte)
	if err != nil {
		return 0, err
	}
	if len(data) <= 1 {
		return 0, errors.New("invalid number format")
	}
	i, err := strconv.ParseInt(string(data[:len(data)-1]), 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (m *unmarshal) unmarshalString() (string, error) {
	data, err := m.buf.ReadBytes(colonByte)
	if err != nil {
		return "", err
	}
	if len(data) <= 1 {
		return "", errors.New("invalid number format")
	}
	number, err := strconv.ParseInt(string(data[:len(data)-1]), 10, 64)
	if err != nil {
		return "", err
	}
	data = make([]byte, number)
	n, err := io.ReadFull(m.buf, data)
	if err != nil {
		return "", err
	}
	if int64(n) != number {
		return "", errors.New("invalid string length format")
	}
	return string(data), nil
}

func (m *unmarshal) unmarshalList() ([]interface{}, error) {
	var data []interface{}
	_, err := m.buf.Discard(1)
	if err != nil {
		return nil, err
	}
	for {
		value, err := m.unMarshal()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		data = append(data, value)
		b, err := m.next()
		if err != nil {
			return nil, err
		}
		if b == commonEndByte {
			m.buf.Discard(1)
			break
		}
	}

	return data, nil
}

func (m *unmarshal) unmarshalMap() (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	_, err := m.buf.Discard(1)
	if err != nil {
		return nil, err
	}
	for {
		key, err := m.unmarshalString()
		if err != nil {
			return nil, err
		}
		val, err := m.unMarshal()
		if err != nil {
			return nil, err
		}
		data[key] = val
		b, err := m.next()
		if err != nil {
			return nil, err
		}
		if b == commonEndByte {
			m.buf.Discard(1)
			break
		}
	}

	return data, nil
}

func UnMarshal(s string) (interface{}, error) {
	m := new(unmarshal)
	m.r = bytes.NewBufferString(s)
	m.buf = bufio.NewReader(m.r)
	return m.unMarshal()
}

func Marshal(v interface{}) string {
	m := new(marshal)
	return m.marshal(v)
}
