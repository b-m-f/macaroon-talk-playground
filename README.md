# Macaroons playground

This repository is a playground for a talk on macaroons.

It will implement an asset server, authentification service and a frontend for 2 users.
If users authenticate successfully via the authentification service, they should receive a macaroon which lets them access pictures on the asset server.

The authentification service will receive a macaroon to access all pictures and adds third party caveats on top of this to delegate certain permissions to the specific users on the frontend.

## Flow
- user goes to asset and gets a macaroon
- needs to get discharge from auth
- once auth gives the discharge he can send both to asset and retrieve the picture
