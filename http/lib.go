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


func DoRequestX(ctx context.Context,req *stdhttp.Request,optionsHeader stdhttp.Header)(resBodyBytes []byte,err error){
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

func DoRequest(ctx context.Context, method string, url string, reqBody io.Reader, optionsHeader stdhttp.Header) (resBodyBytes []byte, err error) {
	req, err := stdhttp.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Build http request object failed: %w", err)
	}
	return DoRequestX(ctx,req,optionsHeader)
}

func DoRequestThenJsonUnMarshalX(
	ctx context.Context,
	req *stdhttp.Request,
	resEntityToUnMarshal interface{},
	optionsHeader stdhttp.Header,
	ifPrintResBody bool,
)(error){
	resBodyBytes, err := DoRequestX(ctx, req,optionsHeader)
	if err != nil {
		return err
	}
	if ifPrintResBody{
		fmt.Println(string(resBodyBytes))
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
