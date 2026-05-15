package scheduling

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type BenchmarkArtifactManifest struct {
	GeneratedAt  time.Time                        `json:"thoi_diem_tao"`
	OutputDir    string                           `json:"thu_muc_dau_ra"`
	ScenarioList []BenchmarkArtifactScenarioEntry `json:"danh_sach_kich_ban"`
}

type BenchmarkArtifactScenarioEntry struct {
	ScenarioName string                         `json:"ten_kich_ban"`
	ScenarioDir  string                         `json:"thu_muc_kich_ban"`
	Solvers      []BenchmarkArtifactSolverEntry `json:"danh_sach_thuat_toan"`
}

type BenchmarkArtifactSolverEntry struct {
	SolverKey string   `json:"khoa_thuat_toan"`
	SolverDir string   `json:"thu_muc_thuat_toan"`
	RunDirs   []string `json:"danh_sach_thu_muc_lan_chay"`
}

type benchmarkScenarioSnapshot struct {
	Name         string      `json:"ten_kich_ban"`
	Description  string      `json:"mo_ta"`
	Iterations   int         `json:"so_lan_chay"`
	ClassCount   int         `json:"so_lop"`
	TeacherCount int         `json:"so_giao_vien"`
	RoomCount    int         `json:"so_phong"`
	ShiftCount   int         `json:"so_ca"`
	SessionCount int         `json:"tong_so_buoi_yeu_cau"`
	Input        SolverInput `json:"du_lieu_dau_vao"`
}

type benchmarkStudyArtifactView struct {
	GeneratedAt     time.Time                       `json:"thoi_diem_tao"`
	ScenarioReports []benchmarkScenarioArtifactView `json:"bao_cao_kich_ban"`
	Recommendation  benchmarkRecommendationView     `json:"khuyen_nghi"`
}

type benchmarkScenarioArtifactView struct {
	Name         string                        `json:"ten_kich_ban"`
	Description  string                        `json:"mo_ta"`
	Iterations   int                           `json:"so_lan_chay"`
	ClassCount   int                           `json:"so_lop"`
	TeacherCount int                           `json:"so_giao_vien"`
	RoomCount    int                           `json:"so_phong"`
	ShiftCount   int                           `json:"so_ca"`
	SessionCount int                           `json:"tong_so_buoi_yeu_cau"`
	Solvers      []benchmarkSolverArtifactView `json:"ket_qua_theo_thuat_toan"`
}

type benchmarkSolverArtifactView struct {
	Key                      string                       `json:"khoa_thuat_toan"`
	Label                    string                       `json:"ten_thuat_toan"`
	Runs                     int                          `json:"so_lan_chay"`
	AvgFeasibilityRate       float64                      `json:"ty_le_xep_thanh_cong_trung_binh"`
	AvgHardViolationCount    float64                      `json:"so_vi_pham_cung_trung_binh"`
	AvgSoftScore             float64                      `json:"diem_mem_trung_binh"`
	AvgRuntimeMs             float64                      `json:"thoi_gian_chay_trung_binh_mili_giay"`
	StdDevHardViolationCount float64                      `json:"do_lech_chuan_vi_pham_cung"`
	StdDevSoftScore          float64                      `json:"do_lech_chuan_diem_mem"`
	StdDevRuntimeMs          float64                      `json:"do_lech_chuan_thoi_gian_chay_mili_giay"`
	MinRuntimeMs             int64                        `json:"thoi_gian_chay_nho_nhat_mili_giay"`
	MaxRuntimeMs             int64                        `json:"thoi_gian_chay_lon_nhat_mili_giay"`
	StableAcrossRuns         bool                         `json:"on_dinh_qua_nhieu_lan_chay"`
	RepresentativeStatus     string                       `json:"trang_thai_dai_dien"`
	RepresentativeSummary    benchmarkSummaryArtifactView `json:"tom_tat_dai_dien"`
	RunRecords               []benchmarkRunArtifactView   `json:"nhat_ky_tung_lan_chay,omitempty"`
}

