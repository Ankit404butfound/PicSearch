import pika
import clip_processor
import face_encoder


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
                    on_message_callback=face_encoder.encode)
    channel.start_consuming()


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print('Interrupted')