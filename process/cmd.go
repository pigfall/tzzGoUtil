package process

import(
    "os/exec"
    "bytes"
)


// 执行 命令 并返回 标准输出和标准错误输出
func ExeOutput(name string,args ...string)(out string,errOut string,err error){
    cmd :=exec.Command(name,args...)
    var outBuffer bytes.Buffer
    var errBuffer  bytes.Buffer
    cmd.Stdout=&outBuffer
    cmd.Stderr = &errBuffer
    err = cmd.Run()
    return outBuffer.String(),errBuffer.String(),err
}
