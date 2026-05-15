# Bộ Dữ Liệu Khởi Tạo Benchmark Xếp Lịch

Thư mục này chứa bộ dữ liệu khởi tạo phục vụ đo đạc đối sánh tính năng xếp lịch.

## Thứ tự nạp dữ liệu

1. Chạy `00_nen_chung.sql`.
2. Chọn đúng một kịch bản quy mô:
   - `01_kich_ban_nho.sql`
   - `02_kich_ban_trung_binh.sql`
   - `03_kich_ban_lon.sql`
3. Nếu cần kiểm thử luật hoặc tiền xử lý, chạy thêm đúng một gói `K1` đến `K9`.

## Mục đích từng nhóm tệp

- `00_nen_chung.sql`
  - tạo cơ sở, ca học và thời gian di chuyển liên cơ sở dùng chung.
- `01_kich_ban_nho.sql`
  - dữ liệu quy mô nhỏ để kiểm thử logic và benchmark nhanh.
- `02_kich_ban_trung_binh.sql`
  - dữ liệu quy mô trung bình cho benchmark chuẩn chi nhánh trung tâm.
- `03_kich_ban_lon.sql`
  - dữ liệu quy mô lớn để đo khả năng mở rộng.
- `11_...` đến `19_...`
  - các gói kiểm thử luật `K1` đến `K9`, áp dụng chồng lên một kịch bản đã nạp.

## Tạo gói seed tự động

Có thể dùng script [tao_goi_seed_xep_lich.sh](/Users/hant/golang/doan/scripts/tao_goi_seed_xep_lich.sh) để ghép thành một tệp SQL duy nhất.

Ví dụ:

```bash
./scripts/tao_goi_seed_xep_lich.sh nho
./scripts/tao_goi_seed_xep_lich.sh trung_binh K1,K4,K7
./scripts/tao_goi_seed_xep_lich.sh lon K8,K9
```

Tệp kết quả sẽ được ghi vào thư mục `artifacts/goi_seed_xep_lich/`.

## Tiền tố mã dữ liệu

- `XLB-CHUNG-...`: dữ liệu nền dùng chung.
- `XLB-NHO-...`: dữ liệu kịch bản nhỏ.
- `XLB-TB-...`: dữ liệu kịch bản trung bình.
- `XLB-LON-...`: dữ liệu kịch bản lớn.

## Ghi chú

- Các tệp được viết theo hướng có thể chạy lặp lại nhiều lần.
- Các câu lệnh chỉ dọn và ghi đè dữ liệu có cùng tiền tố benchmark, không đụng vào dữ liệu vận hành khác.
- Các gói `K1` đến `K9` giả định bạn đã nạp xong `00_nen_chung.sql` và một kịch bản quy mô tương ứng.
