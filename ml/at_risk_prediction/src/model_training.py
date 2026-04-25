from __future__ import annotations

import json
from dataclasses import dataclass
from datetime import datetime
from pathlib import Path
from typing import Any

import joblib
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns
from sklearn.ensemble import RandomForestClassifier
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import (
    accuracy_score,
    classification_report,
    confusion_matrix,
    f1_score,
    precision_score,
    recall_score,
)
from sklearn.pipeline import Pipeline
from sklearn.preprocessing import StandardScaler

from .config import FIGURES_DIR, MODELS_DIR, REPORTS_DIR
from .data_pipeline import build_dataset_summary, save_dataset_artifacts, split_dataset
from .dataset_schema import (
    DEFAULT_DATASET_DEFINITION,
    FEATURE_COLUMNS,
    LABEL_AT_RISK,
    LABEL_NOT_AT_RISK,
)

POSITIVE_LABEL = LABEL_AT_RISK
NEGATIVE_LABEL = LABEL_NOT_AT_RISK


@dataclass(frozen=True)
class TrainingConfig:
    dataset_name: str
    source: str
    test_size: float = 0.2
    seed: int = 42


def train_and_evaluate_models(
    dataset: pd.DataFrame,
    config: TrainingConfig,
) -> dict[str, Any]:
    _validate_dataset(dataset)

    summary = build_dataset_summary(
        dataset,
        DEFAULT_DATASET_DEFINITION,
        source=config.source,
    )
    dataset_paths = save_dataset_artifacts(
        dataset=dataset,
        summary=summary,
        dataset_name=config.dataset_name,
        save_split=True,
        test_size=config.test_size,
        seed=config.seed,
    )

    train_df, test_df = split_dataset(
        dataset,
        test_size=config.test_size,
        seed=config.seed,
    )

    X_train = train_df[FEATURE_COLUMNS].astype(float)
    y_train = train_df["label"].astype(str)
    X_test = test_df[FEATURE_COLUMNS].astype(float)
    y_test = test_df["label"].astype(str)

    if y_train.nunique() < 2:
        raise ValueError(
            "Tap train chi co mot nhan. Khong the huan luyen Logistic Regression hoac Random Forest."
        )

    rule_model = _build_rule_based_model()
    logistic_model = _build_logistic_regression_model(config.seed)
    random_forest_model = _build_random_forest_model(config.seed)

    logistic_model.fit(X_train, y_train)
    random_forest_model.fit(X_train, y_train)

    model_bundle = {
        "rule_based": {
            "model": rule_model,
            "pred": _predict_rule_based(rule_model, X_test),
            "score": _predict_rule_based_proba(rule_model, X_test),
            "model_type": "rule_based",
        },
        "logistic_regression": {
            "model": logistic_model,
            "pred": logistic_model.predict(X_test),
            "score": _positive_label_scores(logistic_model, X_test),
            "model_type": "logistic_regression",
        },
        "random_forest": {
            "model": random_forest_model,
            "pred": random_forest_model.predict(X_test),
            "score": _positive_label_scores(random_forest_model, X_test),
            "model_type": "random_forest",
        },
    }

    metrics_by_model: dict[str, Any] = {}
    for model_name, bundle in model_bundle.items():
        metrics_by_model[model_name] = _evaluate_predictions(
            y_true=y_test,
            y_pred=bundle["pred"],
            y_score=bundle["score"],
            model_type=bundle["model_type"],
        )

    selected_for_confusion = _pick_best_model(metrics_by_model)

    MODELS_DIR.mkdir(parents=True, exist_ok=True)
    REPORTS_DIR.mkdir(parents=True, exist_ok=True)
    FIGURES_DIR.mkdir(parents=True, exist_ok=True)

    rule_based_path = MODELS_DIR / "rule_based.json"
    rule_based_path.write_text(
        json.dumps(rule_model, indent=2, ensure_ascii=False),
        encoding="utf-8",
    )

    logistic_path = MODELS_DIR / "logistic_regression.joblib"
    joblib.dump(
        {
            "model": logistic_model,
            "feature_columns": FEATURE_COLUMNS,
            "positive_label": POSITIVE_LABEL,
            "negative_label": NEGATIVE_LABEL,
        },
        logistic_path,
    )

    random_forest_path = MODELS_DIR / "random_forest.joblib"
    joblib.dump(
        {
            "model": random_forest_model,
            "feature_columns": FEATURE_COLUMNS,
            "positive_label": POSITIVE_LABEL,
            "negative_label": NEGATIVE_LABEL,
        },
        random_forest_path,
    )

    model_metadata = {
        "generated_at": _utc_now(),
        "dataset_name": config.dataset_name,
        "source": config.source,
        "seed": config.seed,
        "test_size": config.test_size,
        "feature_columns": FEATURE_COLUMNS,
        "label_mapping": {
            "positive_label": POSITIVE_LABEL,
            "negative_label": NEGATIVE_LABEL,
        },
        "models": {
            "rule_based": {
                "artifact": str(rule_based_path),
                "type": "rule_based",
            },
            "logistic_regression": {
                "artifact": str(logistic_path),
                "type": "sklearn_pipeline",
            },
            "random_forest": {
                "artifact": str(random_forest_path),
                "type": "sklearn_classifier",
            },
        },
        "evaluation": {
            "selected_for_confusion_matrix": selected_for_confusion,
            "best_by_f1": _pick_best_model(metrics_by_model, by="f1"),
            "best_by_recall": _pick_best_model(metrics_by_model, by="recall"),
        },
    }
    metadata_path = MODELS_DIR / "model_metadata.json"
    metadata_path.write_text(
        json.dumps(model_metadata, indent=2, ensure_ascii=False),
        encoding="utf-8",
    )

    metrics_payload = {
        "generated_at": _utc_now(),
        "dataset_summary": summary,
        "training_config": {
            "dataset_name": config.dataset_name,
            "source": config.source,
            "test_size": config.test_size,
            "seed": config.seed,
            "train_size": int(len(train_df)),
            "test_size_rows": int(len(test_df)),
        },
        "models": metrics_by_model,
    }
    metrics_path = REPORTS_DIR / "metrics.json"
    metrics_path.write_text(
        json.dumps(metrics_payload, indent=2, ensure_ascii=False),
        encoding="utf-8",
    )

    report_path = REPORTS_DIR / "classification_report.md"
    report_path.write_text(
        _build_markdown_report(
            dataset_summary=summary,
            training_config=config,
            metrics_by_model=metrics_by_model,
            selected_for_confusion=selected_for_confusion,
            rule_based=rule_model,
        ),
        encoding="utf-8",
    )

    confusion_path = FIGURES_DIR / "confusion_matrix.png"
    _plot_confusion_matrix(
        y_true=y_test,
        y_pred=model_bundle[selected_for_confusion]["pred"],
        model_name=selected_for_confusion,
        output_path=confusion_path,
    )

    feature_importance_path = FIGURES_DIR / "feature_importance.png"
    _plot_feature_importance(
        logistic_model=logistic_model,
        random_forest_model=random_forest_model,
        output_path=feature_importance_path,
    )

    return {
        "dataset_paths": dataset_paths,
        "metrics_json": str(metrics_path),
        "classification_report_md": str(report_path),
        "confusion_matrix_png": str(confusion_path),
        "feature_importance_png": str(feature_importance_path),
        "rule_based_json": str(rule_based_path),
        "logistic_regression_joblib": str(logistic_path),
        "random_forest_joblib": str(random_forest_path),
        "model_metadata_json": str(metadata_path),
        "selected_for_confusion_matrix": selected_for_confusion,
        "best_by_f1": model_metadata["evaluation"]["best_by_f1"],
        "best_by_recall": model_metadata["evaluation"]["best_by_recall"],
    }


