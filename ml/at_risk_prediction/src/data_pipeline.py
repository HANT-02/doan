from __future__ import annotations

import json
from dataclasses import asdict
from datetime import datetime, timedelta, timezone
from pathlib import Path
from typing import Any

import pandas as pd
from sklearn.model_selection import train_test_split
from sqlalchemy import text
from sqlalchemy.engine import Engine

from .config import PROCESSED_DATA_DIR, REPORTS_DIR
from .dataset_schema import (
    ATTENDANCE_STATUS_ABSENT,
    DATASET_COLUMNS,
    DEFAULT_DATASET_DEFINITION,
    FEATURE_COLUMNS,
    LABEL_AT_RISK,
    LABEL_NOT_AT_RISK,
    PRESENT_LIKE_STATUSES,
    DatasetDefinition,
)


def _ensure_datetime(value: Any) -> datetime | None:
    if value is None or value == "":
        return None
    if isinstance(value, datetime):
        parsed = value
    else:
        parsed = pd.to_datetime(value).to_pydatetime()
    if parsed.tzinfo is not None:
        return parsed.astimezone(timezone.utc).replace(tzinfo=None)
    return parsed


def _safe_float(value: Any, default: float) -> float:
    if value is None or value == "":
        return default
    try:
        return float(value)
    except (TypeError, ValueError):
        return default


def _safe_str(value: Any) -> str:
    if value is None:
        return ""
    return str(value)


def load_source_tables(engine: Engine) -> dict[str, pd.DataFrame]:
    queries = {
        "students": """
            SELECT id, code, full_name, email, grade_level, status, created_at
            FROM students
            WHERE status = 'ACTIVE'
        """,
        "enrollments": """
            SELECT
                e.id,
                e.student_id,
                e.class_id,
                e.status,
                e.approved_at,
                e.created_at,
                c.code AS class_code,
                c.name AS class_name,
                c.start_date AS class_start_date,
                c.end_date AS class_end_date
            FROM enrollments e
            JOIN classes c ON c.id = e.class_id
            WHERE e.status = 'ENROLLED'
        """,
        "lessons": """
            SELECT
                l.id,
                l.class_id,
                l.teacher_id,
                l.room_id,
                l.date_start,
                l.date_end,
                c.code AS class_code,
                c.name AS class_name
            FROM lessons l
            JOIN classes c ON c.id = l.class_id
            ORDER BY l.date_start ASC
        """,
        "attendance": """
            SELECT id, lesson_id, student_id, status, marked_at
            FROM attendances
        """,
        "lesson_summaries": """
            SELECT id, lesson_id, created_at, homework_deadline
            FROM lesson_summaries
        """,
        "academic_records": """
            SELECT
                id,
                lesson_summary_id,
                student_id,
                homework_completed,
                total_score,
                is_completed,
                created_at
            FROM academic_records
        """,
        "leave_requests": """
            SELECT
                id,
                student_id,
                class_id,
                apply_date,
                status
            FROM leave_requests
            WHERE status = 'APPROVED'
        """,
    }

    with engine.connect() as connection:
        return {
            name: pd.read_sql(text(sql), connection)
            for name, sql in queries.items()
        }


def build_dataset_from_db(
    engine: Engine,
    definition: DatasetDefinition = DEFAULT_DATASET_DEFINITION,
) -> pd.DataFrame:
    tables = load_source_tables(engine)
    return build_dataset_from_tables(tables, definition)


def build_inference_dataset_from_db(
    engine: Engine,
    definition: DatasetDefinition = DEFAULT_DATASET_DEFINITION,
) -> pd.DataFrame:
    tables = load_source_tables(engine)
    return build_inference_dataset_from_tables(tables, definition)


