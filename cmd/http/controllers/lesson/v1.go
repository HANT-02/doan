package lesson

import (
	"doan/cmd/http/rest"
	"doan/internal/usecases/lesson"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	listLessonsUseCase lesson.ListLessonsUseCase
	getLessonUseCase   lesson.GetLessonUseCase
}

func NewLessonControllerV1(
	listLessonsUseCase lesson.ListLessonsUseCase,
	getLessonUseCase lesson.GetLessonUseCase,
) *ControllerV1 {
	return &ControllerV1{
		listLessonsUseCase: listLessonsUseCase,
		getLessonUseCase:   getLessonUseCase,
	}
}

func (ctrl *ControllerV1) ListLessons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	classID := c.Query("class_id")
	teacherID := c.Query("teacher_id")
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")

	var dateFrom, dateTo *time.Time
	if dateFromStr != "" {
		t, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			dateFrom = &t
		}
	}
	if dateToStr != "" {
		t, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			// End of day
			end := t.Add(24*time.Hour - time.Second)
			dateTo = &end
		}
	}

	output, err := ctrl.listLessonsUseCase.Execute(c.Request.Context(), lesson.ListLessonsInput{
		ClassID:   classID,
		TeacherID: teacherID,
		DateFrom:  dateFrom,
		DateTo:    dateTo,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to list lessons", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Lessons retrieved successfully", output)
}

func (ctrl *ControllerV1) GetLesson(c *gin.Context) {
	id := c.Param("id")

	output, err := ctrl.getLessonUseCase.Execute(c.Request.Context(), lesson.GetLessonInput{ID: id})
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "Lesson not found", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Lesson retrieved successfully", output.Lesson)
}
