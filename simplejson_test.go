package simplejson

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"io/ioutil"
	"log"
	"strconv"
	"testing"
)

func TestSimplejson(t *testing.T) {
	var ok bool
	var err error

	log.SetOutput(ioutil.Discard)

	js, err := NewJson([]byte(`{ 
		"test": { 
			"string_array": ["asdf", "ghjk", "zxcv"],
			"int_array": [1, 4, 7],
			"array": [1, "2", 3],
			"arraywithsubs": [{"subkeyone": 1},
			{"subkeytwo": 2, "subkeythree": 3}],
			"int": 10,
			"float": 5.150,
			"bignum": 9223372036854775807,
			"string": "simplejson",
			"bool": true 
		}
	}`))

	assert.NotEqual(t, nil, js)
	assert.Equal(t, nil, err)

	_, ok = js.CheckGet("test")
	assert.Equal(t, true, ok)

	_, ok = js.CheckGet("missing_key")
	assert.Equal(t, false, ok)

	arr, _ := js.Get("test").Get("array").Array()
	assert.NotEqual(t, nil, arr)
	for i, v := range arr {
		var iv int
		switch v.(type) {
		case float64:
			iv = int(v.(float64))
		case string:
			iv, _ = strconv.Atoi(v.(string))
		}
		assert.Equal(t, i+1, iv)
	}

	aws := js.Get("test").Get("arraywithsubs")
	assert.NotEqual(t, nil, aws)
	var awsval int
	awsval, _ = aws.GetIndex(0).Get("subkeyone").Int()
	assert.Equal(t, 1, awsval)
	awsval, _ = aws.GetIndex(1).Get("subkeytwo").Int()
	assert.Equal(t, 2, awsval)
	awsval, _ = aws.GetIndex(1).Get("subkeythree").Int()
	assert.Equal(t, 3, awsval)

	i, _ := js.Get("test").Get("int").Int()
	assert.Equal(t, 10, i)

	f, _ := js.Get("test").Get("float").Float64()
	assert.Equal(t, 5.150, f)

	s, _ := js.Get("test").Get("string").String()
	assert.Equal(t, "simplejson", s)

	b, _ := js.Get("test").Get("bool").Bool()
	assert.Equal(t, true, b)

	mi := js.Get("test").Get("int").MustInt()
	assert.Equal(t, 10, mi)

	mi2 := js.Get("test").Get("missing_int").MustInt(5150)
	assert.Equal(t, 5150, mi2)

	ms := js.Get("test").Get("string").MustString()
	assert.Equal(t, "simplejson", ms)

	ms2 := js.Get("test").Get("missing_string").MustString("fyea")
	assert.Equal(t, "fyea", ms2)

	strs, err := js.Get("test").Get("string_array").StringArray()
	assert.Equal(t, err, nil)
	assert.Equal(t, strs[0], "asdf")
	assert.Equal(t, strs[1], "ghjk")
	assert.Equal(t, strs[2], "zxcv")

	ints, erri := js.Get("test").Get("int_array").StringArray()
	assert.Equal(t, erri, nil)
	assert.Equal(t, ints[0], 1)
	assert.Equal(t, ints[1], 4)
	assert.Equal(t, ints[2], 7)
}

func TestStdlibInterfaces(t *testing.T) {
	val := new(struct {
		Name   string `json:"name"`
		Params *Json  `json:"params"`
	})
	val2 := new(struct {
		Name   string `json:"name"`
		Params *Json  `json:"params"`
	})

	raw := `{"name":"myobject","params":{"string":"simplejson"}}`

	assert.Equal(t, nil, json.Unmarshal([]byte(raw), val))

	assert.Equal(t, "myobject", val.Name)
	assert.NotEqual(t, nil, val.Params.data)
	s, _ := val.Params.Get("string").String()
	assert.Equal(t, "simplejson", s)

	p, err := json.Marshal(val)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, json.Unmarshal(p, val2))
	assert.Equal(t, val, val2) // stable
}
