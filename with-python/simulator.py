import json
import random
import urllib.request
import urllib.error
import threading

from datetime import datetime, timedelta


def generate_message():
    # Generate random data
    random_id = random.randint(1000, 9999)
    random_date = datetime.now() - timedelta(days=random.randint(0, 365))
    locations = ["New York", "London", "Tokyo", "Paris", "Sydney", "Berlin", "Moscow", "Dubai"]
    random_location = random.choice(locations)

    # Create the JSON payload
    payload = {
        "units": random_id,
        "date": random_date.strftime("%Y-%m-%d"),
        "location": random_location
    }

    # Convert the payload to a JSON string
    json_payload = json.dumps(payload)
    return json_payload

def send_message(message):
    try:
        url = "http://localhost:8001"
        req = urllib.request.Request(url, data=message.encode('utf-8'), method='POST')
        req.add_header('Content-Type', 'application/json')
        
        with urllib.request.urlopen(req) as response:
            if response.status == 200:
                print(f"Message sent successfully. Status code: {response.status}")
            else:
                print(f"Message sent, but received unexpected status code: {response.status}")
    except urllib.error.URLError as e:
        print(f"Failed to send message. Error: {e.reason}")
    except Exception as e:
        print(f"An unexpected error occurred: {str(e)}")

def daemon():
    message = generate_message()
    send_message(message)

def set_interval(func,interval):
    func()
    threading.Timer(interval,set_interval , [func,interval]).start()

# Start the periodic printing
set_interval(daemon,5)
