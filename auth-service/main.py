from flask import Flask, request, jsonify, Response
import urllib.request
import json
import hashlib
import datetime
from pymacaroons import Macaroon, Verifier


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
    # save hash for every user here
    data = request.data.decode("utf-8")
    json_data = json.loads(data)
    # hash_string = hashlib.md5(data.encode("utf-8")).hexdigest()
    identifiers["alice"] = json_data
    return "alice"


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
        print(macaroon.inspect())
        # create a macaroon that confirms authentification
        # by using the caveat_key that was
        # agreed on with the asset server when it started
        # and added the third party caveat
        # to do this we first need to get the identifier out
        # of the macaroon,
        # usually a service that wants to satisy these
        # would go through all 3rd party caveats
        # check their location
        # and then authorize at the correct service
        discharge = Macaroon(
            location="http://localhost:2222/", key="alice", identifier="alice"
        )

        discharge.add_first_party_caveat("time < " + now_plus_1)
        bound_macaroon = macaroon.prepare_for_request(discharge)
        print(bound_macaroon.inspect())
        verifier = Verifier(discharge_macaroons=[macaroon])
        is_verified = verifier.verify(macaroon, "key-for-alice")
        print(is_verified)

        response = jsonify({"macaroon": bound_macaroon.serialize()})
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
    app.run(host="1.0.0.0", port=2222)
