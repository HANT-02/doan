package entities

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type Teacher struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Code            string         `gorm:"type:varchar(50);unique" json:"code"`
	FullName        string         `gorm:"type:varchar(255)" json:"fullName"`
	Email           string         `gorm:"type:varchar(255)" json:"email"`
	Phone           string         `gorm:"type:varchar(20)" json:"phone"`
	IsSchoolTeacher bool           `gorm:"default:false" json:"isSchoolTeacher"`
	SchoolName      string         `gorm:"type:varchar(255)" json:"schoolName"`
	EmploymentType  string         `gorm:"type:varchar(50);default:'PART_TIME'" json:"employmentType"`
	Status          string         `gorm:"type:varchar(50);default:'ACTIVE'" json:"status"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `gorm:"default:now()" json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

// Table 3.3
type Student struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Code          string         `gorm:"type:varchar(50);unique" json:"code"`
	FullName      string         `gorm:"type:varchar(255)" json:"fullName"`
	Email         string         `gorm:"type:varchar(255)" json:"email"`
	Phone         string         `gorm:"type:varchar(20)" json:"phone"`
	GuardianPhone string         `gorm:"type:varchar(20)" json:"guardianPhone"`
	GradeLevel    string         `gorm:"type:varchar(50)" json:"gradeLevel"`
	Status        string         `gorm:"type:varchar(50);default:'ACTIVE'" json:"status"`
	DateOfBirth   *time.Time     `json:"dateOfBirth"`
	Gender        string         `gorm:"type:varchar(20)" json:"gender"`
	Address       string         `gorm:"type:text" json:"address"`
	CreatedAt     time.Time      `gorm:"default:now()" json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

// Table 3.4
type Room struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Capacity  int       `json:"capacity"`
	Address   string    `gorm:"type:text" json:"address"`
	CreatedAt time.Time `gorm:"default:now()" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Table 3.5
type Course struct {
	ID                     uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Code                   string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name                   string         `gorm:"type:varchar(255)" json:"name"`
	Description            string         `gorm:"type:text" json:"description"`
	GradeLevel             string         `gorm:"type:varchar(50)" json:"gradeLevel"`
	Subject                string         `gorm:"type:varchar(255)" json:"subject"`
	SessionCount           int            `json:"sessionCount"`
	SessionDurationMinutes int            `json:"sessionDurationMinutes"`
	TotalHours             float64        `gorm:"type:numeric(8,2)" json:"totalHours"`
	Price                  float64        `gorm:"type:numeric(10,2)" json:"price"`
	Status                 string         `gorm:"type:varchar(50);default:'ACTIVE'" json:"status"`
	CreatedAt              time.Time      `gorm:"default:now()" json:"createdAt"`
	UpdatedAt              time.Time      `json:"updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

// Table 3.6
type Program struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Code          string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name          string         `gorm:"type:varchar(255)" json:"name"`
	Track         string         `gorm:"type:varchar(50)" json:"track"`
	EffectiveFrom *time.Time     `json:"effectiveFrom"`
	EffectiveTo   *time.Time     `json:"effectiveTo"`
	CreatedByID   *uint          `json:"createdById"`
	CreatedBy     User           `gorm:"foreignKey:CreatedByID" json:"createdBy"`
	ApprovedByID  *uint          `json:"approvedById"`
	ApprovedBy    User           `gorm:"foreignKey:ApprovedByID" json:"approvedBy"`
	ApprovalNote  string         `gorm:"type:text" json:"approvalNote"`
	PublishedAt   *time.Time     `json:"publishedAt"`
	ArchivedAt    *time.Time     `json:"archivedAt"`
	CreatedAt     time.Time      `gorm:"default:now()" json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

// Table 3.7
type ProgramCourse struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProgramID uint    `gorm:"not null" json:"programId"`
	Program   Program `gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE" json:"-"`
	CourseID  uint    `gorm:"not null" json:"courseId"`
	Course    Course  `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course"`
}

// Table 3.8
type Objective struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string  `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name      string  `gorm:"type:text;not null" json:"name"`
	ProgramID uint    `gorm:"not null" json:"programId"`
	Program   Program `gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE" json:"-"`
}

// Table 3.9
type Outcome struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name        string    `gorm:"type:text;not null" json:"name"`
	ProgramID   uint      `gorm:"not null" json:"programId"`
	Program     Program   `gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE" json:"-"`
	ObjectiveID *uint     `json:"objectiveId"`
	Objective   Objective `gorm:"foreignKey:ObjectiveID;constraint:OnDelete:SET NULL" json:"objective"`
}

