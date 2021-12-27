package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	stdhttp "net/http"
	"net/url"
	"strings"
)

func ReadResBody(res *stdhttp.Response)([]byte,error){
	defer res.Body.Close()
	if res.StatusCode != stdhttp.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			bodyBytes = []byte(fmt.Errorf("Read body data failed: %w", err).Error())
		}
		return nil, fmt.Errorf(string(bodyBytes))
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body data failed: %w", err)
	}
	return bodyBytes, nil
}


func doRequestX(ctx context.Context,req *stdhttp.Request,optionsHeader stdhttp.Header,options ...OptionFiller)(resBodyBytes []byte,err error){
	ops := newOptions(options...)
	for k, v := range optionsHeader {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	res, err := stdhttp.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != stdhttp.StatusOK && !ops.StatusCodeOk(res.StatusCode){
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			bodyBytes = []byte(fmt.Errorf("Read body data failed: %w", err).Error())
		}
		return bodyBytes, fmt.Errorf("httpStatusCode: %v, %s",res.StatusCode,string(bodyBytes))
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body data failed: %w", err)
	}
	return bodyBytes, nil
}

func DoRequestX(ctx context.Context,req *stdhttp.Request,optionsHeader stdhttp.Header)(resBodyBytes []byte,err error){
	return doRequestX(ctx,req,optionsHeader)
}

func DoRequest(ctx context.Context, method string, url string, reqBody io.Reader, optionsHeader stdhttp.Header) (resBodyBytes []byte, err error) {
	req, err := stdhttp.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Build http request object failed: %w", err)
	}
	return DoRequestX(ctx,req,optionsHeader)
}

func DoRequestThenJsonUnMarshalAntReturnResBodyDataWithUnmarshal(
	ctx context.Context,
	unmarshal func([]byte,interface{})error,
	req *stdhttp.Request,
	resEntityToUnMarshal interface{},
	optionsHeader stdhttp.Header,
	ifPrintResBody bool,
	options ...OptionFiller,
)(resBodyBytes []byte,err error){
	resBodyBytes, err = doRequestX(ctx, req,optionsHeader,options...)
	if ifPrintResBody{
		fmt.Println("Response body content:\n ",string(resBodyBytes))
	}
	if err != nil {
		return resBodyBytes,err
	}
	err = unmarshal(resBodyBytes, resEntityToUnMarshal)
	if err != nil{
		return resBodyBytes,err
	}
	return resBodyBytes,nil
}

func DoRequestThenJsonUnMarshalAntReturnResBodyData(
	ctx context.Context,
	req *stdhttp.Request,
	resEntityToUnMarshal interface{},
	optionsHeader stdhttp.Header,
	ifPrintResBody bool,
	options ...OptionFiller,
)(resBodyData []byte,err error){
	resBodyBytes, err := doRequestX(ctx, req,optionsHeader,options...)
	if ifPrintResBody{
		fmt.Println("Response body content:\n ",string(resBodyBytes))
	}
	if err != nil {
		return resBodyBytes,err
	}
	err = json.Unmarshal(resBodyBytes, resEntityToUnMarshal)
	if err != nil{
		return resBodyBytes,err
	}
	return resBodyBytes,nil
}


func DoRequestThenJsonUnMarshalX(
	ctx context.Context,
	req *stdhttp.Request,
	resEntityToUnMarshal interface{},
	optionsHeader stdhttp.Header,
	ifPrintResBody bool,
)(error){
	resBodyBytes, err := DoRequestX(ctx, req,optionsHeader)
	if ifPrintResBody{
		fmt.Println("Response body content:\n ",string(resBodyBytes))
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(resBodyBytes, resEntityToUnMarshal)
}


func DoRequestThenJsonUnMarshal(
	ctx context.Context,
	method string,
	url string,
	reqBody io.Reader,
	resEntityToUnMarshal interface{},
	optionsHeader stdhttp.Header,
) (err error) {
	resBodyBytes, err := DoRequest(ctx, method, url, reqBody, optionsHeader)
	if err != nil {
		return err
	}
	return json.Unmarshal(resBodyBytes, resEntityToUnMarshal)
}

func DoPostRequestWithXWWWFormUrlEncoded(
	ctx context.Context, url string, params url.Values, optionsHeader stdhttp.Header,
) (resBody []byte, err error) {
	optionsHeader.Set("Content-Type", "application/x-www-form-urlencoded")
	return DoRequest(
		ctx, stdhttp.MethodPost, url, strings.NewReader(params.Encode()), optionsHeader,
	)
}

func DoPostRequestUrlEncode_ThenUnMarshalJson(
	ctx context.Context, url string, params url.Values, optionsHeader stdhttp.Header, resEntityToUnMarshal interface{},
) error {
	resBytes, err := DoPostRequestWithXWWWFormUrlEncoded(ctx, url, params, optionsHeader)
	if err != nil {
		return err
	}

	return json.Unmarshal(resBytes, resEntityToUnMarshal)
}
