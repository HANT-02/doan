package x_error

// Common x-error codes
const (
	InternalServer    = "INTERNAL_SERVER_ERROR"
	CodeDataNotFound  = "DATA_NOT_FOUND"
	CodeDataMissing   = "DATA_MISSING"
	DataInvalid       = "DATA_INVALID"
	ConfigKeyNotFound = "CONFIG_KEY_NOT_FOUND"
	// ErrorInternalServerError is a constant for x-error code when internal server x-error
	ErrorInternalServerError = "INTERNAL_SERVER_ERROR"
	// ErrorCodeDataInvalid is a constant for x-error code when data is invalid
	ErrorCodeDataInvalid = "DATA_INVALID"
	// ErrorCodeDataNotFound is a constant for x-error code when data is not found
	ErrorCodeDataNotFound      = "DATA_NOT_FOUND"
	ErrorCodeDataMissing       = "DATA_MISSING"
	ErrorCodeConfigKeyNotFound = "CONFIG_KEY_NOT_FOUND"
	ErrorCodeUseCaseProcessing = "USE_CASE_PROCESSING"
	ErrorCodeDataDuplicate     = "DATA_DUPLICATE"
	ErrorInvalidInput          = "INVALID_INPUT"
)

const (
	KafkaBrokersRequired               = "KAFKA_BROKERS_REQUIRED"
	KafkaTopicCreateFailed             = "KAFKA_TOPIC_CREATE_FAILED"
	KafkaPublishFailed                 = "KAFKA_PUBLISH_FAILED"
	KafkaTopicOptionRequired           = "KAFKA_TOPIC_OPTION_REQUIRED"
	KafkaTopicNameRequired             = "KAFKA_TOPIC_NAME_REQUIRED"
	KafkaTopicNumPartitionsInvalid     = "KAFKA_TOPIC_NUM_PARTITIONS_INVALID"
	KafkaTopicReplicationFactorInvalid = "KAFKA_TOPIC_REPLICATION_FACTOR_INVALID"
	KafkaTopicAlreadyConsumed          = "KAFKA_TOPIC_ALREADY_CONSUMED"
)

const (
	OSRequired          = "OS_REQUIRED"
	DeviceIDRequired    = "DEVICE_ID_REQUIRED"
	OSVersionRequired   = "OS_VERSION_REQUIRED"
	AppVersionRequired  = "APP_VERSION_REQUIRED"
	DeviceModelRequired = "DEVICE_MODEL_REQUIRED"
)

