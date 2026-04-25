from __future__ import annotations

import json
from datetime import datetime
from pathlib import Path
from typing import Any

import joblib
import pandas as pd

from .config import MODELS_DIR, REPORTS_DIR
from .dataset_schema import FEATURE_COLUMNS, LABEL_AT_RISK, LABEL_NOT_AT_RISK

POSITIVE_LABEL = LABEL_AT_RISK
NEGATIVE_LABEL = LABEL_NOT_AT_RISK
RISK_BAND_THRESHOLDS = (
    ("HIGH", 0.8),
    ("MEDIUM", 0.6),
    ("LOW", 0.0),
)

MODEL_TIE_BREAK_PRIORITY = {
    "logistic_regression": 3,
    "random_forest": 2,
    "rule_based": 1,
}

FEATURE_REASON_CONFIG = {
    "attendance_rate_28d": {
        "risk_when": "lower",
        "high": 0.70,
        "medium": 0.80,
        "label": "Tỷ lệ chuyên cần thấp",
        "detail": "Tỷ lệ tham dự trong 28 ngày gần đây thấp hơn ngưỡng an toàn.",
    },
    "absence_count_28d": {
        "risk_when": "higher",
        "high": 4.0,
        "medium": 2.0,
        "label": "Số buổi vắng tăng cao",
        "detail": "Học viên vắng nhiều buổi trong 28 ngày gần đây.",
    },
    "average_total_score_28d": {
        "risk_when": "lower",
        "high": 5.0,
        "medium": 6.0,
        "label": "Điểm trung bình thấp",
        "detail": "Điểm tổng kết trung bình đang thấp hơn mức kỳ vọng.",
    },
    "homework_completion_rate_28d": {
        "risk_when": "lower",
        "high": 0.60,
        "medium": 0.75,
        "label": "Tỷ lệ hoàn thành bài tập thấp",
        "detail": "Mức độ hoàn thành bài tập chưa đạt ngưỡng an toàn.",
    },
    "active_enrollment_count_28d": {
        "risk_when": "higher",
        "high": 3.0,
        "medium": 2.0,
        "label": "Tải học phân tán",
        "detail": "Học viên theo học nhiều lớp cùng lúc, có thể ảnh hưởng khả năng theo bài.",
    },
    "weekly_lesson_load_28d": {
        "risk_when": "higher",
        "high": 5.0,
        "medium": 3.5,
        "label": "Cường độ học cao",
        "detail": "Tải buổi học theo tuần cao, dễ gây quá tải.",
    },
    "approved_leave_count_28d": {
        "risk_when": "higher",
        "high": 3.0,
        "medium": 1.0,
        "label": "Tần suất xin phép tăng",
        "detail": "Số đơn xin phép được duyệt trong 28 ngày gần đây tăng cao.",
    },
    "days_since_last_lesson": {
        "risk_when": "higher",
        "high": 10.0,
        "medium": 5.0,
        "label": "Gián đoạn học tập",
        "detail": "Khoảng cách từ buổi học gần nhất đến snapshot đang kéo dài.",
    },
}


def load_json(path: str | Path) -> dict[str, Any]:
    return json.loads(Path(path).read_text(encoding="utf-8"))


def select_primary_model(
    metrics_payload: dict[str, Any],
) -> dict[str, Any]:
    metrics_by_model = metrics_payload["models"]
    ranked = sorted(
        metrics_by_model.items(),
        key=lambda item: (
            float(item[1]["recall"]),
            float(item[1]["f1"]),
            float(item[1]["precision"]),
            float(item[1]["accuracy"]),
            MODEL_TIE_BREAK_PRIORITY.get(item[0], 0),
        ),
        reverse=True,
    )
    selected_name, selected_metrics = ranked[0]
    tied_models = [
        model_name
        for model_name, metrics in metrics_by_model.items()
        if _same_metrics(metrics, selected_metrics)
    ]

    if len(tied_models) > 1 and selected_name == "logistic_regression":
        rationale = (
            "Các mô hình có metric tương đương trên tập kiểm thử hiện tại; "
            "ưu tiên Logistic Regression vì vẫn là mô hình học máy chính thức, "
            "nhẹ khi huấn luyện/suy luận và có khả năng giải thích tốt hơn Random Forest."
        )
    else:
        rationale = (
            "Mô hình được chọn dựa trên Recall, F1-score, Precision và Accuracy theo thứ tự ưu tiên."
        )

    return {
        "selected_model": selected_name,
        "selected_metrics": selected_metrics,
        "tied_models": tied_models,
        "selection_criteria": ["recall", "f1", "precision", "accuracy", "explainability", "lightweight"],
        "rationale": rationale,
        "selected_at": _utc_now(),
    }


