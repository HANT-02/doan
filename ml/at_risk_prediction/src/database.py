from __future__ import annotations

import os
from urllib.parse import quote_plus

from dotenv import load_dotenv
from sqlalchemy import create_engine
from sqlalchemy.engine import Engine


def build_database_url() -> str:
    load_dotenv()

    user = os.getenv("DB_USER", "postgres")
    password = quote_plus(os.getenv("DB_PASSWORD", "postgres"))
    host = os.getenv("DB_HOST", "localhost")
    port = os.getenv("DB_PORT", "5432")
    name = os.getenv("DB_NAME", "doan")

    return f"postgresql+psycopg://{user}:{password}@{host}:{port}/{name}"


def create_db_engine() -> Engine:
    return create_engine(build_database_url(), future=True)
