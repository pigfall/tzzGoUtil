package errors

import (
	"fmt"
)

type Result struct {
	err error
}

func NewResult() *Result {
	return &Result{}
}

func (this *Result) IsErr() bool {
	return this.err != nil
}

func (this *Result) AddErr(err error) *Result {
	if this.err == nil {
		this.err = err
		return this
	}
	this.err = fmt.Errorf("%v, %w", this.err, err)
	return this
}

func (this *Result) ToErr() error {
	return this.err
}
