package material

import (
	"io"
	"net/http"

	"doan/cmd/http/rest"
	usecasematerial "doan/internal/usecases/material"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	uploadUseCase   usecasematerial.UploadMaterialUseCase
	listUseCase     usecasematerial.ListMaterialsUseCase
	getUseCase      usecasematerial.GetMaterialUseCase
	downloadUseCase usecasematerial.DownloadMaterialUseCase
	reviewUseCase   usecasematerial.ReviewMaterialUseCase
}

func NewMaterialControllerV1(
	uploadUseCase usecasematerial.UploadMaterialUseCase,
	listUseCase usecasematerial.ListMaterialsUseCase,
	getUseCase usecasematerial.GetMaterialUseCase,
	downloadUseCase usecasematerial.DownloadMaterialUseCase,
	reviewUseCase usecasematerial.ReviewMaterialUseCase,
) *ControllerV1 {
	return &ControllerV1{
		uploadUseCase:   uploadUseCase,
		listUseCase:     listUseCase,
		getUseCase:      getUseCase,
		downloadUseCase: downloadUseCase,
		reviewUseCase:   reviewUseCase,
	}
}

func (ctrl *ControllerV1) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "File upload is required", err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to open uploaded file", err)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to read uploaded file", err)
		return
	}

	output, err := ctrl.uploadUseCase.Execute(c.Request.Context(), usecasematerial.UploadMaterialInput{
		TeacherID:   c.PostForm("teacher_id"),
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		FileName:    fileHeader.Filename,
		FileType:    fileHeader.Header.Get("Content-Type"),
		Content:     content,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to upload material", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusAccepted, "Material uploaded and audited successfully", output)
}

func (ctrl *ControllerV1) List(c *gin.Context) {
	output, err := ctrl.listUseCase.Execute(c.Request.Context(), usecasematerial.ListMaterialsInput{
		TeacherID: c.Query("teacher_id"),
		Status:    c.Query("status"),
		Queue:     c.Query("queue"),
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to list materials", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Materials retrieved successfully", output)
}

func (ctrl *ControllerV1) ListFlagged(c *gin.Context) {
	output, err := ctrl.listUseCase.Execute(c.Request.Context(), usecasematerial.ListMaterialsInput{
		Status: "AI_REVIEWED",
		Queue:  "flagged",
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to list flagged materials", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Flagged materials retrieved successfully", output)
}

func (ctrl *ControllerV1) Get(c *gin.Context) {
	output, err := ctrl.getUseCase.Execute(c.Request.Context(), usecasematerial.GetMaterialInput{
		ID: c.Param("id"),
	})
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "Material not found", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Material retrieved successfully", output)
}

func (ctrl *ControllerV1) Download(c *gin.Context) {
	output, err := ctrl.downloadUseCase.Execute(c.Request.Context(), usecasematerial.DownloadMaterialInput{
		ID: c.Param("id"),
	})
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "Material file not found", err)
		return
	}

	if output.FileType != "" {
		c.Header("Content-Type", output.FileType)
	}
	c.FileAttachment(output.AbsolutePath, output.FileName)
}

func (ctrl *ControllerV1) Review(c *gin.Context) {
	var req ReviewMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.reviewUseCase.Execute(c.Request.Context(), usecasematerial.ReviewMaterialInput{
		MaterialID:          c.Param("id"),
		ComplianceOfficerID: req.ComplianceOfficerID,
		Approved:            req.Approved,
		RejectReason:        req.RejectReason,
		Notes:               req.Notes,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to review material", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Material reviewed successfully", output)
}
