# Predictive Go -> Python Migration Audit

**Ngày rà soát:** 2026-04-25  
**Phạm vi:** `Task K0` trong [PROJECT_TASKS.md](/Users/hant/golang/doan/PROJECT_TASKS.md)

## 1. Mục tiêu của K0

Khóa rõ các thành phần của nhánh predictive analytics hiện đang viết bằng Go để chuẩn bị chuyển sang kiến trúc mới:

- Python chịu trách nhiệm:
  - xây dựng dataset,
  - huấn luyện mô hình,
  - đánh giá metric,
  - sinh prediction artifact,
  - sinh model metadata/report cho báo cáo.
- Go backend chịu trách nhiệm:
  - đọc artifact do Python sinh ra,
  - phục vụ API cho admin/student,
  - giữ ổn định contract FE-BE và route hệ thống.

`Task K0` **không xóa code ngay**.  
Nó chỉ audit và chốt hướng xử lý để các bước `K1 -> K8` triển khai nhất quán.

---

## 2. Kết luận kỹ thuật ngắn

Nhánh predictive hiện tại của repo được chia thành 4 cụm:

1. **Cụm định nghĩa dữ liệu và logic train/evaluate bằng Go**
   - đang dùng để train nhẹ trong runtime hoặc CLI.
   - đây là cụm **sẽ bị loại bỏ dần** khi project Python hoàn chỉnh.

2. **Cụm API/usecase admin predictive**
   - hiện đang phụ thuộc trực tiếp vào `AtRiskService` của Go.
   - cụm này **không xóa route**, nhưng phải **refactor nguồn dữ liệu** sang artifact Python.

3. **Cụm student portal**
   - route `/api/v1/student/at-risk` đã được FE dùng thật.
   - cụm này **giữ lại về mặt business contract**, nhưng phải **đổi dependency** từ runtime train Go sang artifact reader.

4. **Cụm tài liệu**
   - có tài liệu còn giá trị như dataset definition.
   - có tài liệu mô tả hướng “minimal Go ML” đã lỗi thời, cần viết lại hoặc loại khỏi báo cáo.

---

## 3. Phân loại file Go hiện tại

### 3.1. `internal/services/predictive`

| File | Vai trò hiện tại | Kết luận K0 | Hướng xử lý ở các task sau |
|---|---|---|---|
| `internal/services/predictive/definition.go` | Khóa dataset definition, source/label/feature contract | **Giữ ý tưởng, không giữ implementation Go làm chuẩn cuối** | `K4`: chuyển nội dung khoa học sang `ml/at_risk_prediction/docs/AT_RISK_RESEARCH_SPEC.md`; có thể giữ file Go tạm thời nếu student/admin API vẫn cần enum/field name trong giai đoạn chuyển tiếp |
| `internal/services/predictive/definition_test.go` | Test cho dataset definition Go | **Bỏ sau khi đặc tả Python hoàn chỉnh** | Xóa khi phần spec Python và artifact contract đã ổn định |
| `internal/services/predictive/db_loader.go` | Query DB và dựng `TrainingRow` bằng Go | **Thay thế bằng Python data pipeline** | `K3`: port logic sang Python (`export_dataset.py`, `train_from_db.py`); sau đó xóa loader Go khi không còn endpoint nào dùng |
| `internal/services/predictive/db_loader_test.go` | Test DB loader Go | **Bỏ** | Xóa cùng `db_loader.go` |
| `internal/services/predictive/minimal_pipeline.go` | Rule-based + Logistic Regression thuần Go | **Loại bỏ** | `K5`: thay bằng pipeline Python (`pandas`, `scikit-learn`) |
| `internal/services/predictive/minimal_pipeline_test.go` | Test pipeline train/evaluate Go | **Loại bỏ** | Xóa cùng `minimal_pipeline.go` |
| `internal/services/predictive/runtime.go` | Runtime train-in-memory, cache metadata, list predictions | **Không giữ kiến trúc này** | `K7`: thay bằng service Go chỉ đọc `latest_predictions.json` + `model_metadata.json` do Python sinh |

### 3.2. `internal/usecases/predictive`