def _validate_dataset(dataset: pd.DataFrame) -> None:
    if dataset.empty:
        raise ValueError("Dataset rong. Khong the huan luyen mo hinh.")

    missing = [column for column in FEATURE_COLUMNS + ["label"] if column not in dataset.columns]
    if missing:
        raise ValueError(f"Dataset thieu cot bat buoc: {', '.join(missing)}")

    if dataset["label"].nunique() < 2:
        raise ValueError("Dataset chi co mot nhan. Khong du dieu kien cho bai toan classification.")


def _build_rule_based_model() -> dict[str, Any]:
    return {
        "name": "rule_based",
        "thresholds": {
            "attendance_rate_28d": 0.80,
            "average_total_score_28d": 5.00,
            "homework_completion_rate_28d": 0.60,
        },
        "decision_threshold": 0.50,
        "weights": {
            "attendance": 0.40,
            "score": 0.35,
            "homework": 0.25,
        },
        "generated_at": _utc_now(),
    }


def _predict_rule_based(model: dict[str, Any], features: pd.DataFrame) -> list[str]:
    return [
        POSITIVE_LABEL if score >= float(model["decision_threshold"]) else NEGATIVE_LABEL
        for score in _predict_rule_based_proba(model, features)
    ]


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


def _build_logistic_regression_model(seed: int) -> Pipeline:
    return Pipeline(
        steps=[
            ("scaler", StandardScaler()),
            (
                "classifier",
                LogisticRegression(
                    max_iter=1000,
                    random_state=seed,
                    class_weight="balanced",
                ),
            ),
        ]
    )


