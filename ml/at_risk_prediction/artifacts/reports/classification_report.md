# Báo cáo huấn luyện và đánh giá mô hình AT_RISK

- Thời điểm sinh báo cáo: `2026-04-28T07:58:24.665800Z`
- Dataset: `at_risk_dataset_db`
- Nguồn dữ liệu: `database`
- Seed: `42`
- Tỉ lệ test: `0.2`

## 1. Tóm tắt dataset

- Số dòng: `36`
- Số học viên: `6`
- Số lớp: `1`
- Feature columns: `attendance_rate_28d, absence_count_28d, average_total_score_28d, homework_completion_rate_28d, active_enrollment_count_28d, weekly_lesson_load_28d, approved_leave_count_28d, days_since_last_lesson`
- Phân phối nhãn: `AT_RISK=21`, `NOT_AT_RISK=15`

## 2. Bảng so sánh mô hình

| Mô hình | Accuracy | Precision | Recall | F1 | TP | FP | FN | TN |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `rule_based` | 0.3750 | 0.0000 | 0.0000 | 0.0000 | 0 | 0 | 5 | 3 |
| `logistic_regression` | 0.7500 | 1.0000 | 0.6000 | 0.7500 | 3 | 0 | 2 | 3 |
| `random_forest` | 0.7500 | 1.0000 | 0.6000 | 0.7500 | 3 | 0 | 2 | 3 |

## 3. Nhận xét chính

- Mô hình có `F1-score` cao nhất: `logistic_regression`
- Mô hình có `Recall` cao nhất: `logistic_regression`
- Mô hình được dùng để vẽ confusion matrix: `logistic_regression`

## 4. Cấu hình baseline rule-based

```json
{
  "name": "rule_based",
  "thresholds": {
    "attendance_rate_28d": 0.8,
    "average_total_score_28d": 5.0,
    "homework_completion_rate_28d": 0.6
  },
  "decision_threshold": 0.5,
  "weights": {
    "attendance": 0.4,
    "score": 0.35,
    "homework": 0.25
  },
  "generated_at": "2026-04-28T07:58:24.399676Z"
}
```

## 5. Phân tích chi tiết từng mô hình

### 5.1. rule_based

- Accuracy: `0.3750`
- Precision: `0.0000`
- Recall: `0.0000`
- F1-score: `0.0000`
- Confusion matrix: `[[3, 0], [5, 0]]`

### 5.2. logistic_regression

- Accuracy: `0.7500`
- Precision: `1.0000`
- Recall: `0.6000`
- F1-score: `0.7500`
- Confusion matrix: `[[3, 0], [2, 3]]`

### 5.3. random_forest

- Accuracy: `0.7500`
- Precision: `1.0000`
- Recall: `0.6000`
- F1-score: `0.7500`
- Confusion matrix: `[[3, 0], [2, 3]]`

## 6. Kết luận chọn mô hình chính

- Mô hình được chọn: `logistic_regression`
- Các mô hình có metric tương đương: `logistic_regression, random_forest`
- Lý do chọn: Các mô hình có metric tương đương trên tập kiểm thử hiện tại; ưu tiên Logistic Regression vì vẫn là mô hình học máy chính thức, nhẹ khi huấn luyện/suy luận và có khả năng giải thích tốt hơn Random Forest.
- Tiêu chí ưu tiên: `Recall -> F1-score -> Precision -> Accuracy -> explainability -> lightweight`
