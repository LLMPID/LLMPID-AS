from classifier_engine import ClassifierEngine
from fastapi import FastAPI
from pydantic import BaseModel


classifier_eng = ClassifierEngine('model_data/model.onnx')
app = FastAPI()


class ClassifyRequest(BaseModel):
    text: str

@app.post("/classify")
async def classify(request: ClassifyRequest):
    request_data = request.text
    classification_result, _ = classifier_eng.classify(request_data)

    return {"result": classification_result}

@app.get('/health')
async def get_health():
    return {"status": "healthy"}
