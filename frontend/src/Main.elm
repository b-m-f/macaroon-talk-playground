import Browser
import Html exposing (Html, button, div, text)
import Html.Events exposing (onClick)
import Http
import Json.Decode exposing (Decoder, field, string)
import Debug


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
  | Macaroon String
  | Success String
  | Base

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


type Msg = GetImageAlice String | LoginAlice | GotText (Result Http.Error String) | GotMacaroon (Result Http.Error String)

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    LoginAlice ->
      (
        Loading, Http.get
        {
          url = "http://localhost:2222/login/alice"
        , expect = Http.expectJson GotMacaroon stringDecoder
        }
      )
    GetImageAlice macaroon ->
      (
        Loading, Http.get
        {
          url = "http://localhost:1111/image/alice.jpg?macaroon=" ++ macaroon
        , expect = Http.expectString GotText
        }
      )
    GotMacaroon result ->
      case result of
        Ok macaroon->
          (Macaroon macaroon, Cmd.none)
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
        [ button [ onClick LoginAlice ] [ text "Log in Alice" ],
          div [] [ text "Welcome"]
      ]

    Loading ->
      div []
        [
          div [] [ text "Loading Image"]
      ]
    Failure error ->
      div []
        [ button [ onClick LoginAlice ] [ text "Log in Alice" ],
          div [] [ text error ]
      ]
    Success image ->
      div []
          [ div [] [ text image ]
      ]
    Macaroon macaroon ->
      div []
        [ 
          button [ onClick (GetImageAlice macaroon) ] [ text "Get Image for Alice" ],
          div [] [ text macaroon]
      ]
