package predictive

const (
	SourceStudent           = "student"
	SourceAttendance        = "attendance"
	SourceAcademicRecord    = "academic_record"
	SourceEnrollment        = "enrollment"
	SourceLesson            = "lesson"
	SourceLeaveRequest      = "leave_request"
	SourceLessonSummary     = "lesson_summary"
	LabelAtRisk             = "AT_RISK"
	LabelNotAtRisk          = "NOT_AT_RISK"
	AttendanceStatusPresent = 1
	AttendanceStatusAbsent  = 2
	AttendanceStatusExcused = 3
	AttendanceStatusLate    = 4
	AttendanceStatusEarly   = 5
)

type DataSourceDefinition struct {
	Key            string
	Entity         string
	Required       bool
	Availability   string
	JoinKeys       []string
	SelectedFields []string
	ReasonForUse   string
}

type LabelDefinition struct {
	Name                        string
	PositiveClass               string
	NegativeClass               string
	PredictionUnit              string
	ObservationWindowDays       int
	PredictionHorizonDays       int
	MinimumFutureAttendanceRows int
	MinimumFutureAcademicRows   int
	PositiveWhenAny             []string
	ExcludeWhen                 []string
}

type FeatureDefinition struct {
	Key         string
	SourceKeys  []string
	DataType    string
	Aggregation string
	Description string
}

type DatasetDefinition struct {
	Name                  string
	Version               string
	ProblemType           string
	PredictionUnit        string
	ObservationWindowDays int
	PredictionHorizonDays int
	Sources               []DataSourceDefinition
	Label                 LabelDefinition
	Features              []FeatureDefinition
	LeakageGuards         []string
	ExcludedFields        []string
}

