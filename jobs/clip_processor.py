import torch
from transformers import AutoProcessor, CLIPModel
from transformers.image_utils import load_image
import db_client



model = CLIPModel.from_pretrained("openai/clip-vit-base-patch32")
processor = AutoProcessor.from_pretrained("openai/clip-vit-base-patch32")



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


    image = load_image(file_url)

    inputs = processor(
        text=["a photo of a cat", "a photo of a dog"], images=image, return_tensors="pt", padding=True
    )

    with torch.inference_mode():
        outputs = model(**inputs)

    image_embeddings = outputs.image_embeds

    cur.execute("UPDATE jobs SET universal_encoding_status = 'completed' WHERE id = %s", (job_id,))
    conn.commit()

    cur.execute("UPDATE files SET embedding = %s WHERE id = %s", (image_embeddings, job_id,))
    conn.commit()
    return






   






