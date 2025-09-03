package utils

import (
	"net/http"
	"strconv"
)

// GetQueryInt64Slice parses a slice of int64s from a query parameter.
func GetQueryInt64Slice(r *http.Request, key string) ([]int64, error) {
	values, ok := r.URL.Query()[key]
	if !ok {
		return nil, nil
	}

	var result []int64
	for _, v := range values {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, nil
}

// GetQueryInt parses an int from a query parameter.
func GetQueryInt(r *http.Request, key string) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0, nil
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// GetQueryInt64 parses an int64 from a query parameter.
func GetQueryInt64(r *http.Request, key string) (int64, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0, nil
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}