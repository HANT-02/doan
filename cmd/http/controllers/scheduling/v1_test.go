package scheduling

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	usecasescheduling "doan/internal/usecases/scheduling"

	"github.com/gin-gonic/gin"
)

func TestControllerV1_BenchmarkReturnsBenchmarkPayload(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	benchmarkUseCase := &benchmarkUseCaseStub{
		output: &usecasescheduling.BenchmarkOutput{
			GeneratedAt: time.Date(2026, 4, 22, 10, 0, 0, 0, time.UTC),
			Mode:        "ADMIN_BENCHMARK_EXECUTION",
			Solvers: []usecasescheduling.BenchmarkSolverResult{
				{
					Key:             "cp_sat",
					Label:           "CP-SAT",
					ExecutionStatus: "COMPLETED",
				},
			},
		},
	}

	controller := NewSchedulingControllerV1(
		previewUseCaseStub{},
		benchmarkUseCase,
		getPreviewUseCaseStub{},
		commitPreviewUseCaseStub{},
		substituteUseCaseStub{},
		makeupUseCaseStub{},
	)

	router := gin.New()
	router.POST("/benchmark", controller.Benchmark)

	body := bytes.NewBufferString(`{"date_from":"2026-04-13","date_to":"2026-04-20","class_ids":["class-1"]}`)
	request := httptest.NewRequest(http.MethodPost, "/benchmark", body)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	if !benchmarkUseCase.called {
		t.Fatalf("expected benchmark use case to be called")
	}
	if len(benchmarkUseCase.input.ClassIDs) != 1 || benchmarkUseCase.input.ClassIDs[0] != "class-1" {
		t.Fatalf("unexpected benchmark input: %+v", benchmarkUseCase.input)
	}

	var response baseResponseForTest
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Fatalf("expected success response")
	}
	if response.Data == nil {
		t.Fatalf("expected response data")
	}

	var payload usecasescheduling.BenchmarkOutput
	if err := json.Unmarshal(response.Data, &payload); err != nil {
		t.Fatalf("failed to decode response payload: %v", err)
	}

	if payload.Mode != "ADMIN_BENCHMARK_EXECUTION" {
		t.Fatalf("expected benchmark mode, got %s", payload.Mode)
	}
	if len(payload.Solvers) != 1 || payload.Solvers[0].Key != "cp_sat" {
		t.Fatalf("unexpected solver payload: %+v", payload.Solvers)
	}
}

func TestControllerV1_BenchmarkRejectsInvalidDateRange(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	benchmarkUseCase := &benchmarkUseCaseStub{}
	controller := NewSchedulingControllerV1(
		previewUseCaseStub{},
		benchmarkUseCase,
		getPreviewUseCaseStub{},
		commitPreviewUseCaseStub{},
		substituteUseCaseStub{},
		makeupUseCaseStub{},
	)

	router := gin.New()
	router.POST("/benchmark", controller.Benchmark)

	body := bytes.NewBufferString(`{"date_from":"2026-04-20","date_to":"2026-04-13"}`)
	request := httptest.NewRequest(http.MethodPost, "/benchmark", body)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d with body %s", recorder.Code, recorder.Body.String())
	}
	if benchmarkUseCase.called {
		t.Fatalf("did not expect benchmark use case to be called on invalid request")
	}
}

type baseResponseForTest struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}

type previewUseCaseStub struct{}

func (previewUseCaseStub) Execute(context.Context, usecasescheduling.PreviewInput) (*usecasescheduling.PreviewResult, error) {
	return nil, nil
}

type benchmarkUseCaseStub struct {
	called bool
	input  usecasescheduling.BenchmarkInput
	output *usecasescheduling.BenchmarkOutput
	err    error
}

func (s *benchmarkUseCaseStub) Execute(_ context.Context, input usecasescheduling.BenchmarkInput) (*usecasescheduling.BenchmarkOutput, error) {
	s.called = true
	s.input = input
	return s.output, s.err
}

type getPreviewUseCaseStub struct{}

func (getPreviewUseCaseStub) Execute(context.Context, usecasescheduling.GetPreviewInput) (*usecasescheduling.PreviewResult, error) {
	return nil, nil
}

func (getPreviewUseCaseStub) GetLatest(context.Context) (*usecasescheduling.PreviewResult, error) {
	return nil, nil
}

type commitPreviewUseCaseStub struct{}

func (commitPreviewUseCaseStub) Execute(context.Context, usecasescheduling.CommitPreviewInput) (*usecasescheduling.CommitPreviewOutput, error) {
	return nil, nil
}

type substituteUseCaseStub struct{}

func (substituteUseCaseStub) SuggestSubstituteTeachers(context.Context, usecasescheduling.Actor, string) ([]usecasescheduling.SubstituteSuggestion, error) {
	return nil, nil
}

func (substituteUseCaseStub) AssignSubstitute(context.Context, usecasescheduling.Actor, string, string, string) error {
	return nil
}

type makeupUseCaseStub struct{}

func (makeupUseCaseStub) FindMakeupSpots(context.Context, usecasescheduling.FindMakeupSpotsInput) (*usecasescheduling.FindMakeupSpotsOutput, error) {
	return nil, nil
}
