package structhash

import (
	"fmt"
	"testing"
)

type First struct {
	Bool   bool    `version:"1"`
	String string  `version:"2"`
	Int    int     `version:"1" lastversion:"1"`
	Float  float64 `version:"1"`
	Struct *Second `version:"1"`
	Uint   uint    `version:"1"`
}

type Second struct {
	Map   map[string]string `version:"1"`
	Slice []int             `version:"1"`
}

type Tags1 struct {
	Int int    `hash:"-"`
	Str string `hash:"name:Foo version:1 lastversion:2"`
	Bar string `hash:"version:1"`
}

type Tags2 struct {
	Foo string
	Bar string
}

type Tags3 struct {
	Bar string
}

type Nils struct {
	Str   *string
	Int   *int
	Bool  *bool
	Map   map[string]string
	Slice []string
}

type unexportedTags struct {
	foo  string
	bar  string
	aMap map[string]string
}

func dataSetup() *First {
	tmpmap := make(map[string]string)
	tmpmap["foo"] = "bar"
	tmpmap["baz"] = "go"
	tmpslice := make([]int, 3)
	tmpslice[0] = 0
	tmpslice[1] = 1
	tmpslice[2] = 2
	return &First{
		Bool:   true,
		String: "test",
		Int:    123456789,
		Float:  65.3458,
		Struct: &Second{
			Map:   tmpmap,
			Slice: tmpslice,
		},
		Uint: 1,
	}
}

func TestHash(t *testing.T) {
	v1Hash := "v1_e8e67581aee36d7237603381a9cbd9fc"
	v2Hash := "v2_5e51490d7c24c4b7a9e63c04f55734eb"

	data := dataSetup()
	v1, err := Hash(data, 1)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(string(Dump(data, 1)))
	if v1 != v1Hash {
		t.Errorf("%s is not %s", v1, v1Hash)
	}
	v2, err := Hash(data, 2)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(string(Dump(data, 2)))
	if v2 != v2Hash {
		t.Errorf("%s is not %s", v2, v2Hash)
	}

	v1md5 := fmt.Sprintf("v1_%x", Md5(data, 1))
	if v1md5 != v1Hash {
		t.Errorf("%s is not %s", v1md5, v1Hash[3:])
	}
	v2md5 := fmt.Sprintf("v2_%x", Md5(data, 2))
	if v2md5 != v2Hash {
		t.Errorf("%s is not %s", v2md5, v2Hash[3:])
	}
}

func TestTags(t *testing.T) {
	t1 := Tags1{11, "foo", "bar"}
	t1x := Tags1{22, "foo", "bar"}
	t2 := Tags2{"foo", "bar"}
	t3 := Tags3{"bar"}

	t1_dump := string(Dump(t1, 1))
	t1x_dump := string(Dump(t1x, 1))
	if t1_dump != t1x_dump {
		t.Errorf("%s is not %s", t1_dump, t1x_dump)
	}

	t2_dump := string(Dump(t2, 1))
	if t1_dump != t2_dump {
		t.Errorf("%s is not %s", t1_dump, t2_dump)
	}

	t1v3_dump := string(Dump(t1, 3))
	t3v3_dump := string(Dump(t3, 3))
	if t1v3_dump != t3v3_dump {
		t.Errorf("%s is not %s", t1v3_dump, t3v3_dump)
	}
}

func TestNils(t *testing.T) {
	s1 := Nils{
		Str:   nil,
		Int:   nil,
		Bool:  nil,
		Map:   nil,
		Slice: nil,
	}

	s2 := Nils{
		Str:   new(string),
		Int:   new(int),
		Bool:  new(bool),
		Map:   make(map[string]string),
		Slice: make([]string, 0),
	}

	s1_dump := string(Dump(s1, 1))
	s2_dump := string(Dump(s2, 1))
	if s1_dump != s2_dump {
		t.Errorf("%s is not %s", s1_dump, s2_dump)
	}
}

func TestUnexportedFields(t *testing.T) {
	v1Hash := "v1_750efb7c919caf87f2ab0d119650c87d"
	data := unexportedTags{
		foo: "foo",
		bar: "bar",
		aMap: map[string]string{
			"key1": "val",
		},
	}
	v1, err := Hash(data, 1)
	if err != nil {
		t.Error(err)
	}
	if v1 != v1Hash {
		t.Errorf("%s is not %s", v1, v1Hash)
	}
	v1md5 := fmt.Sprintf("v1_%x", Md5(data, 1))
	if v1md5 != v1Hash {
		t.Errorf("%s is not %s", v1md5, v1Hash[3:])
	}
}
