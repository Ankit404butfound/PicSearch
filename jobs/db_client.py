import os
import psycopg2

def get_conn():
    # Get config from environment
    conn = psycopg2.connect(
        os.getenv("DSN")
    )
    
    # Create cursor
    cur = conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor)
    
    return conn, cur