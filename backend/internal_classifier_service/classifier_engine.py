import onnxruntime as ort
import numpy as np
from transformers import BertTokenizer

class ClassifierEngine():
    def __init__(self, model_path: str):
        self.__tokenizer = BertTokenizer.from_pretrained("bert-base-uncased")  # Change if needed
        self.__onnx_session = ort.InferenceSession(model_path, providers=["CPUExecutionProvider"])

    def classify(self, input_text: str):
        tokens = self.__tokenizer(input_text, padding="max_length", truncation=True, max_length=128, return_tensors="np")

        # Convert the input to int64 tensors
        input_ids = tokens["input_ids"].astype(np.int64)
        attention_mask = tokens["attention_mask"].astype(np.float32)

        # Extract input names 
        input_names = {inp.name: inp for inp in self.__onnx_session.get_inputs()}

        # Prepare model inputs
        feed_dict = {
            "input_ids": input_ids,
            "attention_mask": attention_mask
        }

        # Check if token_type_ids is required
        if "token_type_ids" in input_names:
            feed_dict["token_type_ids"] = tokens["token_type_ids"].astype(np.int64)

        # Run inference
        output = self.__onnx_session.run(None, feed_dict)
        output_array = output[0]  # Extract first output tensor

        # Ensure the output is a scalar or a single-value array
        if output_array.ndim > 1:
            logit = output_array[0][0]  # Assuming shape is (1,1)
        else:
            logit = output_array[0]  # Assuming shape is (1,)

        # Convert logit to probability (if needed)
        probability = 1 / (1 + np.exp(-logit))  # Sigmoid

        # Classification decision
        threshold = 0.6
        is_normal = probability > threshold

        return "Normal" if is_normal else "Injection", probability