def build_dataset_from_tables(
    tables: dict[str, pd.DataFrame],
    definition: DatasetDefinition = DEFAULT_DATASET_DEFINITION,
) -> pd.DataFrame:
    students = tables["students"].copy()
    enrollments = tables["enrollments"].copy()
    lessons = tables["lessons"].copy()
    attendance = tables["attendance"].copy()
    lesson_summaries = tables["lesson_summaries"].copy()
    academic_records = tables["academic_records"].copy()
    leave_requests = tables["leave_requests"].copy()

    if lessons.empty:
        raise ValueError(
            "Khong tim thay lesson trong database. Hay commit scheduling preview "
            "hoac seed lesson truoc khi xay dung dataset predictive."
        )

    for frame, column_names in (
        (lessons, ["date_start", "date_end"]),
        (enrollments, ["approved_at", "created_at", "class_start_date", "class_end_date"]),
        (attendance, ["marked_at"]),
        (lesson_summaries, ["created_at", "homework_deadline"]),
        (academic_records, ["created_at"]),
        (leave_requests, ["apply_date"]),
    ):
        for column in column_names:
            if column in frame.columns:
                frame[column] = pd.to_datetime(frame[column], errors="coerce")

    lesson_by_id = {
        row["id"]: row
        for row in lessons.to_dict("records")
    }
    lessons_by_class: dict[str, list[dict[str, Any]]] = {}
    for row in lessons.to_dict("records"):
        lessons_by_class.setdefault(row["class_id"], []).append(row)
    for class_id in lessons_by_class:
        lessons_by_class[class_id].sort(key=lambda item: item["date_start"])

    summaries_by_id = {
        row["id"]: row
        for row in lesson_summaries.to_dict("records")
    }

    students_by_id = {
        row["id"]: row
        for row in students.to_dict("records")
    }

    enrollments_by_student: dict[str, list[dict[str, Any]]] = {}
    for row in enrollments.to_dict("records"):
        enrollments_by_student.setdefault(row["student_id"], []).append(row)

    attendance_by_student_class: dict[tuple[str, str], list[dict[str, Any]]] = {}
    for row in attendance.to_dict("records"):
        lesson = lesson_by_id.get(row["lesson_id"])
        if not lesson:
            continue
        key = (row["student_id"], lesson["class_id"])
        attendance_by_student_class.setdefault(key, []).append(
            {
                "at": lesson["date_start"],
                "status": int(row["status"]),
            }
        )

    academic_by_student_class: dict[tuple[str, str], list[dict[str, Any]]] = {}
    for row in academic_records.to_dict("records"):
        summary = summaries_by_id.get(row["lesson_summary_id"])
        if not summary:
            continue
        lesson = lesson_by_id.get(summary["lesson_id"])
        if not lesson:
            continue
        key = (row["student_id"], lesson["class_id"])
        academic_by_student_class.setdefault(key, []).append(
            {
                "at": lesson["date_start"],
                "homework_completed": bool(row["homework_completed"]),
                "total_score": _safe_float(row["total_score"], 0.0),
                "is_completed": bool(row["is_completed"]),
            }
        )

    leaves_by_student: dict[str, list[dict[str, Any]]] = {}
    for row in leave_requests.to_dict("records"):
        leaves_by_student.setdefault(row["student_id"], []).append(
            {
                "at": row["apply_date"],
                "class_id": row["class_id"],
            }
        )

    records: list[dict[str, Any]] = []
    for student_id, student_enrollments in enrollments_by_student.items():
        student = students_by_id.get(student_id)
        if not student:
            continue

        for enrollment in student_enrollments:
            class_lessons = lessons_by_class.get(enrollment["class_id"], [])
            if not class_lessons:
                continue

            key = (student_id, enrollment["class_id"])
            class_attendance = attendance_by_student_class.get(key, [])
            class_academic = academic_by_student_class.get(key, [])

            for snapshot_lesson in class_lessons:
                snapshot_at = _ensure_datetime(snapshot_lesson["date_start"])
                if snapshot_at is None:
                    continue
                if not _is_enrollment_active_at(enrollment, snapshot_at):
                    continue

                observation_start = snapshot_at - timedelta(days=definition.observation_window_days)
                future_end = snapshot_at + timedelta(days=definition.prediction_horizon_days)

                obs_attendance = _filter_between(class_attendance, "at", observation_start, snapshot_at, include_start=True, include_end=False)
                future_attendance = _filter_between(class_attendance, "at", snapshot_at, future_end, include_start=False, include_end=False)
                obs_academic = _completed_records(_filter_between(class_academic, "at", observation_start, snapshot_at, include_start=True, include_end=False))
                future_academic = _completed_records(_filter_between(class_academic, "at", snapshot_at, future_end, include_start=False, include_end=False))

                if (
                    len(future_attendance) < definition.minimum_future_attendance_rows
                    and len(future_academic) < definition.minimum_future_academic_rows
                ):
                    continue

                active_class_ids = _active_enrollment_class_ids(student_enrollments, snapshot_at)
                records.append(
                    {
                        "snapshot_id": _make_snapshot_id(student_id, enrollment["class_id"], snapshot_at),
                        "student_id": _safe_str(student_id),
                        "student_code": _safe_str(student.get("code", "")),
                        "student_name": _safe_str(student.get("full_name", "")),
                        "class_id": _safe_str(enrollment["class_id"]),
                        "class_code": _safe_str(enrollment.get("class_code", "")),
                        "class_name": _safe_str(enrollment.get("class_name", "")),
                        "snapshot_at": snapshot_at.isoformat(),
                        "attendance_rate_28d": _attendance_rate(obs_attendance),
                        "absence_count_28d": float(_absence_count(obs_attendance)),
                        "average_total_score_28d": _average_score(obs_academic),
                        "homework_completion_rate_28d": _homework_completion_rate(obs_academic),
                        "active_enrollment_count_28d": float(len(active_class_ids)),
                        "weekly_lesson_load_28d": _weekly_lesson_load(
                            lessons_by_class,
                            active_class_ids,
                            observation_start,
                            snapshot_at,
                            definition.observation_window_days,
                        ),
                        "approved_leave_count_28d": float(
                            _approved_leave_count(
                                leaves_by_student.get(student_id, []),
                                observation_start,
                                snapshot_at,
                                enrollment["class_id"],
                            )
                        ),
                        "days_since_last_lesson": _days_since_last_lesson(
                            lessons_by_class,
                            active_class_ids,
                            snapshot_at,
                            definition.observation_window_days,
                        ),
                        "label": _derive_label(future_attendance, future_academic),
                    }
                )

    dataset = pd.DataFrame.from_records(records)
    if dataset.empty:
        raise ValueError("Khong tao duoc dong dataset nao tu nguon du lieu hien tai.")

    dataset = dataset.drop_duplicates(subset=["snapshot_id"]).sort_values("snapshot_id").reset_index(drop=True)
    return dataset[DATASET_COLUMNS]