def save_primary_model_selection(
    metadata_path: str | Path | None = None,
    metrics_path: str | Path | None = None,
) -> dict[str, Any]:
    metadata_file = Path(metadata_path or MODELS_DIR / "model_metadata.json")
    metrics_file = Path(metrics_path or REPORTS_DIR / "metrics.json")
    metadata = load_json(metadata_file)
    metrics_payload = load_json(metrics_file)
    selection = select_primary_model(metrics_payload)
    metadata["selection"] = selection
    metadata_file.write_text(json.dumps(metadata, indent=2, ensure_ascii=False), encoding="utf-8")
    return selection


def append_selection_conclusion_to_report(
    report_path: str | Path | None = None,
    selection: dict[str, Any] | None = None,
) -> Path:
    report_file = Path(report_path or REPORTS_DIR / "classification_report.md")
    if selection is None:
        selection = load_json(MODELS_DIR / "model_metadata.json").get("selection", {})

    existing_content = report_file.read_text(encoding="utf-8").rstrip()
    marker = "\n## 6. Kết luận chọn mô hình chính\n"
    if marker in existing_content:
        existing_content = existing_content.split(marker, 1)[0].rstrip()

    content = existing_content + "\n\n"
    content += "## 6. Kết luận chọn mô hình chính\n\n"
    content += f"- Mô hình được chọn: `{selection.get('selected_model', 'unknown')}`\n"
    tied_models = selection.get("tied_models") or []
    if tied_models:
        content += f"- Các mô hình có metric tương đương: `{', '.join(tied_models)}`\n"
    content += f"- Lý do chọn: {selection.get('rationale', '')}\n"
    content += "- Tiêu chí ưu tiên: `Recall -> F1-score -> Precision -> Accuracy -> explainability -> lightweight`\n"
    report_file.write_text(content, encoding="utf-8")
    return report_file


def predict_dataset(
    dataset: pd.DataFrame,
    model_name: str | None = None,
    metadata_path: str | Path | None = None,
    output_path: str | Path | None = None,
    source: str | None = None,
) -> dict[str, Any]:
    metadata_file = Path(metadata_path or MODELS_DIR / "model_metadata.json")
    metadata = load_json(metadata_file)
    selection = metadata.get("selection") or {}
    selected_model_name = model_name or selection.get("selected_model") or "logistic_regression"

    if dataset.empty:
        payload = _build_prediction_payload(
            predictions=[],
            model_name=selected_model_name,
            metadata=metadata,
            source=source or metadata.get("source", "unknown"),
            note="Khong co dong dataset nao de sinh prediction.",
        )
        output_file = Path(output_path or REPORTS_DIR / "latest_predictions.json")
        output_file.write_text(json.dumps(payload, indent=2, ensure_ascii=False), encoding="utf-8")
        return {"payload": payload, "output_path": str(output_file)}

    prediction_rows = _predict_rows(dataset, selected_model_name, metadata)
    payload = _build_prediction_payload(
        predictions=prediction_rows,
        model_name=selected_model_name,
        metadata=metadata,
        source=source or metadata.get("source", "unknown"),
        note=None,
    )
    output_file = Path(output_path or REPORTS_DIR / "latest_predictions.json")
    output_file.write_text(json.dumps(payload, indent=2, ensure_ascii=False), encoding="utf-8")
    return {"payload": payload, "output_path": str(output_file)}


def _predict_rows(
    dataset: pd.DataFrame,
    model_name: str,
    metadata: dict[str, Any],
) -> list[dict[str, Any]]:
    features = dataset[FEATURE_COLUMNS].astype(float)
    scores, labels = _run_model_prediction(model_name, features, metadata)

    rows: list[dict[str, Any]] = []
    for index, (_, row) in enumerate(dataset.iterrows()):
        score = float(scores[index])
        label = labels[index]
        feature_summary, top_features, primary_reason = _build_feature_explanations(row)
        rows.append(
            {
                "snapshot_id": row["snapshot_id"],
                "student_id": row["student_id"],
                "student_code": row.get("student_code", ""),
                "student_name": row.get("student_name", ""),
                "class_id": row["class_id"],
                "class_code": row.get("class_code", ""),
                "class_name": row.get("class_name", ""),
                "snapshot_at": row.get("snapshot_at", ""),
                "risk_label": label,
                "risk_score": round(score, 6),
                "risk_band": _derive_risk_band(score),
                "primary_reason": primary_reason,
                "top_features": top_features,
                "feature_summary": feature_summary,
                "model_version": metadata.get("selection", {}).get("selected_model", model_name),
                "predicted_at": _utc_now(),
            }
        )
    return rows


