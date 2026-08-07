import json
import requests
import time

BASE_AUTH = "http://localhost:8082/api/v1/auth"
BASE_MGMT = "http://localhost:8080/api/v1"


def auth_register(email, password, role):
    payload = {"email": email, "password": password, "role": role}
    resp = requests.post(f"{BASE_AUTH}/register", json=payload)
    print("register", email, resp.status_code, resp.text)
    return resp


def auth_login(email, password):
    payload = {"email": email, "password": password}
    resp = requests.post(f"{BASE_AUTH}/login", json=payload)
    print("login", email, resp.status_code, resp.text)
    data = resp.json()
    return data["access_token"], data["refresh_token"]


def mgmt_post(path, payload, token):
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    resp = requests.post(f"{BASE_MGMT}{path}", json=payload, headers=headers)
    print(path, resp.status_code, resp.text)
    return resp


def main():
    admin_email = "admin_teste@cisterna.com"
    admin_password = "admin123"
    auth_register(admin_email, admin_password, "ADMIN")
    token, _ = auth_login(admin_email, admin_password)

    for i in range(3):
        cistern_payload = {
            "uuid": f"cisterna-{i+1}",
            "name": f"Cisterna {i+1}",
            "address": f"Rua {i+1}",
            "capacity": 1000 + i * 100,
            "status": "ATIVA"
        }
        mgmt_post("/cisterns/", cistern_payload, token)

    for i in range(3):
        truck_payload = {
            "uuid": f"caminhao-{i+1}",
            "plate": f"ABC{i+1:02d}",
            "model": "Truck X",
            "capacity": 5000 + i * 100,
            "status": "DISPONIVEL"
        }
        mgmt_post("/trucks/", truck_payload, token)

    for i in range(3):
        pipeiro_payload = {
            "uuid": f"pipeiro-{i+1}",
            "name": f"Pipeiro {i+1}",
            "cpf": f"{10000000000 + i}",
            "phone": f"99999{1000 + i}",
            "status": "ATIVO"
        }
        mgmt_post("/pipeiros/", pipeiro_payload, token)

    for i in range(3):
        delivery_payload = {
            "uuid": f"entrega-{i+1}",
            "cistern_uuid": f"cisterna-{i+1}",
            "truck_uuid": f"caminhao-{i+1}",
            "pipeiro_uuid": f"pipeiro-{i+1}",
            "status": "PENDENTE"
        }
        mgmt_post("/deliveries/", delivery_payload, token)


if __name__ == "__main__":
    main()
