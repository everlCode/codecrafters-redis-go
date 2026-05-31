package database

import "slices"

type Zset struct {
	Keys map[string]*Zvalue
	Set []Zvalue
}

type Zvalue struct {
	Id float64
	Value string
}

func NewZset() *Zset {
	return &Zset{
		Keys: make(map[string]*Zvalue),
	}
}

func (z *Zset) Add(id float64, value string) int {
	var v *Zvalue
	var isNewValue int = 0
	v, ok := z.Keys[value]
	if !ok {
		v = &Zvalue{Id: id, Value: value}
		isNewValue = 1
	}
	
	z.Keys[value] = v
	z.Set = append(z.Set, *v)
	slices.SortFunc(z.Set, func(a, b Zvalue) int {
		if a.Id < b.Id {
			return -1
		} else if a.Id == b.Id {
			return 0
		} else {
			return 1
		}
	})

	return isNewValue
}