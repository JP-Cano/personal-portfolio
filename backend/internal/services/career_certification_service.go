package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/JuanPabloCano/personal-portfolio/backend/internal/handlers/dto"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/models"
	"github.com/JuanPabloCano/personal-portfolio/backend/internal/repository"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/constants"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/logger"
	"github.com/JuanPabloCano/personal-portfolio/backend/pkg/utils"
	"github.com/google/uuid"
)

var (
	DefaultMaxWorkers = 3
)

// UploadResult represents the outcome of an individual file upload process, including its success status and any errors encountered.
type UploadResult struct {
	Certification *models.CareerCertification
	OriginalName  string
	Error         error
	Success       bool
}

// FileWithMetadata represents a file upload along with its associated metadata.
type FileWithMetadata struct {
	File     *multipart.FileHeader
	Metadata *dto.CertificationMetadata
}

// CareerCertificationService represents the service interface for managing career certifications.
//
// The interface provides methods to store, retrieve, and delete career certification records. It supports batch
// uploading of certification files with metadata and handles storage and CRUD operations.
type CareerCertificationService interface {
	StoreBatch(ctx context.Context, filesWithMetadata []FileWithMetadata, maxWorkers int, baseURL string) []UploadResult
	GetAll() ([]models.CareerCertification, error)
	GetByID(id uint) (*models.CareerCertification, error)
	Delete(id uint) error
}

// careerCertificationService provides methods for managing career certifications, including file handling and database operations.
type careerCertificationService struct {
	uploadDir string
	repo      repository.CareerCertificationRepository
}

// NewCareerCertificationService initializes and returns a new CareerCertificationService implementation.
// It creates the upload directory for career certifications if it does not exist and utilizes the provided repository.
func NewCareerCertificationService(repo repository.CareerCertificationRepository) CareerCertificationService {
	service := &careerCertificationService{
		uploadDir: constants.CareerCertificationsDir,
		repo:      repo,
	}

	if err := os.MkdirAll(service.uploadDir, 0755); err != nil {
		logger.Fatal("Failed to create upload directory: %v", err)
		panic(err)
	}

	return service
}

// StoreBatch uploads multiple files concurrently, utilizing a worker pool. Returns a slice of UploadResult for each file.
func (c *careerCertificationService) StoreBatch(ctx context.Context, filesWithMetadata []FileWithMetadata, maxWorkers int, baseURL string) []UploadResult {
	if len(filesWithMetadata) == 0 {
		return []UploadResult{}
	}

	if maxWorkers <= 0 {
		maxWorkers = DefaultMaxWorkers
	}

	if maxWorkers > len(filesWithMetadata) {
		maxWorkers = len(filesWithMetadata)
	}

	logger.Info("Starting batch upload of %d files with %d workers", len(filesWithMetadata), maxWorkers)

	jobs := make(chan FileWithMetadata, len(filesWithMetadata))
	results := make(chan UploadResult, len(filesWithMetadata))

	var wg sync.WaitGroup

	// Start workers' goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go c.uploadWorker(ctx, i+1, jobs, results, &wg, baseURL)
	}

	for _, fwm := range filesWithMetadata {
		jobs <- fwm
	}
	close(jobs)

	// Wait for all workers to complete and close the results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	uploadResults := make([]UploadResult, 0, len(filesWithMetadata))
	for result := range results {
		uploadResults = append(uploadResults, result)
	}

	successCount := 0
	for _, r := range uploadResults {
		if r.Success {
			successCount++
		}
	}

	logger.Info("Finished batch upload: %d/%d successful", successCount, len(filesWithMetadata))
	return uploadResults
}

