package errors

// ErrorType định nghĩa loại lỗi
type ErrorType uint

const (
	// ErrorTypeUnknown lỗi không xác định
	ErrorTypeUnknown ErrorType = iota
	// ErrorTypeValidation lỗi validation
	ErrorTypeValidation
	// ErrorTypeNotFound lỗi không tìm thấy
	ErrorTypeNotFound
	// ErrorTypeConflict lỗi xung đột
	ErrorTypeConflict
	// ErrorTypeAuthorization lỗi phân quyền
	ErrorTypeAuthorization
	// ErrorTypePermission lỗi quyền truy cập
	ErrorTypePermission
	// ErrorTypeTimeout lỗi timeout
	ErrorTypeTimeout
	// ErrorTypeInternal lỗi hệ thống
	ErrorTypeInternal
	// ErrorTypeExternal lỗi từ dịch vụ bên ngoài
	ErrorTypeExternal
	// ErrorTypeTokenExpired lỗi token hết hạn
	ErrorTypeTokenExpired
	// ErrorTypeForbidden lỗi truy cập bị cấm
	ErrorTypeForbidden
	// ErrorTypeGone lỗi đã bị xóa
	ErrorTypeGone
	// ErrorTypeUnprocessableEntity lỗi không thể xử lý
	ErrorTypeUnprocessableEntity
)

// String trả về tên của loại lỗi dưới dạng chuỗi
func (e ErrorType) String() string {
	names := map[ErrorType]string{
		ErrorTypeUnknown:             "unknown",
		ErrorTypeValidation:          "validation",
		ErrorTypeNotFound:            "not_found",
		ErrorTypeConflict:            "conflict",
		ErrorTypeAuthorization:       "authorization",
		ErrorTypePermission:          "permission",
		ErrorTypeTimeout:             "timeout",
		ErrorTypeInternal:            "internal",
		ErrorTypeExternal:            "external",
		ErrorTypeTokenExpired:        "token_expired",
		ErrorTypeForbidden:           "forbidden",
		ErrorTypeGone:                "gone",
		ErrorTypeUnprocessableEntity: "unprocessable_entity",
	}

	if name, ok := names[e]; ok {
		return name
	}
	return "unknown"
}
