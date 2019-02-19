import Browser
import Html exposing (Html, button, div, text)
import Html.Events exposing (onClick)
import Http


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
  | Success String

init : () -> (Model, Cmd Msg)
init _ =
 (
   Loading,
   Cmd.none
 )

-- UPDATE

type Msg = GetImageAlice | LoginAlice | GotText (Result Http.Error String)

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    LoginAlice ->
      (
        Loading, Http.get
        {
          url = "http://localhost:2222/login/alice"
        , expect = Http.expectString GotText
        }
      )
    GetImageAlice ->
      (
        Loading, Http.get
        {
          url = "http://localhost:1111/image/alice.jpg"
        , expect = Http.expectString GotText
        }
      )
    GotText result ->
      case result of
        Ok image ->
          (Success image, Cmd.none)
        Err error ->
          case error of
            Http.NetworkError ->
              (Failure (Debug.toString error), Cmd.none)
            _ ->
              (Failure "failed", Cmd.none)

-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none


-- VIEW

view : Model -> Html Msg
view model =
  case model of
    Loading ->
      div []
        [ button [ onClick LoginAlice ] [ text "Log in Alice" ],
          button [ onClick GetImageAlice ] [ text "Get Image for Alice" ],
        div [] [ text "Loading Image"]
      ]
    Failure error ->
      div []
        [ button [ onClick LoginAlice ] [ text "Log in Alice" ],
          button [ onClick GetImageAlice ] [ text "Get Image for Alice" ],
          div [] [ text error ]
      ]
    Success image ->
      div []
        [ button [ onClick LoginAlice ] [ text "Log in Alice" ],
          button [ onClick GetImageAlice ] [ text "Get Image for Alice" ],
          div [] [ text image ]
      ]
