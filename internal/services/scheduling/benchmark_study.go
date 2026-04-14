package scheduling

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"doan/internal/entities"
)

type BenchmarkScenario struct {
	Name        string
	Description string
	Input       SolverInput
	Iterations  int
}

type BenchmarkStudy struct {
	GeneratedAt     time.Time
	ScenarioReports []BenchmarkScenarioReport
	Recommendation  BenchmarkRecommendation
}

type BenchmarkScenarioReport struct {
	Name         string
	Description  string
	Iterations   int
	ClassCount   int
	TeacherCount int
	RoomCount    int
	ShiftCount   int
	SessionCount int
	Solvers      []BenchmarkScenarioSolverReport
}

type BenchmarkScenarioSolverReport struct {
	Key                   string
	Label                 string
	Runs                  int
	AvgFeasibilityRate    float64
	AvgHardViolationCount float64
	AvgSoftScore          float64
	AvgRuntimeMs          float64
	MinRuntimeMs          int64
	MaxRuntimeMs          int64
	StableAcrossRuns      bool
	RepresentativeStatus  string
	RepresentativeSummary SolverSummary
}

type BenchmarkRecommendation struct {
	SelectedSolverKey   string
	SelectedSolverLabel string
	Rationale           []string
}

type benchmarkAggregate struct {
	runs                  int
	feasibilityTotal      float64
	hardViolationTotal    float64
	softScoreTotal        float64
	runtimeTotalMs        float64
	minRuntimeMs          int64
	maxRuntimeMs          int64
	signatures            map[string]struct{}
	representativeStatus  string
	representativeSummary SolverSummary
}

func BuildDefaultBenchmarkScenarios() []BenchmarkScenario {
	baseDate := time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC)

	return []BenchmarkScenario{
		buildSyntheticScenario(syntheticScenarioConfig{
			Name:             "small",
			Description:      "6 lop, 4 giao vien, 3 phong, 2 buoi moi lop. Muc tieu la do feasibility va runtime baseline.",
			DateFrom:         baseDate,
			Days:             14,
			ClassCount:       6,
			TeacherCount:     4,
			RoomCount:        3,
			SessionCount:     2,
			Iterations:       7,
			PreferredRoomMod: 2,
		}),
		buildSyntheticScenario(syntheticScenarioConfig{
			Name:             "medium",
			Description:      "10 lop, 5 giao vien, 4 phong, 3 buoi moi lop. Muc tieu la do chat luong nghiem khi rang buoc bat dau tang.",
			DateFrom:         baseDate,
			Days:             14,
			ClassCount:       10,
			TeacherCount:     5,
			RoomCount:        4,
			SessionCount:     3,
			Iterations:       7,
			PreferredRoomMod: 3,
		}),
		buildSyntheticScenario(syntheticScenarioConfig{
			Name:             "large",
			Description:      "16 lop, 7 giao vien, 5 phong, 4 buoi moi lop. Muc tieu la do scalability va do on dinh khi so bien tang.",
			DateFrom:         baseDate,
			Days:             14,
			ClassCount:       16,
			TeacherCount:     7,
			RoomCount:        5,
			SessionCount:     4,
			Iterations:       7,
			PreferredRoomMod: 4,
		}),
	}
}

