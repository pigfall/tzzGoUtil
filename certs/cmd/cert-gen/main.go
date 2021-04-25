package main

import(
    "os"
    "fmt"
    "io/ioutil"
)


func main(){
    if len(os.Args) < 2{
        printHelpMsg()
        os.Exit(0)
    }

    subCmd := os.Args[1]
    var cfgPath string = "config.json"
    if len(os.Args) >=3 {
        cfgPath = os.Args[2]
    }

    cfgContent,err := ioutil.ReadFile(cfgPath)
    if err != nil{
        fmt.Printf("< Load config file %s failed >: %v\n",cfgPath,err)
        os.Exit(1)
    }


    switch subCmd{
    case  SUB_CMD_NAME_RSA:
        err = subCmdRsa(cfgContent)
    default:
        fmt.Printf("< Unknown subCmd %s >\n",subCmd)
        os.Exit(1)
    }
    if err != nil{
        fmt.Println("< Generate failed >: %w",err)
        os.Exit(1)
    }
}

func printHelpMsg(){
    fmt.Println("unimpl")
}
