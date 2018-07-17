package JSON

import (
	"runtime"
	"strconv"
)

// JObject estrutura para json files
type JObject struct {
	data  []byte
	pos   int
	size  int
	base  int
	block int
}

// JArray estrutura de array com json
type JArray struct {
	data      []byte
	base      int
	pos       int
	size      int
	block     int
	array     int
	lastIndex int
	current   JObject
}

var (
	blocks int
	arrays int
)

func skip(data []byte, pos int) int {
	for i := pos; i < len(data); i++ {
		if data[i] != ' ' && data[i] != '\t' && data[i] != '\r' && data[i] != '\n' {
			return i
		}
	}
	return -1
}

func (j *JObject) expect(b byte) bool {
	j.pos = skip(j.data, j.pos)
	// if pos == idx {
	// pos++
	// }
	if j.data[j.pos] == b {
		j.pos++
		return true
	}
	return false
}

func index(s []byte, b string) int {
	return indexFrom(s, b, 0)
}

func indexFrom(s []byte, b string, p int) int {
	if string(s) == b {
		return 0
	}
	l := len(s)
	lb := len(b)
	if p >= l || l < lb {
		return -1
	}
	pos := 0
	block := 0
	array := 0
	for i := p; i < l; i++ {
		if s[i] == '{' {
			block++
			continue
		} else if s[i] == '}' {
			block--
		} else if s[i] == '[' {
			array++
			continue
		} else if s[i] == ']' {
			array--
		}
		if block > 0 || array > 0 {
			continue
		}
		if array < 0 || block < 0 {
			return -1
		}
		if s[i] == b[pos] {
			pos++
			if pos == lb {
				return i + 1
			}
		} else {
			pos = 0
		}
	}
	return -1
}

func (j *JObject) find(s string) bool {
	idx := indexFrom(j.data, "\""+s+"\"", j.base)
	if idx == -1 {
		return false
	}
	j.pos = idx
	return true
}

func isNum(b byte) bool {
	if (b >= '0') && (b <= '9') {
		return true
	}
	return false
}

func Load(data []byte) JObject {
	var j JObject
	j.data = data
	j.base = skip(j.data, 0)
	if j.data[j.pos] != '{' {
		println("Erro no json")
	}
	blocks = 1
	j.base++
	j.pos = j.base
	j.size = len(data)
	return j
}

func (j *JObject) Free() {
	runtime.GC()
}

func (j *JObject) Position() int {
	return j.pos
}

func (j *JObject) GetInt(s string) int {
	if j.find(s) == false {
		return 0
	}
	if j.expect(':') == false {
		return 0
	}
	idx := skip(j.data, j.pos)
	if ok := isNum(j.data[idx]); ok {
		num := make([]byte, 50)
		num[0] = j.data[idx]
		k := 1
		for i := idx + 1; i < j.size; i++ {
			if ok := isNum(j.data[i]); !ok {
				n, _ := strconv.Atoi(string(num[:k]))
				return n
			}
			num[k] = j.data[i]
			k++
		}
	}
	return 0
}

func (j *JObject) GetFloat(s string) float64 {
	if j.find(s) == false {
		return 0
	}
	if j.expect(':') == false {
		return 0
	}
	idx := skip(j.data, j.pos)
	if ok := isNum(j.data[idx]); ok {
		num := make([]byte, 100)
		num[0] = j.data[idx]
		k := 1
		comma := false
		for i := idx + 1; i < j.size; i++ {
			if j.data[i] == '.' {
				if comma == false {
					comma = true
					num[k] = '.'
					k++
					continue
				} else {
					return 0
				}
			}
			if ok := isNum(j.data[i]); !ok {
				n, _ := strconv.ParseFloat(string(num[:k]), 64)
				return n
			}
			num[k] = j.data[i]
			k++
		}
	}
	return 0
}

// func (j *JObject) GetString(s string) (string, bool) {
// 	if j.find(s) == false {
// 		return "", false
// 	}
// 	if j.expect(':') == false {
// 		return "", false
// 	}
// 	if j.expect('"') {
// 		idx := j.pos
// 		ret := indexFrom(j.data, "\"", j.pos)
// 		if ret == -1 {
// 			return "", false
// 		}
// 		j.pos = ret
// 		return string(j.data[idx : ret-1]), true
// 	}
// 	return "", false
// }

func (j *JObject) GetString(s string) string {
	if j.find(s) == false {
		return ""
	}
	if j.expect(':') == false {
		return ""
	}
	if j.expect('"') {
		idx := j.pos
		ret := indexFrom(j.data, "\"", j.pos)
		if ret == -1 {
			return ""
		}
		j.pos = ret
		return string(j.data[idx : ret-1])
	}
	return ""
}

func (j *JObject) GetObject(s string) (JObject, bool) {
	if j.find(s) == false {
		return JObject{}, false
	}
	if j.expect(':') == false {
		return JObject{}, false
	}
	if j.expect('{') {
		blocks++
		var o JObject
		o.block = blocks
		o.pos = j.pos
		o.base = j.pos
		o.data = j.data
		o.size = j.size
		return o, true
	}
	return JObject{}, false
}

func (j *JObject) GetArray(s string) (JArray, bool) {
	if j.find(s) == false {
		return JArray{}, false
	}
	if j.expect(':') == false {
		return JArray{}, false
	}
	if j.expect('[') {
		arrays++
		var o JArray
		o.block = blocks
		o.array = arrays
		o.pos = j.pos
		o.base = j.pos
		o.data = j.data
		o.size = j.size
		o.lastIndex = -1
		return o, true
	}
	return JArray{}, false
}

func (j *JArray) expect(b byte) bool {
	j.pos = skip(j.data, j.pos)
	// if pos == idx {
	// pos++
	// }
	if j.data[j.pos] == b {
		j.pos++
		return true
	}
	return false
}

func (j *JArray) Next() bool {
	if c, ok := j.Get(j.lastIndex + 1); ok {
		j.current = c
		return ok
	}
	return false
}

func (j *JArray) Current() JObject {
	return j.current
}

func (j *JArray) Get(idx int) (JObject, bool) {
	// if idx > a.lastIndex {
	// }
	if idx == 0 {
		j.pos = j.base
		if j.expect('{') == false {
			return JObject{}, false
		}
		j.lastIndex = 0
		blocks++
		var o JObject
		o.block = blocks
		o.pos = j.pos
		o.base = j.pos
		o.data = j.data
		o.size = j.size
		return o, true
	} else if idx > j.lastIndex {
		b := 1
		a := 1
		ind := j.lastIndex
		for i := j.pos; i < j.size; i++ {
			if j.data[i] == '{' {
				if b < 1 {
					if ind+1 < idx {
						ind++
						b++
						continue
					}
					j.pos = i + 1
					j.lastIndex = idx
					blocks++
					var o JObject
					o.block = blocks
					o.pos = j.pos
					o.base = j.pos
					o.data = j.data
					o.size = j.size
					return o, true
				}
				b++
			} else if j.data[i] == '}' {
				b--
			} else if j.data[i] == '[' {
				a++
			} else if j.data[i] == ']' {
				if a < 1 {
					return JObject{}, false
				}
				a--
			}
		}
	} else {

	}
	return JObject{}, false
}
