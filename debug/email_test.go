package debug

import(
    "testing"
)

func TestSendMail(t *testing.T){
    req :=&MailRequest{
        Addr:"smtp.163.com:smtp",
        User:"tzzNotify@163.com",
        Passwd:"tzzNotify999",
        SendTo:[]string{"tzzNotify@163.com"},
        Content:"你怎么知道的",
    }
    err := SendMail(req)
    if err != nil{
        t.Error(err)
    }
}
