from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

PROJECT_ROOT = Path(__file__).resolve().parents[1]
if str(PROJECT_ROOT) not in sys.path:
    sys.path.insert(0, str(PROJECT_ROOT))

from src.data_pipeline import load_dataset_from_csv
from src.model_training import TrainingConfig, train_and_evaluate_models


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Train and evaluate AT_RISK models from CSV."
    )
    parser.add_argument("--input", required=True, help="Path to normalized AT_RISK CSV dataset.")
    parser.add_argument("--dataset-name", default="at_risk_dataset_csv", help="Prefix name for exported files.")
    parser.add_argument("--test-size", type=float, default=0.2, help="Test split ratio.")
    parser.add_argument("--seed", type=int, default=42, help="Random seed.")
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    input_path = Path(args.input).expanduser().resolve()
    dataset = load_dataset_from_csv(input_path)
    outputs = train_and_evaluate_models(
        dataset=dataset,
        config=TrainingConfig(
            dataset_name=args.dataset_name,
            source=f"csv:{input_path.name}",
            test_size=args.test_size,
            seed=args.seed,
        ),
    )
    print(json.dumps(outputs, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
