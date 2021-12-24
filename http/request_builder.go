package http

import(
	"strings"
	"encoding/json"
	"io"
	"context"
	stdhttp "net/http"
)

type RequestBuilder struct{
//	request *stdhttp.Request
	method string
	url string
	paramsToPutBody  paramsToPutBody
	paramsToPutUrl map[string]string
	headerToSet map[string]string
}

func NewRequestBuilder() *RequestBuilder{
	return &RequestBuilder{}
}

type paramsToPutBody interface{
	GetContentType() string
	GetParamReader() (io.Reader,error)
}

type jsonParam struct{
	entity interface{}
	marshal func(v interface{})([]byte,error)
}

func newJsonParam(entity interface{}) paramsToPutBody{
	return &jsonParam{entity:entity}
}

func withJsonParam(entity interface{},marshalFunc func(v interface{})([]byte,error))paramsToPutBody{
	return &jsonParam{
		entity:entity,
		marshal:marshalFunc,
	}
}

func (this *jsonParam) GetContentType()string{
	return "application/json"
}
func (this *jsonParam)GetParamReader()(io.Reader,error){
	var bytes []byte
	var err error
	var marshal = json.Marshal
	if this.marshal != nil{
		marshal = this.marshal
	}
	bytes,err = marshal(this.entity)
	if err != nil{
		return nil,err
	}
	
	return strings.NewReader(string(bytes)),nil
}

func (this *RequestBuilder) MethodPost() *RequestBuilder{
	this.method= stdhttp.MethodPost
	return this
}

func (this *RequestBuilder) MethodPatch() *RequestBuilder{
	this.method = stdhttp.MethodPatch
	return this
}

func (this *RequestBuilder) MethodDelete() *RequestBuilder{
	this.method = stdhttp.MethodDelete
	return this
}

func (this *RequestBuilder) MethodGet() *RequestBuilder{
	this.method = stdhttp.MethodGet
	return this
}

func (this *RequestBuilder) MethodPut() *RequestBuilder{
	this.method = stdhttp.MethodPut
	return this
}


func (this *RequestBuilder) URL(url string)*RequestBuilder{
	this.url = url
	return this
}

func (this *RequestBuilder) JsonParam(entity interface{})*RequestBuilder{
	this.paramsToPutBody= newJsonParam(entity)
	return this
}

func (this *RequestBuilder) JsonParamWithMarshalFunc(entity interface{},marshal func(v interface{})([]byte,error))*RequestBuilder{
	this.paramsToPutBody = withJsonParam(entity,marshal)
	return this
}

func (this *RequestBuilder) GetParamToPutBody()interface{}{
	return this.paramsToPutBody
}

func (this *RequestBuilder) SetHeader(key string,value string)(*RequestBuilder){
	if this.headerToSet == nil{
		this.headerToSet = make(map[string]string)
	}
	this.headerToSet[key] = value
	return this
}

func (this *RequestBuilder) PutParamsToUrl (params map[string]string)*RequestBuilder{
	if this.paramsToPutUrl == nil{
		this.paramsToPutUrl = params
		return this
	}
	for k,v := range params{
		this.paramsToPutUrl[k]=v
	}
	return this
}

func (this *RequestBuilder) Build(ctx context.Context)(*stdhttp.Request,error){
	var body io.Reader
	var err error
	if this.paramsToPutBody!= nil{
		body,err = this.paramsToPutBody.GetParamReader()
		if err != nil{
			return nil,err
		}
	}
	req,err := stdhttp.NewRequestWithContext(ctx,this.method,this.url,body)
	if err != nil{
		return nil,err
	}
	if this.paramsToPutBody!= nil{
		req.Header.Set("Content-Type",this.paramsToPutBody.GetContentType())
	}
	if this.headerToSet != nil{
		for k,v := range this.headerToSet{
			req.Header.Add(k,v)
		}
	}

	if this.paramsToPutUrl !=nil{
		query := req.URL.Query()
		for k,v := range this.paramsToPutUrl{
			query.Add(k,v)
		}
		req.URL.RawQuery = query.Encode()
	} 

	return req,nil
}
