from __future__ import annotations

from dataclasses import dataclass
from typing import Final


LABEL_AT_RISK: Final[str] = "AT_RISK"
LABEL_NOT_AT_RISK: Final[str] = "NOT_AT_RISK"

ATTENDANCE_STATUS_PRESENT: Final[int] = 1
ATTENDANCE_STATUS_ABSENT: Final[int] = 2
ATTENDANCE_STATUS_EXCUSED: Final[int] = 3
ATTENDANCE_STATUS_LATE: Final[int] = 4
ATTENDANCE_STATUS_EARLY: Final[int] = 5

PRESENT_LIKE_STATUSES: Final[set[int]] = {
    ATTENDANCE_STATUS_PRESENT,
    ATTENDANCE_STATUS_LATE,
    ATTENDANCE_STATUS_EARLY,
}

DATASET_ID_COLUMNS: Final[list[str]] = [
    "snapshot_id",
    "student_id",
    "student_code",
    "student_name",
    "class_id",
    "class_code",
    "class_name",
    "snapshot_at",
]

FEATURE_COLUMNS: Final[list[str]] = [
    "attendance_rate_28d",
    "absence_count_28d",
    "average_total_score_28d",
    "homework_completion_rate_28d",
    "active_enrollment_count_28d",
    "weekly_lesson_load_28d",
    "approved_leave_count_28d",
    "days_since_last_lesson",
]

DATASET_COLUMNS: Final[list[str]] = DATASET_ID_COLUMNS + FEATURE_COLUMNS + ["label"]


@dataclass(frozen=True)
class DatasetDefinition:
    name: str = "student_at_risk_classification"
    prediction_unit: str = "student_class_snapshot"
    observation_window_days: int = 28
    prediction_horizon_days: int = 28
    minimum_future_attendance_rows: int = 4
    minimum_future_academic_rows: int = 2


DEFAULT_DATASET_DEFINITION = DatasetDefinition()