func DefaultAtRiskDatasetDefinition() DatasetDefinition {
	return DatasetDefinition{
		Name:                  "student_at_risk_classification",
		Version:               "v1",
		ProblemType:           "classification",
		PredictionUnit:        "student_class_snapshot",
		ObservationWindowDays: 28,
		PredictionHorizonDays: 28,
		Sources: []DataSourceDefinition{
			{
				Key:          SourceStudent,
				Entity:       "entities.Student",
				Required:     true,
				Availability: "implemented_entity_and_crud",
				JoinKeys:     []string{"student_id"},
				SelectedFields: []string{
					"id",
					"grade_level",
					"date_of_birth",
					"gender",
					"status",
				},
				ReasonForUse: "Nguon profile co san cho age band, gender, grade level va active status.",
			},
			{
				Key:          SourceEnrollment,
				Entity:       "entities.Enrollment",
				Required:     true,
				Availability: "implemented_entity_and_crud_supporting_class_roster",
				JoinKeys:     []string{"student_id", "class_id"},
				SelectedFields: []string{
					"id",
					"student_id",
					"class_id",
					"status",
					"approved_at",
					"created_at",
				},
				ReasonForUse: "Dung de xac dinh student dang hoc lop nao, tai thoi diem nao va muc do tai hoc hien tai.",
			},
			{
				Key:          SourceLesson,
				Entity:       "entities.Lesson",
				Required:     true,
				Availability: "implemented_from_scheduling_commit",
				JoinKeys:     []string{"class_id", "lesson_id"},
				SelectedFields: []string{
					"id",
					"class_id",
					"teacher_id",
					"date_start",
					"date_end",
				},
				ReasonForUse: "Dung de suy ra cuong do hoc theo tuan, khoang cach giua cac buoi va gan attendance/summary vao dung snapshot.",
			},
			{
				Key:          SourceAttendance,
				Entity:       "entities.Attendance",
				Required:     true,
				Availability: "implemented_entity_only_status_mapping_required",
				JoinKeys:     []string{"student_id", "lesson_id"},
				SelectedFields: []string{
					"id",
					"student_id",
					"lesson_id",
					"status",
					"marked_at",
				},
				ReasonForUse: "Nguon chinh cho attendance rate, absence count va attendance incident features.",
			},
			{
				Key:          SourceAcademicRecord,
				Entity:       "entities.AcademicRecord",
				Required:     true,
				Availability: "implemented_entity_only",
				JoinKeys:     []string{"student_id", "lesson_summary_id"},
				SelectedFields: []string{
					"id",
					"student_id",
					"lesson_summary_id",
					"homework_completed",
					"homework_score",
					"attitude_rating",
					"participation_score",
					"total_score",
					"is_completed",
					"created_at",
				},
				ReasonForUse: "Nguon grade chinh cho bai toan AT_RISK, dung de tao score trend va homework completion features.",
			},
			{
				Key:          SourceLeaveRequest,
				Entity:       "entities.LeaveRequest",
				Required:     false,
				Availability: "implemented_entity_only",
				JoinKeys:     []string{"student_id", "lesson_id", "class_id"},
				SelectedFields: []string{
					"id",
					"student_id",
					"leave_type",
					"apply_date",
					"late_minutes",
					"early_minutes",
					"status",
				},
				ReasonForUse: "Tin hieu van hanh bo sung cho vang co phep, di tre, ve som va tan suat xin nghi.",
			},
			{
				Key:          SourceLessonSummary,
				Entity:       "entities.LessonSummary",
				Required:     false,
				Availability: "implemented_entity_only",
				JoinKeys:     []string{"lesson_id"},
				SelectedFields: []string{
					"id",
					"lesson_id",
					"homework_deadline",
					"created_at",
				},
				ReasonForUse: "Nguon bo sung de tinh ty le buoi hoc da co summary va lien ket hoc vu hoan chinh.",
			},
		},
		Label: LabelDefinition{
			Name:                        "student_at_risk_label",
			PositiveClass:               LabelAtRisk,
			NegativeClass:               LabelNotAtRisk,
			PredictionUnit:              "student_class_snapshot",
			ObservationWindowDays:       28,
			PredictionHorizonDays:       28,
			MinimumFutureAttendanceRows: 4,
			MinimumFutureAcademicRows:   2,
			PositiveWhenAny: []string{
				"future_attendance_rate < 0.80 over the next 28 days",
				"future_average_total_score < 5.00 over completed academic records in the next 28 days",
				"future_homework_completion_rate < 0.60 over completed academic records in the next 28 days",
			},
			ExcludeWhen: []string{
				"student is not ACTIVE at snapshot date",
				"student has no active enrollment at snapshot date",
				"future horizon has fewer than 4 attendance rows and fewer than 2 academic records",
			},
		},
		Features: []FeatureDefinition{
			{Key: "student_grade_level", SourceKeys: []string{SourceStudent}, DataType: "categorical", Aggregation: "snapshot", Description: "Khoi lop cua hoc vien tai thoi diem snapshot."},
			{Key: "student_age_years", SourceKeys: []string{SourceStudent}, DataType: "numeric", Aggregation: "snapshot", Description: "Tuoi tai thoi diem snapshot, suy ra tu date_of_birth."},
			{Key: "student_gender", SourceKeys: []string{SourceStudent}, DataType: "categorical", Aggregation: "snapshot", Description: "Gioi tinh neu co du lieu."},
			{Key: "active_enrollment_count_28d", SourceKeys: []string{SourceEnrollment}, DataType: "numeric", Aggregation: "count", Description: "So lop dang hoc trong cua so quan sat."},
			{Key: "weekly_lesson_load_28d", SourceKeys: []string{SourceEnrollment, SourceLesson}, DataType: "numeric", Aggregation: "average_per_week", Description: "Cuong do hoc trung binh moi tuan trong 28 ngay quan sat."},
			{Key: "attendance_rate_28d", SourceKeys: []string{SourceAttendance}, DataType: "numeric", Aggregation: "ratio", Description: "Ty le tham du trong 28 ngay truoc snapshot."},
			{Key: "absence_count_28d", SourceKeys: []string{SourceAttendance}, DataType: "numeric", Aggregation: "count", Description: "Tong so buoi vang trong 28 ngay truoc snapshot."},
			{Key: "late_or_early_incident_count_28d", SourceKeys: []string{SourceAttendance, SourceLeaveRequest}, DataType: "numeric", Aggregation: "count", Description: "So lan di tre/ve som trong cua so quan sat."},
			{Key: "approved_leave_count_28d", SourceKeys: []string{SourceLeaveRequest}, DataType: "numeric", Aggregation: "count", Description: "So don nghi da duoc duyet trong 28 ngay truoc snapshot."},
			{Key: "average_total_score_28d", SourceKeys: []string{SourceAcademicRecord}, DataType: "numeric", Aggregation: "mean", Description: "Diem tong ket trung binh cua academic record trong 28 ngay truoc snapshot."},
			{Key: "minimum_total_score_28d", SourceKeys: []string{SourceAcademicRecord}, DataType: "numeric", Aggregation: "min", Description: "Diem tong ket thap nhat trong cua so quan sat."},
			{Key: "homework_completion_rate_28d", SourceKeys: []string{SourceAcademicRecord}, DataType: "numeric", Aggregation: "ratio", Description: "Ty le hoan thanh bai tap trong 28 ngay truoc snapshot."},
			{Key: "average_homework_score_28d", SourceKeys: []string{SourceAcademicRecord}, DataType: "numeric", Aggregation: "mean", Description: "Diem bai tap trung binh trong cua so quan sat."},
			{Key: "average_participation_score_28d", SourceKeys: []string{SourceAcademicRecord}, DataType: "numeric", Aggregation: "mean", Description: "Diem tham gia trung binh trong cua so quan sat."},
			{Key: "average_attitude_rating_28d", SourceKeys: []string{SourceAcademicRecord}, DataType: "numeric", Aggregation: "mean", Description: "Thai do hoc tap trung binh trong cua so quan sat."},
			{Key: "completed_record_ratio_28d", SourceKeys: []string{SourceAcademicRecord, SourceLessonSummary}, DataType: "numeric", Aggregation: "ratio", Description: "Ty le ban ghi academic da completed trong cua so quan sat."},
			{Key: "days_since_last_lesson", SourceKeys: []string{SourceLesson}, DataType: "numeric", Aggregation: "snapshot", Description: "So ngay ke tu lesson gan nhat truoc snapshot."},
		},
		LeakageGuards: []string{
			"Khong dua future attendance, future academic score hoac bat ky field nao nam sau snapshot vao feature set.",
			"Khong dua total_score, homework_completed, attendance status cua prediction horizon vao feature engineering.",
			"Khong dua free-text fields nhu personal_comment, teacher_notes, class_feedback vao baseline F2 de giam do phuc tap va tranh leakage qua ghi chu hau kiem.",
		},
		ExcludedFields: []string{
			"student.full_name",
			"student.phone",
			"student.guardian_phone",
			"student.address",
			"academic_record.personal_comment",
			"lesson_summary.teacher_notes",
			"lesson_summary.lesson_content",
			"leave_request.reason",
		},
	}
}
