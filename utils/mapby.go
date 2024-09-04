package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// MapBy represents a custom struct tag
const MapByTag = "mapby"

// FillFromMapByTags fills fields of the first struct from the second struct based on mapby tags
func FillFromMapByTags(dest, src interface{}) error {
	destVal := reflect.ValueOf(dest).Elem()
	srcVal := reflect.ValueOf(src)

	if destVal.Kind() != reflect.Struct || srcVal.Kind() != reflect.Struct {
		return fmt.Errorf("both arguments must be structs")
	}

	destType := destVal.Type()
	// srcType := srcVal.Type()

	for i := 0; i < destVal.NumField(); i++ {
		destField := destVal.Field(i)
		destFieldType := destType.Field(i)
		mapByTag := destFieldType.Tag.Get(MapByTag)

		if mapByTag != "" {
			srcField := srcVal.FieldByName(mapByTag)
			if !srcField.IsValid() {
				return fmt.Errorf("field %s not found in the source struct", mapByTag)
			}

			if !destField.CanSet() {
				return fmt.Errorf("field %s in destination struct cannot be set", destFieldType.Name)
			}

			destField.Set(srcField)
		}
	}

	return nil
}

// FillByFieldNameAndType fills fields of the destination struct from the source struct based on field names and types
func FillByFieldNameAndType(dest, src interface{}) error {
	destVal := reflect.ValueOf(dest).Elem()
	srcVal := reflect.ValueOf(src)

	if destVal.Kind() != reflect.Struct || srcVal.Kind() != reflect.Struct {
		return fmt.Errorf("both arguments must be structs")
	}

	destType := destVal.Type()

	for i := 0; i < destVal.NumField(); i++ {
		destField := destVal.Field(i)
		destFieldType := destType.Field(i)
		destFieldName := destFieldType.Name

		srcField := srcVal.FieldByName(destFieldName)
		if !srcField.IsValid() {
			continue // Skip if the source struct does not have the same field name
		}

		srcFieldType := srcField.Type()
		if destField.Type() != srcFieldType {
			continue // Skip if the types of the fields are different
		}

		if !destField.CanSet() {
			return fmt.Errorf("field %s in destination struct cannot be set", destFieldName)
		}

		destField.Set(srcField)
	}

	return nil
}

func FillByFieldNameAndTypeSlice(destSlice, srcSlice interface{}) error {
	destSliceValue := reflect.ValueOf(destSlice)
	srcSliceValue := reflect.ValueOf(srcSlice)

	if destSliceValue.Kind() != reflect.Slice || srcSliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("both arguments must be slices")
	}

	if destSliceValue.Len() != srcSliceValue.Len() {
		return fmt.Errorf("slices must have the same length")
	}

	for i := 0; i < destSliceValue.Len(); i++ {
		destElem := destSliceValue.Index(i)
		srcElem := srcSliceValue.Index(i)

		if err := FillByFieldNameAndType(destElem.Addr().Interface(), srcElem.Interface()); err != nil {
			return err
		}
	}

	return nil
}

// new mapping
func MapByTagComplex(dest, src interface{}) error {
	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)

	if srcVal.Kind() != reflect.Ptr || destVal.Kind() != reflect.Ptr {
		return fmt.Errorf("both arguments must be pointers")
	}

	srcElem := srcVal.Elem()
	destElem := destVal.Elem()

	switch srcElem.Kind() {
	case reflect.Struct:
		return fillFromMapByTagsRecursive(destElem, srcElem)
	case reflect.Slice:
		srcSlice := srcElem
		destSlice := destElem

		if destSlice.Kind() != reflect.Slice {
			return fmt.Errorf("destination must be a slice")
		}

		destType := destSlice.Type().Elem()

		for i := 0; i < srcSlice.Len(); i++ {
			srcItem := srcSlice.Index(i)
			destItem := reflect.New(destType).Elem()

			err := fillFromMapByTagsRecursive(destItem, srcItem)
			if err != nil {
				return err
			}

			destSlice.Set(reflect.Append(destSlice, destItem))
		}

		return nil
	default:
		return fmt.Errorf("unsupported source type: %s", srcElem.Kind())
	}
}

func fillFromMapByTagsRecursive(destVal reflect.Value, srcVal reflect.Value) error {
	if destVal.Kind() != reflect.Struct || srcVal.Kind() != reflect.Struct {
		return fmt.Errorf("both arguments must be structs")
	}

	destType := destVal.Type()

	for i := 0; i < destVal.NumField(); i++ {
		destField := destVal.Field(i)
		destFieldType := destType.Field(i)
		mapByTag := destFieldType.Tag.Get(MapByTag)

		if mapByTag != "" {
			srcField := getNestedField(srcVal, mapByTag)
			if !srcField.IsValid() {
				return fmt.Errorf("field %s not found in the source struct", mapByTag)
			}

			if !destField.CanSet() {
				return fmt.Errorf("field %s in destination struct cannot be set", destFieldType.Name)
			}

			// Check if destination field is a pointer and source field is not
			if destField.Kind() == reflect.Ptr && srcField.Kind() != reflect.Ptr {
				newValue := reflect.New(destField.Type().Elem()).Elem()
				newValue.Set(srcField)
				destField.Set(newValue.Addr())
			} else if destField.Kind() == reflect.Ptr && srcField.Kind() == reflect.Ptr {
				if !srcField.IsNil() {
					destField.Set(srcField)
				}
			} else {
				destField.Set(srcField)
			}
		} else {
			srcFieldName := destFieldType.Name
			srcField := srcVal.FieldByName(srcFieldName)

			if srcField.IsValid() && srcField.Type() == destFieldType.Type {
				if !destField.CanSet() {
					return fmt.Errorf("field %s in destination struct cannot be set", destFieldType.Name)
				}
				destField.Set(srcField)
			}
		}

		if destFieldType.Type.Kind() == reflect.Struct {
			err := fillFromMapByTagsRecursive(destField, srcVal)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getNestedField(srcVal reflect.Value, mapByTag string) reflect.Value {
	if mapByTag == "" {
		return reflect.Value{}
	}

	tags := strings.Split(mapByTag, ".")
	currentField := srcVal

	for _, tag := range tags {
		currentField = currentField.FieldByName(tag)
		if !currentField.IsValid() {
			break
		}
	}

	return currentField
}
