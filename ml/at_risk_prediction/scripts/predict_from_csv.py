from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

PROJECT_ROOT = Path(__file__).resolve().parents[1]
if str(PROJECT_ROOT) not in sys.path:
    sys.path.insert(0, str(PROJECT_ROOT))

from src.data_pipeline import load_dataset_from_csv
from src.prediction_inference import (
    append_selection_conclusion_to_report,
    predict_dataset,
    save_primary_model_selection,
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Generate AT_RISK prediction artifact from CSV."
    )
    parser.add_argument("--input", required=True, help="Path to normalized AT_RISK CSV dataset.")
    parser.add_argument(
        "--model",
        default=None,
        choices=["rule_based", "logistic_regression", "random_forest"],
        help="Override selected model from model_metadata.json.",
    )
    parser.add_argument(
        "--output",
        default=None,
        help="Optional path for latest_predictions.json.",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    input_path = Path(args.input).expanduser().resolve()
    dataset = load_dataset_from_csv(input_path)
    selection = save_primary_model_selection()
    append_selection_conclusion_to_report(selection=selection)
    result = predict_dataset(
        dataset=dataset,
        model_name=args.model,
        output_path=args.output,
        source=f"csv:{input_path.name}",
    )
    print(json.dumps(result, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
