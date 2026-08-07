import json
import time
import requests

BASE_AUTH = "http://localhost:8082/api/v1/auth"
BASE_TRACKING = "http://localhost:8080/api/v1/gps"


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
    return data["access_token"]


def load_positions(path):
    with open(path, "r", encoding="utf-8") as fh:
        return json.load(fh)


def send_positions(token, positions):
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    for idx, pos in enumerate(positions):
        payload = {
            "truck_id": "truck-sim-001",
            "latitude": pos["lat"],
            "longitude": pos["lon"],
            "timestamp": pos["timestamp"]
        }
        resp = requests.post(BASE_TRACKING, json=payload, headers=headers)
        print(idx, resp.status_code, resp.text)
        time.sleep(2)


def main():
    email = "pipeiro_teste@cisterna.com"
    password = "pipeiro123"
    auth_register(email, password, "PIPEIRO")
    token = auth_login(email, password)
    positions = load_positions("test-harness/gps_points.json")
    send_positions(token, positions)


if __name__ == "__main__":
    main()
