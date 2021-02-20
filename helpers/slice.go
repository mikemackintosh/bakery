package helpers

import "github.com/mikemackintosh/bakery/pantry"

// Append will append to a slice
func Append(data interface{}, toAppend interface{}) interface{} {
	switch toAppend.(type) {
	case []string:
		for _, v := range toAppend.([]string) {
			data = append(data.([]string), v)
		}
		return data
	case []int:
		for _, v := range toAppend.([]int) {
			data = append(data.([]int), v)
		}
		return data
	case []bool:
		for _, v := range toAppend.([]bool) {
			data = append(data.([]bool), v)
		}
		return data
	case []pantry.PantryItem:
		for _, v := range toAppend.([]pantry.PantryItem) {
			data = append(data.([]pantry.PantryItem), v)
		}
		return data
	}
	return nil
}

// Prepend will prepend to a slice
// Append will append to a slice
func Prepend(toPrepend interface{}, data interface{}) interface{} {
	switch toPrepend.(type) {
	case []string:
		data = append(toPrepend.([]string), data.([]string)...)
		return data
	case []int:
		data = append(toPrepend.([]int), data.([]int)...)
		return data
	case []bool:
		data = append(toPrepend.([]bool), data.([]bool)...)
		return data
	case []pantry.PantryItem:
		data = append(toPrepend.([]pantry.PantryItem), data.([]pantry.PantryItem)...)
		return data
	}
	return nil
}