func RunBenchmarkStudy(ctx context.Context, catalog SolverCatalog, scenarios []BenchmarkScenario) (*BenchmarkStudy, error) {
	descriptors := catalog.BenchmarkSolvers()
	reports := make([]BenchmarkScenarioReport, 0, len(scenarios))

	for _, scenario := range scenarios {
		solverReports := make([]BenchmarkScenarioSolverReport, 0, len(descriptors))
		for _, descriptor := range descriptors {
			solver, ok := catalog.GetSolver(descriptor.Key)
			if !ok {
				return nil, fmt.Errorf("solver %s not registered in catalog", descriptor.Key)
			}

			aggregate, err := runScenarioForSolver(ctx, solver, scenario)
			if err != nil {
				return nil, fmt.Errorf("run benchmark scenario %s with solver %s: %w", scenario.Name, descriptor.Key, err)
			}

			solverReports = append(solverReports, BenchmarkScenarioSolverReport{
				Key:                   descriptor.Key,
				Label:                 descriptor.Label,
				Runs:                  aggregate.runs,
				AvgFeasibilityRate:    aggregate.feasibilityTotal / float64(aggregate.runs),
				AvgHardViolationCount: aggregate.hardViolationTotal / float64(aggregate.runs),
				AvgSoftScore:          aggregate.softScoreTotal / float64(aggregate.runs),
				AvgRuntimeMs:          aggregate.runtimeTotalMs / float64(aggregate.runs),
				MinRuntimeMs:          aggregate.minRuntimeMs,
				MaxRuntimeMs:          aggregate.maxRuntimeMs,
				StableAcrossRuns:      len(aggregate.signatures) == 1,
				RepresentativeStatus:  aggregate.representativeStatus,
				RepresentativeSummary: aggregate.representativeSummary,
			})
		}

		sort.Slice(solverReports, func(i, j int) bool {
			return compareScenarioSolverReports(solverReports[i], solverReports[j])
		})

		reports = append(reports, BenchmarkScenarioReport{
			Name:         scenario.Name,
			Description:  scenario.Description,
			Iterations:   scenario.Iterations,
			ClassCount:   len(scenario.Input.Classes),
			TeacherCount: uniqueTeacherCount(scenario.Input.Classes),
			RoomCount:    len(scenario.Input.Rooms),
			ShiftCount:   len(scenario.Input.Shifts),
			SessionCount: totalRequestedSessions(scenario.Input.Classes),
			Solvers:      solverReports,
		})
	}

	return &BenchmarkStudy{
		GeneratedAt:     time.Now(),
		ScenarioReports: reports,
		Recommendation:  selectBenchmarkRecommendation(reports),
	}, nil
}

func RenderBenchmarkStudyMarkdown(study *BenchmarkStudy) string {
	var builder strings.Builder

	builder.WriteString("# Scheduling Benchmark Study\n\n")
	builder.WriteString(fmt.Sprintf("- Generated at: %s\n", study.GeneratedAt.Format(time.RFC3339)))
	builder.WriteString(fmt.Sprintf("- Recommended solver: `%s` (%s)\n\n", study.Recommendation.SelectedSolverKey, study.Recommendation.SelectedSolverLabel))

	for _, scenario := range study.ScenarioReports {
		builder.WriteString(fmt.Sprintf("## Scenario `%s`\n\n", scenario.Name))
		builder.WriteString(fmt.Sprintf("- Description: %s\n", scenario.Description))
		builder.WriteString(fmt.Sprintf("- Dataset: %d classes, %d teachers, %d rooms, %d shifts, %d requested sessions\n", scenario.ClassCount, scenario.TeacherCount, scenario.RoomCount, scenario.ShiftCount, scenario.SessionCount))
		builder.WriteString(fmt.Sprintf("- Iterations: %d\n\n", scenario.Iterations))
		builder.WriteString("| Solver | Avg feasibility | Avg hard violations | Avg soft score | Avg runtime (ms) | Runtime range (ms) | Stable | Status |\n")
		builder.WriteString("| --- | ---: | ---: | ---: | ---: | --- | --- | --- |\n")
		for _, solver := range scenario.Solvers {
			builder.WriteString(fmt.Sprintf(
				"| %s | %.3f | %.3f | %.3f | %.3f | %d-%d | %t | %s |\n",
				solver.Label,
				solver.AvgFeasibilityRate,
				solver.AvgHardViolationCount,
				solver.AvgSoftScore,
				solver.AvgRuntimeMs,
				solver.MinRuntimeMs,
				solver.MaxRuntimeMs,
				solver.StableAcrossRuns,
				solver.RepresentativeStatus,
			))
		}
		builder.WriteString("\n")
	}

	builder.WriteString("## Recommendation\n\n")
	builder.WriteString(fmt.Sprintf("- Selected solver: `%s` (%s)\n", study.Recommendation.SelectedSolverKey, study.Recommendation.SelectedSolverLabel))
	for _, rationale := range study.Recommendation.Rationale {
		builder.WriteString(fmt.Sprintf("- %s\n", rationale))
	}

	return builder.String()
}

