package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	stdhttp "net/http"
)

func DoRequest(ctx context.Context, method string, url string, reqBody io.Reader, optionsHeader stdhttp.Header) (resBodyBytes []byte, err error) {
	req, err := stdhttp.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Build http request object failed: %w", err)
	}
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