def _build_random_forest_model(seed: int) -> RandomForestClassifier:
    return RandomForestClassifier(
        n_estimators=300,
        max_depth=8,
        min_samples_leaf=2,
        random_state=seed,
        class_weight="balanced_subsample",
    )


def _positive_label_scores(model: Any, features: pd.DataFrame) -> list[float]:
    probabilities = model.predict_proba(features)
    classes = list(model.classes_)
    if POSITIVE_LABEL not in classes:
        raise ValueError(f"Khong tim thay positive label `{POSITIVE_LABEL}` trong model.classes_.")
    positive_index = classes.index(POSITIVE_LABEL)
    return probabilities[:, positive_index].tolist()


def _evaluate_predictions(
    y_true: pd.Series,
    y_pred: list[str] | pd.Series,
    y_score: list[float] | Any,
    model_type: str,
) -> dict[str, Any]:
    labels = [NEGATIVE_LABEL, POSITIVE_LABEL]
    matrix = confusion_matrix(y_true, y_pred, labels=labels)
    report = classification_report(
        y_true,
        y_pred,
        labels=labels,
        output_dict=True,
        zero_division=0,
    )
    return {
        "model_type": model_type,
        "accuracy": round(float(accuracy_score(y_true, y_pred)), 6),
        "precision": round(float(precision_score(y_true, y_pred, pos_label=POSITIVE_LABEL, zero_division=0)), 6),
        "recall": round(float(recall_score(y_true, y_pred, pos_label=POSITIVE_LABEL, zero_division=0)), 6),
        "f1": round(float(f1_score(y_true, y_pred, pos_label=POSITIVE_LABEL, zero_division=0)), 6),
        "support": int(len(y_true)),
        "confusion_matrix": {
            "labels": labels,
            "matrix": matrix.tolist(),
            "tn": int(matrix[0][0]),
            "fp": int(matrix[0][1]),
            "fn": int(matrix[1][0]),
            "tp": int(matrix[1][1]),
        },
        "prediction_score_summary": {
            "min": round(float(min(y_score)), 6),
            "max": round(float(max(y_score)), 6),
            "avg": round(float(sum(y_score) / len(y_score)), 6),
        },
        "classification_report": report,
    }


def _pick_best_model(metrics_by_model: dict[str, Any], by: str = "f1") -> str:
    ranking_priority = ("recall", "f1", "precision", "accuracy")
    sorted_models = sorted(
        metrics_by_model.items(),
        key=lambda item: tuple(float(item[1][metric]) for metric in ranking_priority),
        reverse=True,
    )
    if by in ranking_priority:
        sorted_models = sorted(
            metrics_by_model.items(),
            key=lambda item: (
                float(item[1][by]),
                float(item[1]["f1"]),
                float(item[1]["recall"]),
                float(item[1]["precision"]),
                float(item[1]["accuracy"]),
            ),
            reverse=True,
        )
    return sorted_models[0][0]


