from app.config.env import get_env
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

env = get_env()
app = FastAPI(
    title="Machine Learning API",
    description="Machine Learning API",
    version="0.1.0",
    docs_url="/docs",
    root_path=env.ML_PREFIX_API,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/")
async def root():
    return {"message": "Welcome to the Machine Learning API"}
