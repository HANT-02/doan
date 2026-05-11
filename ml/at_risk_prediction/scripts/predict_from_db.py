from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

import pandas as pd

PROJECT_ROOT = Path(__file__).resolve().parents[1]
if str(PROJECT_ROOT) not in sys.path:
    sys.path.insert(0, str(PROJECT_ROOT))

from src.data_pipeline import build_inference_dataset_from_db
from src.database import create_db_engine
from src.dataset_schema import DEFAULT_DATASET_DEFINITION
from src.prediction_inference import (
    append_selection_conclusion_to_report,
    predict_dataset,
    save_primary_model_selection,
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Generate AT_RISK prediction artifact directly from PostgreSQL."
    )
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
    selection = save_primary_model_selection()
    append_selection_conclusion_to_report(selection=selection)

    try:
        engine = create_db_engine()
        dataset = build_inference_dataset_from_db(engine, DEFAULT_DATASET_DEFINITION)
    except ValueError as exc:
        result = predict_dataset(
            dataset=pd.DataFrame(columns=[]),
            model_name=args.model,
            output_path=args.output,
            source="database",
        )
        result["warning"] = str(exc)
        print(json.dumps(result, indent=2, ensure_ascii=False))
        return

    result = predict_dataset(
        dataset=dataset,
        model_name=args.model,
        output_path=args.output,
        source="database",
    )
    print(json.dumps(result, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