def build_inference_dataset_from_tables(
    tables: dict[str, pd.DataFrame],
    definition: DatasetDefinition = DEFAULT_DATASET_DEFINITION,
) -> pd.DataFrame:
    students = tables["students"].copy()
    enrollments = tables["enrollments"].copy()
    lessons = tables["lessons"].copy()
    attendance = tables["attendance"].copy()
    lesson_summaries = tables["lesson_summaries"].copy()
    academic_records = tables["academic_records"].copy()
    leave_requests = tables["leave_requests"].copy()

    for frame, column_names in (
        (lessons, ["date_start", "date_end"]),
        (enrollments, ["approved_at", "created_at", "class_start_date", "class_end_date"]),
        (attendance, ["marked_at"]),
        (lesson_summaries, ["created_at", "homework_deadline"]),
        (academic_records, ["created_at"]),
        (leave_requests, ["apply_date"]),
    ):
        for column in column_names:
            if column in frame.columns:
                frame[column] = pd.to_datetime(frame[column], errors="coerce")

    lesson_by_id = {
        row["id"]: row
        for row in lessons.to_dict("records")
    }
    lessons_by_class: dict[str, list[dict[str, Any]]] = {}
    for row in lessons.to_dict("records"):
        lessons_by_class.setdefault(row["class_id"], []).append(row)
    for class_id in lessons_by_class:
        lessons_by_class[class_id].sort(key=lambda item: item["date_start"])

    summaries_by_id = {
        row["id"]: row
        for row in lesson_summaries.to_dict("records")
    }
    students_by_id = {
        row["id"]: row
        for row in students.to_dict("records")
    }

    enrollments_by_student: dict[str, list[dict[str, Any]]] = {}
    for row in enrollments.to_dict("records"):
        enrollments_by_student.setdefault(row["student_id"], []).append(row)

    attendance_by_student_class: dict[tuple[str, str], list[dict[str, Any]]] = {}
    for row in attendance.to_dict("records"):
        lesson = lesson_by_id.get(row["lesson_id"])
        if not lesson:
            continue
        key = (row["student_id"], lesson["class_id"])
        attendance_by_student_class.setdefault(key, []).append(
            {
                "at": lesson["date_start"],
                "status": int(row["status"]),
            }
        )

    academic_by_student_class: dict[tuple[str, str], list[dict[str, Any]]] = {}
    for row in academic_records.to_dict("records"):
        summary = summaries_by_id.get(row["lesson_summary_id"])
        if not summary:
            continue
        lesson = lesson_by_id.get(summary["lesson_id"])
        if not lesson:
            continue
        key = (row["student_id"], lesson["class_id"])
        academic_by_student_class.setdefault(key, []).append(
            {
                "at": lesson["date_start"],
                "homework_completed": bool(row["homework_completed"]),
                "total_score": _safe_float(row["total_score"], 0.0),
                "is_completed": bool(row["is_completed"]),
            }
        )

    leaves_by_student: dict[str, list[dict[str, Any]]] = {}
    for row in leave_requests.to_dict("records"):
        leaves_by_student.setdefault(row["student_id"], []).append(
            {
                "at": row["apply_date"],
                "class_id": row["class_id"],
            }
        )

    now = datetime.utcnow()
    records: list[dict[str, Any]] = []
    for student_id, student_enrollments in enrollments_by_student.items():
        student = students_by_id.get(student_id)
        if not student:
            continue

        for enrollment in student_enrollments:
            snapshot_at = _pick_inference_snapshot_at(lessons_by_class.get(enrollment["class_id"], []), enrollment, now)
            if snapshot_at is None:
                continue
            if not _is_enrollment_active_at(enrollment, snapshot_at):
                continue

            observation_start = snapshot_at - timedelta(days=definition.observation_window_days)
            key = (student_id, enrollment["class_id"])
            obs_attendance = _filter_between(
                attendance_by_student_class.get(key, []),
                "at",
                observation_start,
                snapshot_at,
                include_start=True,
                include_end=True,
            )
            obs_academic = _completed_records(
                _filter_between(
                    academic_by_student_class.get(key, []),
                    "at",
                    observation_start,
                    snapshot_at,
                    include_start=True,
                    include_end=True,
                )
            )
            active_class_ids = _active_enrollment_class_ids(student_enrollments, snapshot_at)

            records.append(
                {
                    "snapshot_id": _make_snapshot_id(student_id, enrollment["class_id"], snapshot_at),
                    "student_id": _safe_str(student_id),
                    "student_code": _safe_str(student.get("code", "")),
                    "student_name": _safe_str(student.get("full_name", "")),
                    "class_id": _safe_str(enrollment["class_id"]),
                    "class_code": _safe_str(enrollment.get("class_code", "")),
                    "class_name": _safe_str(enrollment.get("class_name", "")),
                    "snapshot_at": snapshot_at.isoformat(),
                    "attendance_rate_28d": _attendance_rate(obs_attendance),
                    "absence_count_28d": float(_absence_count(obs_attendance)),
                    "average_total_score_28d": _average_score(obs_academic),
                    "homework_completion_rate_28d": _homework_completion_rate(obs_academic),
                    "active_enrollment_count_28d": float(len(active_class_ids)),
                    "weekly_lesson_load_28d": _weekly_lesson_load(
                        lessons_by_class,
                        active_class_ids,
                        observation_start,
                        snapshot_at,
                        definition.observation_window_days,
                    ),
                    "approved_leave_count_28d": float(
                        _approved_leave_count(
                            leaves_by_student.get(student_id, []),
                            observation_start,
                            snapshot_at,
                            enrollment["class_id"],
                        )
                    ),
                    "days_since_last_lesson": _days_since_last_lesson(
                        lessons_by_class,
                        active_class_ids,
                        snapshot_at,
                        definition.observation_window_days,
                    ),
                    "label": LABEL_NOT_AT_RISK,
                }
            )

    dataset = pd.DataFrame.from_records(records)
    if dataset.empty:
        raise ValueError("Khong tao duoc dong inference nao tu nguon du lieu hien tai.")

    dataset = dataset.drop_duplicates(subset=["snapshot_id"]).sort_values("snapshot_id").reset_index(drop=True)
    return dataset[DATASET_COLUMNS]


