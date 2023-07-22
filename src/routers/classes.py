from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from src.db.engine import get_db
from src.db.class_table import ClassTableRepository

router = APIRouter(prefix="/classes", tags=["Classes"])


@router.get("/all")
async def all_classes(
    page: int = 1, page_size: int = 10, db: AsyncSession = Depends(get_db)
):
    return await ClassTableRepository(db).get_all(page=page, page_size=page_size)
