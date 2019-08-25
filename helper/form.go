package helper

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ParseFormToArray(values url.Values,key string) []string {
	if val,ok := values[key];ok && len(val) > 0{
		return val
	}
	if val,ok := values[key+"[]"];ok && len(val) > 0{
		return val
	}
	return []string{}
}

func ToInt(input string) int {
	out,_ := strconv.Atoi(strings.TrimSpace(input))
	return out
}

func ToArrayInt(values url.Values,key string) []int {
	valueStrs := ParseFormToArray(values,key)
	data := []int{}
	for _,val := range valueStrs {
		data = append(data,ToInt(val))
	}
	return data
}

func Boolean(form string) bool {
	return (form == "true" || form == "1")
}

func ToFloat(input string) float64 {
	out,_ := strconv.ParseFloat(strings.TrimSpace(input),64)
	return out
}

func ToDate(val string) time.Time {
	t,err := time.Parse("02/01/2006",val)
	if err == nil {
		return t
	}
	return time.Now()
}

func ToDateWithErr(val string) (time.Time,error) {
	t,err := time.Parse("1/2/2006",val)
	return t,err
}