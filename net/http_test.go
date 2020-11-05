package net

import(
    "strings"
    "context"
    "testing"
    "net/http"
    "net/url"
    "io/ioutil"
)

func TestHttpGet(t *testing.T){
    var reqUrl = "http://www.baidu.com"
    res,err := http.Get(reqUrl)
    if err != nil{
        t.Fatal(err)
    }
    defer res.Body.Close()
    data,err := ioutil.ReadAll(res.Body)
    if err != nil{
        t.Fatal(err)
    }
    if res.StatusCode !=http.StatusOK{
        // handle res err
        t.Error("res error:",string(data))
    }else{
        // handle ok 
        t.Log("res ok: ",string(data))
    }

    // 自己构造一个  get 请求
    req,err := http.NewRequestWithContext(context.Background(),http.MethodGet,reqUrl,nil)
    // 使用默认的 http client
    res,err = http.DefaultClient.Do(req)
    if err != nil{
        t.Fatal(err)
    }
    defer res.Body.Close()
    data,err = ioutil.ReadAll(res.Body)
    if err != nil{
        t.Fatal(err)
    }
    if res.StatusCode !=http.StatusOK{
        // handle res err
        t.Error("res error:",string(data))
    }else{
        // handle ok 
        t.Log("res ok: ",string(data))
    }
}

func HttpPostForm(t *testing.T){
    req,err := http.NewRequestWithContext(context.Background(),http.MethodPost,"testurl",strings.NewReader(url.Values{"key":[]string{"value"}}.Encode()))
    if err != nil{
        t.Fatal(err)
    }
    // 使用 form表单需要特别设置
    req.Header.Set("Content-Type","application/x-www-form-urlencoded")
    _,err = http.DefaultClient.Do(req)
    if err != nil{
        t.Fatal(err)
    }
}
