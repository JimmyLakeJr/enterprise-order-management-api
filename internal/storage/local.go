package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"enterprise-order-management-api/internal/pkg/apperror"
)

type LocalFileStorage struct {
	rootDir   string
	publicURL string
}

func NewLocalFileStorage(rootDir string, publicURL string) *LocalFileStorage {
	return &LocalFileStorage{
		rootDir:   strings.TrimSpace(rootDir),
		publicURL: strings.TrimRight(strings.TrimSpace(publicURL), "/"),
	}
}

func (s *LocalFileStorage) SaveImage(file *multipart.FileHeader, subdir string, maxBytes int64) (string, error) {
	return s.save(file, subdir, maxBytes, map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
		"image/gif":  ".gif",
		"image/avif": ".avif",
	})
}

func (s *LocalFileStorage) SaveVideo(file *multipart.FileHeader, subdir string, maxBytes int64) (string, error) {
	return s.save(file, subdir, maxBytes, map[string]string{
		"video/mp4":  ".mp4",
		"video/webm": ".webm",
	})
}

func (s *LocalFileStorage) DeleteManagedFile(fileURL string) error {
	relativePath, ok := s.relativeUploadPath(fileURL)
	if !ok {
		return nil
	}
	if relativePath == "" {
		return nil
	}
	if err := os.Remove(filepath.Join(s.rootDir, filepath.FromSlash(relativePath))); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *LocalFileStorage) save(file *multipart.FileHeader, subdir string, maxBytes int64, allowed map[string]string) (string, error) {
	if file == nil {
		return "", apperror.BadRequest("File is required")
	}
	if file.Size <= 0 {
		return "", apperror.BadRequest("File is empty")
	}
	if maxBytes > 0 && file.Size > maxBytes {
		return "", apperror.BadRequest(fmt.Sprintf("File exceeds maximum size of %d bytes", maxBytes))
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	header := make([]byte, 512)
	n, err := src.Read(header)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := http.DetectContentType(header[:n])
	extension, ok := allowed[contentType]
	if !ok {
		return "", apperror.BadRequest("File type is not supported")
	}

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	targetDir := filepath.Join(s.rootDir, filepath.FromSlash(strings.Trim(subdir, "/")))
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	targetPath := filepath.Join(targetDir, filename)

	dst, err := os.Create(targetPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	relativePath := filepath.ToSlash(filepath.Join(strings.Trim(subdir, "/"), filename))
	publicPath := "/uploads/" + strings.TrimLeft(relativePath, "/")
	if s.publicURL == "" {
		return publicPath, nil
	}
	return s.publicURL + publicPath, nil
}

func (s *LocalFileStorage) relativeUploadPath(fileURL string) (string, bool) {
	cleanURL := strings.TrimSpace(fileURL)
	if cleanURL == "" {
		return "", false
	}

	if strings.HasPrefix(cleanURL, "/uploads/") {
		return strings.TrimPrefix(cleanURL, "/uploads/"), true
	}

	prefix := s.publicURL + "/uploads/"
	if s.publicURL != "" && strings.HasPrefix(cleanURL, prefix) {
		return strings.TrimPrefix(cleanURL, prefix), true
	}

	return "", false
}
