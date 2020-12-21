package main

import "reflect"

func isFieldExported(typeOfField reflect.StructField) bool {
	return typeOfField.PkgPath == ""
}
