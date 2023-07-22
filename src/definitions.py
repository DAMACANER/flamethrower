from pathlib import Path
import sqlalchemy
from sqlalchemy.ext.declarative import declarative_base

ROOT_DIR = Path(__file__).parent.parent.__str__()


metadata = sqlalchemy.MetaData()

Base = declarative_base(metadata=metadata)
