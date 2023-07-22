from fastapi import FastAPI
from src.db.engine import engine
from src.definitions import Base
from src.routers.classes import router as class_router

app = FastAPI()
app.include_router(class_router)


@app.get("/")
async def health():
    return {"status": "I feel... not too bad. not too good too. At least I am alive."}


@app.on_event("startup")
async def init_tables():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)


@app.on_event("shutdown")
async def shutdown():
    await engine.dispose()
