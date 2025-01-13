package api

import (
	"simplebank/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldlevel validator.FieldLevel) bool {
	if currency,ok:=fieldlevel.Field().Interface().(string);ok{
		return util.IsSupportCurrency(currency)
	}
	return false
}