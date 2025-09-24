package hpfs_test

import (
	"os"
	"testing"
)

const dependenceDir = "testdata/dependence"

type Dependence struct{}

var dependence Dependence

func TestMain(m *testing.M) {
	if skip := os.Getenv("SKIP_HPFS_TEST"); skip != "" {
		return
	}
	// 执行测试用例
	code := m.Run()
	os.Exit(code)
}