def _run_model_prediction(
    model_name: str,
    features: pd.DataFrame,
    metadata: dict[str, Any],
) -> tuple[list[float], list[str]]:
    model_info = metadata["models"][model_name]
    artifact_path = Path(model_info["artifact"])
    if model_name == "rule_based":
        rule_model = load_json(artifact_path)
        scores = _predict_rule_based_proba(rule_model, features)
        labels = [
            POSITIVE_LABEL if score >= float(rule_model["decision_threshold"]) else NEGATIVE_LABEL
            for score in scores
        ]
        return scores, labels

    artifact = joblib.load(artifact_path)
    model = artifact["model"]
    probabilities = model.predict_proba(features)
    classes = list(model.classes_)
    positive_index = classes.index(POSITIVE_LABEL)
    scores = probabilities[:, positive_index].tolist()
    labels = [POSITIVE_LABEL if score >= 0.5 else NEGATIVE_LABEL for score in scores]
    return scores, labels


def _predict_rule_based_proba(model: dict[str, Any], features: pd.DataFrame) -> list[float]:
    thresholds = model["thresholds"]
    weights = model["weights"]
    scores: list[float] = []
    for _, row in features.iterrows():
        risk_score = 0.0
        if float(row["attendance_rate_28d"]) < float(thresholds["attendance_rate_28d"]):
            risk_score += float(weights["attendance"])
        if float(row["average_total_score_28d"]) < float(thresholds["average_total_score_28d"]):
            risk_score += float(weights["score"])
        if float(row["homework_completion_rate_28d"]) < float(thresholds["homework_completion_rate_28d"]):
            risk_score += float(weights["homework"])
        scores.append(min(1.0, risk_score))
    return scores


def _build_feature_explanations(row: pd.Series) -> tuple[dict[str, float], list[dict[str, Any]], str]:
    feature_summary = {column: round(float(row[column]), 6) for column in FEATURE_COLUMNS}
    feature_signals: list[dict[str, Any]] = []

    for feature_name in FEATURE_COLUMNS:
        config = FEATURE_REASON_CONFIG[feature_name]
        value = float(row[feature_name])
        severity = _compute_feature_severity(value, config)
        if severity <= 0:
            continue
        feature_signals.append(
            {
                "feature": feature_name,
                "value": round(value, 6),
                "severity": round(severity, 6),
                "label": config["label"],
                "detail": config["detail"],
            }
        )

    feature_signals.sort(key=lambda item: item["severity"], reverse=True)
    if feature_signals:
        primary_reason = feature_signals[0]["label"]
    else:
        primary_reason = "Chưa ghi nhận tín hiệu rủi ro nổi bật trong cửa sổ quan sát."

    return feature_summary, feature_signals[:3], primary_reason


def _compute_feature_severity(value: float, config: dict[str, Any]) -> float:
    if config["risk_when"] == "lower":
        if value < float(config["high"]):
            return 1.0
        if value < float(config["medium"]):
            return 0.5
        return 0.0

    if value >= float(config["high"]):
        return 1.0
    if value >= float(config["medium"]):
        return 0.5
    return 0.0


def _derive_risk_band(score: float) -> str:
    for band, threshold in RISK_BAND_THRESHOLDS:
        if score >= threshold:
            return band
    return "LOW"


def _build_prediction_payload(
    predictions: list[dict[str, Any]],
    model_name: str,
    metadata: dict[str, Any],
    source: str,
    note: str | None,
) -> dict[str, Any]:
    return {
        "generated_at": _utc_now(),
        "source": source,
        "model_version": model_name,
        "selection": metadata.get("selection", {}),
        "feature_columns": metadata.get("feature_columns", FEATURE_COLUMNS),
        "prediction_count": len(predictions),
        "note": note,
        "predictions": predictions,
    }


def _same_metrics(left: dict[str, Any], right: dict[str, Any], tolerance: float = 1e-9) -> bool:
    for metric_name in ("recall", "f1", "precision", "accuracy"):
        if abs(float(left[metric_name]) - float(right[metric_name])) > tolerance:
            return False
    return True


def _utc_now() -> str:
    return datetime.utcnow().isoformat() + "Z"
