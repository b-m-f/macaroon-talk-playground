# Macaroons playground

This repository is a playground for a talk on macaroons.

It will implement an asset server, authentification service and a frontend.
If users authenticate successfully via the authentification service, they should receive a macaroon which lets them access pictures on the asset server. 

For simplicity the Frontend will first get the undischarged Macaroon from the Asset server.


## Running the code
You will need both [golang]() and [elm](https://guide.elm-lang.org/install.html) installed.
Go was used due to its rapid development capabilities and Types providing good compiler feedback for understanding the Macaroon library implementation better.
Elm lang has very good Refactoring capabilities, which made it a good fit for exploring and experimenting.


Run both the asset and auth server with `go run`.
The Frontend can be build with `elm make src/Main.elm` which provides an `index.html `or loaded for development with `elm reactor`, in which case it will be at `http://localhost:8000/src/Main.elm`.


## Flow
- user goes to asset and gets a macaroon
- needs to get discharge from auth
- once auth gives the discharge he can send both to asset and retrieve the picture


## Missing
The third party caveat should be securely negotiated between the asset and auth server. At the moment simple hardcoded strings are used, that do not provide real security.
To improve this, the Asset server can let the Auth server know about the ThirdPartyCaveat and its predicate by either sending a request for which the auth server will return an identifier, or just directly set the Identifier and inform the Auth server of this securely (f.e. by encrypting this information with the Auth serves public key).
