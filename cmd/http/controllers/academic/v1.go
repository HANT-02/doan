package academic

import (
	"net/http"

	"doan/cmd/http/rest"
	lessonrecord "doan/internal/usecases/lessonrecord"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	listMyAcademicRecordsUseCase lessonrecord.ListMyAcademicRecordsUseCase
}

func NewAcademicControllerV1(
	listMyAcademicRecordsUseCase lessonrecord.ListMyAcademicRecordsUseCase,
) *ControllerV1 {
	return &ControllerV1{listMyAcademicRecordsUseCase: listMyAcademicRecordsUseCase}
}

func (ctrl *ControllerV1) ListMyAcademicRecords(c *gin.Context) {
	userRole, _ := c.Get("user_role")
	userEmail, _ := c.Get("user_email")
	role, _ := userRole.(string)
	email, _ := userEmail.(string)

	output, err := ctrl.listMyAcademicRecordsUseCase.Execute(c.Request.Context(), lessonrecord.ListMyAcademicRecordsInput{
		Actor: lessonrecord.LessonActor{
			Role:  role,
			Email: email,
		},
		ClassID: c.Query("class_id"),
	})
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to list academic records", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Academic records retrieved successfully", output)
}
