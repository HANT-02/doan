package constants

const (
	// RegexUsername format: 3-20 characters, only letters, numbers, and underscores
	RegexUsername = "^[a-zA-Z0-9_]{3,20}$"

	// RegexPassword format: 8-50 characters, at least one uppercase letter, one lowercase letter, one number, and one special character
	RegexPassword = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{8,50}$"

	// RegexEmail format: email address
	RegexEmail = "^[a-zA-Z0-9](\\.?[a-zA-Z0-9_%+-])*@[a-zA-Z0-9-]+(\\.[a-zA-Z0-9-]+)*\\.[a-zA-Z]{2,}$"

	// RegexPhoneNumber format: phone number, starting with 0 or +84, 9-10 digits
	RegexPhoneNumber = "^(?:\\+84|0)\\d{9,10}$"

	// RegexURL format: URL, starting with api:// or https:// and ending with a domain
	RegexURL = "^(https?):\\/\\/([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,}(?:\\/[a-zA-Z0-9\\-._~:/@!$&'()*+,;=%]*)?(?:\\?[a-zA-Z0-9\\-._~!$&'()*+,;=:@/?%]*)?(?:\\#[a-zA-Z0-9\\-._~!$&'()*+,;=:@/?%]*)?$"

	// RegexDate format: date, dd/MM/yyyy
	RegexDate = "^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/\\d{4}$"

	// RegexTime format: time, HH:mm:ss
	RegexTime = "^([01]\\d|2[0-3]):([0-5]\\d):([0-5]\\d)$"

	// RegexDateTime format: ISO 8601 date and time, yyyy-MM-ddTHH:mm:ssZ
	RegexDateTime = "^(19|20)\\d{2}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01])T([01]\\d|2[0-3]):([0-5]\\d):([0-5]\\d)(?:\\.\\d+)?(?:Z|[+-]\\d{2}:\\d{2})$"

	RegexDomainName = `^(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`
)
