import io
import json
import torch
import requests
from PIL import Image
from transformers import AutoProcessor, CLIPModel
import db_client



model = CLIPModel.from_pretrained("openai/clip-vit-base-patch32")
processor = AutoProcessor.from_pretrained("openai/clip-vit-base-patch32")


def encode_image(image):
    inputs = processor(images=image, return_tensors="pt")

    with torch.inference_mode():
        outputs = model.get_image_features(**inputs)

    return outputs[0].detach().cpu().numpy().tolist()


def generate_encoding_for_channel(ch, method, properties, body):
    payload = json.loads(body)
    job_id = int(payload['job_id'])
    image_bytes = payload['image_bytes']
    redis_client = db_client.get_redis_client()

    embeddings = encode_image(Image.open(io.BytesIO(image_bytes)))
    if not embeddings:
        print(f"Failed to generate embeddings for job id {job_id}")
        redis_client.publish('encoding_results', json.dumps({
            'job_id': job_id,
            'status': 'failed',
            'embeddings': None
        }))
        return
    
    redis_client.publish('encoding_results', json.dumps({
        'job_id': job_id,
        'status': 'completed',
        'embeddings': embeddings
    }))


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






   






