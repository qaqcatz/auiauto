package putils

import (
	"auiauto/perrorx"
	"bufio"
	"bytes"
	"os/exec"
)

// 调用/bin/bash -c执行myCmdStr, 出现异常时默认合并输出流和错误流到error中
func Sh(myCmdStr string) (string, *perrorx.ErrorX) {
	exeCmd := exec.Command("/bin/bash", "-c", myCmdStr)

	stdout, _ := exeCmd.StdoutPipe()
	stderr, _ := exeCmd.StderrPipe()
	defer stdout.Close()
	defer stderr.Close()

	if err := exeCmd.Start(); err != nil {
		return "", perrorx.NewErrorXShellStart(myCmdStr, err.Error(), nil)
	}
	// 输出流
	scanner := bufio.NewScanner(stdout)
	var outBuf bytes.Buffer
	for scanner.Scan() {
		outBuf.WriteString(scanner.Text()+"\n")
	}
	myOut := outBuf.String()
	// 错误流
	scanner = bufio.NewScanner(stderr)
	var errBuf bytes.Buffer
	for scanner.Scan() {
		errBuf.WriteString(scanner.Text()+"\n")
	}
	myErr := errBuf.String()

	if myErr != "" {
		// 异常时默认合并错误流和输出流
		return "", perrorx.NewErrorXShellExecute(myCmdStr, myErr+"\n"+myOut, nil)
	}
	return myOut, nil
}