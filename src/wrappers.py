from typing import TypeVar, Generic

from pydantic import BaseModel

DataT = TypeVar("DataT")


class ResponseWrapper(BaseModel, Generic[DataT]):
    data: DataT
