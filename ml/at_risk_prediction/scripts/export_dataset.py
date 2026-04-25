from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

PROJECT_ROOT = Path(__file__).resolve().parents[1]
if str(PROJECT_ROOT) not in sys.path:
    sys.path.insert(0, str(PROJECT_ROOT))

from src.data_pipeline import build_dataset_from_db, build_dataset_summary, save_dataset_artifacts
from src.database import create_db_engine
from src.dataset_schema import DEFAULT_DATASET_DEFINITION


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Export training-ready AT_RISK dataset from PostgreSQL.")
    parser.add_argument("--dataset-name", default="at_risk_dataset", help="Prefix name for exported files.")
    parser.add_argument("--test-size", type=float, default=0.2, help="Test split ratio for exported train/test CSV.")
    parser.add_argument("--seed", type=int, default=42, help="Random seed used for train/test split.")
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    engine = create_db_engine()
    dataset = build_dataset_from_db(engine, DEFAULT_DATASET_DEFINITION)
    summary = build_dataset_summary(dataset, DEFAULT_DATASET_DEFINITION, source="database")
    paths = save_dataset_artifacts(
        dataset=dataset,
        summary=summary,
        dataset_name=args.dataset_name,
        save_split=True,
        test_size=args.test_size,
        seed=args.seed,
    )

    print(json.dumps({"summary": summary, "paths": paths}, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
