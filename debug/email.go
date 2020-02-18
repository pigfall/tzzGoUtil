package debug

import(
    "net/smtp"
    "strings"

)

type MailRequest struct{
    Addr string
    User string
    Passwd string
    SendTo []string
    Content string
}

func SendMail(req *MailRequest)error{
    host := strings.Split(req.Addr,":")
    auth := smtp.PlainAuth("",req.User,req.Passwd,host[0])
    content_type := "Content-Type: text/plain" + "; charset=UTF-8"
    msg := []byte("To: " +strings.Join(req.SendTo,";") + "\r\nFrom: " + req.User+ ">\r\nSubject: " + "\r\n" + content_type + "\r\n\r\n" +req.Content)
    err := smtp.SendMail(req.Addr,auth,req.User,req.SendTo,msg)
    return err
}


