package entity

import (
	"fmt"
	"reflect"
	"strings"
)

func entityColumns(entity interface{}) []string {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var fields []string
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i).Tag.Get("db")
			if field != "" {
				fields = append(fields, field)
			}
		}
	} else if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			fields = append(fields, key.String())
		}
	}

	return fields
}

func UpdateQueryString(entity interface{}, tableName string) string {
	columnNames := entityColumns(entity)
	params := make([]string, len(columnNames))
	for i, columnName := range columnNames {
		params[i] = fmt.Sprintf("%s = :%s", columnName, columnName)
	}
	paramsString := strings.Join(params, ", ")
	return fmt.Sprintf("UPDATE %s SET %s WHERE id = :id",
		tableName, paramsString)
}

func InsertQueryString(entity interface{}, tableName string) string {
	columnNames := entityColumns(entity)
	values := make([]string, len(columnNames))
	for i, columnName := range columnNames {
		values[i] = fmt.Sprintf(":%s", columnName)
	}
	valuesString := strings.Join(values, ", ")
	columnsString := strings.Join(columnNames, ", ")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		tableName, columnsString, valuesString)
}
