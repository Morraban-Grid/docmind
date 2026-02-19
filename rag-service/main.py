# import classes anf functions from http.server module
from http.server import BaseHTTPRequestHandler, HTTPServer

"""
This is a simple HTTP server that listens on port 8000
and responds to GET requests on the /health endpoint with
a message indicating that the RAG service is running.
For any other endpoint, it responds with a 404 Not Found status.
"""
class Handler(BaseHTTPRequestHandler):
    # Override the do_GET method to handle GET requests
    def do_GET(self):
        if self.path == "/health":
            self.send_response(200)
            self.send_header("Content-type", "text/plain")
            self.end_headers()
            self.wfile.write(b"DocMind RAG Service is running")
        else:
            self.send_response(404)
            self.end_headers()

def run():
    server_address = ('', 8000)
    httpd = HTTPServer(server_address, Handler)
    print("RAG service running on port 8000")
    httpd.serve_forever()

if __name__ == "__main__":
    run()
