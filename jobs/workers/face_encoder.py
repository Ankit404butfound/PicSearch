import os
import io
import secrets
import PIL.Image as Image
from networkx import second_order_centrality
import requests
import face_recognition

import db_client


def encode_face(image):
    # Load the image
    image = face_recognition.load_image_file(image)

    # Encode the face(s)
    face_encodings = face_recognition.face_encodings(image, num_jitters=10, model="large")
    face_locations = face_recognition.face_locations(image)
    return face_encodings, face_locations if face_encodings else ([], [])


def process_image(ch, method, properties, body):
    print("Received message for processing image", body)
    job_id = int(body)
    conn, cur = db_client.get_conn()
    # Fetch job details
    cur.execute(
        "SELECT files.id, files.url FROM files \
            JOIN jobs ON files.id = jobs.file_id \
            WHERE jobs.id = %s AND jobs.face_encoding_status = 'pending'",
        (job_id,)
    )
    job = cur.fetchone()
    if not job:
        print(f"No pending job found with id {job_id}")
        return
    file_url = job.get('url')

    # Load image and compute face encoding
    response = requests.get(file_url)
    image = io.BytesIO(response.content)
    face_encodings, face_locations = encode_face(image)
    if not face_encodings:
        # Update job status to failed if no faces found
        cur.execute("UPDATE jobs SET face_encoding_status = 'failed' WHERE id = %s", (job_id,))
        conn.commit()
        print(f"No faces found in image for job id {job_id}")
        return
    
    print(f"Found {len(face_encodings)} face(s) in image for job id {job_id}")
    # Store face encodings and locations in the database
    for encoding, location in zip(face_encodings, face_locations):
        # Find closest existing encoding from unique faces
        cur.execute(
            "SELECT id, embedding <-> %s AS distance FROM unique_faces \
                WHERE embedding <-> %s < %s \
                ORDER BY distance ASC LIMIT 1",
            (str(encoding.tolist()), str(encoding.tolist()), float(os.getenv("FACE_ENCODING_THRESHOLD")))
        )
        result = cur.fetchone()
        if result:
            unique_face_id = result['id']
        else:
            # Insert new unique face
            # First crop the face image for storage
            top, right, bottom, left = location
            image = io.BytesIO(response.content)
            image = Image.open(image)
            face_image = image.crop((left, top, right, bottom))
            face_image_io = io.BytesIO()
            face_image.save(face_image_io, format='JPEG')
            face_image_io.seek(0)
            with open(f'../uploads/faces/face_{job_id}.jpg', 'wb') as f:
                f.write(face_image_io.read())
            url = f'{os.getenv("SERVER_HOST")}/api/files/download/faces/face_{job_id}/{secrets.token_hex(16)}.jpg'
            cur.execute(
                "INSERT INTO unique_faces (embedding, image_url) VALUES (%s, %s) RETURNING id",
                (str(encoding.tolist()), url)
            )
            unique_face_id = cur.fetchone()['id']

        # Insert face instance
        cur.execute(
            "INSERT INTO faces (file_id, unique_face_id, coordinates) VALUES (%s, %s, %s)",
            (job['id'], unique_face_id, list(location))
        )

    # Update job status to completed
    cur.execute("UPDATE jobs SET face_encoding_status = 'completed' WHERE id = %s", (job_id,))
    conn.commit()
    print(f"Successfully processed job id {job_id}")

    