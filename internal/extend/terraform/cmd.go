package terraform

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// ApplyResource terraform apply
func ApplyResource(fileDir string) (outputs map[string]struct {
	Value string `json:"value"`
}, err error) {
	// 应用配置
	var out []byte
	out, err = terraformApply(fileDir)
	if err != nil {
		return
	}

	// 获取输出值
	out, err = terraformOutput(fileDir)
	if err != nil {
		return
	}

	// 解析输出
	if err = json.Unmarshal(out, &outputs); err != nil {
		err = fmt.Errorf("output parsing failed: %s", err)
		return
	}
	return
}

// DestroyResource terraform destroy
func DestroyResource(fileDir string) error {
	_, err := terraformDestroy(fileDir)
	return err
}

func terraformApply(fileDir string) (out []byte, err error) {
	cmd := exec.Command("terraform", "apply", "-auto-approve", "-input=false")
	cmd.Dir = fileDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("apply failed: %s\n%s", err, string(out))
	}
	return
}

func terraformOutput(fileDir string) (out []byte, err error) {
	cmd := exec.Command("terraform", "output", "-json")
	cmd.Dir = fileDir
	out, err = cmd.Output()
	if err != nil {
		err = fmt.Errorf("output failed: %s", err)
	}
	return
}

func terraformDestroy(fileDir string) (out []byte, err error) {
	cmd := exec.Command("terraform", "destroy", "-auto-approve", "-input=false")
	cmd.Dir = fileDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("destroy failed: %s", err)
	}
	return
}
