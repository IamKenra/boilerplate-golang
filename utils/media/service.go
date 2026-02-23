package media

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"mime"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"boilerplate-golang/utils/storage"
)

var (
	ErrStorageUnavailable = errors.New("media storage is not configured")
	ErrFileRequired       = errors.New("file is required")
	ErrFileEmpty          = errors.New("file is empty")
)

type Category string

const (
	CategoryAssetImage    Category = "asset_image"
	CategoryAssetQRCode   Category = "asset_qrcode"
	CategoryEmployeeImage Category = "employee_image"

	VariantMain      = "main"
	VariantThumbnail = "thumbnail"

	defaultPresignExpiry = 15 * time.Minute
	commandTimeout       = 15 * time.Second

	formatPNG  outputFormat = "png"
	formatWebP outputFormat = "webp"
	formatJPEG outputFormat = "jpeg"
)

type UploadPayload struct {
	File        io.ReadSeeker
	Size        int64
	Filename    string
	ContentType string
	ReferenceID string
}

type StoredObject struct {
	Variant      string
	ObjectName   string
	ContentType  string
	Size         int64
	URL          string
	TemporaryURL string
}

type UploadResult struct {
	Bucket  string
	Objects []StoredObject
}

type Service struct {
	storage storage.FileStorage
}

func NewService(fs storage.FileStorage) *Service {
	return &Service{storage: fs}
}

func (s *Service) Upload(ctx context.Context, category Category, payload UploadPayload) (UploadResult, error) {
	var result UploadResult

	if s.storage == nil {
		return result, ErrStorageUnavailable
	}

	rule, ok := categoryRules[category]
	if !ok {
		return result, fmt.Errorf("unsupported media category: %s", category)
	}

	data, contentType, err := prepareInput(rule, payload)
	if err != nil {
		return result, err
	}

	img, err := decodeImage(data, contentType)
	if err != nil {
		return result, err
	}

	commonMetadata := buildBaseMetadata(payload)
	var uploaded []StoredObject

	for _, variant := range rule.Variants {
		processed := img
		if variant.MaxWidth > 0 || variant.MaxHeight > 0 {
			processed = resizeImage(processed, variant.MaxWidth, variant.MaxHeight)
		}

		encoded, err := encodeImage(processed, variant.Format)
		if err != nil {
			s.cleanup(ctx, uploaded)
			return result, fmt.Errorf("failed to encode %s variant: %w", variant.Name, err)
		}

		if variant.MaxBytes > 0 && int64(len(encoded)) > variant.MaxBytes {
			s.cleanup(ctx, uploaded)
			return result, fmt.Errorf("%s variant exceeds %dMB limit", variant.Name, variant.MaxBytes/1024/1024)
		}

		objectName := buildObjectName(variant, payload.ReferenceID)
		metadata := cloneMetadata(commonMetadata)
		if metadata == nil {
			metadata = map[string]string{}
		}
		metadata["variant"] = variant.Name
		if payload.ReferenceID != "" {
			metadata["reference-id"] = payload.ReferenceID
		}

		if err := s.storage.Upload(ctx, objectName, bytes.NewReader(encoded), int64(len(encoded)), variant.ContentType, metadata); err != nil {
			s.cleanup(ctx, uploaded)
			return result, fmt.Errorf("failed to upload %s variant: %w", variant.Name, err)
		}

		stored := StoredObject{
			Variant:     variant.Name,
			ObjectName:  objectName,
			ContentType: variant.ContentType,
			Size:        int64(len(encoded)),
			URL:         s.storage.ObjectURL(objectName),
		}

		if stored.URL == "" {
			if url, err := s.storage.PresignedURL(ctx, objectName, defaultPresignExpiry); err == nil {
				stored.TemporaryURL = url
			}
		}

		uploaded = append(uploaded, stored)
	}

	result.Bucket = s.storage.BucketName()
	result.Objects = uploaded
	return result, nil
}

func (s *Service) Delete(ctx context.Context, objectNames ...string) error {
	if s.storage == nil {
		return ErrStorageUnavailable
	}

	for _, name := range objectNames {
		if strings.TrimSpace(name) == "" {
			continue
		}
		if err := s.storage.Delete(ctx, name); err != nil {
			return fmt.Errorf("failed to delete object %s: %w", name, err)
		}
	}

	return nil
}

