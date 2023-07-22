from fastapi import FastAPI

app = FastAPI()



@app.get("/healthcheck")
async def health():
    return {"status": "I feel... not too bad. not too good too. At least I am alive."}

