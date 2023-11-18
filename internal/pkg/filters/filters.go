package handlers

import (
	"cards/internal/models"
	"errors"
	"reflect"
	"strings"
)

var cardFields = getCardFields()

func getCardFields() []string {
    var field []string
    v := reflect.ValueOf(models.Card{})
    for i := 0; i < v.Type().NumField(); i++ {
        field = append(field, v.Type().Field(i).Tag.Get("json"))
    }
    return field
}

func StringInSlice(strSlice []string, s string) bool {
    for _, v := range strSlice {
        if v == s {
            return true
        }
    }
    return false
}

func ValidateAndReturnFilterMap(filter string) (map[string]string, error) {
    splits := strings.Split(filter, ".")
    if len(splits) != 2 {
        return nil, errors.New("malformed sortBy query parameter, should be field.orderdirection")
    }
    field, value := splits[0], splits[1]
    if !StringInSlice(cardFields, field) {
        return nil, errors.New("unknown field in filter query parameter")
    }

    return map[string]string{field: value}, nil
}