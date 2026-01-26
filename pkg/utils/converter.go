package utils

import (
	apperrors "doan/pkg/error"
	"encoding/json"
	"reflect"
	"strings"
)

// --- Các hàm tiện ích khác của bạn (giữ nguyên) ---

func BytesToStruct[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func BytesToStructWithDefault[T any](data []byte, defaultValue T) T {
	result, err := BytesToStruct[T](data)
	if err != nil {
		return defaultValue
	}
	return result
}

func StructToBytes[T any](data T) ([]byte, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func StructToBytesWithDefault[T any](data T, defaultValue []byte) []byte {
	result, err := StructToBytes(data)
	if err != nil {
		return defaultValue
	}
	return result
}

func StructToString[T any](data T) (string, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func StructToStringWithDefault[T any](data T, defaultValue string) string {
	result, err := StructToString(data)
	if err != nil {
		return defaultValue
	}
	return result
}

func StringToStruct[T any](data string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func StringToStructWithDefault[T any](data string, defaultValue T) T {
	result, err := StringToStruct[T](data)
	if err != nil {
		return defaultValue
	}
	return result
}

func StringToBytes(data string) ([]byte, error) {
	return []byte(data), nil
}

func StringToBytesWithDefault(data string, defaultValue []byte) []byte {
	result, err := StringToBytes(data)
	if err != nil {
		return defaultValue
	}
	return result
}

func BytesToString(data []byte) (string, error) {
	return string(data), nil
}

func BytesToStringWithDefault(data []byte, defaultValue string) string {
	result, err := BytesToString(data)
	if err != nil {
		return defaultValue
	}
	return result
}
func ConvertMapClaimToStruct(data map[string]interface{}, output interface{}) error {
	val := reflect.ValueOf(output).Elem()
	for key, value := range data {
		field := val.FieldByNameFunc(func(fieldName string) bool {
			field, _ := val.Type().FieldByName(fieldName)
			return field.Tag.Get("json") == key
		})

		if field.IsValid() && field.CanSet() {
			v := reflect.ValueOf(value)

			if field.Type() == reflect.TypeOf([]string{}) && v.Type() == reflect.TypeOf([]interface{}{}) {
				var stringSlice []string
				for _, item := range value.([]interface{}) {
					stringSlice = append(stringSlice, item.(string))
				}
				v = reflect.ValueOf(stringSlice)
			} else if field.Type() != v.Type() {
				v = v.Convert(field.Type())
			}
			field.Set(v)
		}
	}
	return nil
}

func ConvertByJSONTag[T any, U any](from T) (U, error) {
	var to U
	bytes, err := json.Marshal(from)
	if err != nil {
		return to, err
	}
	err = json.Unmarshal(bytes, &to)
	if err != nil {
		return to, err
	}
	return to, nil
}

// --- PHIÊN BẢN HOÀN CHỈNH CỦA MapStructByJSONTag ---

func MapStructByJSONTag(src any, dst any) error {
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return apperrors.ErrInternalServer.WithOp("MapStructByJSONTag").WithUserMessage("dst must be a non-nil pointer to a struct")
	}

	dstElem := dstVal.Elem()
	if dstElem.Kind() != reflect.Struct {
		return apperrors.ErrInternalServer.WithOp("MapStructByJSONTag").WithUserMessage("dst must be a pointer to a struct")
	}

	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}
	if srcVal.Kind() != reflect.Struct {
		return apperrors.ErrInternalServer.WithOp("MapStructByJSONTag").WithUserMessage("src must be a struct or a pointer to a struct")
	}

	srcFieldMap := make(map[string]reflect.Value)
	buildSrcFieldMapRecursive(srcVal, srcFieldMap)

	// ⭐ Thay đổi: Không cần truyền `srcVal` vào hàm map của dst nữa.
	mapToDstRecursive(dstElem, srcFieldMap)

	return nil
}

func buildSrcFieldMapRecursive(val reflect.Value, fieldMap map[string]reflect.Value) {
	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)

		// Logic này đúng: đi sâu vào struct nhúng để lấy tất cả các trường con
		if fieldTyp.Anonymous && fieldVal.Kind() == reflect.Struct {
			buildSrcFieldMapRecursive(fieldVal, fieldMap)
		} else {
			tag := fieldTyp.Tag.Get("json")
			tagName := strings.Split(tag, ",")[0]
			if tagName != "" && tagName != "-" {
				fieldMap[tagName] = fieldVal
			}
		}
	}
}

// ⭐ Nâng cấp `mapToDstRecursive` để xử lý struct nhúng một cách tổng quát
func mapToDstRecursive(dstVal reflect.Value, srcFieldMap map[string]reflect.Value) {
	if dstVal.Kind() != reflect.Struct {
		return
	}

	typ := dstVal.Type()
	for i := 0; i < dstVal.NumField(); i++ {
		dstField := dstVal.Field(i)
		dstFieldInfo := typ.Field(i)

		// ⭐ LOGIC THỐNG NHẤT:
		// 1. Nếu là struct nhúng, chỉ cần đi sâu vào trong nó để xử lý các trường con.
		// Không cần so khớp tên struct nhúng nữa.
		if dstFieldInfo.Anonymous && dstField.Kind() == reflect.Struct {
			mapToDstRecursive(dstField, srcFieldMap)
			continue // Sau khi đi sâu, tiếp tục vòng lặp với trường tiếp theo của struct cha
		}

		// 2. Xử lý tất cả các trường có tên (không phải nhúng).
		if !dstField.CanSet() {
			continue
		}

		tag := dstFieldInfo.Tag.Get("json")
		tagName := strings.Split(tag, ",")[0]
		if tagName == "" || tagName == "-" {
			continue
		}

		if srcField, ok := srcFieldMap[tagName]; ok {
			// LOGIC XỬ LÝ STRUCT LỒNG NHAU (NESTED STRUCT)
			if srcField.Kind() == reflect.Struct && dstField.Kind() == reflect.Struct {
				if dstField.CanAddr() {
					// Gọi lại hàm map chính cho cặp struct con này
					MapStructByJSONTag(srcField.Interface(), dstField.Addr().Interface())
				}
				continue // Đã xử lý xong, chuyển sang trường tiếp theo
			}

			// --- Logic để xử lý các kiểu dữ liệu khác ---
			if srcField.Type().AssignableTo(dstField.Type()) {
				dstField.Set(srcField)
			} else if srcField.Type().ConvertibleTo(dstField.Type()) {
				dstField.Set(srcField.Convert(dstField.Type()))
			} else if srcField.Kind() == reflect.Ptr && !srcField.IsNil() && srcField.Type().Elem().ConvertibleTo(dstField.Type()) {
				dstField.Set(srcField.Elem().Convert(dstField.Type()))
			} else if dstField.Kind() == reflect.Ptr && srcField.Type().ConvertibleTo(dstField.Type().Elem()) {
				convertedVal := srcField.Convert(dstField.Type().Elem())
				newPtr := reflect.New(dstField.Type().Elem())
				newPtr.Elem().Set(convertedVal)
				dstField.Set(newPtr)
			} else if srcField.Kind() == reflect.Ptr && dstField.Kind() == reflect.Ptr && !srcField.IsNil() {
				if srcField.Type().Elem().ConvertibleTo(dstField.Type().Elem()) {
					srcElemVal := srcField.Elem()
					convertedVal := srcElemVal.Convert(dstField.Type().Elem())
					newDstPtr := reflect.New(dstField.Type().Elem())
					newDstPtr.Elem().Set(convertedVal)
					dstField.Set(newDstPtr)
				}
			}
		}
	}
}
