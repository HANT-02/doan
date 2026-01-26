package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AppError đại diện cho lỗi ứng dụng
type AppError struct {
	Type        ErrorType              // Loại lỗi
	Code        string                 // Mã lỗi duy nhất, dạng "module.error_code"
	Message     string                 // Thông điệp lỗi cho dev
	UserMessage string                 // Thông điệp lỗi thân thiện cho người dùng
	Cause       error                  // Lỗi gốc
	Stack       string                 // Stack trace
	Meta        map[string]interface{} // Metadata bổ sung
	Ops         []string               // Operations (function names) đã xử lý lỗi
	Retryable   bool                   // Có thể retry hay không
	Err         error                  // Lỗi wrapped (để tương thích với code cũ)
	GRPCStatus  *status.Status         // Status code cho GRPC
}

func (e AppError) Error() string {
	message := ""
	if e.Message != "" {
		message = e.Message
	}
	if e.Err != nil {
		message = e.Err.Error()
	}
	if e.UserMessage != "" {
		message = fmt.Sprintf("%s (%s)", message, e.UserMessage)
	}
	if message == "" {
		message = "unknown error"
	}
	return message
}

// Unwrap trả về lỗi gốc
func (e AppError) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return e.Err
}

// HTTPStatusCode trả về HTTP status code phù hợp với loại lỗi
func (e AppError) HTTPStatusCode() int {
	switch e.Type {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeAuthorization:
		return http.StatusUnauthorized
	case ErrorTypePermission:
		return http.StatusForbidden
	case ErrorTypeTimeout:
		return http.StatusGatewayTimeout
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	case ErrorTypeExternal:
		return http.StatusBadGateway
	case ErrorTypeTokenExpired:
		return HTTPStatusTokenExpired
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeGone:
		return http.StatusGone
	case ErrorTypeUnprocessableEntity:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

// WithMeta thêm metadata cho lỗi
func (e AppError) WithMeta(key string, value interface{}) AppError {
	if e.Meta == nil {
		e.Meta = make(map[string]interface{})
	}
	e.Meta[key] = value
	return e
}

// WithUserMessage thêm thông điệp thân thiện cho người dùng
func (e AppError) WithUserMessage(message string) AppError {
	e.UserMessage = message
	return e
}

// AsRetryable đánh dấu lỗi có thể retry
func (e AppError) AsRetryable() AppError {
	e.Retryable = true
	return e
}

// WithOp thêm operation vào chuỗi operations
func (e AppError) WithOp(op string) AppError {
	if e.Ops == nil {
		e.Ops = []string{}
	}
	e.Ops = append(e.Ops, op)
	return e
}

// Is kiểm tra target có phải là AppError và có cùng Type không
func (e AppError) Is(target error) bool {
	var t AppError
	if errors.As(target, &t) {
		return e.Type == t.Type
	}
	return errors.Is(e.Err, target) || errors.Is(e.Cause, target)
}

// New tạo một AppError mới
func New(errType ErrorType, code, message string) AppError {
	return AppError{
		Type:    errType,
		Code:    code,
		Message: message,
		Stack:   captureStack(),
		Ops:     []string{},
	}
}

// Errorf tạo một AppError với format
func Errorf(errType ErrorType, code, format string, args ...interface{}) AppError {
	return AppError{
		Type:    errType,
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Stack:   captureStack(),
		Ops:     []string{},
	}
}

// Wrap bọc một error với thông tin bổ sung
func Wrap(err error, op, message string) error {
	if err == nil {
		return nil
	}

	var appErr AppError
	if errors.As(err, &appErr) {
		// Nếu đã là AppError, thêm op vào ops
		appErr.Ops = append(appErr.Ops, op)
		if message != "" {
			appErr.Message = message + ": " + appErr.Message
		}
		return appErr
	}

	// Tạo AppError mới nếu là error thông thường
	return AppError{
		Type:    ErrorTypeUnknown,
		Message: message,
		Cause:   err,
		Stack:   captureStack(),
		Ops:     []string{op},
	}
}

// WrapWithType bọc một error với type cụ thể
func WrapWithType(err error, op string, message string, errType ErrorType, code string) error {
	if err == nil {
		return nil
	}

	var appErr AppError
	if errors.As(err, &appErr) {
		// Nếu đã là AppError, giữ Type
		appErr.Ops = append(appErr.Ops, op)
		if code != "" {
			appErr.Code = code
		}
		if message != "" {
			appErr.Message = message + ": " + appErr.Message
		}
		return appErr
	}

	// Tạo AppError mới
	return AppError{
		Type:    errType,
		Code:    code,
		Message: message,
		Cause:   err,
		Stack:   captureStack(),
		Ops:     []string{op},
	}
}

func WrapWithMeta(err error, op string, message string, key string, value string) error {
	if err == nil {
		return nil
	}

	var appErr AppError
	if !errors.As(err, &appErr) {
		appErr = AppError{
			Type:  ErrorTypeUnknown,
			Cause: err,
			Stack: captureStack(),
		}
	}

	appErr.Ops = append(appErr.Ops, op)
	if message != "" {
		appErr.Message = message + ": " + appErr.Message
	}
	return appErr.WithMeta(key, value)
}

// captureStack bắt stack trace
func captureStack() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var builder strings.Builder
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&builder, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return builder.String()
}

// Các hàm tạo lỗi tiện ích - giữ lại từ phiên bản cũ để tương thích ngược

// NewUnknownError tạo lỗi không xác định
func NewUnknownError(message string, cause error) AppError {
	return AppError{
		Type:    ErrorTypeUnknown,
		Message: message,
		Cause:   cause,
		Stack:   captureStack(),
		Code:    "system.unknown_error",
	}
}

// NewValidationError tạo lỗi validation
func NewValidationError(message string, cause error) AppError {
	return AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Cause:   cause,
		Stack:   captureStack(),
		Code:    "validation.invalid_input",
	}
}

