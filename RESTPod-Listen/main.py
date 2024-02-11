import json
import os
from http.server import BaseHTTPRequestHandler, HTTPServer
from kubernetes import client, config
from datetime import datetime
import logging


class RequestHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path != "/send-data":
            self.send_error(405, "Method Not Allowed")
            return

        content_length = int(self.headers["Content-Length"])
        payload_data = self.rfile.read(content_length)
        try:
            payload = json.loads(payload_data.decode("utf-8"))
            data_received(payload)
        except json.JSONDecodeError:
            self.send_error(400, "Bad Request")
            return

        response = {"message": "Data Received"}
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(json.dumps(response).encode("utf-8"))

    def send_error(self, code, message=None):
        self.send_response(code)
        self.send_header("Content-Type", "text/plain")
        self.end_headers()
        if message:
            self.wfile.write(message.encode("utf-8"))

    def log_message(self, format, *args):
        pass


def data_received(payload):
    # Afficher les données reçues
    logging.info(f"Data Received: {payload}")

    # Charger les clés de surveillance depuis la variable d'environnement
    data_to_monitor = get_data_to_monitor()

    # Configurer le client Kubernetes
    config.load_incluster_config()  # Utilise la configuration du pod

    # Créer l'objet client
    custom_api = client.CustomObjectsApi()

    # Spécifier les informations sur la Custom Resource
    group = "closedlooppooc.closedloop.io"
    version = "v1"
    plural = "monitoringv2s"
    namespace = os.getenv("CLOSED_LOOP_MONITOR_NAMESPACE")

    # Récupérer la Custom Resource "Monitoringv2" spécifique
    name = os.getenv("CLOSED_LOOP_MONITOR_NAME")
    cr = custom_api.get_namespaced_custom_object(group, version, namespace, plural, name)
    logging.info(cr)
    cr["spec"]["time"] = str(datetime.now())

    # Effacer le contenu du champ spec.data
    if "spec" in cr:
        cr["spec"]["data"] = {}

    # Mettre à jour le champ spec.data avec les données de payload
    if "spec" in cr and "data" in cr["spec"]:
        for key, value in payload.items():
            if key in data_to_monitor.values():
                cr["spec"]["data"][key] = str(value)

    # Mettre à jour la Custom Resource
    updated_cr = custom_api.replace_namespaced_custom_object(group, version, namespace, plural, name, cr)
    logging.info(updated_cr)


def get_data_to_monitor():
    data_to_monitor = {}
    raw_data = os.getenv("CLOSED_LOOP_DATA_TO_MONITOR")
    if raw_data == "":
        raise ValueError("La variable d'environnement CLOSED_LOOP_DATA_TO_MONITOR est vide")

    try:
        data_to_monitor = json.loads(raw_data)
    except json.JSONDecodeError as e:
        raise ValueError(f"Erreur lors de la lecture de la variable d'environnement CLOSED_LOOP_DATA_TO_MONITOR : {e}")

    return data_to_monitor


def main():
    server_address = ("", 80)
    httpd = HTTPServer(server_address, RequestHandler)
    logging.basicConfig(format='%(asctime)s :: %(levelname)s :: %(message)s',level=logging.NOTSET)
    logging.info("Server running on port 80...")

    httpd.serve_forever()


if __name__ == "__main__":
    main()
