import os
import io
import db_client
import requests
import face_recognition


def encode_face(image):
    # Load the image
    image = face_recognition.load_image_file(image)

    # Encode the face(s)
    face_encodings = face_recognition.face_encodings(image, num_jitters=10, model="large")
    face_locations = face_recognition.face_locations(image)
    return face_encodings, face_locations if face_encodings else ([], [])


def encode(ch, method, properties, body):
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
    file_url = job.get('file_url')

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
    
    # Store face encodings and locations in the database
    for encoding, location in zip(face_encodings, face_locations):
        # Find closest existing encoding from unique faces
        cur.execute(
            "SELECT id, embedding <=> %s AS distance FROM unique_faces \
                WHERE distance < %s \
                ORDER BY distance ASC LIMIT 1",
            (list(encoding), float(os.getenv("FACE_ENCODING_THRESHOLD")))
        )
        result = cur.fetchone()
        if result:
            unique_face_id = result['id']
        else:
            # Insert new unique face
            cur.execute(
                "INSERT INTO unique_faces (embedding) VALUES (%s) RETURNING id",
                (list(encoding),)
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

    