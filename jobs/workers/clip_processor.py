import io
import json
import torch
import requests
from PIL import Image
from transformers import CLIPProcessor, CLIPModel
import db_client



model = CLIPModel.from_pretrained("openai/clip-vit-base-patch32")
processor = CLIPProcessor.from_pretrained("openai/clip-vit-base-patch32")
redis_client = db_client.get_redis_client()


def encode_image(image):
    inputs = processor(images=image, return_tensors="pt")
    outputs = model.get_image_features(**inputs)
    return outputs[0].detach().cpu().numpy().tolist()

def encode_text(text):
    inputs = processor(text=[text], return_tensors="pt", padding=True)
    outputs = model.get_text_features(**inputs)
    return outputs[0].detach().cpu().numpy().tolist()


def generate_encoding_for_channel(ch, method, properties, body):
    # try:
        print("Received message for generating CLIP encoding", body)

        embeddings = encode_text(text=body.decode('utf-8'))
        if not embeddings:
            print(f"Failed to generate embeddings for text: {body.decode('utf-8')}")
            redis_client.publish('encoding_results', json.dumps({
                'embeddings': None
            }))
            return
        
        print(embeddings)
        redis_client.publish('encoding_results', json.dumps({
            'embeddings': embeddings
        }))
    # except Exception as e:
    #     print(f"Error generating CLIP encoding: {e}")
    #     redis_client.publish('encoding_results', json.dumps({
    #         'embeddings': None
    #     }))


def process_image(ch, method, properties, body):

    job_id = int(body)
    conn, cur = db_client.get_conn()
    # Fetch job details
    cur.execute(
        "SELECT id, file_url, FROM files \
            JOIN jobs ON files.id = jobs.file_id \
            WHERE jobs.id = %s AND jobs.face_encoding_status = 'pending'",
        (job_id,)
    )
    job = cur.fetchone()
    if not job:
        print(f"No pending job found with id {job_id}")
        return
    
    file_url = job['file_url']
    image_response = requests.get(file_url)
    image_bytes = image_response.content
    image = Image.open(io.BytesIO(image_bytes))
    image_embeddings = encode_image(image)

    cur.execute("UPDATE jobs SET universal_encoding_status = 'completed' WHERE id = %s", (job_id,))
    conn.commit()

    cur.execute("UPDATE files SET embedding = %s WHERE id = %s", (image_embeddings, job_id,))
    conn.commit()
    return






   






