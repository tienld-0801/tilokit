package python

const (
	FastAPIRequirements = `fastapi==0.104.0
uvicorn[standard]==0.24.0
python-decouple==3.8
pydantic==2.4.0`

	FastAPIMain = `from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from decouple import config

app = FastAPI(title="{{.ProjectName}}", version="1.0.0")

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
async def read_root():
    return {"message": "Hello World from {{.ProjectName}}!"}

@app.get("/health")
async def health_check():
    return {"status": "ok"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)`
)