type benchmarkRunArtifactView struct {
	RunIndex           int              `json:"chi_so_lan_chay"`
	StartedAt          time.Time        `json:"thoi_diem_bat_dau"`
	FinishedAt         time.Time        `json:"thoi_diem_ket_thuc"`
	RuntimeMs          int64            `json:"thoi_gian_chay_mili_giay"`
	Status             string           `json:"trang_thai"`
	ScheduledLessons   int              `json:"so_buoi_xep_duoc"`
	UnscheduledLessons int              `json:"so_buoi_chua_xep_duoc"`
	ConflictCount      int              `json:"so_xung_dot"`
	SoftScore          int              `json:"diem_mem"`
	Telemetry          *SolverTelemetry `json:"ghi_do_dac"`
}

type benchmarkSummaryArtifactView struct {
	RequestedClasses           int     `json:"so_lop_yeu_cau"`
	RequestedSessions          int     `json:"so_buoi_yeu_cau"`
	ScheduledLessons           int     `json:"so_buoi_xep_duoc"`
	UnscheduledLessons         int     `json:"so_buoi_chua_xep_duoc"`
	ConflictCount              int     `json:"so_xung_dot"`
	SoftScore                  int     `json:"diem_mem"`
	ScheduleChangeCount        int     `json:"so_lan_doi_lich"`
	TeacherChangeCount         int     `json:"so_lan_doi_giao_vien"`
	RoomChangeCount            int     `json:"so_lan_doi_phong"`
	AverageCapacityUtilization float64 `json:"ty_le_su_dung_suc_chua_trung_binh"`
}

type benchmarkRecommendationView struct {
	SelectedSolverKey   string   `json:"khoa_thuat_toan_duoc_chon"`
	SelectedSolverLabel string   `json:"ten_thuat_toan_duoc_chon"`
	Rationale           []string `json:"ly_do_lua_chon"`
}

func WriteBenchmarkArtifacts(baseDir string, study *BenchmarkStudy, scenarios []BenchmarkScenario) (string, error) {
	if study == nil {
		return "", fmt.Errorf("benchmark study is nil")
	}
	if baseDir == "" {
		baseDir = filepath.Join("artifacts", "do_dac_xep_lich")
	}

	rootDir := filepath.Join(baseDir, time.Now().Format("20060102_150405"))
	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		return "", fmt.Errorf("create artifact root dir: %w", err)
	}

	if err := writeJSONFile(filepath.Join(rootDir, "tong_hop.json"), toBenchmarkStudyArtifactView(study)); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(rootDir, "tong_hop.md"), []byte(RenderBenchmarkStudyMarkdown(study)), 0o644); err != nil {
		return "", fmt.Errorf("write markdown summary: %w", err)
	}
	if err := writeScenarioSummaryCSV(filepath.Join(rootDir, "bang_tong_hop_theo_thuat_toan.csv"), study); err != nil {
		return "", err
	}
	if err := writeRawRunsCSV(filepath.Join(rootDir, "so_lieu_tho_tung_lan_chay.csv"), study); err != nil {
		return "", err
	}

	scenariosByName := make(map[string]BenchmarkScenario, len(scenarios))
	for _, scenario := range scenarios {
		scenariosByName[scenario.Name] = scenario
	}

	manifest := BenchmarkArtifactManifest{
		GeneratedAt:  time.Now(),
		OutputDir:    rootDir,
		ScenarioList: make([]BenchmarkArtifactScenarioEntry, 0, len(study.ScenarioReports)),
	}

	for _, report := range study.ScenarioReports {
		scenarioDirName := "kich_ban_" + slugifyForPath(report.Name)
		scenarioDir := filepath.Join(rootDir, scenarioDirName)
		if err := os.MkdirAll(scenarioDir, 0o755); err != nil {
			return "", fmt.Errorf("create scenario dir: %w", err)
		}

		if scenario, ok := scenariosByName[report.Name]; ok {
			snapshot := benchmarkScenarioSnapshot{
				Name:         report.Name,
				Description:  report.Description,
				Iterations:   report.Iterations,
				ClassCount:   report.ClassCount,
				TeacherCount: report.TeacherCount,
				RoomCount:    report.RoomCount,
				ShiftCount:   report.ShiftCount,
				SessionCount: report.SessionCount,
				Input:        scenario.Input,
			}
			if err := writeJSONFile(filepath.Join(scenarioDir, "du_lieu_dau_vao.json"), snapshot); err != nil {
				return "", err
			}
		}

		if err := writeJSONFile(filepath.Join(scenarioDir, "tong_quan_kich_ban.json"), toBenchmarkScenarioArtifactView(report)); err != nil {
			return "", err
		}

		manifestScenario := BenchmarkArtifactScenarioEntry{
			ScenarioName: report.Name,
			ScenarioDir:  scenarioDir,
			Solvers:      make([]BenchmarkArtifactSolverEntry, 0, len(report.Solvers)),
		}

		for _, solver := range report.Solvers {
			solverDirName := "thuat_toan_" + slugifyForPath(solver.Key)
			solverDir := filepath.Join(scenarioDir, solverDirName)
			if err := os.MkdirAll(solverDir, 0o755); err != nil {
				return "", fmt.Errorf("create solver dir: %w", err)
			}

			if err := writeJSONFile(filepath.Join(solverDir, "tong_quan_thuat_toan.json"), toBenchmarkSolverArtifactView(solver, false)); err != nil {
				return "", err
			}

			manifestSolver := BenchmarkArtifactSolverEntry{
				SolverKey: solver.Key,
				SolverDir: solverDir,
				RunDirs:   make([]string, 0, len(solver.RunRecords)),
			}

			for _, run := range solver.RunRecords {
				runDir := filepath.Join(solverDir, fmt.Sprintf("lan_chay_%03d", run.RunIndex))
				if err := os.MkdirAll(runDir, 0o755); err != nil {
					return "", fmt.Errorf("create run dir: %w", err)
				}

				if err := writeJSONFile(filepath.Join(runDir, "so_lieu_tho.json"), toBenchmarkRunArtifactView(run)); err != nil {
					return "", err
				}
				if err := writeJSONFile(filepath.Join(runDir, "thong_tin_telemetry.json"), run.Telemetry); err != nil {
					return "", err
				}

				manifestSolver.RunDirs = append(manifestSolver.RunDirs, runDir)
			}

			manifestScenario.Solvers = append(manifestScenario.Solvers, manifestSolver)
		}

		manifest.ScenarioList = append(manifest.ScenarioList, manifestScenario)
	}

	if err := writeJSONFile(filepath.Join(rootDir, "bang_ke_tap_tin.json"), manifest); err != nil {
		return "", err
	}

	return rootDir, nil
}

