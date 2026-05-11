package predictive

import (
	"net/http"
	"strconv"

	"doan/cmd/http/rest"
	predictiveuc "doan/internal/usecases/predictive"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	listPredictionsUseCase  predictiveuc.ListStudentPredictionsUseCase
	getModelMetadataUseCase predictiveuc.GetModelMetadataUseCase
	trainFromDBUseCase      predictiveuc.TrainAtRiskFromDBUseCase
}

func NewPredictiveControllerV1(
	listPredictionsUseCase predictiveuc.ListStudentPredictionsUseCase,
	getModelMetadataUseCase predictiveuc.GetModelMetadataUseCase,
	trainFromDBUseCase predictiveuc.TrainAtRiskFromDBUseCase,
) *ControllerV1 {
	return &ControllerV1{
		listPredictionsUseCase:  listPredictionsUseCase,
		getModelMetadataUseCase: getModelMetadataUseCase,
		trainFromDBUseCase:      trainFromDBUseCase,
	}
}

func (ctrl *ControllerV1) ListAtRiskStudents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	refresh, _ := strconv.ParseBool(c.DefaultQuery("refresh", "false"))
	onlyAtRisk, _ := strconv.ParseBool(c.DefaultQuery("only_at_risk", "false"))

	output, err := ctrl.listPredictionsUseCase.Execute(c.Request.Context(), predictiveuc.ListStudentPredictionsInput{
		Search:     c.Query("search"),
		OnlyAtRisk: onlyAtRisk,
		Page:       page,
		Limit:      limit,
		Refresh:    refresh,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to list at-risk predictions", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "At-risk predictions retrieved successfully", output)
}

func (ctrl *ControllerV1) GetModelMetadata(c *gin.Context) {
	refresh, _ := strconv.ParseBool(c.DefaultQuery("refresh", "false"))
	output, err := ctrl.getModelMetadataUseCase.Execute(c.Request.Context(), refresh)
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to get predictive model metadata", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Predictive model metadata retrieved successfully", output)
}

func (ctrl *ControllerV1) TrainAtRiskFromDB(c *gin.Context) {
	output, err := ctrl.trainFromDBUseCase.Execute(c.Request.Context())
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to train predictive model from DB", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Predictive model trained from DB successfully", output)
}