// NewNotFoundError tạo lỗi không tìm thấy
func NewNotFoundError(message string, cause error) AppError {
	return AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Cause:   cause,
		Stack:   captureStack(),
		Code:    ErrCodeDataNotFound,
	}
}

// NewForbiddenError tạo lỗi không có quyền truy cập
func NewForbiddenError(message string, cause error) AppError {
	return AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Cause:   cause,
		Stack:   captureStack(),
		Code:    ErrCodeForbidden,
	}
}

// NewConflictError tạo lỗi xung đột
func NewConflictError(message string, cause error) AppError {
	return AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Cause:   cause,
		Stack:   captureStack(),
		Code:    "entity.conflict",
	}
}

// NewAuthorizationError tạo lỗi phân quyền
func NewAuthorizationError(message string, cause error) AppError {
	return AppError{
		Type:    ErrorTypeAuthorization,
		Message: message,
		Cause:   cause,
		Stack:   captureStack(),
		Code:    "auth.unauthorized",
	}
}

// NewInternalError tạo lỗi hệ thống
func NewInternalError(message string, cause error) AppError {
	return AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Cause:   cause,
		Stack:   captureStack(),
		Code:    "system.internal_error",
	}
}

// ErrorTypeToGRPCCode chuyển đổi ErrorType sang gRPC code
func ErrorTypeToGRPCCode(errType ErrorType) codes.Code {
	switch errType {
	case ErrorTypeValidation:
		return codes.InvalidArgument
	case ErrorTypeNotFound:
		return codes.NotFound
	case ErrorTypeConflict:
		return codes.AlreadyExists
	case ErrorTypeAuthorization:
		return codes.Unauthenticated
	case ErrorTypePermission:
		return codes.PermissionDenied
	case ErrorTypeTimeout:
		return codes.DeadlineExceeded
	case ErrorTypeExternal:
		return codes.Unavailable
	case ErrorTypeInternal:
		return codes.Internal
	case ErrorTypeTokenExpired:
		return codes.Unauthenticated
	case ErrorTypeForbidden:
		return codes.PermissionDenied
	default:
		return codes.Unknown
	}
}