func writeJSONFile(path string, value interface{}) error {
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json for %s: %w", path, err)
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		return fmt.Errorf("write json file %s: %w", path, err)
	}
	return nil
}

var nonPathSafePattern = regexp.MustCompile(`[^a-z0-9_]+`)

func slugifyForPath(input string) string {
	normalized := strings.ToLower(strings.TrimSpace(input))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")
	normalized = nonPathSafePattern.ReplaceAllString(normalized, "_")
	normalized = strings.Trim(normalized, "_")
	if normalized == "" {
		return "khong_xac_dinh"
	}
	return normalized
}

func toBenchmarkStudyArtifactView(study *BenchmarkStudy) benchmarkStudyArtifactView {
	view := benchmarkStudyArtifactView{
		GeneratedAt:     study.GeneratedAt,
		ScenarioReports: make([]benchmarkScenarioArtifactView, 0, len(study.ScenarioReports)),
		Recommendation: benchmarkRecommendationView{
			SelectedSolverKey:   study.Recommendation.SelectedSolverKey,
			SelectedSolverLabel: study.Recommendation.SelectedSolverLabel,
			Rationale:           append([]string(nil), study.Recommendation.Rationale...),
		},
	}

	for _, scenario := range study.ScenarioReports {
		view.ScenarioReports = append(view.ScenarioReports, toBenchmarkScenarioArtifactView(scenario))
	}
	return view
}

func toBenchmarkScenarioArtifactView(report BenchmarkScenarioReport) benchmarkScenarioArtifactView {
	view := benchmarkScenarioArtifactView{
		Name:         report.Name,
		Description:  report.Description,
		Iterations:   report.Iterations,
		ClassCount:   report.ClassCount,
		TeacherCount: report.TeacherCount,
		RoomCount:    report.RoomCount,
		ShiftCount:   report.ShiftCount,
		SessionCount: report.SessionCount,
		Solvers:      make([]benchmarkSolverArtifactView, 0, len(report.Solvers)),
	}
	for _, solver := range report.Solvers {
		view.Solvers = append(view.Solvers, toBenchmarkSolverArtifactView(solver, true))
	}
	return view
}

