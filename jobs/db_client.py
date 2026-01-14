import os
import redis
import psycopg2


def get_conn():
    # Get config from environment
    conn = psycopg2.connect(
        os.getenv("DSN")
    )
    
    # Create cursor
    cur = conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor)
    
    return conn, cur


def get_redis_client():
    redis_url = os.getenv("REDIS_URL", "redis://localhost:6379/0")
    return redis.from_url(redis_url)