def _build_markdown_report(
    dataset_summary: dict[str, Any],
    training_config: TrainingConfig,
    metrics_by_model: dict[str, Any],
    selected_for_confusion: str,
    rule_based: dict[str, Any],
) -> str:
    lines = [
        "# Báo cáo huấn luyện và đánh giá mô hình AT_RISK",
        "",
        f"- Thời điểm sinh báo cáo: `{_utc_now()}`",
        f"- Dataset: `{training_config.dataset_name}`",
        f"- Nguồn dữ liệu: `{training_config.source}`",
        f"- Seed: `{training_config.seed}`",
        f"- Tỉ lệ test: `{training_config.test_size}`",
        "",
        "## 1. Tóm tắt dataset",
        "",
        f"- Số dòng: `{dataset_summary['row_count']}`",
        f"- Số học viên: `{dataset_summary['student_count']}`",
        f"- Số lớp: `{dataset_summary['class_count']}`",
        f"- Feature columns: `{', '.join(dataset_summary['feature_columns'])}`",
        f"- Phân phối nhãn: `AT_RISK={dataset_summary['label_distribution'][LABEL_AT_RISK]}`, `NOT_AT_RISK={dataset_summary['label_distribution'][LABEL_NOT_AT_RISK]}`",
        "",
        "## 2. Bảng so sánh mô hình",
        "",
        "| Mô hình | Accuracy | Precision | Recall | F1 | TP | FP | FN | TN |",
        "|---|---:|---:|---:|---:|---:|---:|---:|---:|",
    ]

    for model_name, metrics in metrics_by_model.items():
        matrix = metrics["confusion_matrix"]
        lines.append(
            f"| `{model_name}` | {metrics['accuracy']:.4f} | {metrics['precision']:.4f} | "
            f"{metrics['recall']:.4f} | {metrics['f1']:.4f} | {matrix['tp']} | {matrix['fp']} | "
            f"{matrix['fn']} | {matrix['tn']} |"
        )

    best_by_f1 = _pick_best_model(metrics_by_model, by="f1")
    best_by_recall = _pick_best_model(metrics_by_model, by="recall")

    lines.extend(
        [
            "",
            "## 3. Nhận xét chính",
            "",
            f"- Mô hình có `F1-score` cao nhất: `{best_by_f1}`",
            f"- Mô hình có `Recall` cao nhất: `{best_by_recall}`",
            f"- Mô hình được dùng để vẽ confusion matrix: `{selected_for_confusion}`",
            "",
            "## 4. Cấu hình baseline rule-based",
            "",
            "```json",
            json.dumps(rule_based, indent=2, ensure_ascii=False),
            "```",
            "",
            "## 5. Phân tích chi tiết từng mô hình",
            "",
        ]
    )

    for model_name, metrics in metrics_by_model.items():
        matrix = metrics["confusion_matrix"]
        lines.extend(
            [
                f"### 5.{list(metrics_by_model.keys()).index(model_name) + 1}. {model_name}",
                "",
                f"- Accuracy: `{metrics['accuracy']:.4f}`",
                f"- Precision: `{metrics['precision']:.4f}`",
                f"- Recall: `{metrics['recall']:.4f}`",
                f"- F1-score: `{metrics['f1']:.4f}`",
                f"- Confusion matrix: `[[{matrix['tn']}, {matrix['fp']}], [{matrix['fn']}, {matrix['tp']}]]`",
                "",
            ]
        )
    return "\n".join(lines) + "\n"


def _plot_confusion_matrix(
    y_true: pd.Series,
    y_pred: list[str] | pd.Series,
    model_name: str,
    output_path: Path,
) -> None:
    matrix = confusion_matrix(y_true, y_pred, labels=[NEGATIVE_LABEL, POSITIVE_LABEL])
    plt.figure(figsize=(6, 5))
    sns.heatmap(
        matrix,
        annot=True,
        fmt="d",
        cmap="Blues",
        xticklabels=[NEGATIVE_LABEL, POSITIVE_LABEL],
        yticklabels=[NEGATIVE_LABEL, POSITIVE_LABEL],
    )
    plt.title(f"Confusion Matrix - {model_name}")
    plt.xlabel("Predicted label")
    plt.ylabel("True label")
    plt.tight_layout()
    plt.savefig(output_path, dpi=200)
    plt.close()


def _plot_feature_importance(
    logistic_model: Pipeline,
    random_forest_model: RandomForestClassifier,
    output_path: Path,
) -> None:
    logistic_classifier: LogisticRegression = logistic_model.named_steps["classifier"]
    logistic_importance = pd.Series(
        logistic_classifier.coef_[0],
        index=FEATURE_COLUMNS,
    ).sort_values()
    random_forest_importance = pd.Series(
        random_forest_model.feature_importances_,
        index=FEATURE_COLUMNS,
    ).sort_values()

    fig, axes = plt.subplots(1, 2, figsize=(14, 6))

    logistic_importance.plot.barh(ax=axes[0], color="#4c78a8")
    axes[0].set_title("Logistic Regression Coefficients")
    axes[0].set_xlabel("Coefficient")
    axes[0].set_ylabel("Feature")

    random_forest_importance.plot.barh(ax=axes[1], color="#f58518")
    axes[1].set_title("Random Forest Feature Importance")
    axes[1].set_xlabel("Importance")
    axes[1].set_ylabel("Feature")

    plt.tight_layout()
    plt.savefig(output_path, dpi=200)
    plt.close(fig)


def _utc_now() -> str:
    return datetime.utcnow().isoformat() + "Z"
