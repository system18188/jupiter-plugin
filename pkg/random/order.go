package random

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	rd   *rand.Rand
	node *snowflake.Node
	once sync.Once
)

// 时间订单号 年两位，天三位，时分秒6位，毫秒6位 （17位 + 随机数）
func TimeOrderStr(n int) string {
	now := time.Now()
	order := now.Format("06-150405.999999999")
	order = Remove(order, ".") // 删除点号
	order = Replace(order, "-", fmt.Sprintf("%03d", now.YearDay()), 1)
	count := len(order)
	if count < 17 {
		order += Repeat("0", 17-count)
	}
	if n > 0 {
		rdinit()
		return fmt.Sprint(order, rd.Intn(n))
	}
	return order
}

// 生成随机数字
func IntN(n int) int {
	rdinit()
	return rd.Intn(n)
}

func rdinit() {
	once.Do(func() { rd = rand.New(rand.NewSource(time.Now().UnixNano())) })
}

// 格式 1076459955600494592
func NodeStr() string {
	nodeinit()
	return node.Generate().String()
}

func nodeinit() {
	once.Do(func() {
		var err error
		node, err = snowflake.NewNode(1)
		if err != nil {
			panic("snowflake.Node:" + err.Error())
			return
		}
	})
}

// Replace returns a copy of the string s with the first n
// non-overlapping instances of old replaced by new.
// If n < 0, there is no limit on the number of replacements.
func Replace(s, old, new string, n int) string { return strings.Replace(s, old, new, n) }


// IsBlank checks if a string is empty ("") or whitespace.
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Remove removes all occurrences of a substring from within the source string.
func Remove(s, remove string) string {
	if IsBlank(s) || IsBlank(remove) {
		return s
	}
	return strings.Replace(s, remove, "", -1)
}

// Repeat returns a new string consisting of count copies of the string s.
func Repeat(s string, count int) string { return strings.Repeat(s, count) }