// uploadWorker handles the upload and processing of career certification files in a concurrent worker routine.
// It reads file and metadata from the jobs channel, processes the file, and stores the result in the database.
// Results are sent to the results channel, capturing any errors or success status.
// The method stops processing when jobs channel is closed or context cancellation occurs.
// It must be called in a goroutine and signals completion by calling Done on the provided WaitGroup.
func (c *careerCertificationService) uploadWorker(ctx context.Context, workerID int, jobs <-chan FileWithMetadata, results chan<- UploadResult, wg *sync.WaitGroup, baseURL string) {
	defer wg.Done()

	for fwm := range jobs {
		file := fwm.File
		metadata := fwm.Metadata

		select {
		case <-ctx.Done():
			results <- UploadResult{
				OriginalName: file.Filename,
				Error:        ctx.Err(),
				Success:      false,
			}
			continue
		default:
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))

		// Validate file extension
		if !utils.AllowedExtensions[ext] {
			logger.Warn("Worker %d: invalid file type %s for %s", workerID, ext, file.Filename)
			results <- UploadResult{
				OriginalName: file.Filename,
				Error:        fmt.Errorf("invalid file type: only JPG, JPEG, PNG, and WEBP images are allowed"),
				Success:      false,
			}
			continue
		}

		filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), uuid.New().String(), ext)
		filePath := filepath.Join(c.uploadDir, filename)

		if err := c.saveFile(file, filePath); err != nil {
			logger.Error("Worker %d: failed to save %s: %v", workerID, file.Filename, err)
			results <- UploadResult{
				OriginalName: file.Filename,
				Error:        err,
				Success:      false,
			}
			continue
		}

		certification := c.buildCareerCertification(metadata, file, baseURL, filename)

		c.setOptionalFields(metadata, certification)

		if err := c.repo.Create(certification); err != nil {
			logger.Error("Worker %d: failed to save to database: %v", workerID, err)
			os.Remove(filePath)
			results <- UploadResult{
				OriginalName: file.Filename,
				Error:        err,
				Success:      false,
			}
			continue
		}

		logger.Debug("Worker %d: successfully saved %s", workerID, filename)
		results <- UploadResult{
			Certification: certification,
			OriginalName:  file.Filename,
			Success:       true,
		}
	}
}

// setOptionalFields sets optional fields in the CareerCertification model based on provided CertificationMetadata.
// It applies values for ExpiryDate, CredentialID, CredentialURL, and Description if they are present in the metadata.
func (c *careerCertificationService) setOptionalFields(metadata *dto.CertificationMetadata, certification *models.CareerCertification) {
	if metadata.ExpiryDate != "" {
		if expiry, err := utils.ParseDateToPtr(metadata.ExpiryDate); err == nil {
			certification.ExpiryDate = expiry
		}
	}

	if metadata.CredentialID != "" {
		certification.CredentialID = &metadata.CredentialID
	}

	if metadata.CredentialURL != "" {
		certification.CredentialURL = &metadata.CredentialURL
	}

	if metadata.Description != "" {
		certification.Description = metadata.Description
	}
}

// buildCareerCertification constructs a CareerCertification model combining metadata, file details, and generated fields.
func (c *careerCertificationService) buildCareerCertification(metadata *dto.CertificationMetadata, file *multipart.FileHeader, baseURL string, filename string) *models.CareerCertification {
	return &models.CareerCertification{
		Title:        getOrDefault(metadata.Title, file.Filename),
		Issuer:       getOrDefault(metadata.Issuer, "N/A"),
		IssueDate:    parseIssueDate(metadata.IssueDate),
		FileURL:      fmt.Sprintf("%s/%s", baseURL, filename),
		FileName:     filename,
		OriginalName: file.Filename,
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
	}
}

// saveFile saves an uploaded file to the specified destination path. It opens the source file, creates the destination file,
// copies the file data to the destination, and handles cleanup on failure. Returns an error if any operation fails.
func (c *careerCertificationService) saveFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(path)
		return fmt.Errorf("failed to copy uploaded file: %w", err)
	}

	return nil
}

// GetAll retrieves all CareerCertification records from the repository and returns them along with any potential error.
func (c *careerCertificationService) GetAll() ([]models.CareerCertification, error) {
	return c.repo.FindAll()
}

// GetByID retrieves a CareerCertification by its unique ID from the repository and returns it.
func (c *careerCertificationService) GetByID(id uint) (*models.CareerCertification, error) {
	return c.repo.FindByID(id)
}

// Delete removes a career certification by its ID, deletes the corresponding file from disk, and returns an error if any occur.
func (c *careerCertificationService) Delete(id uint) error {
	cert, err := c.repo.FindByID(id)
	if err != nil {
		return err
	}

	if err := c.repo.Delete(id); err != nil {
		return err
	}

	filePath := filepath.Join(c.uploadDir, cert.FileName)
	if err := os.Remove(filePath); err != nil {
		logger.Warn("Failed to delete physical file %s: %v", filePath, err)
	}

	return nil
}

// getOrDefault returns the `value` if it is not an empty string, otherwise it returns the `defaultValue`.
func getOrDefault(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

// parseIssueDate parses a date string and returns a time.Time object; defaults to current time if parsing fails or input is empty.
func parseIssueDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}

	if date, err := utils.ParseDate(dateStr); err == nil {
		return date
	}

	return time.Now()
}
