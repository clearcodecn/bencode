package bencode

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncode(t *testing.T) {
	m := new(marshal)

	s := m.marshalString("spam")
	require.Equal(t, "4:spam", s)

	s = m.marshalNumber(3)
	require.Equal(t, "i3e", s)

	s = m.marshalNumber(-3)
	require.Equal(t, "i-3e", s)

	s = m.marshalNumber(0)
	require.Equal(t, "i0e", s)

	s = m.marshalList([]interface{}{
		"spam", "eggs",
	})
	require.Equal(t, "l4:spam4:eggse", s)

	s = m.marshalMap(map[string]interface{}{
		"cow":  "moo",
		"spam": "eggs",
	})
	require.Equal(t, "d3:cow3:moo4:spam4:eggse", s)
}

func TestBuf(t *testing.T) {
	var x = "aaaaae"
	buf := bufio.NewReader(bytes.NewBufferString(x))
	// include e
	data, _ := buf.ReadBytes('e')
	fmt.Println(string(data))
}

func TestDecode(t *testing.T) {
	var listMap = []interface{}{
		12563214,
		"absadasd",
		map[string]string{
			"hello": "111",
		},
	}

	data := Marshal(listMap)
	fmt.Println(data)

	v, err := UnMarshal(data)
	require.Nil(t, err)

	data2 := Marshal(v)
	fmt.Println(data2)
}

func BenchmarkMarshal(t *testing.B) {
	var listMap = []interface{}{
		12563214,
		"absadasd",
		map[string]string{
			"hello": "111",
		},
		map[string]string{
			"hello": "111",
		},
		map[string]string{
			"hello": "111",
		},
		map[string]string{
			"hello": "111",
		},
		[]interface{}{
			map[string]string{
				"hello": "111",
			},
			map[string]string{
				"hello": "111",
			},
			map[string]string{
				"hello": "111",
			},
			map[string]string{
				"hello": "111",
			},
			12563214,
			"absadasd",
			12563214,
			"absadasd",
			12563214,
			"absadasd",
		},
	}
	data := Marshal(listMap)
	fmt.Println(data)
}

func BenchmarkUnMarshal(t *testing.B) {
	var listMap = []interface{}{
		12563214,
		"absadasd",
		map[string]string{
			"hello": "111",
		},
		map[string]string{
			"hello": "111",
		},
		map[string]string{
			"hello": "111",
		},
		map[string]string{
			"hello": "111",
		},
		[]interface{}{
			map[string]string{
				"hello": "111",
			},
			map[string]string{
				"hello": "111",
			},
			map[string]string{
				"hello": "111",
			},
			map[string]string{
				"hello": "111",
			},
			12563214,
			"absadasd",
			12563214,
			"absadasd",
			12563214,
			"absadasd",
		},
	}
	data := Marshal(listMap)
	fmt.Println(data)

	t.ResetTimer()
	_, err := UnMarshal(data)
	require.Nil(t, err)
}