| File | Vai trò hiện tại | Kết luận K0 | Hướng xử lý |
|---|---|---|---|
| `internal/usecases/predictive/list_predictions.go` | Usecase admin list AT_RISK predictions | **Giữ route/usecase name, refactor nguồn dữ liệu** | `K7`: đổi service dependency từ `AtRiskService` runtime Go sang `PredictionArtifactService` |
| `internal/usecases/predictive/get_model_metadata.go` | Usecase đọc metadata model | **Giữ route/usecase name, refactor nguồn dữ liệu** | `K7`: đọc `model_metadata.json` từ Python artifact |

### 3.3. `cmd/cli/predictive_train`

| File | Vai trò hiện tại | Kết luận K0 | Hướng xử lý |
|---|---|---|---|
| `cmd/cli/predictive_train/main.go` | CLI train/evaluate predictive bằng Go | **Loại bỏ** | `K1-K5`: thay bằng `ml/at_risk_prediction/scripts/train_from_db.py` và `train_from_csv.py` |

### 3.4. `cmd/http/controllers/predictive`

| File | Vai trò hiện tại | Kết luận K0 | Hướng xử lý |
|---|---|---|---|
| `cmd/http/controllers/predictive/controller.go` | Đăng ký route admin predictive | **Giữ route** | Không đổi path; chỉ refactor cách lấy dữ liệu |
| `cmd/http/controllers/predictive/v1.go` | Controller admin AT_RISK list + metadata | **Giữ controller contract, refactor implementation** | `K7`: trả dữ liệu từ artifact Python |

### 3.5. Student portal liên quan predictive

| File | Vai trò hiện tại | Kết luận K0 | Hướng xử lý |
|---|---|---|---|
| `internal/usecases/studentportal/get_my_at_risk_prediction.go` | Lấy prediction của chính học sinh từ `AtRiskService` Go | **Giữ usecase/response shape, refactor service dependency** | `K7`: usecase này sẽ gọi artifact reader service mới thay vì runtime Go training |
| `cmd/http/controllers/studentportal/controller.go` | Route `GET /api/v1/student/at-risk` | **Giữ route** | Không đổi path |
| `cmd/http/controllers/studentportal/v1.go` | Mapping response student AT_RISK | **Giữ controller contract** | Chỉ đổi nguồn dữ liệu nội bộ |

### 3.6. Wire / bootstrap cần refactor

| File | Vai trò hiện tại | Kết luận K0 | Hướng xử lý |
|---|---|---|---|
| `internal/usecases/provider.go` | Cấp provider cho predictive usecase | **Refactor** | Trỏ sang service đọc artifact Python |
| `cmd/http/controllers/provider.go` | Wire controller predictive | **Refactor** | Không đổi public route |
| `cmd/http/main.go` | Mount route predictive | **Giữ** | Không đổi |
| `cmd/http/wire_gen.go` | File sinh tự động đang tạo `predictive.NewAtRiskService(db)` | **Sẽ thay đổi sau khi regenerate wire** | `K7`: regenerate theo service mới |

---

## 4. Ảnh hưởng tới frontend

Frontend hiện **không phải xóa màn**.  
Các page và API client đang dùng vẫn có thể giữ nếu Go backend tiếp tục duy trì route hiện tại.

### 4.1. Admin UI

| File | Trạng thái | Kết luận |
|---|---|---|
| `frontend/src/api/predictiveApi.ts` | Đang gọi `/v1/predictive/at-risk/students` và `/v1/predictive/at-risk/model-metadata` | **Giữ contract API** |
| `frontend/src/pages/admin/PredictiveAlertsPage.tsx` | Dùng dữ liệu cảnh báo thật | **Giữ UI, không cần viết lại** |
| `frontend/src/config/nav.ts` + `frontend/src/App.tsx` | Route/nav admin predictive | **Giữ** |

### 4.2. Student UI

| File | Trạng thái | Kết luận |
|---|---|---|
| `frontend/src/api/studentPortalApi.ts` | Có endpoint `/v1/student/at-risk` | **Giữ** |
| `frontend/src/pages/student/StudentOverview.tsx` | Hiển thị card AT_RISK | **Giữ** |
| `frontend/src/pages/student/StudentResultsPage.tsx` | Hiển thị alert AT_RISK | **Giữ** |