def load_dataset_from_csv(path: str | Path) -> pd.DataFrame:
    dataset = pd.read_csv(path)
    missing = [column for column in DATASET_COLUMNS if column not in dataset.columns]
    if missing:
        raise ValueError(f"CSV dataset thieu cac cot bat buoc: {', '.join(missing)}")
    return dataset[DATASET_COLUMNS].copy()


def split_dataset(
    dataset: pd.DataFrame,
    test_size: float = 0.2,
    seed: int = 42,
) -> tuple[pd.DataFrame, pd.DataFrame]:
    stratify = dataset["label"] if dataset["label"].nunique() > 1 else None
    train_df, test_df = train_test_split(
        dataset,
        test_size=test_size,
        random_state=seed,
        shuffle=True,
        stratify=stratify,
    )
    return (
        train_df.sort_values("snapshot_id").reset_index(drop=True),
        test_df.sort_values("snapshot_id").reset_index(drop=True),
    )


def build_dataset_summary(
    dataset: pd.DataFrame,
    definition: DatasetDefinition = DEFAULT_DATASET_DEFINITION,
    source: str = "unknown",
) -> dict[str, Any]:
    label_counts = dataset["label"].value_counts().to_dict()
    return {
        "dataset_name": definition.name,
        "prediction_unit": definition.prediction_unit,
        "source": source,
        "observation_window_days": definition.observation_window_days,
        "prediction_horizon_days": definition.prediction_horizon_days,
        "row_count": int(len(dataset)),
        "student_count": int(dataset["student_id"].nunique()),
        "class_count": int(dataset["class_id"].nunique()),
        "feature_columns": FEATURE_COLUMNS,
        "label_distribution": {
            LABEL_AT_RISK: int(label_counts.get(LABEL_AT_RISK, 0)),
            LABEL_NOT_AT_RISK: int(label_counts.get(LABEL_NOT_AT_RISK, 0)),
        },
        "generated_at": datetime.utcnow().isoformat() + "Z",
    }


