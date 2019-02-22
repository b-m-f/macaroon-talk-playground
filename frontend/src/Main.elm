import Browser
import Html exposing (Html, button, div, text)
import Html.Events exposing (onClick)
import Http 
import Json.Decode exposing (Decoder, field, string)
import Json.Encode
import Debug
import Url.Builder 

main =
    Browser.element {
      init = init
    , update = update
    , subscriptions = subscriptions
    , view = view
  }


-- MODEL

type Model
  = Failure String
  | Loading
  | DischargeMacaroon String
  | Success String
  | Base
  | Macaroon String

init : () -> (Model, Cmd Msg)
init _ =
 (
   Base,
   Cmd.none
 )

-- UPDATE

stringDecoder : Decoder String
stringDecoder =
  field "macaroon" string



encodeMacaroon: String -> Http.Body
encodeMacaroon macaroon =
    Http.jsonBody <|
        Json.Encode.object
            [ 
             ( "macaroon", Json.Encode.string macaroon)
            ]


type Msg = GetImage String | Login String | GetMacaroon | GotText (Result Http.Error String) | GotMacaroon (Result Http.Error String) | GotDischargeMacaroon (Result Http.Error String)

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    GetMacaroon ->
      (
        Loading, Http.get
        {
          url = "http://localhost:8080/macaroon"
        , expect = Http.expectJson GotMacaroon stringDecoder
        }
      )
    Login macaroon ->
      (
        Loading, Http.post
        {
          url = "http://localhost:9999/login"
        , body = (encodeMacaroon macaroon)
        , expect = Http.expectJson GotDischargeMacaroon stringDecoder
        }
      )
    GetImage macaroon ->
      (
        Loading, Http.get
        {
          url = "http://localhost:8080/get-image" ++ (Url.Builder.toQuery [Url.Builder.string "macaroon" macaroon])
        , expect = Http.expectString GotText
        }
      )
    GotMacaroon result ->
      case result of
        Ok macaroon->
          (Macaroon macaroon, Cmd.none)
        Err error ->
          (Failure "failed", Cmd.none)
    GotDischargeMacaroon result ->
      case result of
        Ok macaroon->
          (DischargeMacaroon macaroon, Cmd.none)
        Err error ->
          (Failure "failed", Cmd.none)
    GotText result ->
      case result of
        Ok image ->
          (Success image, Cmd.none)
        Err error ->
          (Failure "failed", Cmd.none)

-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none


-- VIEW

view : Model -> Html Msg
view model =
  case model of
    Base ->
      div []
        [ button [ onClick GetMacaroon] [ text "Get Macaroon" ],
          div [] [ text "First we need to get a macaroon"]
      ]

    Macaroon macaroon ->
      div []
        [ button [ onClick (Login macaroon)] [ text "Log in Alice" ],
          div [] [ text "Now we have received the asset server macaroon, which has a third party caveat. This caveat has a URL which tells us where we can satisfy this requirement. So lets do it"]
      ]

    Loading ->
      div []
        [
          div [] [ text "Loading Image"]
      ]
    Failure error ->
      div []
        [ button [ onClick GetMacaroon ] [ text "Get Macaroon" ],
          div [] [ text error, text " Try again from the start" ]
      ]
    Success image ->
      div []
          [ div [] [ text image ]
      ]
    DischargeMacaroon macaroon ->
      div []
        [ 
          button [ onClick (GetImage macaroon) ] [ text "Get Image for Alice" ],
          div [] [ text "Perfect, we have successfully authenticated and can now use our discharge macaroon, and the original one, to request an image from the asset server"]
      ]
