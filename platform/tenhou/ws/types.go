package ws

import (
	"encoding/json"
	"strconv"
	"strings"
)

const sep = ","

// "13,4,65" <=> []int{13, 4, 65}
type Ints []int

func (q *Ints) UnmarshalJSON(raw []byte) error {
	*q = (*q)[:0]
	var ss string
	if err := json.Unmarshal(raw, &ss); err != nil {
		return err
	}
	for _, s := range strings.Split(ss, sep) {
		v, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*q = append(*q, v)
	}
	return nil
}

func (q Ints) MarshalJSON() ([]byte, error) {
	ss := []string{}
	for _, v := range q {
		ss = append(ss, strconv.Itoa(v))
	}
	return json.Marshal(strings.Join(ss, sep))
}

// "13.26,4.84,65" <=> []float64{13.26, 4.84, 65}
type Float64s []float64

func (q *Float64s) UnmarshalJSON(raw []byte) error {
	*q = (*q)[:0]
	var ss string
	if err := json.Unmarshal(raw, &ss); err != nil {
		return err
	}
	for _, s := range strings.Split(ss, sep) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*q = append(*q, v)
	}
	return nil
}

func (q Float64s) MarshalJSON() ([]byte, error) {
	ss := []string{}
	for _, v := range q {
		ss = append(ss, strconv.FormatFloat(v, 'f', -1, 64))
	}
	return json.Marshal(strings.Join(ss, sep))
}

// "ac,wa,tle" <=> []string{"ac", "wa", "tle"}
type Strings []string

func (q *Strings) UnmarshalJSON(raw []byte) error {
	*q = (*q)[:0]
	var ss string
	if err := json.Unmarshal(raw, &ss); err != nil {
		return err
	}
	*q = strings.Split(ss, sep)
	return nil
}

func (q Strings) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.Join(q, sep))
}