func runScenarioForSolver(ctx context.Context, solver SchedulingSolver, scenario BenchmarkScenario) (*benchmarkAggregate, error) {
	aggregate := &benchmarkAggregate{
		minRuntimeMs: mathMaxInt64,
		maxRuntimeMs: 0,
		signatures:   make(map[string]struct{}),
	}

	for iteration := 0; iteration < scenario.Iterations; iteration++ {
		startedAt := time.Now()
		output, err := solver.Solve(ctx, scenario.Input)
		runtimeMs := time.Since(startedAt).Milliseconds()
		if err != nil {
			return nil, err
		}

		aggregate.runs++
		aggregate.feasibilityTotal += calculateScenarioFeasibilityRate(output.Summary)
		aggregate.hardViolationTotal += float64(output.Summary.ConflictCount)
		aggregate.softScoreTotal += float64(output.Summary.SoftScore)
		aggregate.runtimeTotalMs += float64(runtimeMs)
		if runtimeMs < aggregate.minRuntimeMs {
			aggregate.minRuntimeMs = runtimeMs
		}
		if runtimeMs > aggregate.maxRuntimeMs {
			aggregate.maxRuntimeMs = runtimeMs
		}
		if aggregate.runs == 1 {
			aggregate.representativeStatus = output.Status
			aggregate.representativeSummary = output.Summary
		}

		signature := fmt.Sprintf("%s|%d|%d|%d|%d", output.Status, output.Summary.ScheduledLessons, output.Summary.UnscheduledLessons, output.Summary.ConflictCount, output.Summary.SoftScore)
		aggregate.signatures[signature] = struct{}{}
	}

	return aggregate, nil
}

func calculateScenarioFeasibilityRate(summary SolverSummary) float64 {
	if summary.RequestedSessions == 0 {
		return 0
	}

	return float64(summary.ScheduledLessons) / float64(summary.RequestedSessions)
}

func compareScenarioSolverReports(left, right BenchmarkScenarioSolverReport) bool {
	if left.AvgFeasibilityRate != right.AvgFeasibilityRate {
		return left.AvgFeasibilityRate > right.AvgFeasibilityRate
	}
	if left.AvgHardViolationCount != right.AvgHardViolationCount {
		return left.AvgHardViolationCount < right.AvgHardViolationCount
	}
	if left.AvgSoftScore != right.AvgSoftScore {
		return left.AvgSoftScore > right.AvgSoftScore
	}
	if left.AvgRuntimeMs != right.AvgRuntimeMs {
		return left.AvgRuntimeMs < right.AvgRuntimeMs
	}
	return left.Key < right.Key
}

