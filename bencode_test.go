package dht

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
	m := new(marshal)

	var listMap = []interface{}{
		[]int{1, 2, 3, 45, 65, 67},
		[]int{1, 2, 3, 45, 65, 67},
		[]int{1, 2, 3, 45, 65, 67},
		[]int{1, 2, 3, 45, 65, 67},
		12563214,
		"absadasd",
		map[string]string{
			"hello": "111",
		},
		map[string]interface{}{
			"hello": "111",
		},
		"asdasdasd",
		"hahaha",
		[]int{1, 2, 3, 45, 65, 67},
		[]int{1, 2, 3, 45, 65, 67},
		[]int{1, 2, 3, 45, 65, 67},
		[]int{1, 2, 3, 45, 65, 67},
		[]int{1, 2, 3, 45, 65, 67},
		[]int{1, 2, 3, 45, 65, 67},
		[]interface{}{
			12563214,
			"absadasd",
			map[string]string{
				"hello": "111",
			},
			map[string]interface{}{
				"hello": "111",
			},
			"asdasdasd",
			"hahaha",
			[]interface{}{
				12563214,
				"absadasd",
				map[string]string{
					"hello": "111",
				},
				map[string]interface{}{
					"hello": "111",
				},
				"asdasdasd",
				"hahaha",
				[]interface{}{
					12563214,
					"absadasd",
					map[string]string{
						"hello": "111",
					},
					map[string]interface{}{
						"hello": "111",
					},
					"asdasdasd",
					"hahaha",
					[]interface{}{
						12563214,
						"absadasd",
						map[string]string{
							"hello": "111",
						},
						map[string]interface{}{
							"hello": "111",
						},
						"asdasdasd",
						"hahaha",
					},
				},
			},
		},
	}

	data := m.marshalList(listMap)
	fmt.Println(data)

	v, err := UnMarshal(data)
	require.Nil(t, err)

	data2 := m.marshal(v)
	fmt.Println(data2)
}
