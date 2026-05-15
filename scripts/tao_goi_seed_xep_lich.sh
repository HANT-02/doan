#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SEED_DIR="$ROOT_DIR/internal/infrastructure/database/postgres/sql_seed/xep_lich_benchmark"
OUTPUT_DIR="$ROOT_DIR/artifacts/goi_seed_xep_lich"

SCENARIO="${1:-}"
PACKAGES="${2:-}"

if [[ -z "$SCENARIO" ]]; then
  echo "Cach dung: $0 <kich_ban> [goi_luat]"
  echo "Vi du:"
  echo "  $0 nho"
  echo "  $0 trung_binh K1,K4,K7"
  exit 1
fi

case "$SCENARIO" in
  nho)
    SCENARIO_FILE="$SEED_DIR/01_kich_ban_nho.sql"
    ;;
  trung_binh)
    SCENARIO_FILE="$SEED_DIR/02_kich_ban_trung_binh.sql"
    ;;
  lon)
    SCENARIO_FILE="$SEED_DIR/03_kich_ban_lon.sql"
    ;;
  *)
    echo "Kich ban khong hop le: $SCENARIO"
    echo "Gia tri hop le: nho, trung_binh, lon"
    exit 1
    ;;
esac

mkdir -p "$OUTPUT_DIR"

TIMESTAMP="$(date +"%Y%m%d_%H%M%S")"
PACKAGE_SLUG="khong_co_goi_luat"
if [[ -n "$PACKAGES" ]]; then
  PACKAGE_SLUG="$(echo "$PACKAGES" | tr ',' '_' | tr '[:upper:]' '[:lower:]')"
fi

OUTPUT_FILE="$OUTPUT_DIR/goi_seed_${SCENARIO}_${PACKAGE_SLUG}_${TIMESTAMP}.sql"

{
  echo "-- Goi seed xep lich duoc tao tu dong"
  echo "-- Thoi diem tao: $(date +"%Y-%m-%d %H:%M:%S %z")"
  echo "-- Kich ban: $SCENARIO"
  echo "-- Goi luat: ${PACKAGES:-khong co}"
  echo
  cat "$SEED_DIR/00_nen_chung.sql"
  echo
  cat "$SCENARIO_FILE"
  echo
} > "$OUTPUT_FILE"

if [[ -n "$PACKAGES" ]]; then
  IFS=',' read -r -a package_array <<< "$PACKAGES"
  for raw_name in "${package_array[@]}"; do
    package_name="$(echo "$raw_name" | tr '[:lower:]' '[:upper:]' | xargs)"
    if [[ -z "$package_name" ]]; then
      continue
    fi

    case "$package_name" in
      K1) package_file="$SEED_DIR/11_kiem_thu_luat_k1_si_so_khong_dat.sql" ;;
      K2) package_file="$SEED_DIR/12_kiem_thu_luat_k2_thieu_giao_vien.sql" ;;
      K3) package_file="$SEED_DIR/13_kiem_thu_luat_k3_sai_ky_nang.sql" ;;
      K4) package_file="$SEED_DIR/14_kiem_thu_luat_k4_thieu_lich_tuan.sql" ;;
      K5) package_file="$SEED_DIR/15_kiem_thu_luat_k5_khong_du_ca_xep.sql" ;;
      K6) package_file="$SEED_DIR/16_kiem_thu_luat_k6_ap_luc_suc_chua_phong.sql" ;;
      K7) package_file="$SEED_DIR/17_kiem_thu_luat_k7_xung_dot_di_chuyen.sql" ;;
      K8) package_file="$SEED_DIR/18_kiem_thu_luat_k8_dieu_chinh_buoi_khoa.sql" ;;
      K9) package_file="$SEED_DIR/19_kiem_thu_luat_k9_dieu_chinh_buoi_nhap.sql" ;;
      *)
        echo "Goi luat khong hop le: $package_name" >&2
        rm -f "$OUTPUT_FILE"
        exit 1
        ;;
    esac

    {
      echo
      echo "-- Bat dau goi luat $package_name"
      cat "$package_file"
      echo "-- Ket thuc goi luat $package_name"
      echo
    } >> "$OUTPUT_FILE"
  done
fi

echo "Da tao goi seed tai: $OUTPUT_FILE"