func toBenchmarkSolverArtifactView(report BenchmarkScenarioSolverReport, includeRuns bool) benchmarkSolverArtifactView {
	view := benchmarkSolverArtifactView{
		Key:                      report.Key,
		Label:                    report.Label,
		Runs:                     report.Runs,
		AvgFeasibilityRate:       report.AvgFeasibilityRate,
		AvgHardViolationCount:    report.AvgHardViolationCount,
		AvgSoftScore:             report.AvgSoftScore,
		AvgRuntimeMs:             report.AvgRuntimeMs,
		StdDevHardViolationCount: report.StdDevHardViolationCount,
		StdDevSoftScore:          report.StdDevSoftScore,
		StdDevRuntimeMs:          report.StdDevRuntimeMs,
		MinRuntimeMs:             report.MinRuntimeMs,
		MaxRuntimeMs:             report.MaxRuntimeMs,
		StableAcrossRuns:         report.StableAcrossRuns,
		RepresentativeStatus:     toVietnameseStatus(report.RepresentativeStatus),
		RepresentativeSummary:    toBenchmarkSummaryArtifactView(report.RepresentativeSummary),
	}
	if includeRuns {
		view.RunRecords = make([]benchmarkRunArtifactView, 0, len(report.RunRecords))
		for _, run := range report.RunRecords {
			view.RunRecords = append(view.RunRecords, toBenchmarkRunArtifactView(run))
		}
	}
	return view
}

func toBenchmarkRunArtifactView(run BenchmarkRunRecord) benchmarkRunArtifactView {
	return benchmarkRunArtifactView{
		RunIndex:           run.RunIndex,
		StartedAt:          run.StartedAt,
		FinishedAt:         run.FinishedAt,
		RuntimeMs:          run.RuntimeMs,
		Status:             toVietnameseStatus(run.Status),
		ScheduledLessons:   run.ScheduledLessons,
		UnscheduledLessons: run.UnscheduledLessons,
		ConflictCount:      run.ConflictCount,
		SoftScore:          run.SoftScore,
		Telemetry:          run.Telemetry,
	}
}

func toVietnameseStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "COMPLETED":
		return "HOAN_THANH"
	case "PARTIAL":
		return "MOT_PHAN"
	case "FAILED":
		return "THAT_BAI"
	default:
		if strings.TrimSpace(status) == "" {
			return "KHONG_XAC_DINH"
		}
		return status
	}
}

func writeScenarioSummaryCSV(path string, study *BenchmarkStudy) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create scenario summary csv %s: %w", path, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"ten_kich_ban",
		"mo_ta_kich_ban",
		"so_lop",
		"so_giao_vien",
		"so_phong",
		"so_ca",
		"tong_so_buoi_yeu_cau",
		"so_lan_chay",
		"khoa_thuat_toan",
		"ten_thuat_toan",
		"ty_le_xep_thanh_cong_trung_binh",
		"so_vi_pham_cung_trung_binh",
		"diem_mem_trung_binh",
		"thoi_gian_chay_trung_binh_mili_giay",
		"do_lech_chuan_thoi_gian_mili_giay",
		"thoi_gian_nho_nhat_mili_giay",
		"thoi_gian_lon_nhat_mili_giay",
		"trang_thai_dai_dien",
		"on_dinh_qua_nhieu_lan_chay",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("write csv header for %s: %w", path, err)
	}

	for _, scenario := range study.ScenarioReports {
		for _, solver := range scenario.Solvers {
			record := []string{
				scenario.Name,
				scenario.Description,
				strconv.Itoa(scenario.ClassCount),
				strconv.Itoa(scenario.TeacherCount),
				strconv.Itoa(scenario.RoomCount),
				strconv.Itoa(scenario.ShiftCount),
				strconv.Itoa(scenario.SessionCount),
				strconv.Itoa(scenario.Iterations),
				solver.Key,
				solver.Label,
				formatFloat(solver.AvgFeasibilityRate),
				formatFloat(solver.AvgHardViolationCount),
				formatFloat(solver.AvgSoftScore),
				formatFloat(solver.AvgRuntimeMs),
				formatFloat(solver.StdDevRuntimeMs),
				strconv.FormatInt(solver.MinRuntimeMs, 10),
				strconv.FormatInt(solver.MaxRuntimeMs, 10),
				toVietnameseStatus(solver.RepresentativeStatus),
				boolToVietnamese(solver.StableAcrossRuns),
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("write scenario summary row for %s: %w", path, err)
			}
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush scenario summary csv %s: %w", path, err)
	}

	return nil
}

