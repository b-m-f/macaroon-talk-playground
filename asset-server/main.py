from flask import Flask, request, jsonify
from pymacaroons import Macaroon, Verifier
import urllib.request
import json

auth_url = "http://localhost:2222/"


keys = {"key-for-alice": "alice_secret"}

alice_macaroon = Macaroon(
    location="cool-assets",
    identifier="key-for-alice",
    key=keys["key-for-alice"],
)

alice_macaroon.add_first_party_caveat("picture_ids = alice.jpg")


caveat_key = "alice"
identifier = "alice"
alice_macaroon.add_third_party_caveat(auth_url, caveat_key, identifier)
print("ALICE MACAROON with 3rd party auth:")
print(alice_macaroon.inspect())


app = Flask(__name__)

# get macaroon for user with third party caveat


@app.route("/get-unauthorized-macaroon/<user>")
def get_access(user):
    try:
        if user == "alice":
            return alice_macaroon.serialize()

    except Exception:
        return "Error"


@app.route("/image/alice.jpg")
def get_alice_image():
    macaroon = Macaroon.deserialize(request.args["macaroon"])
    # print("ALICE MACAROON that is authorized")
    print(macaroon.inspect())
    # print(macaroon.signature_bytes)
    ## verfiy macaroon
    verifier = Verifier(discharge_macaroons=[macaroon])
    is_verified = verifier.verify(
        alice_macaroon, keys[alice_macaroon.identifier]
    )
    print(is_verified)
    # if ok do
    # else throw
    return "True"


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=1111)
