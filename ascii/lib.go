package ascii

import(
    "unicode/utf8"
)

func IsASCII(b byte)bool{
    return b < utf8.RuneSelf
}

func IsAlpha(b byte)bool{
    return   ('A' <= b && b<= 'Z') || ('a'<=b && b<='z')
}

func IsNumber(n byte)bool{
    return  '0' <= n && n<='9'
}

func IsAlphaNum(n byte)bool{
    return IsNumber(n) || IsAlpha(n)
}