**Kết luận FE:**  
Nếu Go backend giữ nguyên response contract, phần FE gần như không cần đổi đáng kể trong `K7`.

---

## 5. Phân loại tài liệu

| File | Giá trị hiện tại | Kết luận K0 | Hướng xử lý |
|---|---|---|---|
| `docs/AT_RISK_DATASET_DEFINITION.md` | Vẫn còn giá trị ở mức dataset source/label/features | **Giữ tạm, viết lại theo Python** | `K4/K8`: chuyển thành bản nghiên cứu chính thức cho Python |
| `docs/AT_RISK_MINIMAL_APPROACH.md` | Mô tả hướng “Go tối giản” | **Lỗi thời** | `K8`: loại khỏi báo cáo chính; có thể xóa sau khi Python README + research spec hoàn chỉnh |
| `docs/assets/benchmark/figure_4_9_at_risk_data_flow.svg` | Hình luồng dữ liệu AT_RISK hiện phản ánh hướng cũ | **Phải cập nhật** | Vẽ lại pipeline Python: `DB/CSV -> Python FE -> Train/Evaluate -> Artifact -> Go API -> UI` |

---

## 6. Danh sách xóa / giữ / refactor sau K0

### 6.1. Nhóm sẽ xóa sau khi Python project chạy ổn

- `internal/services/predictive/minimal_pipeline.go`
- `internal/services/predictive/minimal_pipeline_test.go`
- `internal/services/predictive/db_loader.go`
- `internal/services/predictive/db_loader_test.go`
- `cmd/cli/predictive_train/main.go`
- `internal/services/predictive/runtime.go`
- `docs/AT_RISK_MINIMAL_APPROACH.md`

### 6.2. Nhóm giữ route/contract nhưng refactor implementation

- `internal/usecases/predictive/get_model_metadata.go`
- `internal/usecases/predictive/list_predictions.go`
- `cmd/http/controllers/predictive/controller.go`
- `cmd/http/controllers/predictive/v1.go`
- `internal/usecases/studentportal/get_my_at_risk_prediction.go`
- `cmd/http/controllers/studentportal/controller.go`
- `cmd/http/controllers/studentportal/v1.go`
- `internal/usecases/provider.go`
- `cmd/http/controllers/provider.go`
- `cmd/http/wire_gen.go`

### 6.3. Nhóm giữ nội dung nhưng chuyển hóa thành tài liệu/spec Python

- `internal/services/predictive/definition.go`
- `internal/services/predictive/definition_test.go`
- `docs/AT_RISK_DATASET_DEFINITION.md`

---

## 7. Kiến trúc đích sau khi chuyển sang Python

```text
PostgreSQL / CSV
        |
        v
Python ML project
  - export_dataset.py
  - train_from_db.py / train_from_csv.py
  - feature engineering
  - Logistic Regression / Random Forest / Rule-based baseline
  - metrics + confusion matrix + feature importance
  - latest_predictions.json
  - model_metadata.json
        |
        v
Go backend
  - PredictionArtifactService
  - /api/v1/predictive/at-risk/students
  - /api/v1/predictive/at-risk/model-metadata
  - /api/v1/student/at-risk
        |
        v
Frontend admin + student
```

---

## 8. Quyết định chốt sau K0

1. **Không tiếp tục phát triển nhánh train/evaluate bằng Go.**
2. **Python là hướng chính thức cho predictive analytics trong đồ án.**
3. **Go backend chỉ là lớp tích hợp artifact và phục vụ API/UI.**
4. **Route FE-BE hiện tại sẽ được giữ tối đa để tránh phát sinh sửa UI không cần thiết.**
5. **Báo cáo đồ án phải loại bỏ mô tả “pipeline ML Go” khỏi phần chính thức sau khi K8 hoàn tất.**

---

## 9. Đầu ra của K0

`Task K0` được xem là hoàn tất khi:

- đã audit toàn bộ file predictive Go hiện có,
- đã phân loại rõ `xóa / giữ / refactor / viết lại`,
- đã chỉ ra tài liệu lỗi thời,
- đã khóa kiến trúc đích Python -> artifact -> Go API -> FE.
