from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

PROJECT_ROOT = Path(__file__).resolve().parents[1]
if str(PROJECT_ROOT) not in sys.path:
    sys.path.insert(0, str(PROJECT_ROOT))

from src.data_pipeline import build_dataset_from_db
from src.database import create_db_engine
from src.dataset_schema import DEFAULT_DATASET_DEFINITION
from src.model_training import TrainingConfig, train_and_evaluate_models


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Train and evaluate AT_RISK models directly from PostgreSQL."
    )
    parser.add_argument("--dataset-name", default="at_risk_dataset_db", help="Prefix name for exported files.")
    parser.add_argument("--test-size", type=float, default=0.2, help="Test split ratio.")
    parser.add_argument("--seed", type=int, default=42, help="Random seed.")
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    engine = create_db_engine()
    dataset = build_dataset_from_db(engine, DEFAULT_DATASET_DEFINITION)
    outputs = train_and_evaluate_models(
        dataset=dataset,
        config=TrainingConfig(
            dataset_name=args.dataset_name,
            source="database",
            test_size=args.test_size,
            seed=args.seed,
        ),
    )
    print(json.dumps(outputs, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
