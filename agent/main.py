from fastapi import FastAPI
import uvicorn

app = FastAPI(title="Hackathon Judge Agent")

@app.get("/")
async def root():
    return {"message": "Welcome to the Hackathon Judge Agent"}

@app.get("/health")
async def health():
    return {"status": "healthy"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
