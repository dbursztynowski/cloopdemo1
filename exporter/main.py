import json
import time
import random
import subprocess

command = "/etc/init.d/nginx start"
subprocess.run(command, shell=True)

def generer_objet_json():
    cpu = random.randint(1, 10)
    memory = random.randint(1, 100)
    disk = random.randint(1, 1000)
    return {"cpu": cpu, "memory": memory, "disk": disk}

liste_objets = []

while True:
    nouvel_objet = generer_objet_json()

    liste_objets.append(nouvel_objet)

    with open('/usr/share/nginx/html/index.html', 'w') as file:
        file.write(json.dumps(liste_objets))

    print("Données ajoutées avec succès !")

    time.sleep(30)
