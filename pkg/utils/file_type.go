package utils

import (
	"doan/pkg/constants"
	"doan/pkg/types"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	AllowedFileTypes = map[string]bool{
		// Documents
		constants.MimeTypePDF:            true,
		constants.MimeTypeMSWord:         true,
		constants.MimeTypeWordDocx:       true,
		constants.MimeTypeMSExcel:        true,
		constants.MimeTypeExcelXlsx:      true,
		constants.MimeTypeMSPowerPoint:   true,
		constants.MimeTypePowerPointPptx: true,
		constants.MimeTypeRTF:            true,
		constants.MimeTypeTextPlain:      true,
		constants.MimeTypeTextCSV:        true,
		constants.MimeTypeJSON:           true,
		constants.MimeTypeXML:            true,
		constants.MimeTypeHTML:           true,
		constants.MimeTypeMarkdown:       true,
		constants.MimeTypeYAML:           true,
		constants.MimeTyeExe:             true,
		constants.MimeTypeSQL:            true,

		// Images
		constants.MimeTypeImageJPEG: true,
		constants.MimeTypeImagePNG:  true,
		constants.MimeTypeImageGIF:  true,
		constants.MimeTypeImageWebP: true,
		constants.MimeTypeImageBMP:  true,
		constants.MimeTypeImageTIFF: true,
		constants.MimeTypeImageSVG:  true,
		constants.MimeTypeIco:       true,
		constants.MimeTypeHEIC:      true,
		constants.MimeTypeHEIF:      true,

		// Audio
		constants.MimeTypeAudioMP3: true,
		constants.MimeTypeAudioWAV: true,
		constants.MimeTypeAudioOGG: true,

		// Video
		constants.MimeTypeVideoMP4:  true,
		constants.MimeTypeVideoWebM: true,
		constants.MimeTypeVideoMOV:  true,

		// Archives
		constants.MimeTypeZip:  true,
		constants.MimeTypeGzip: true,
		constants.MimeTypeTar:  true,
		constants.MimeTypeRar:  true,

		// Generic fallback
		constants.MimeTypeOctetStream: true,
	}

	ExtensionMap = map[string]string{
		constants.MimeTypePDF:            constants.ExtPDF,
		constants.MimeTypeMSWord:         constants.ExtDOC,
		constants.MimeTypeWordDocx:       constants.ExtDOCX,
		constants.MimeTypeMSExcel:        constants.ExtXLS,
		constants.MimeTypeExcelXlsx:      constants.ExtXLSX,
		constants.MimeTypeMSPowerPoint:   constants.ExtPPT,
		constants.MimeTypePowerPointPptx: constants.ExtPPTX,
		constants.MimeTypeTextPlain:      constants.ExtTXT,
		constants.MimeTypeTextCSV:        constants.ExtCSV,
		constants.MimeTypeJSON:           constants.ExtJSON,
		constants.MimeTypeXML:            constants.ExtXML,
		constants.MimeTypeImageJPEG:      constants.ExtJPG,
		constants.MimeTypeImagePNG:       constants.ExtPNG,
		constants.MimeTypeImageGIF:       constants.ExtGIF,
		constants.MimeTypeImageWebP:      constants.ExtWebP,
		constants.MimeTypeImageBMP:       constants.ExtBMP,
		constants.MimeTypeImageTIFF:      constants.ExtTIFF,
		constants.MimeTypeImageSVG:       constants.ExtSVG,
		constants.MimeTypeVideoMP4:       constants.ExtMP4,
		constants.MimeTypeVideoWebM:      constants.ExtWebM,
		constants.MimeTypeVideoMOV:       constants.ExtMOV,
		constants.MimeTypeAudioMP3:       constants.ExtMP3,
		constants.MimeTypeAudioWAV:       constants.ExtWAV,
		constants.MimeTypeAudioOGG:       constants.ExtOGG,
		constants.MimeTypeZip:            constants.ExtZIP,
		constants.MimeTypeGzip:           constants.ExtGZ,
		constants.MimeTypeRar:            constants.ExtRAR,
		constants.MimeTypeTar:            constants.ExtTAR,
		constants.MimeTypeOctetStream:    constants.ExtBIN,
		constants.MimeTyeExe:             constants.ExtEXE,
		constants.MimeTypeHTML:           constants.ExtHtml,
		constants.MimeTypeMarkdown:       constants.ExtMd,
		constants.MimeTypeYAML:           constants.ExtYaml,
		constants.MimeTypeSQL:            constants.ExtSQL,
		constants.MimeTypeIco:            constants.ExtICO,
		constants.MimeTypeHEIC:           constants.ExtHEIC,
		constants.MimeTypeHEIF:           constants.ExtHEIF,
	}
)

