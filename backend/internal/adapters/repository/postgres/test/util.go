package test

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
)

func EntityColumns(entity interface{}) []string {
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

func EntityValues(entity interface{}) []driver.Value {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var values []driver.Value
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			value := v.Field(i).Interface()
			values = append(values, value)
		}
	}

	return values
}

func UpdateQueryString(entity interface{}, tableName string) string {
	columnNames := EntityColumns(entity)
	params := make([]string, len(columnNames))
	for i, columnName := range columnNames {
		params[i] = fmt.Sprintf("%s = $%d", columnName, i+1)
	}
	paramsString := strings.Join(params, ", ")
	return fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		tableName, paramsString, len(columnNames)+1)
}

func InsertQueryString(entity interface{}, tableName string) string {
	columnNames := EntityColumns(entity)
	values := make([]string, len(columnNames))
	for i, _ := range columnNames {
		values[i] = fmt.Sprintf("$%d", i+1)
	}
	valuesString := strings.Join(values, ", ")
	columnsString := strings.Join(columnNames, ", ")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		tableName, columnsString, valuesString)
}
