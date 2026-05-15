package scheduling

import (
	"sort"
	"time"
)

type BenchmarkOptions struct {
	ScenarioName string        `json:"ten_kich_ban,omitempty"`
	RunIndex     int           `json:"chi_so_lan_chay,omitempty"`
	TotalRuns    int           `json:"tong_so_lan_chay,omitempty"`
	TimeBudget   time.Duration `json:"ngan_sach_thoi_gian,omitempty"`
	RandomSeed   int64         `json:"hat_giong_ngau_nhien,omitempty"`
	IsWarmupRun  bool          `json:"la_lan_lam_nong,omitempty"`
}

type DomainStats struct {
	VariableCount     int     `json:"so_bien"`
	EmptyDomainCount  int     `json:"so_mien_rong"`
	TotalDomainValues int     `json:"tong_so_gia_tri_mien"`
	MinDomainSize     int     `json:"kich_thuoc_mien_nho_nhat"`
	AvgDomainSize     float64 `json:"kich_thuoc_mien_trung_binh"`
	MedianDomainSize  float64 `json:"kich_thuoc_mien_trung_vi"`
	MaxDomainSize     int     `json:"kich_thuoc_mien_lon_nhat"`
}

type SolverTelemetry struct {
	SolverKey                      string        `json:"khoa_thuat_toan,omitempty"`
	SolverLabel                    string        `json:"nhan_thuat_toan,omitempty"`
	ScenarioName                   string        `json:"ten_kich_ban,omitempty"`
	RunIndex                       int           `json:"chi_so_lan_chay,omitempty"`
	TotalRuns                      int           `json:"tong_so_lan_chay,omitempty"`
	RandomSeed                     int64         `json:"hat_giong_ngau_nhien,omitempty"`
	IsWarmupRun                    bool          `json:"la_lan_lam_nong,omitempty"`
	StartedAt                      time.Time     `json:"thoi_diem_bat_dau,omitempty"`
	FinishedAt                     time.Time     `json:"thoi_diem_ket_thuc,omitempty"`
	MeasuredAt                     time.Time     `json:"thoi_diem_ghi_nhan,omitempty"`
	Runtime                        time.Duration `json:"thoi_gian_chay,omitempty"`
	InputClassCount                int           `json:"so_lop_dau_vao,omitempty"`
	RequestedSessionCount          int           `json:"so_buoi_yeu_cau,omitempty"`
	DomainStats                    DomainStats   `json:"thong_ke_mien"`
	CandidateEvaluatedCount        int           `json:"so_phuong_an_da_danh_gia,omitempty"`
	CandidateRejectedConflictCount int           `json:"so_phuong_an_bi_loai_do_xung_dot,omitempty"`
	SlotGroupEvaluatedCount        int           `json:"so_nhom_slot_da_xet,omitempty"`
	FirstPassAssignedCount         int           `json:"so_gan_thanh_cong_luot_dau,omitempty"`
	NodesVisited                   int           `json:"so_nut_da_duyet,omitempty"`
	PrunedBranchCount              int           `json:"so_nhanh_bi_cat,omitempty"`
	BestSolutionUpdateCount        int           `json:"so_lan_cap_nhat_nghiem_tot_nhat,omitempty"`
	LeafSolutionCount              int           `json:"so_nghiem_la,omitempty"`
	HitMaxNodeLimit                bool          `json:"cham_nguong_so_nut_toi_da,omitempty"`
	IterationsExecuted             int           `json:"so_vong_lap_da_thuc_hien,omitempty"`
	InitialPenalty                 int           `json:"diem_phat_ban_dau,omitempty"`
	BestPenalty                    int           `json:"diem_phat_tot_nhat,omitempty"`
	AcceptedMoveCount              int           `json:"so_buoc_chuyen_duoc_chap_nhan,omitempty"`
	TabuRejectedMoveCount          int           `json:"so_buoc_chuyen_bi_cam_loai_bo,omitempty"`
	RepairAssignmentCount          int           `json:"so_gan_duoc_sua_buoc_phuc_hoi,omitempty"`
}

func newSolverTelemetry(solverKey, solverLabel string, input SolverInput, problem preparedSchedulingProblem) *SolverTelemetry {
	telemetry := &SolverTelemetry{
		SolverKey:             solverKey,
		SolverLabel:           solverLabel,
		InputClassCount:       len(input.Classes),
		RequestedSessionCount: len(problem.variables),
		DomainStats:           computeDomainStats(problem.domains, problem.variables),
	}

	if input.BenchmarkOptions != nil {
		telemetry.ScenarioName = input.BenchmarkOptions.ScenarioName
		telemetry.RunIndex = input.BenchmarkOptions.RunIndex
		telemetry.TotalRuns = input.BenchmarkOptions.TotalRuns
		telemetry.RandomSeed = input.BenchmarkOptions.RandomSeed
		telemetry.IsWarmupRun = input.BenchmarkOptions.IsWarmupRun
	}

	return telemetry
}

func finalizeSolverTelemetry(telemetry *SolverTelemetry, startedAt time.Time, finishedAt time.Time) *SolverTelemetry {
	if telemetry == nil {
		return nil
	}
	telemetry.StartedAt = startedAt
	telemetry.FinishedAt = finishedAt
	telemetry.MeasuredAt = finishedAt
	telemetry.Runtime = finishedAt.Sub(startedAt)
	return telemetry
}

func computeDomainStats(domains map[string][]DomainValue, variables []Variable) DomainStats {
	stats := DomainStats{VariableCount: len(variables)}
	if len(variables) == 0 {
		return stats
	}

	sizes := make([]int, 0, len(variables))
	minSize := 0
	maxSize := 0
	totalSize := 0

	for _, variable := range variables {
		size := len(domains[variable.ID])
		sizes = append(sizes, size)
		totalSize += size
		if size == 0 {
			stats.EmptyDomainCount++
		}
		if len(sizes) == 1 || size < minSize {
			minSize = size
		}
		if size > maxSize {
			maxSize = size
		}
	}

	sort.Ints(sizes)
	stats.TotalDomainValues = totalSize
	stats.MinDomainSize = minSize
	stats.MaxDomainSize = maxSize
	stats.AvgDomainSize = float64(totalSize) / float64(len(sizes))
	stats.MedianDomainSize = medianIntSlice(sizes)
	return stats
}

func medianIntSlice(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	mid := len(values) / 2
	if len(values)%2 == 1 {
		return float64(values[mid])
	}
	return float64(values[mid-1]+values[mid]) / 2
}