def save_dataset_artifacts(
    dataset: pd.DataFrame,
    summary: dict[str, Any],
    dataset_name: str,
    save_split: bool = True,
    test_size: float = 0.2,
    seed: int = 42,
) -> dict[str, str]:
    PROCESSED_DATA_DIR.mkdir(parents=True, exist_ok=True)
    REPORTS_DIR.mkdir(parents=True, exist_ok=True)

    full_path = PROCESSED_DATA_DIR / f"{dataset_name}_full.csv"
    dataset.to_csv(full_path, index=False)

    paths: dict[str, str] = {"full_dataset_csv": str(full_path)}

    if save_split:
        train_df, test_df = split_dataset(dataset, test_size=test_size, seed=seed)
        train_path = PROCESSED_DATA_DIR / f"{dataset_name}_train.csv"
        test_path = PROCESSED_DATA_DIR / f"{dataset_name}_test.csv"
        train_df.to_csv(train_path, index=False)
        test_df.to_csv(test_path, index=False)
        paths["train_dataset_csv"] = str(train_path)
        paths["test_dataset_csv"] = str(test_path)

        summary["train_size"] = int(len(train_df))
        summary["test_size"] = int(len(test_df))

    summary_path = REPORTS_DIR / f"{dataset_name}_dataset_summary.json"
    summary_path.write_text(json.dumps(summary, indent=2, ensure_ascii=False), encoding="utf-8")
    paths["dataset_summary_json"] = str(summary_path)
    return paths


def _is_enrollment_active_at(enrollment: dict[str, Any], snapshot_at: datetime) -> bool:
    effective_start = _ensure_datetime(enrollment.get("created_at")) or snapshot_at
    approved_at = _ensure_datetime(enrollment.get("approved_at"))
    if approved_at and approved_at > effective_start:
        effective_start = approved_at
    if effective_start > snapshot_at:
        return False

    class_start = _ensure_datetime(enrollment.get("class_start_date"))
    class_end = _ensure_datetime(enrollment.get("class_end_date"))
    if class_start and class_start > snapshot_at:
        return False
    if class_end and class_end < snapshot_at:
        return False
    return True


def _filter_between(
    items: list[dict[str, Any]],
    field: str,
    start: datetime,
    end: datetime,
    *,
    include_start: bool,
    include_end: bool,
) -> list[dict[str, Any]]:
    filtered: list[dict[str, Any]] = []
    for item in items:
        value = _ensure_datetime(item.get(field))
        if value is None:
            continue

        lower_ok = value >= start if include_start else value > start
        upper_ok = value <= end if include_end else value < end
        if lower_ok and upper_ok:
            filtered.append(item)
    return filtered


def _completed_records(items: list[dict[str, Any]]) -> list[dict[str, Any]]:
    return [item for item in items if bool(item.get("is_completed"))]


