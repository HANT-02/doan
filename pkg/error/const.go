package errors

const (
	ErrCodeDataInvalid         = "data.invalid"
	ErrCodeDataConflict        = "data.conflict"
	ErrCodeInternalServer      = "internal.server_error"
	ErrCodeTokenExpired        = "token.expired"
	ErrCodeUnauthorized        = "unauthorized"
	ErrCodeDataNotFound        = "data.not_found"
	ErrCodeForbidden           = "forbidden"
	ErrCodeUnprocessableEntity = "unprocessable_entity"
)

const (
	HTTPStatusTokenExpired = 469
)

var (
	ErrDataInvalid = New(
		ErrorTypeValidation,
		ErrCodeDataInvalid,
		"Data is invalid",
	).WithUserMessage("Dữ liệu không hợp lệ. Vui lòng kiểm tra lại.")
	ErrDataNotFound = New(
		ErrorTypeNotFound,
		ErrCodeDataNotFound,
		"Data not found",
	).WithUserMessage("Dữ liệu không tìm thấy. Vui lòng kiểm tra lại.")
	ErrDataConflict = New(
		ErrorTypeConflict,
		ErrCodeDataConflict,
		"Data conflict",
	).WithUserMessage("Dữ liệu bị xung đột. Vui lòng thử lại với thông tin khác.")
	ErrInternalServer = New(
		ErrorTypeInternal,
		ErrCodeInternalServer,
		"Internal server error",
	).WithUserMessage("Đã xảy ra lỗi hệ thống. Vui lòng thử lại sau.")
	ErrTokenExpired = New(
		ErrorTypeTokenExpired,
		ErrCodeTokenExpired,
		"Token expired",
	).WithUserMessage("Token đã hết hạn. Vui lòng đăng nhập lại.")
	ErrUnauthorized = New(
		ErrorTypeAuthorization,
		ErrCodeUnauthorized,
		"Unauthorized",
	).WithUserMessage("unauthorized")
	ErrForbidden = New(
		ErrorTypeForbidden,
		ErrCodeForbidden,
		"Forbidden",
	).WithUserMessage("Bạn không có quyền thực hiện hành động này")
	ErrUnprocessableEntity = New(
		ErrorTypeUnprocessableEntity,
		ErrCodeUnprocessableEntity,
		"Unprocessable entity",
	).WithUserMessage("Yêu cầu không thể xử lý do dữ liệu không hợp lệ hoặc thiếu thông tin cần thiết.")
)
