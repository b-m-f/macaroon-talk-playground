from flask import Flask
from pymacaroons import Macaroon, Verifier

keys = {
    "key-for-alice": "WOW-a-very-secret-signing-key",
    "key-for-bob": "WOW-another-very-secret-signing-key",
}

macaroon1 = Macaroon(
    location="cool-picture-service.example.com",
    identifier="key-for-alice",
    key=keys["key-for-alice"],
)
macaroon2 = Macaroon(
    location="cool-picture-service.example.com",
    identifier="key-for-bob",
    key=keys["key-for-bob"],
)

macaroon1.add_first_party_caveat("picture_id = picture_for_alice.jpg")
macaroon2.add_first_party_caveat("picture_id = picture_for_bob.jpg")

print(macaroon1.inspect())
print(macaroon2.inspect())

serialized1 = macaroon1.serialize()
serialized2 = macaroon1.serialize()

print(serialized1)
print(serialized2)

app = Flask(__name__)


@app.route("/macaroon-for-alice")
def macaroon_1():
    return serialized1


@app.route("/macaroon-for-bob")
def macaroon_2():
    return serialized2


if __name__ == "__main__":
    app.run(host="0.0.0.0")
