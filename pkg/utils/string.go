package utils

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func CheckLength(s string, min, max int) bool {
	if utf8.RuneCountInString(s) < min || utf8.RuneCountInString(s) > max {
		return false
	}
	return true
}

// Check if string contains special characters
func ContainsStrangeSpecialChar(s string) bool {
	re := regexp.MustCompile(`[@!#$%^&*()_+\-=\[\]{};':"\\|,.<>/?~]`)
	return re.MatchString(s)
}

func ContainsPercentSign(s string) bool {
	re := regexp.MustCompile(`%`)
	return re.MatchString(s)
}

func IsValidContainsStrangeSpecialCharForTextInput(s string) bool {
	re := regexp.MustCompile(`^[\p{L}0-9.\-\_\+\*\&@\(\),;:\?!\/' ]+$`)
	return re.MatchString(s)
}

// Check if string contains special characters
func ContainsStrangeSpecialCharForCode(s string) bool {
	re := regexp.MustCompile(`[@!#$%^&*+=\[\]{};':"\\|,.<>/?~]`)
	return re.MatchString(s)
}

// Check if string contains special characters minus (& - . _)
func ContainsStrangeSpecialCharExceptAllowed(s string) bool {
	re := regexp.MustCompile(`[@!#$%^*()+=\[\]{};':"\\|,<>/?~]`) // Exclude & - . _
	return re.MatchString(s)
}

// Check if string contains emoji
func ContainsEmoji(s string) bool {
	for _, r := range s {
		// Kiểm tra các Unicode ranges chứa emoji
		if isEmoji(r) {
			return true
		}
	}
	return false
}

// isEmoji kiểm tra xem một rune có phải là emoji không
func isEmoji(r rune) bool {
	// Emoji ranges trong Unicode
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
		(r >= 0x1F1E0 && r <= 0x1F1FF) || // Regional indicator symbols
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) || // Dingbats
		(r >= 0xFE00 && r <= 0xFE0F) || // Variation Selectors
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1F018 && r <= 0x1F270) || // Various symbols
		(r >= 0x238C && r <= 0x2454) || // Misc symbols
		(r >= 0x20D0 && r <= 0x20FF) // Combining Diacritical Marks for Symbols
}

func IsValidInputForCode(s string) bool {
	re := regexp.MustCompile(`^[\p{L}0-9.\-_\+]+$`)
	return re.MatchString(s)
}

// Check if string is all digits
func IsAllDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Check if string is all uppercase
func IsAllUpper(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// Check if string is all lowercase
func IsAllLower(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsLower(r) {
			return false
		}
	}
	return true
}

// Check if string is a script code (e.g., "Latn", "Cyrl", "Hans")
// ContainsScriptRegex sử dụng regex để detect script (comprehensive)
func ContainsScriptRegex(s string) bool {
	// regex patterns for detecting script codes
	patterns := []string{
		`(?i)<script[^>]*>.*?</script>`,             // Script tags
		`(?i)<script[^>]*>`,                         // Opening script tag
		`(?i)javascript\s*:`,                        // Javascript protocol
		`(?i)vbscript\s*:`,                          // VBScript protocol
		`(?i)on\w+\s*=\s*["\'].*?["\']`,             // Event handlers
		`(?i)on\w+\s*=\s*\w+`,                       // Event handlers without quotes
		`(?i)<iframe[^>]*>`,                         // Iframe tags
		`(?i)<object[^>]*>`,                         // Object tags
		`(?i)<embed[^>]*>`,                          // Embed tags
		`(?i)<form[^>]*>`,                           // Form tags
		`(?i)expression\s*\(`,                       // CSS expression
		`(?i)url\s*\(\s*javascript:`,                // CSS with javascript
		`(?i)@import`,                               // CSS import
		`(?i)<meta[^>]*http-equiv`,                  // Meta refresh
		`(?i)<link[^>]*href\s*=\s*["\']javascript:`, // Link with javascript
	}

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, s)
		if err == nil && matched {
			return true
		}
	}

	return false
}

func ContainsSpecialCharsExceptComma(s string) bool {
	forbiddenChars := "!@#$%^&*()_+-=[]{}|\\:;\"'<>?/~`"

	for _, char := range s {
		if strings.ContainsRune(forbiddenChars, char) {
			return true
		}
	}

	return false
}

// Hàm kiểm tra chuỗi chỉ chứa ký tự đặc biệt
func IsAllSpecialChars(s string) bool {
	if len(s) == 0 {
		return false // chuỗi rỗng không tính là đặc biệt
	}
	for _, r := range s {
		// Nếu là chữ cái hoặc số → false
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Check if string contains Vietnamese accented characters
func ContainsVietnameseAccent(s string) bool {
	re := regexp.MustCompile(`[àáảãạầấẩẫậằắẳẵặèéẻẽẹềếểễệìíỉĩịòóỏõọồốổỗộờớởỡợùúủũụừứửữựỳýỷỹÀÁẢÃẠẦẤẨẪẬẰẮẲẴẶÈÉẺẼẸỀẾỂỄỆÌÍỈĨỊÒÓỎÕỌỒỐỔỖỘỜỚỞỠỢÙÚỦŨỤỪỨỬỮỰỲÝỶỸỴ]`)
	return re.MatchString(s)
}

// Check if string contains special characters minus (- _ , ( ))
func ContainsStrangeSpecialCharExceptAllowedSet(s string) bool {
	re := regexp.MustCompile(`[@!#$%^&*.+=\[\]{};':"\\|,<>/?~]`) // Exclude ( ) - _
	return re.MatchString(s)
}

// Check if string contains space between words
func ContainsSpaceBetweenWords(s string) bool {
	re := regexp.MustCompile(`\s+`)
	return re.MatchString(s)
}

func IsInterfaceArrayString(value interface{}) bool {
	valueArray, ok := value.([]interface{})
	if !ok {
		return false
	}
	for _, item := range valueArray {
		if _, ok = item.(string); !ok {
			return false
		}
	}
	return true
}