// User x-error codes
const (
	// CodeInvalidUsernameOrPassword is a constant for x-error code when username or password is invalid
	CodeInvalidUsernameOrPassword = "INVALID_USERNAME_OR_PASSWORD"
	// CodeUserNotFound is a constant for x-error code when user-management is not found
	CodeUserNotFound    = "USER_NOT_FOUND"
	CodeProfileNotFound = "PROFILE_NOT_FOUND"
	// CodeUserAlreadyExists is a constant for x-error code when user-management already exists
	CodeUserAlreadyExists = "USER_ALREADY_EXISTS"
	// CodeProfileAlreadyExists is a constant for x-error code when profile-management already exists
	CodeProfileAlreadyExists = "Profile_ALREADY_EXISTS"
	// CodeInvalidToken is a constant for x-error code when token is invalid
	CodeInvalidToken = "INVALID_TOKEN"
	// CodeUserHasNoRole is a constant for x-error code when user-management has no role
	CodeUserHasNoRole = "USER_HAS_NO_ROLE"
	// CodeUnexpectedSigningMethod is a constant for x-error code when unexpected signing method
	CodeUnexpectedSigningMethod = "UNEXPECTED_SIGNING_METHOD"
	// CodeUsernameExists is a constant for x-error code when username exists
	CodeUsernameExists = "USERNAME_EXISTS"
	// CodeUsernameRequired is a constant for x-error code when username is required
	CodeUsernameRequired = "USERNAME_REQUIRED"
	// CodeCreatedByRequired is a constant for x-error code when created by is required
	CodeCreatedByRequired = "CREATED_BY_REQUIRED"
	// CodeCreatedByNotExists is a constant for x-error code when created by is not exists
	CodeCreatedByNotExists = "CREATED_BY_NOT_EXISTS"
	// CodeRoleIdsRequired is a constant for x-error code when role ids is required
	CodeRoleIdsRequired = "ROLE_IDS_REQUIRED"
	// CodeRoleIdsNotExists is a constant for x-error code when role ids is not exists
	CodeRoleIdsNotExists = "ROLE_IDS_NOT_EXISTS"
	// CodePasswordRequired is a constant for x-error code when password is required
	CodePasswordRequired = "PASSWORD_REQUIRED"
	// CodeAccessTokenExpired is a constant for x-error code when token is expired
	CodeAccessTokenExpired = "TOKEN_EXPIRED"
	// CodeUnauthorized is a constant for x-error code when unauthorized
	CodeUnauthorized = "UNAUTHORIZED"
	// CodeDoNotPermissionCreateUploadFolder is a constant for x-error code when user-management is not active
	CodeDoNotPermissionCreateUploadFolder = "DO_NOT_PERMISSION_CREATE_UPLOAD_FOLDER"
	// DoNotPermissionCreateNewFile is a constant for x-error code when user-management is not rejected
	DoNotPermissionCreateNewFile = "DO_NOT_PERMISSION_CREATE_NEW_FILE"
	// CanNotReadFileContent is a constant for x-error code when user-management is not deactivated
	CanNotReadFileContent = "CAN_NOT_READ_FILE_CONTENT"
	// DoNotPermissionToUploadFile is a constant for x-error code when user-management is not deactivated
	DoNotPermissionToUploadFile = "DO_NOT_PERMISSION_TO_UPLOAD_FILE"
	// Forbidden is a constant for x-error code when user-management is not deactivated
	Forbidden = "FORBIDDEN"
	// EmailRequired is a constant for x-error code when email is required
	EmailRequired = "EMAIL_REQUIRED"
	// OTPNotMatch is a constant for x-error code when OTP is not match
	OTPNotMatch = "OTP_NOT_MATCH"
	// AccountNeedToVerify is a constant for x-error code when account need to verify
	AccountNeedToVerify = "ACCOUNT_NEED_TO_VERIFY"
	// RefreshTokenNotFound is a constant for x-error code when refresh token is not found
	RefreshTokenNotFound = "REFRESH_TOKEN_NOT_FOUND"
	// RefreshTokenExpired is a constant for x-error code when refresh token is expired
	RefreshTokenExpired = "REFRESH_TOKEN_EXPIRED"
	// RefreshTokenInvalid is a constant for x-error code when refresh token is invalid
	RefreshTokenInvalid = "REFRESH_TOKEN_INVALID"
	// DataExisted is a constant for x-error code when data is existed
	DataExisted = "DATA_EXISTED"
)

const (
	SMSTemplateNotFound = "SMS_TEMPLATE_NOT_FOUND"
)
const (
	CantOpenFile    = "CANT_OPEN_FILE"
	NoFileUploaded  = "NO_FILE_UPLOADED"
	FileNotFound    = "FILE_NOT_FOUND"
	ImageNotFound   = "IMAGE_NOT_FOUND"
	InvalidFileType = "INVALID_FILE_TYPE"
	FileReadError   = "FILE_READ_ERROR"
)

// Review x-error codes
const (
	ReviewExisted = "REVIEW_EXISTED"
)

// UserFollow x-error codes
const (
	UserFollowExisted = "USER_FOLLOW_EXISTED"
)

const (
	MeiliSearchCreateCollectionFailed = "MEILI_SEARCH_CREATE_COLLECTION_FAILED"
	MeiliSearchCreateDocumentsFailed  = "MEILI_SEARCH_CREATE_DOCUMENTS_FAILED"
	MeiliSearchDeleteDocumentsFailed  = "MEILI_SEARCH_DELETE_DOCUMENTS_FAILED"
	MeiliSearchSearchFailed           = "MEILI_SEARCH_SEARCH_FAILED"
)

const (
	IndexNotFound = "INDEX_NOT_FOUND"
)

const (
	CategoryExisted = "CATEGORY_EXITED"
)
