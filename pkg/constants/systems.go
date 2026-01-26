package constants

type DateTimeFormat string

const (
	DatetimeFormat    DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat        DateTimeFormat = "2006-01-02"
	TimeFormat        DateTimeFormat = "15:04:05"
	DateTimeIOSFormat DateTimeFormat = "2006-01-02T15:04:05Z"
)

type SystemENV string

const (
	EnvDebug   SystemENV = "debug"
	EnvRelease SystemENV = "release"
)

const (
	ConfigTypeFile   = "file"
	ConfigTypeRemote = "remote"
)

type HeaderKey string

const (
	HeaderKeyDeviceID      HeaderKey = "Device-ID"
	HeaderKeyDeviceModel   HeaderKey = "Device-Model"
	HeaderKeyAppVersion    HeaderKey = "App-Version"
	HeaderKeyOsVersion     HeaderKey = "Os-Version"
	HeaderKeyOs            HeaderKey = "Os"
	HeaderKeyAuthorization HeaderKey = "Authorization"
)

const (
	AuthenticationScheme = "Bearer "
)

const (
	ContextKeyDBTransaction = "context_db_transaction"
)
