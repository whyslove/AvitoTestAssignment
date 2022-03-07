import json
from multiprocessing.context import assert_spawning
import pytest
import requests
from datetime import datetime

# Добавить юзера с id и проверить получение с валютой
def test_add_user_and_valutes_check():
    requests.post(
        "http://localhost:8000/api/balance/operation", json={"id": 1, "amount_money": 100}
    )  # add user with id = 1 and his money

    r = requests.get("http://localhost:8000/api/balance/get/1?currency=RUB")
    assert r.json() == {"id": "1", "balance": 100}  # В рублях

    r = requests.get("http://localhost:8000/api/balance/get/1?currency=USD")
    assert (
        r.json().get("balance") > 0
    )  # не могу проверить точно, по скольку не ясен курс доллара, но главное, что оно выводит доллары


def test_transactions():
    requests.post(
        "http://localhost:8000/api/balance/operation", json={"id": 2, "amount_money": 200}
    )  # add new user and his money

    r = requests.post(
        "http://localhost:8000/api/balance/transfer",
        json={"sender_id": 1, "receiver_id": 2, "amount_money": 99},
    )
    assert r.json() == {"status": "ok"}  # he has enough money

    # next give them him back
    r = requests.post(
        "http://localhost:8000/api/balance/transfer",
        json={"sender_id": 2, "receiver_id": 1, "amount_money": 99},
    )
    # and try same request with currency in usd
    r = requests.post(
        "http://localhost:8000/api/balance/transfer?currency=USD",
        json={"sender_id": 1, "receiver_id": 2, "amount_money": 99},
    )
    assert r.json() == {"message": "id=1 has not enough money "}


def test_draw_off_money():
    r = requests.post(
        "http://localhost:8000/api/balance/operation", json={"id": 2, "amount_money": -400}
    )
    assert r.json() == {"message": "not enough money for this operation"}
    r = requests.post(
        "http://localhost:8000/api/balance/operation", json={"id": 2, "amount_money": 400}
    )
    assert r.json() == {"status": "ok"}


def test_transactions_stat():
    # simple story
    r = requests.get("http://localhost:8000/api/operations/1")
    # i check two fields because i can not check all jsons due to the executed_at. I can not specify this value
    assert r.json()[0]["id"] == 1
    assert r.json()[1]["amount_of_money"] == 99
    # Full output
    # [
    #     {
    #         "id": 1,
    #         "main_subject_id": 1,
    #         "other_subject_id": null,
    #         "executed_at": "0000-01-01T19:10:47.713932Z",
    #         "amount_of_money": 100,
    #     },
    #     {
    #         "id": 3,
    #         "main_subject_id": 1,
    #         "other_subject_id": 2,
    #         "executed_at": "0000-01-01T19:10:47.734199Z",
    #         "amount_of_money": 99,
    #     },
    #     {
    #         "id": 6,
    #         "main_subject_id": 1,
    #         "other_subject_id": 2,
    #         "executed_at": "0000-01-01T19:10:47.741527Z",
    #         "amount_of_money": -99,
    #     },
    # ]

    # lets add some more operations
    for i in range(10):
        r = requests.post(
            "http://localhost:8000/api/balance/operation", json={"id": 1, "amount_money": 10}
        )
    # and test pagingation
    r = requests.get("http://localhost:8000/api/operations/1?page=1")
    assert len(r.json()) == 10  # can not check all json's due to problem metioned bellow
    r = requests.get("http://localhost:8000/api/operations/1?page=2")
    assert len(r.json()) == 3


def test_transactions_sort():
    r = requests.get("http://localhost:8000/api/operations/1?sort=summ&page=1")
    r = r.json()

    i = 0
    for i in range(len(r) - 2):
        assert r[i]["amount_of_money"] >= r[i + 1]["amount_of_money"]

    r = requests.get("http://localhost:8000/api/operations/1?sort=date&page=2")  # use 2 page
    r = r.json()
    i = 0
    for i in range(len(r) - 2):
        assert datetime.strptime(r[i]["executed_at"], "%Y-%m-%dT%H:%M:%S.%fZ") >= datetime.strptime(
            r[i + 1]["executed_at"], "%Y-%m-%dT%H:%M:%S.%fZ"
        )
