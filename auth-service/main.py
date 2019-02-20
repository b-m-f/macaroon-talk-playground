from flask import Flask, request, jsonify, Response
from pymacaroons import Macaroon, Verifier
import urllib.request
import json
import hashlib
import datetime


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
    now = datetime.datetime.now()
    now_plus_1 = (now + datetime.timedelta(minutes=1)).strftime(
        "%Y-%m-%d %H:%M:%S"
    )
    if user == "alice":
        # get unauthorized macaroon
        macaroon = Macaroon.deserialize(
            urllib.request.urlopen(
                "http://localhost:1111/get-unauthorized-macaroon/alice"
            ).read()
        )
        # create a macaroon that confirms authentification
        # by using the caveat_key that was
        # agreed on with the asset server when it started
        # and added the third party caveat
        discharge = Macaroon(
            "http://localhost:2222/",
            identifiers["needs_auth"]["caveat_key"],
            "needs_auth",
        )
        discharge.add_first_party_caveat("time < " + now_plus_1)
        macaroon.prepare_for_request(discharge)
        response = jsonify({"macaroon": macaroon.serialize()})
        response.headers.add("Access-Control-Allow-Origin", "*")
        return response
    if user == "bob":
        #        macaroon = urllib.request.urlopen(
        #            "http://localhost:1111/get-access/bob"
        #        ).read()
        response = jsonify({"status": "success"})
        response.headers.add("Access-Control-Allow-Origin", "*")
        return response
    response = jsonify({"status": "failed"})
    response.headers.add("Access-Control-Allow-Origin", "*")
    return response


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=2222)
