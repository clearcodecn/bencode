# bencode
bencode implement


### install 

```
go get -u github.com/clearcodecn/bencode
```


### import 

```
import "github.com/clearcodecn/bencode"

```

### Usage 

```go 
	var listMap = []interface{}{
		12563214,
		"absadasd",
		map[string]string{
			"hello": "111",
		},
	}
	data := bencode.Marshal(listMap)
	fmt.Println(data)

	v, err := bencode.UnMarshal(data)
	require.Nil(t, err)

```

### benchMarks

```
goos: darwin
goarch: amd64
pkg: github.com/clearcodecn/bencode
BenchmarkMarshal
BenchmarkMarshal-8   	1000000000	         0.000045 ns/op
PASS


goos: darwin
goarch: amd64
pkg: github.com/clearcodecn/bencode
BenchmarkUnMarshal
BenchmarkUnMarshal-8   	1000000000	         0.000016 ns/op
PASS

```

### protocol

see [BEP-3](http://www.bittorrent.org/beps/bep_0003.html)