func DetectFileType(headerBytes []byte, fileName string) (*types.FileTypeInfo, error) {
	result := &types.FileTypeInfo{}

	// Use a Go standard library http.DetectContentType
	detectedMime := http.DetectContentType(headerBytes)
	// Refine detection for specific types that http.DetectContentType might miss
	refinedMime := refineDetection(headerBytes, fileName, detectedMime)
	result.DetectedMime = refinedMime
	result.Extension = getExtensionForMime(refinedMime)
	return result, nil
}

func refineDetection(header []byte, filename, httpDetectedMime string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	// Check for Office documents (ZIP-based)
	if strings.HasPrefix(httpDetectedMime, constants.MimeTypeZip) {
		switch ext {
		case constants.ExtDOCX:
			return constants.MimeTypeWordDocx
		case constants.ExtXLSX:
			return constants.MimeTypeExcelXlsx
		case constants.ExtPPTX:
			return constants.MimeTypePowerPointPptx
		}
	}

	// Check for specific video formats
	if strings.HasPrefix(httpDetectedMime, constants.MimeTypeOctetStream) {
		// MP4 detection
		if len(header) > 8 && string(header[4:8]) == "ftyp" {
			return constants.MimeTypeVideoMP4
		}
		// Check file extension for common cases
		switch ext {
		case constants.ExtMP4:
			return constants.MimeTypeVideoMP4
		case constants.ExtWebM:
			return constants.MimeTypeVideoWebM
		case constants.ExtMOV:
			return constants.MimeTypeVideoMOV
		case constants.ExtMP3:
			return constants.MimeTypeAudioMP3
		case constants.ExtWAV:
			return constants.MimeTypeAudioWAV
		case constants.ExtPDF:
			// Double-check PDF
			if len(header) > 4 && string(header[:4]) == "%PDF" {
				return constants.MimeTypePDF
			}
		case constants.ExtOGG:
			return constants.MimeTypeAudioOGG
		case constants.ExtDOC:
			return constants.MimeTypeMSWord
		case constants.ExtXLS:
			return constants.MimeTypeMSExcel
		case constants.ExtPPT:
			return constants.MimeTypeMSPowerPoint
		case constants.MimeTypeHEIC:
			return constants.MimeTypeHEIC
		case constants.MimeTypeHEIF:
			return constants.MimeTypeHEIF
		}
	}

	// Text files with specific extensions
	if strings.HasPrefix(httpDetectedMime, constants.MimeTypeTextPlain) {
		switch ext {
		case constants.ExtCSV:
			return constants.MimeTypeTextCSV
		case constants.ExtJSON:
			return constants.MimeTypeJSON
		case constants.ExtXML:
			return constants.MimeTypeXML
		case constants.ExtTXT:
			return constants.MimeTypeTextPlain
		case constants.ExtSVG:
			return constants.MimeTypeImageSVG
		case constants.ExtSQL:
			return constants.MimeTypeSQL
		}
	}

	if strings.HasPrefix(httpDetectedMime, constants.MimeTypeXML) {
		switch ext {
		case constants.ExtSVG:
			return constants.MimeTypeImageSVG
		}
	}

	if strings.HasPrefix(httpDetectedMime, constants.MimeTypeImageGIF) {
		switch ext {
		case constants.ExtGIF:
			return constants.MimeTypeImageGIF
		case constants.ExtTAR:
			return constants.MimeTypeTar
		}
	}

	if strings.HasPrefix(httpDetectedMime, constants.MimeTypeVideoMP4) {
		switch ext {
		case constants.MimeTypeHEIC:
			return constants.MimeTypeHEIC
		case constants.MimeTypeHEIF:
			return constants.MimeTypeHEIF
		}
	}

	return httpDetectedMime
}

func getExtensionForMime(mimeType string) string {
	baseMimeType := strings.Split(mimeType, ";")[0]
	if ext, exists := ExtensionMap[baseMimeType]; exists {
		return ext
	}

	// Fallback to mime package
	if exts, err := mime.ExtensionsByType(baseMimeType); err == nil && len(exts) > 0 {
		return exts[0]
	}

	return constants.ExtBIN
}
