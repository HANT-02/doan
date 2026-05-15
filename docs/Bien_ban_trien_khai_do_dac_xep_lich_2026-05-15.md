# Biên Bản Triển Khai Đo Đạc Xếp Lịch Ngày 15-05-2026

## 1. Phạm vi đã triển khai

Đã hoàn thành các hạng mục sau:

- dựng bộ dữ liệu khởi tạo quy mô `nhỏ`, `trung_bình`, `lớn`;
- dựng đủ các gói kiểm thử luật `K1` đến `K9`;
- chạy đo đạc đối sánh cho đủ `3` thuật toán, mỗi kịch bản `10` lần chạy;
- ghi số liệu thô theo từng lần chạy kèm thời điểm bắt đầu, kết thúc và thời gian chạy;
- sinh thêm tệp tổng hợp để tiện đưa vào bảng biểu của báo cáo.

## 2. Bộ dữ liệu khởi tạo đã sẵn sàng

Thư mục dữ liệu:

- [xep_lich_benchmark](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark)

Tệp nền và kịch bản:

- [00_nen_chung.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/00_nen_chung.sql)
- [01_kich_ban_nho.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/01_kich_ban_nho.sql)
- [02_kich_ban_trung_binh.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/02_kich_ban_trung_binh.sql)
- [03_kich_ban_lon.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/03_kich_ban_lon.sql)

Gói kiểm thử luật:

- [11_kiem_thu_luat_k1_si_so_khong_dat.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/11_kiem_thu_luat_k1_si_so_khong_dat.sql)
- [12_kiem_thu_luat_k2_thieu_giao_vien.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/12_kiem_thu_luat_k2_thieu_giao_vien.sql)
- [13_kiem_thu_luat_k3_sai_ky_nang.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/13_kiem_thu_luat_k3_sai_ky_nang.sql)
- [14_kiem_thu_luat_k4_thieu_lich_tuan.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/14_kiem_thu_luat_k4_thieu_lich_tuan.sql)
- [15_kiem_thu_luat_k5_khong_du_ca_xep.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/15_kiem_thu_luat_k5_khong_du_ca_xep.sql)
- [16_kiem_thu_luat_k6_ap_luc_suc_chua_phong.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/16_kiem_thu_luat_k6_ap_luc_suc_chua_phong.sql)
- [17_kiem_thu_luat_k7_xung_dot_di_chuyen.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/17_kiem_thu_luat_k7_xung_dot_di_chuyen.sql)
- [18_kiem_thu_luat_k8_dieu_chinh_buoi_khoa.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/18_kiem_thu_luat_k8_dieu_chinh_buoi_khoa.sql)
- [19_kiem_thu_luat_k9_dieu_chinh_buoi_nhap.sql](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/19_kiem_thu_luat_k9_dieu_chinh_buoi_nhap.sql)

Tệp hướng dẫn:

- [README.md](/Users/hant/golang/doan/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark/README.md)

## 3. Gói seed ghép tự động

Tệp lệnh ghép:

- [tao_goi_seed_xep_lich.sh](/Users/hant/golang/doan/scripts/tao_goi_seed_xep_lich.sh)

Ví dụ đã tạo:

- [goi_seed_nho_k1_k4_20260515_204311.sql](/Users/hant/golang/doan/artifacts/goi_seed_xep_lich/goi_seed_nho_k1_k4_20260515_204311.sql)

## 4. Lần đo đạc đã chạy

Lệnh đã chạy:

```bash
env GOCACHE=/Users/hant/golang/doan/.gocache go run ./cmd/cli/scheduling_benchmark -kich-ban nho,trung_binh,lon -so-lan-chay 10
```

Thư mục minh chứng đã sinh:

- [20260515_204909](/Users/hant/golang/doan/artifacts/do_dac_xep_lich/20260515_204909)

Tệp tổng hợp:

- [tong_hop.md](/Users/hant/golang/doan/artifacts/do_dac_xep_lich/20260515_204909/tong_hop.md)
- [tong_hop.json](/Users/hant/golang/doan/artifacts/do_dac_xep_lich/20260515_204909/tong_hop.json)
- [bang_ke_tap_tin.json](/Users/hant/golang/doan/artifacts/do_dac_xep_lich/20260515_204909/bang_ke_tap_tin.json)
- [bang_tong_hop_theo_thuat_toan.csv](/Users/hant/golang/doan/artifacts/do_dac_xep_lich/20260515_204909/bang_tong_hop_theo_thuat_toan.csv)
- [so_lieu_tho_tung_lan_chay.csv](/Users/hant/golang/doan/artifacts/do_dac_xep_lich/20260515_204909/so_lieu_tho_tung_lan_chay.csv)

Ví dụ tệp số liệu thô theo lần chạy:

- [so_lieu_tho.json](/Users/hant/golang/doan/artifacts/do_dac_xep_lich/20260515_204909/kich_ban_nho/thuat_toan_cp_sat/lan_chay_001/so_lieu_tho.json)
- [thong_tin_telemetry.json](/Users/hant/golang/doan/artifacts/do_dac_xep_lich/20260515_204909/kich_ban_nho/thuat_toan_cp_sat/lan_chay_001/thong_tin_telemetry.json)

## 5. Danh sách chỉ số đo đạc đã ghi

Trong tệp tổng hợp và tệp thô hiện đã có đầy đủ:

- tên kịch bản, khóa thuật toán, lần chạy;
- thời điểm bắt đầu, thời điểm kết thúc, thời gian chạy;
- số buổi xếp được, số buổi chưa xếp được, số xung đột, điểm mềm;
- số phương án đã đánh giá, số phương án bị loại do xung đột;
- số nút đã duyệt, số nhánh bị cắt;
- tốc độ đánh giá phương án mỗi giây;
- trung bình, độ lệch chuẩn, giá trị nhỏ nhất, giá trị lớn nhất theo từng thuật toán.
