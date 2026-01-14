import pika
import jobs.workers.clip_processor as clip_processor
import jobs.workers.face_encoder as face_encoder


def main():
    connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    channel = connection.channel()

    channel.queue_declare(queue='clip_processor')
    channel.basic_consume(queue='clip_processor',
                    auto_ack=True,
                    on_message_callback=clip_processor.process_image)
    
    channel.queue_declare(queue='face_encoder')
    channel.basic_consume(queue='face_encoder',
                    auto_ack=True,
                    on_message_callback=face_encoder.process_image)
    
    channel.queue_declare(queue='generate_clip_encoding')
    channel.basic_consume(queue='generate_clip_encoding',
                    auto_ack=True,
                    on_message_callback=clip_processor.generate_encoding_for_channel)
    channel.start_consuming()


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print('Interrupted')