func selectBenchmarkRecommendation(reports []BenchmarkScenarioReport) BenchmarkRecommendation {
	type totalScore struct {
		key                    string
		label                  string
		feasibilityTotal       float64
		hardViolationTotal     float64
		softScoreTotal         float64
		runtimeTotal           float64
		stableScenarioCount    int
		completedScenarioCount int
	}

	indexed := make(map[string]*totalScore)
	for _, scenario := range reports {
		for _, solver := range scenario.Solvers {
			item, ok := indexed[solver.Key]
			if !ok {
				item = &totalScore{
					key:   solver.Key,
					label: solver.Label,
				}
				indexed[solver.Key] = item
			}

			item.feasibilityTotal += solver.AvgFeasibilityRate
			item.hardViolationTotal += solver.AvgHardViolationCount
			item.softScoreTotal += solver.AvgSoftScore
			item.runtimeTotal += solver.AvgRuntimeMs
			if solver.StableAcrossRuns {
				item.stableScenarioCount++
			}
			if solver.RepresentativeStatus == "COMPLETED" {
				item.completedScenarioCount++
			}
		}
	}

	items := make([]*totalScore, 0, len(indexed))
	for _, item := range indexed {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].completedScenarioCount != items[j].completedScenarioCount {
			return items[i].completedScenarioCount > items[j].completedScenarioCount
		}
		if items[i].feasibilityTotal != items[j].feasibilityTotal {
			return items[i].feasibilityTotal > items[j].feasibilityTotal
		}
		if items[i].hardViolationTotal != items[j].hardViolationTotal {
			return items[i].hardViolationTotal < items[j].hardViolationTotal
		}
		if items[i].softScoreTotal != items[j].softScoreTotal {
			return items[i].softScoreTotal > items[j].softScoreTotal
		}
		if items[i].stableScenarioCount != items[j].stableScenarioCount {
			return items[i].stableScenarioCount > items[j].stableScenarioCount
		}
		if items[i].runtimeTotal != items[j].runtimeTotal {
			return items[i].runtimeTotal < items[j].runtimeTotal
		}
		return items[i].key < items[j].key
	})

	if len(items) == 0 {
		return BenchmarkRecommendation{}
	}

	selected := items[0]
	return BenchmarkRecommendation{
		SelectedSolverKey:   selected.key,
		SelectedSolverLabel: selected.label,
		Rationale: []string{
			fmt.Sprintf("Uu tien solver co so scenario COMPLETED nhieu nhat: %d/%d.", selected.completedScenarioCount, len(reports)),
			fmt.Sprintf("Tong feasibility trung binh cao nhat hoac dong hang cao nhat: %.3f.", selected.feasibilityTotal),
			fmt.Sprintf("Tong hard violations thap nhat sau khi uu tien feasibility: %.3f.", selected.hardViolationTotal),
			fmt.Sprintf("Tong runtime trung binh duoc dung lam tie-break sau cung: %.3f ms.", selected.runtimeTotal),
		},
	}
}

func totalRequestedSessions(classes []entities.Class) int {
	total := 0
	for _, classEntity := range classes {
		total += classEntity.Course.SessionCount
	}
	return total
}

func uniqueTeacherCount(classes []entities.Class) int {
	teacherIDs := make(map[string]struct{})
	for _, classEntity := range classes {
		if classEntity.TeacherID == nil || *classEntity.TeacherID == "" {
			continue
		}
		teacherIDs[*classEntity.TeacherID] = struct{}{}
	}
	return len(teacherIDs)
}

type syntheticScenarioConfig struct {
	Name             string
	Description      string
	DateFrom         time.Time
	Days             int
	ClassCount       int
	TeacherCount     int
	RoomCount        int
	SessionCount     int
	Iterations       int
	PreferredRoomMod int
}

