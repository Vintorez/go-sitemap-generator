package stm

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/imdario/mergo"
)

// SetBuilderElementValue if it will change to struct from map if the future's
// author is feeling a bothersome in this function.
func SetBuilderElementValue(elm *etree.Element, data map[string]interface{}, basekey string) bool {
	key := basekey
	ts, tk := spaceDecompose(elm.Tag)
	_, sk := spaceDecompose(elm.Space)

	if elm.Tag != "" && ts != "" && tk != "" {
		key = fmt.Sprintf("%s:%s", elm.Space, basekey)
	} else if sk != "" {
		key = fmt.Sprintf("%s:%s", sk, basekey)
	}

	if values, ok := data[basekey]; ok {
		switch value := values.(type) {
		case nil:
		default:
			child := elm.CreateElement(key)
			child.SetText(fmt.Sprint(value))
		case int:
			child := elm.CreateElement(key)
			child.SetText(fmt.Sprint(value))
		case string:
			child := elm.CreateElement(key)
			child.SetText(value)
		case float64, float32:
			child := elm.CreateElement(key)
			child.SetText(fmt.Sprint(value))
		case time.Time:
			child := elm.CreateElement(key)
			child.SetText(value.Format(time.RFC3339))
		case bool:
			_ = elm.CreateElement(fmt.Sprintf("%s:%s", key, key))
		case []int:
			for _, v := range value {
				child := elm.CreateElement(key)
				child.SetText(fmt.Sprint(v))
			}
		case []string:
			for _, v := range value {
				child := elm.CreateElement(key)
				child.SetText(v)
			}
		case interface{}:
			var childkey string
			if sk == "" {
				childkey = fmt.Sprintf("%s:%s", key, key)
			} else {
				childkey = fmt.Sprint(key)
			}

			switch value := values.(type) {
			case []URL:
				for _, v := range value {
					child := elm.CreateElement(childkey)
					for ck := range v {
						SetBuilderElementValue(child, v, ck)
					}
				}
			case URL:
				child := elm.CreateElement(childkey)
				for ck := range value {
					SetBuilderElementValue(child, value, ck)
				}
			}
		}

		return true
	}
	return false
}

// MergeMap TODO: Slow function: It wants to change fast function
func MergeMap(src, dst map[string]interface{}) map[string]interface{} {
	mergo.MapWithOverwrite(&dst, src)
	return dst
}

// ToLowerString converts lower strings from including capital or upper strings.
func ToLowerString(befores []string) (afters []string) {
	for _, name := range befores {
		afters = append(afters, strings.ToLower(name))
	}
	return afters
}

// URLJoin TODO: Too slowly
func URLJoin(src string, joins ...string) string {
	var u *url.URL
	lastnum := len(joins)
	base, _ := url.Parse(src)

	for i, j := range joins {
		if !strings.HasSuffix(j, "/") && lastnum > (i+1) {
			j = j + "/"
		}

		u, _ = url.Parse(j)
		base = base.ResolveReference(u)
	}

	return base.String()
}

// spaceDecompose is separating strings for the SetBuilderElementValue
func spaceDecompose(str string) (space, key string) {
	colon := strings.IndexByte(str, ':')
	if colon == -1 {
		return "", str
	}
	return str[:colon], str[colon+1:]
}