// GRPCCodeToErrorType chuyển đổi gRPC code sang ErrorType
func GRPCCodeToErrorType(code codes.Code) ErrorType {
	switch code {
	case codes.InvalidArgument:
		return ErrorTypeValidation
	case codes.NotFound:
		return ErrorTypeNotFound
	case codes.AlreadyExists:
		return ErrorTypeConflict
	case codes.Unauthenticated:
		return ErrorTypeAuthorization
	case codes.PermissionDenied:
		return ErrorTypePermission
	case codes.DeadlineExceeded:
		return ErrorTypeTimeout
	case codes.Unavailable:
		return ErrorTypeExternal
	case codes.Internal:
		return ErrorTypeInternal
	default:
		return ErrorTypeUnknown
	}
}

func GRPCCodeToAppErrCode(code codes.Code) string {
	switch code {
	case codes.InvalidArgument:
		return ErrCodeDataInvalid
	case codes.NotFound:
		return ErrCodeDataNotFound
	case codes.AlreadyExists:
		return ErrCodeDataConflict
	case codes.Unauthenticated:
		return ErrCodeForbidden
	case codes.PermissionDenied:
		return ErrCodeUnauthorized
	case codes.Internal:
		return ErrCodeInternalServer
	default:
		return ErrCodeInternalServer
	}
}

// MakeGRPCStatus để sử dụng mã code tương ứng với loại lỗi
func (e AppError) MakeGRPCStatus() AppError {
	st := status.New(ErrorTypeToGRPCCode(e.Type), e.Message)

	data, _ := json.Marshal(e.Meta)

	metadata := map[string]string{
		"user_message": e.UserMessage,
		"message":      e.Message,
		"stack":        e.Stack,
		"ops":          strings.Join(e.Ops, ","),
		"retryable":    fmt.Sprintf("%t", e.Retryable),
		"meta":         string(data),
		"type":         string(e.Type),
		"code":         e.Code,
	}

	// Xử lý trường hợp e.Cause có thể là nil
	if e.Cause != nil {
		metadata["cause"] = e.Cause.Error()
	}

	st, _ = st.WithDetails(&errdetails.ErrorInfo{
		Reason:   e.Code,
		Domain:   "eduone.service",
		Metadata: metadata,
	})

	e.GRPCStatus = st
	return e
}

// AsGRPCError trả về lỗi dạng gRPC
func (e AppError) AsGRPCError() error {
	err := e
	if e.GRPCStatus == nil {
		err = e.MakeGRPCStatus()
	}
	return err.GRPCStatus.Err()
}

// FromGRPCError chuyển đổi gRPC error thành AppError
func FromGRPCError(err error) (AppError, bool) {
	if err == nil {
		return AppError{}, false
	}

	st := status.Convert(err)

	appErr := AppError{
		Type:       GRPCCodeToErrorType(st.Code()),
		Message:    st.Message(),
		GRPCStatus: st,
		Code:       GRPCCodeToAppErrCode(st.Code()),
	}

	// Trích xuất thông tin từ details
	for _, detail := range st.Details() {
		if errInfo, ok := detail.(*errdetails.ErrorInfo); ok {
			appErr.Code = errInfo.Reason
			appErr.UserMessage = errInfo.Metadata["user_message"]
			appErr.Stack = errInfo.Metadata["stack"]

			// Phục hồi ops
			if opsStr, ok := errInfo.Metadata["ops"]; ok && opsStr != "" {
				appErr.Ops = strings.Split(opsStr, ",")
			}

			// Phục hồi retryable
			if retryableStr, ok := errInfo.Metadata["retryable"]; ok {
				appErr.Retryable = retryableStr == "true"
			}

			// Phục hồi meta
			if metaStr, ok := errInfo.Metadata["meta"]; ok && metaStr != "" {
				json.Unmarshal([]byte(metaStr), &appErr.Meta)
			}

			return appErr, true
		}
	}

	// Nếu không có details, trả về lỗi cơ bản
	return appErr, true
}

func HandleError(err error) AppError {
	if err == nil {
		return AppError{}
	}

	var appErr AppError
	if errors.As(err, &appErr) {
		return appErr.MakeGRPCStatus()

	}

	// Nếu không phải AppError, tạo một AppError mới với thông tin cơ bản
	appErr = AppError{
		Type:    ErrorTypeUnknown,
		Message: err.Error(),
		Cause:   err,
		Stack:   captureStack(),
		Code:    "system.unknown_error",
	}
	return appErr.MakeGRPCStatus()
}
