package api

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getIntQueryParameter(c *gin.Context, parameter string, defaultValue int) int {
	queryParameters := c.Request.URL.Query()

	value := defaultValue
	valueStr, hasValue := queryParameters[parameter]
	if hasValue {
		value, _ = strconv.Atoi(valueStr[0])
	}

	return value
}

func getInt64QueryParameter(c *gin.Context, parameter string, defaultValue int64) int64 {
	queryParameters := c.Request.URL.Query()

	value := defaultValue
	valueStr, hasValue := queryParameters[parameter]
	if hasValue {
		value, _ = strconv.ParseInt(valueStr[0], 10, 64)
	}

	return value
}

func getInt64ListQueryParameter(c *gin.Context, parameter string) []int64 {
	queryParameters := c.Request.URL.Query()

	ret := make([]int64, 0)

	valueStr, hasValue := queryParameters[parameter]
	if hasValue {
		stringList := strings.Split(valueStr[0], ",")
		for _, s := range stringList {
			value, err := strconv.ParseInt(s, 10, 64)
			if err == nil {
				ret = append(ret, value)
			}
		}
	}

	return ret
}

func getStringQueryParameter(c *gin.Context, parameter string) (string, error) {
	queryParameters := c.Request.URL.Query()

	valueStr, hasValue := queryParameters[parameter]
	if hasValue {
		return url.QueryUnescape(valueStr[0])
	}

	return "", errors.New("required parameter '" + parameter + "' not provided in query string")
}
