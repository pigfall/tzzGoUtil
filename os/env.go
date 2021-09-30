package os

import(
	"fmt"
		"os"
)



func MustFindEnvVarAndNotEmpty(envVarName string)(value string,err error) {
	value,exist:= os.LookupEnv(envVarName)
	if !exist {
		return "",fmt.Errorf("EnvVar %s not exists",envVarName)
	}

	if len(value) == 0{
		return "",fmt.Errorf("Found EnvVar %s but length is 0",envVarName)
	}

	return value,nil
}
