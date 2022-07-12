package jisilu

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

func MarshalAndUnmarshal(dst interface{}, src interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

// 查找第一个子匹配
func FindFirstSubMatch(expr, str string) string {

	sub := regexp.MustCompile(expr).FindStringSubmatch(str)
	if len(sub) >= 1 {
		return strings.TrimSpace(sub[1])
	}
	return ""
}

// 保留小数点后n位
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

func ParseDate(date string) string {
	if len(date) == 0 {
		return ""
	}
	t, _ := time.Parse("2006-01-02", date)
	return fmt.Sprintf("%d月%d日", t.Month(), t.Day())
}

func GetMarket(bondID string) string {
	if strings.HasPrefix(bondID, "11") {
		return "沪"
	} else if strings.HasPrefix(bondID, "12") {
		return "深"
	}
	return ""
}
