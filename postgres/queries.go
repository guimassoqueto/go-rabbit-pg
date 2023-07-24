package pg

import (
	"fmt"
	"grp/types"
	"reflect"
	"strings"
)

var Test string = "abc"

func UpsertQuery(table string, product types.Product) string {
	var insert string
	var values string
	var onConflict string

	insert += fmt.Sprintf("INSERT INTO %s(", table)
	values += "VALUES("
	onConflict += "ON CONFLICT (id)\nDO UPDATE SET "

	structType := reflect.TypeOf(product)
	structValue := reflect.ValueOf(product)
	
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldValue := structValue.Field(i)

		key := strings.ToLower(fieldType.Name)
		value := fieldValue.Interface()

		insert += fmt.Sprintf("%s,", key)

		if reflect.TypeOf(value).Kind() == reflect.String {	
			values += fmt.Sprintf("'%s',", value)
			onConflict += fmt.Sprintf("%s = '%s',", key, value)
		} else if reflect.TypeOf(value).Kind() == reflect.Float64 {
			values += fmt.Sprintf("%.2f,", value)
			onConflict += fmt.Sprintf("%s = %.2f,", key, value)
		} else {
			values += fmt.Sprintf("%d,", value)
			onConflict += fmt.Sprintf("%s = %d,", key, value)
		}
	}

	insert += "updated_at)\n"
	values += "NOW())\n"
	onConflict += "updated_at = NOW();"

	return fmt.Sprintf("%s %s %s", insert, values, onConflict)
}