def _attendance_rate(items: list[dict[str, Any]]) -> float:
    if not items:
        return 1.0
    attended = sum(1 for item in items if int(item.get("status", 0)) in PRESENT_LIKE_STATUSES)
    return attended / len(items)


def _absence_count(items: list[dict[str, Any]]) -> int:
    return sum(1 for item in items if int(item.get("status", 0)) == ATTENDANCE_STATUS_ABSENT)


def _average_score(items: list[dict[str, Any]]) -> float:
    if not items:
        return 7.0
    return sum(_safe_float(item.get("total_score"), 0.0) for item in items) / len(items)


def _homework_completion_rate(items: list[dict[str, Any]]) -> float:
    if not items:
        return 1.0
    completed = sum(1 for item in items if bool(item.get("homework_completed")))
    return completed / len(items)


def _derive_label(future_attendance: list[dict[str, Any]], future_academic: list[dict[str, Any]]) -> str:
    if future_attendance and _attendance_rate(future_attendance) < 0.80:
        return LABEL_AT_RISK
    if future_academic and _average_score(future_academic) < 5.0:
        return LABEL_AT_RISK
    if future_academic and _homework_completion_rate(future_academic) < 0.60:
        return LABEL_AT_RISK
    return LABEL_NOT_AT_RISK


def _active_enrollment_class_ids(enrollments: list[dict[str, Any]], snapshot_at: datetime) -> list[str]:
    class_ids: list[str] = []
    seen: set[str] = set()
    for enrollment in enrollments:
        class_id = enrollment["class_id"]
        if class_id in seen:
            continue
        if _is_enrollment_active_at(enrollment, snapshot_at):
            seen.add(class_id)
            class_ids.append(class_id)
    return class_ids


def _weekly_lesson_load(
    lessons_by_class: dict[str, list[dict[str, Any]]],
    class_ids: list[str],
    start: datetime,
    end: datetime,
    observation_window_days: int,
) -> float:
    if not class_ids or observation_window_days <= 0:
        return 0.0

    lesson_count = 0
    for class_id in class_ids:
        for lesson in lessons_by_class.get(class_id, []):
            lesson_at = _ensure_datetime(lesson.get("date_start"))
            if lesson_at is None:
                continue
            if lesson_at >= start and lesson_at < end:
                lesson_count += 1

    weeks = observation_window_days / 7.0
    return lesson_count / weeks if weeks else 0.0


def _approved_leave_count(
    items: list[dict[str, Any]],
    start: datetime,
    end: datetime,
    class_id: str,
) -> int:
    total = 0
    for item in items:
        item_at = _ensure_datetime(item.get("at"))
        if item_at is None:
            continue
        if item.get("class_id") not in ("", None, class_id):
            continue
        if item_at >= start and item_at < end:
            total += 1
    return total


def _days_since_last_lesson(
    lessons_by_class: dict[str, list[dict[str, Any]]],
    class_ids: list[str],
    snapshot_at: datetime,
    observation_window_days: int,
) -> float:
    last_lesson_at: datetime | None = None
    for class_id in class_ids:
        for lesson in lessons_by_class.get(class_id, []):
            lesson_at = _ensure_datetime(lesson.get("date_start"))
            if lesson_at is None or lesson_at >= snapshot_at:
                continue
            if last_lesson_at is None or lesson_at > last_lesson_at:
                last_lesson_at = lesson_at
    if last_lesson_at is None:
        return float(observation_window_days + 1)
    return (snapshot_at - last_lesson_at).total_seconds() / 86400.0


def _make_snapshot_id(student_id: str, class_id: str, snapshot_at: datetime) -> str:
    return f"{student_id}:{class_id}:{snapshot_at.strftime('%Y%m%d')}"


def _pick_inference_snapshot_at(
    class_lessons: list[dict[str, Any]],
    enrollment: dict[str, Any],
    now: datetime,
) -> datetime | None:
    lesson_times = [
        lesson_at
        for lesson_at in (_ensure_datetime(lesson.get("date_start")) for lesson in class_lessons)
        if lesson_at is not None
    ]
    past_or_current = [lesson_at for lesson_at in lesson_times if lesson_at <= now]
    if past_or_current:
        return max(past_or_current)
    if lesson_times:
        return min(lesson_times)

    class_start = _ensure_datetime(enrollment.get("class_start_date"))
    if class_start and class_start > now:
        return class_start
    return now
