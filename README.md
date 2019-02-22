# Macaroons playground

This repository is a playground for a talk on macaroons.

It will implement an asset server, authentification service and a frontend.
If users authenticate successfully via the authentification service, they should receive a macaroon which lets them access pictures on the asset server. 

For simplicity the Frontend will first get the undischarged Macaroon from the Asset server.



## Flow
- user goes to asset and gets a macaroon
- needs to get discharge from auth
- once auth gives the discharge he can send both to asset and retrieve the picture


## Missing
The third party caveat should be securely negotiated between the asset and auth server. At the moment simple hardcoded strings are used, that do not provide real security.
To improve this, the Asset server can let the Auth server know about the ThirdPartyCaveat and its predicate by either sending a request for which the auth server will return an identifier, or just directly set the Identifier and inform the Auth server of this securely (f.e. by encrypting this information with the Auth serves public key).