// Table 3.10 (Structure corresponds to Class in PDF)
type Class struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Notes       string         `gorm:"type:text" json:"notes"`
	StartDate   time.Time      `gorm:"not null" json:"startDate"`
	EndDate     *time.Time     `json:"endDate"`
	MaxStudents int            `json:"maxStudents"`
	Status      string         `gorm:"type:varchar(50);default:'OPEN'" json:"status"`
	Price       float64        `gorm:"type:numeric(10,2)" json:"price"`
	ProgramID   *uint          `json:"programId"`
	Program     Program        `gorm:"foreignKey:ProgramID" json:"program"`
	CourseID    *uint          `json:"courseId"`
	Course      Course         `gorm:"foreignKey:CourseID" json:"course"`
	TeacherID   *uint          `json:"teacherId"`
	Teacher     Teacher        `gorm:"foreignKey:TeacherID" json:"teacher"`
	CreatedAt   time.Time      `gorm:"default:now()" json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

// Table 3.11
type Lesson struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassID   uint      `gorm:"not null" json:"classId"`
	Class     Class     `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"class"`
	DateStart time.Time `gorm:"not null" json:"dateStart"`
	DateEnd   time.Time `gorm:"not null" json:"dateEnd"`
	RoomID    *uint     `json:"roomId"`
	Room      Room      `gorm:"foreignKey:RoomID" json:"room"`
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `gorm:"default:now()" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Table 3.12
type Enrollment struct {
	ID         string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClassID    uint       `gorm:"not null" json:"classId"`
	Class      Class      `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"class"`
	StudentID  uint       `gorm:"not null" json:"studentId"`
	Student    Student    `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student"`
	Status     string     `gorm:"type:varchar(50);default:'APPLIED'" json:"status"`
	ApprovedAt *time.Time `json:"approvedAt"`
	RejectedAt *time.Time `json:"rejectedAt"`
	CreatedAt  time.Time  `gorm:"default:now()" json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// Table 3.13
type Attendance struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LessonID  uint      `gorm:"not null" json:"lessonId"`
	Lesson    Lesson    `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"lesson"`
	StudentID uint      `gorm:"not null" json:"studentId"`
	Student   Student   `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student"`
	Status    int       `gorm:"not null" json:"status"`
	Note      string    `gorm:"type:text" json:"note"`
	MarkedAt  time.Time `json:"markedAt"`
	CreatedAt time.Time `gorm:"default:now()" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Table 3.14
type ClassSchedule struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassID   uint   `gorm:"not null" json:"classId"`
	Class     Class  `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"-"`
	DayOfWeek string `gorm:"type:varchar(20);not null" json:"dayOfWeek"`
	StartTime string `gorm:"type:varchar(10);not null" json:"startTime"`
	EndTime   string `gorm:"type:varchar(10);not null" json:"endTime"`
	RoomID    *uint  `json:"roomId"`
	Room      Room   `gorm:"foreignKey:RoomID" json:"room"`
}

// Table 3.15
type LessonSummary struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LessonID         uint      `gorm:"unique;not null" json:"lessonId"`
	Lesson           Lesson    `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"lesson"`
	Topic            string    `gorm:"type:text" json:"topic"`
	LessonContent    string    `gorm:"type:text" json:"lessonContent"`
	ClassFeedback    string    `gorm:"type:text" json:"classFeedback"`
	Homework         string    `gorm:"type:text" json:"homework"`
	HomeworkDeadline time.Time `json:"homeworkDeadline"`
	TeacherNotes     string    `gorm:"type:text" json:"teacherNotes"`
	CreatedByID      *uint     `json:"createdById"`
	CreatedBy        User      `gorm:"foreignKey:CreatedByID" json:"createdBy"`
	CreatedAt        time.Time `gorm:"default:now()" json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// Table 3.16
type AcademicRecord struct {
	ID                 uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	LessonSummaryID    uint          `gorm:"not null" json:"lessonSummaryId"`
	LessonSummary      LessonSummary `gorm:"foreignKey:LessonSummaryID;constraint:OnDelete:CASCADE" json:"-"`
	StudentID          uint          `gorm:"not null" json:"studentId"`
	Student            Student       `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student"`
	HomeworkCompleted  bool          `gorm:"default:false" json:"homeworkCompleted"`
	HomeworkScore      float64       `gorm:"type:numeric(5,2)" json:"homeworkScore"`
	AttitudeRating     int           `json:"attitudeRating"`
	ParticipationScore float64       `gorm:"type:numeric(5,2)" json:"participationScore"`
	PersonalComment    string        `gorm:"type:text" json:"personalComment"`
	TotalScore         float64       `gorm:"type:numeric(5,2)" json:"totalScore"`
	IsCompleted        bool          `gorm:"default:false" json:"isCompleted"`
	CreatedAt          time.Time     `gorm:"default:now()" json:"createdAt"`
	UpdatedAt          time.Time     `json:"updatedAt"`
}

// Table 3.17
type Consultation struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName   string    `gorm:"type:varchar(255);not null" json:"fullName"`
	Phone      string    `gorm:"type:varchar(20);not null" json:"phone"`
	GradeLevel string    `gorm:"type:varchar(50);not null" json:"gradeLevel"`
	Notes      string    `gorm:"type:text" json:"notes"`
	Status     string    `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	CreatedAt  time.Time `gorm:"default:now()" json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Table 3.18
type LeaveRequest struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID       uint           `gorm:"not null" json:"studentId"`
	Student         Student        `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student"`
	LeaveType       string         `gorm:"type:varchar(50);not null" json:"leaveType"` // LEAVE, LATE, EARLY
	ApplyDate       time.Time      `gorm:"not null" json:"applyDate"`
	LateMinutes     int            `json:"lateMinutes"`
	EarlyMinutes    int            `json:"earlyMinutes"`
	Reason          string         `gorm:"type:text;not null" json:"reason"`
	Documents       pq.StringArray `gorm:"type:text[]" json:"documents"`
	ClassID         *uint          `json:"classId"`
	Class           Class          `gorm:"foreignKey:ClassID;constraint:OnDelete:SET NULL" json:"class"`
	LessonID        *uint          `json:"lessonId"`
	Lesson          Lesson         `gorm:"foreignKey:LessonID;constraint:OnDelete:SET NULL" json:"lesson"`
	Subject         string         `gorm:"type:varchar(255)" json:"subject"`
	Status          string         `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	ApprovedByID    *uint          `json:"approvedBy"`
	ApprovedBy      User           `gorm:"foreignKey:ApprovedByID" json:"approver"`
	ApprovedAt      *time.Time     `json:"approvedAt"`
	RejectionReason string         `gorm:"type:text" json:"rejectionReason"`
	CreatedAt       time.Time      `gorm:"default:now()" json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}
