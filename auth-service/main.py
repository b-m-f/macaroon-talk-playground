from flask import Flask, request
from pymacaroons import Macaroon, Verifier
import urllib.request
import json
import hashlib


identifiers = {}

# alice_macaroon = Macaroon.deserialize(macaroon)
#
# alice_macaroon.add_third_party_caveat("picture_id = picture_for_alice.jpg")
#
# print(alice_macaroon)
#
# alice = alice_macaroon.serialize()
#
# print(alice)

app = Flask(__name__)


@app.route("/get-identifier", methods=["POST"])
def add_identifier():
    data = json.loads(request.data)
    identifiers["needs_auth"] = data
    return "needs_auth"


@app.route("/login/<user>")
def login(user):
    if user == "alice":
        return identifiers["asset"]
    if user == "bob":
        return identifiers["asset"]
    return "Login failed"


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=2222)
