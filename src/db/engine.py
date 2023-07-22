from typing import AsyncGenerator
from src.definitions import ROOT_DIR
from sqlalchemy.ext.asyncio import (
    create_async_engine,
    async_sessionmaker,
    AsyncSession,
    AsyncConnection,
)

DB_LOC = ROOT_DIR + "/dnd35.db"
DATABASE_URL = f"sqlite+aiosqlite:///{DB_LOC}"

engine = create_async_engine(DATABASE_URL, echo=True, pool_pre_ping=True)

async_session = async_sessionmaker(engine, expire_on_commit=False)


async def get_db() -> AsyncGenerator[AsyncSession, None]:
    async with async_session() as session:
        yield session


async def get_engine() -> AsyncGenerator[AsyncConnection, None]:
    async with engine.begin() as conn:
        yield conn
