import requests

BASE_AUTH = "http://localhost:8082/api/v1/auth"


def auth_register(email, password, role):
    payload = {"email": email, "password": password, "role": role}
    resp = requests.post(f"{BASE_AUTH}/register", json=payload)
    print("register", email, resp.status_code, resp.text)
    return resp


def auth_login(email, password):
    payload = {"email": email, "password": password}
    resp = requests.post(f"{BASE_AUTH}/login", json=payload)
    print("login", email, resp.status_code, resp.text)
    return resp.json()


def main():
    email = "cidadao_teste@cisterna.com"
    password = "cidadao123"
    auth_register(email, password, "CIDADAO")
    auth_login(email, password)


if __name__ == "__main__":
    main()
