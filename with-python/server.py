import http.server
import socketserver
import json
from datetime import datetime

PORT = 8001

def write_to_file(filename,data):
    try:
        with open(filename, 'w') as file:
            file.write(f"Location: {data['location']}\n")
            file.write(f"Units: {data['units']}\n")
            file.write(f"Date: {data['date']}\n")
        return True
    except Exception as e:
        print(f"Error writing to file: {e}")
        return False
    
def clean_message(message):
    if message['location'] is None:
        return None
    if message['units'] is None:
        return None
    if message['date'] is None:
        return None
    
    clean_message = {
        'location': message['location'].upper(),
        'units': message['units'],
        'date': message['date']
    }
    return clean_message

class RequestHandler(http.server.SimpleHTTPRequestHandler):
    def do_POST(self):
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        
        try:
            json_data = json.loads(post_data.decode('utf-8'))
            print("Received JSON payload:", json_data)
            cleaned = clean_message(json_data)
            if cleaned:
                filename = "tmp/"+datetime.utcnow().strftime("%Y-%m-%d-%H-%M-%S")+".txt"
                if write_to_file(filename,cleaned):
                    self.send_response(200)
                    self.send_header('Content-type', 'text/plain')
                    self.end_headers()
                    self.wfile.write(b"JSON received successfully")
                else:
                    self.send_response(500)
                    self.send_header('Content-type', 'text/plain')
                    self.end_headers()
                    self.wfile.write(b"Error writing to file")
        except json.JSONDecodeError:
            self.send_response(400)
            self.send_header('Content-type', 'text/plain')
            self.end_headers()
            self.wfile.write(b"Invalid JSON payload")

with socketserver.TCPServer(("", PORT), RequestHandler) as httpd:
    print(f"Serving at port {PORT}")
    httpd.serve_forever()