func (s *Service) cleanup(ctx context.Context, objects []StoredObject) {
	for _, obj := range objects {
		_ = s.storage.Delete(ctx, obj.ObjectName)
	}
}

type outputFormat string

type variantRule struct {
	Name        string
	Prefix      string
	Format      outputFormat
	ContentType string
	MaxBytes    int64
	MaxWidth    int
	MaxHeight   int
	Extension   string
}

type categoryRule struct {
	AllowedTypes map[string]struct{}
	MaxInput     int64
	Variants     []variantRule
}

var imageContentTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
	"image/gif":  {},
}

var categoryRules = map[Category]categoryRule{
	CategoryAssetImage: {
		AllowedTypes: imageContentTypes,
		MaxInput:     15 * 1024 * 1024,
		Variants: []variantRule{
			{
				Name:        VariantMain,
				Prefix:      "assets/img/",
				Format:      formatPNG,
				ContentType: "image/png",
				MaxBytes:    2 * 1024 * 1024,
				Extension:   ".png",
			},
			{
				Name:        VariantThumbnail,
				Prefix:      "assets/img/thumbnail/",
				Format:      formatWebP,
				ContentType: "image/webp",
				MaxBytes:    512 * 1024,
				MaxWidth:    320,
				MaxHeight:   320,
				Extension:   ".webp",
			},
		},
	},
	CategoryAssetQRCode: {
		AllowedTypes: imageContentTypes,
		MaxInput:     10 * 1024 * 1024,
		Variants: []variantRule{
			{
				Name:        VariantMain,
				Prefix:      "assets/qr/",
				Format:      formatPNG,
				ContentType: "image/png",
				MaxBytes:    1 * 1024 * 1024,
				Extension:   ".png",
			},
		},
	},
	CategoryEmployeeImage: {
		AllowedTypes: imageContentTypes,
		MaxInput:     10 * 1024 * 1024,
		Variants: []variantRule{
			{
				Name:        VariantMain,
				Prefix:      "employees/img/",
				Format:      formatJPEG,
				ContentType: "image/jpeg",
				MaxBytes:    2 * 1024 * 1024,
				MaxWidth:    1024,
				MaxHeight:   1024,
				Extension:   ".jpg",
			},
			{
				Name:        VariantThumbnail,
				Prefix:      "employees/img/thumbnail/",
				Format:      formatWebP,
				ContentType: "image/webp",
				MaxBytes:    256 * 1024,
				MaxWidth:    256,
				MaxHeight:   256,
				Extension:   ".webp",
			},
		},
	},
}

func prepareInput(rule categoryRule, payload UploadPayload) ([]byte, string, error) {
	if payload.File == nil {
		return nil, "", ErrFileRequired
	}
	if payload.Size <= 0 {
		return nil, "", ErrFileEmpty
	}

	if _, err := payload.File.Seek(0, io.SeekStart); err != nil {
		return nil, "", fmt.Errorf("failed to reset file reader: %w", err)
	}

	limit := rule.MaxInput
	if limit == 0 && payload.Size > 0 {
		limit = payload.Size
	}
	if limit == 0 {
		limit = 10 * 1024 * 1024
	}

	data, err := io.ReadAll(io.LimitReader(payload.File, limit+1))
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %w", err)
	}

	if int64(len(data)) == 0 {
		return nil, "", ErrFileEmpty
	}

	if rule.MaxInput > 0 && int64(len(data)) > rule.MaxInput {
		return nil, "", fmt.Errorf("file size exceeds %dMB", rule.MaxInput/1024/1024)
	}

	contentType := strings.TrimSpace(payload.ContentType)
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	if _, ok := rule.AllowedTypes[contentType]; !ok {
		if ext := strings.ToLower(filepath.Ext(payload.Filename)); ext != "" {
			if guessed := mime.TypeByExtension(ext); guessed != "" {
				contentType = guessed
			}
		}
	}

	if _, ok := rule.AllowedTypes[contentType]; !ok {
		return nil, "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	return data, contentType, nil
}

