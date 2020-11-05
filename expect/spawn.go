package expect

import(
    "regexp"
)

type Expecter struct{
}


func (this *Expecter) Cmd()*Expecter{

    return nil
}

func (this *Expecter) Expect(reg *regexp.Regexp)(*Expecter){

    return nil
}


func (this *Expecter) Run()(error){

    return nil
}