func writeRawRunsCSV(path string, study *BenchmarkStudy) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create raw runs csv %s: %w", path, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"ten_kich_ban",
		"khoa_thuat_toan",
		"ten_thuat_toan",
		"chi_so_lan_chay",
		"thoi_diem_bat_dau",
		"thoi_diem_ket_thuc",
		"thoi_gian_chay_mili_giay",
		"trang_thai",
		"so_buoi_xep_duoc",
		"so_buoi_chua_xep_duoc",
		"so_xung_dot",
		"diem_mem",
		"so_phuong_an_da_danh_gia",
		"so_phuong_an_bi_loai_do_xung_dot",
		"so_nut_da_duyet",
		"so_nhanh_bi_cat",
		"toc_do_danh_gia_phuong_an_moi_giay",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("write csv header for %s: %w", path, err)
	}

	for _, scenario := range study.ScenarioReports {
		for _, solver := range scenario.Solvers {
			for _, run := range solver.RunRecords {
				var (
					candidateEvaluated      int
					candidateRejected       int
					nodesVisited            int
					prunedBranchCount       int
					calculateCountPerSecond float64
				)

				if run.Telemetry != nil {
					candidateEvaluated = run.Telemetry.CandidateEvaluatedCount
					candidateRejected = run.Telemetry.CandidateRejectedConflictCount
					nodesVisited = run.Telemetry.NodesVisited
					prunedBranchCount = run.Telemetry.PrunedBranchCount
					if run.RuntimeMs > 0 {
						calculateCountPerSecond = float64(run.Telemetry.CandidateEvaluatedCount) / (float64(run.RuntimeMs) / 1000)
					}
				}

				record := []string{
					scenario.Name,
					solver.Key,
					solver.Label,
					strconv.Itoa(run.RunIndex),
					run.StartedAt.Format(time.RFC3339Nano),
					run.FinishedAt.Format(time.RFC3339Nano),
					strconv.FormatInt(run.RuntimeMs, 10),
					toVietnameseStatus(run.Status),
					strconv.Itoa(run.ScheduledLessons),
					strconv.Itoa(run.UnscheduledLessons),
					strconv.Itoa(run.ConflictCount),
					strconv.Itoa(run.SoftScore),
					strconv.Itoa(candidateEvaluated),
					strconv.Itoa(candidateRejected),
					strconv.Itoa(nodesVisited),
					strconv.Itoa(prunedBranchCount),
					formatFloat(calculateCountPerSecond),
				}
				if err := writer.Write(record); err != nil {
					return fmt.Errorf("write raw run row for %s: %w", path, err)
				}
			}
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush raw runs csv %s: %w", path, err)
	}

	return nil
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 6, 64)
}

func boolToVietnamese(value bool) string {
	if value {
		return "Co"
	}
	return "Khong"
}

func toBenchmarkSummaryArtifactView(summary SolverSummary) benchmarkSummaryArtifactView {
	return benchmarkSummaryArtifactView{
		RequestedClasses:           summary.RequestedClasses,
		RequestedSessions:          summary.RequestedSessions,
		ScheduledLessons:           summary.ScheduledLessons,
		UnscheduledLessons:         summary.UnscheduledLessons,
		ConflictCount:              summary.ConflictCount,
		SoftScore:                  summary.SoftScore,
		ScheduleChangeCount:        summary.ScheduleChangeCount,
		TeacherChangeCount:         summary.TeacherChangeCount,
		RoomChangeCount:            summary.RoomChangeCount,
		AverageCapacityUtilization: summary.AverageCapacityUtilization,
	}
}