func buildBaseMetadata(payload UploadPayload) map[string]string {
	meta := map[string]string{}

	if trimmed := strings.TrimSpace(payload.Filename); trimmed != "" {
		meta["original-name"] = trimmed
	}
	if trimmed := strings.TrimSpace(payload.ContentType); trimmed != "" {
		meta["source-content-type"] = trimmed
	}

	if len(meta) == 0 {
		return nil
	}
	return meta
}

func cloneMetadata(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	clone := make(map[string]string, len(src))
	for key, value := range src {
		clone[key] = value
	}
	return clone
}

func decodeImage(data []byte, contentType string) (image.Image, error) {
	reader := bytes.NewReader(data)
	if contentType == "image/webp" {
		converted, err := decodeWebPWithBinary(data)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(converted)
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}

func encodeImage(img image.Image, format outputFormat) ([]byte, error) {
	switch format {
	case formatPNG:
		var buffer bytes.Buffer
		if err := png.Encode(&buffer, img); err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	case formatJPEG:
		var buffer bytes.Buffer
		if err := jpeg.Encode(&buffer, img, &jpeg.Options{Quality: 95}); err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	case formatWebP:
		return encodeWebPWithBinary(img)
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}

func resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	if maxWidth <= 0 && maxHeight <= 0 {
		return img
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width == 0 || height == 0 {
		return img
	}

	targetWidth := width
	targetHeight := height

	if maxWidth > 0 && width > maxWidth {
		scale := float64(maxWidth) / float64(width)
		targetWidth = maxWidth
		targetHeight = int(math.Round(float64(height) * scale))
	}

	if maxHeight > 0 && targetHeight > maxHeight {
		scale := float64(maxHeight) / float64(targetHeight)
		targetHeight = maxHeight
		targetWidth = int(math.Round(float64(targetWidth) * scale))
	}

	if targetWidth == width && targetHeight == height {
		return img
	}

	dst := image.NewNRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	srcBounds := img.Bounds()

	for y := 0; y < targetHeight; y++ {
		srcY := srcBounds.Min.Y + int(float64(y)*float64(srcBounds.Dy())/float64(targetHeight))
		if srcY >= srcBounds.Max.Y {
			srcY = srcBounds.Max.Y - 1
		}

		for x := 0; x < targetWidth; x++ {
			srcX := srcBounds.Min.X + int(float64(x)*float64(srcBounds.Dx())/float64(targetWidth))
			if srcX >= srcBounds.Max.X {
				srcX = srcBounds.Max.X - 1
			}
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	return dst
}

func decodeWebPWithBinary(data []byte) ([]byte, error) {
	if _, err := exec.LookPath("dwebp"); err != nil {
		return nil, errors.New("dwebp binary not found for webp decoding")
	}

	var out bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "dwebp", "-", "-o", "-", "-quiet")
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to decode webp: %w", err)
	}

	return out.Bytes(), nil
}

func encodeWebPWithBinary(img image.Image) ([]byte, error) {
	if _, err := exec.LookPath("cwebp"); err != nil {
		return nil, errors.New("cwebp binary not found for thumbnail encoding")
	}

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return nil, fmt.Errorf("failed to prepare buffer for webp encoding: %w", err)
	}

	var out bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "cwebp", "-lossless", "-quiet", "-", "-o", "-")
	cmd.Stdin = bytes.NewReader(buffer.Bytes())
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail to webp: %w", err)
	}

	return out.Bytes(), nil
}

func buildObjectName(variant variantRule, referenceID string) string {
	base := variant.Prefix
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}

	if ref := sanitizeReference(referenceID); ref != "" {
		base += ref + "/"
	}

	return base + uuid.NewString() + variant.Extension
}

func sanitizeReference(reference string) string {
	if reference == "" {
		return ""
	}

	reference = strings.TrimSpace(strings.ToLower(reference))
	var builder strings.Builder
	for _, r := range reference {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			builder.WriteRune(r)
		case r == '-' || r == '_':
			builder.WriteRune(r)
		default:
			builder.WriteRune('-')
		}
		if builder.Len() >= 48 {
			break
		}
	}

	return strings.Trim(builder.String(), "-_")
}
