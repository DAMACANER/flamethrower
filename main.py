from fastapi import FastAPI
from src.db.engine import engine
from src.definitions import Base

app = FastAPI()


@app.get("/healthcheck")
async def health():
    return {"status": "I feel... not too bad. not too good too. At least I am alive."}


@app.on_event("startup")
async def init_tables():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