func buildSyntheticScenario(cfg syntheticScenarioConfig) BenchmarkScenario {
	teachers := make([]entities.Teacher, 0, cfg.TeacherCount)
	for teacherIndex := 0; teacherIndex < cfg.TeacherCount; teacherIndex++ {
		teachers = append(teachers, entities.Teacher{
			ID:       fmt.Sprintf("teacher-%02d", teacherIndex+1),
			Code:     fmt.Sprintf("GV-%02d", teacherIndex+1),
			FullName: fmt.Sprintf("Giao vien %02d", teacherIndex+1),
			Status:   "ACTIVE",
		})
	}

	rooms := make([]entities.Room, 0, cfg.RoomCount)
	for roomIndex := 0; roomIndex < cfg.RoomCount; roomIndex++ {
		rooms = append(rooms, entities.Room{
			ID:       fmt.Sprintf("room-%02d", roomIndex+1),
			Code:     fmt.Sprintf("P-%02d", roomIndex+1),
			Name:     fmt.Sprintf("Phong %02d", roomIndex+1),
			Capacity: 22 + roomIndex*6,
		})
	}

	shifts := []entities.Shift{
		{ID: "shift-1", Code: "S1", Name: "Ca sang", StartTime: "08:00", EndTime: "10:00", DurationMinutes: 120, SessionType: "MORNING", IsActive: true},
		{ID: "shift-2", Code: "S2", Name: "Ca chieu", StartTime: "13:30", EndTime: "15:30", DurationMinutes: 120, SessionType: "AFTERNOON", IsActive: true},
		{ID: "shift-3", Code: "S3", Name: "Ca toi", StartTime: "18:00", EndTime: "20:00", DurationMinutes: 120, SessionType: "EVENING", IsActive: true},
	}

	dayPatterns := [][]string{
		{"MONDAY", "WEDNESDAY"},
		{"TUESDAY", "THURSDAY"},
		{"WEDNESDAY", "FRIDAY"},
		{"MONDAY", "FRIDAY"},
	}

	classes := make([]entities.Class, 0, cfg.ClassCount)
	for classIndex := 0; classIndex < cfg.ClassCount; classIndex++ {
		teacher := teachers[classIndex%len(teachers)]
		courseID := fmt.Sprintf("course-%02d", classIndex+1)
		course := entities.Course{
			ID:                     courseID,
			Code:                   fmt.Sprintf("KH-%02d", classIndex+1),
			Name:                   fmt.Sprintf("Khoa hoc %02d", classIndex+1),
			SessionCount:           cfg.SessionCount,
			SessionDurationMinutes: 120,
			Status:                 "ACTIVE",
		}

		maxStudents := 18 + (classIndex % 4 * 4)
		pattern := dayPatterns[classIndex%len(dayPatterns)]
		primaryShift := shifts[classIndex%len(shifts)]
		secondaryShift := shifts[(classIndex+1)%len(shifts)]

		var preferredRoomID *string
		if cfg.PreferredRoomMod > 0 && classIndex%cfg.PreferredRoomMod == 0 {
			roomID := rooms[classIndex%len(rooms)].ID
			preferredRoomID = &roomID
		}

		classSchedules := []entities.ClassSchedule{
			{
				ID:        fmt.Sprintf("schedule-%02d-a", classIndex+1),
				ClassID:   fmt.Sprintf("class-%02d", classIndex+1),
				DayOfWeek: pattern[0],
				ShiftID:   primaryShift.ID,
				Shift:     primaryShift,
				RoomID:    preferredRoomID,
			},
			{
				ID:        fmt.Sprintf("schedule-%02d-b", classIndex+1),
				ClassID:   fmt.Sprintf("class-%02d", classIndex+1),
				DayOfWeek: pattern[1],
				ShiftID:   secondaryShift.ID,
				Shift:     secondaryShift,
				RoomID:    preferredRoomID,
			},
		}

		classes = append(classes, entities.Class{
			ID:             fmt.Sprintf("class-%02d", classIndex+1),
			Code:           fmt.Sprintf("L-%02d", classIndex+1),
			Name:           fmt.Sprintf("Lop %02d", classIndex+1),
			StartDate:      cfg.DateFrom,
			MaxStudents:    maxStudents,
			Status:         "OPEN",
			CourseID:       &courseID,
			Course:         course,
			TeacherID:      &teacher.ID,
			Teacher:        teacher,
			RoomID:         preferredRoomID,
			ClassSchedules: classSchedules,
		})
	}

	return BenchmarkScenario{
		Name:        cfg.Name,
		Description: cfg.Description,
		Iterations:  cfg.Iterations,
		Input: SolverInput{
			DateFrom: cfg.DateFrom,
			DateTo:   cfg.DateFrom.AddDate(0, 0, cfg.Days-1),
			Classes:  classes,
			Rooms:    rooms,
			Shifts:   shifts,
		},
	}
}

const mathMaxInt64 = int64(^uint64(0) >> 1)
