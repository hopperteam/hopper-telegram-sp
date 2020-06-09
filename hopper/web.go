package hopper

import (
	"github.com/hopperteam/hopper-telegram-sp/model"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

var logoBytes = readLogo()

func StartWebEndpoint() {
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/logo.svg", serveLogo)
	err := http.ListenAndServe(":80", nil)
	log.Fatal().Err(err).Msg("Error while starting web server")
}

func readLogo() []byte {
	dat, err := ioutil.ReadFile("res/logo.svg")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not read logo")
	}
	return dat
}

func serveLogo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	_, _ = w.Write(logoBytes)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["chat"]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request"))
		return
	}

	key := keys[0]
	chat, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request"))
		return
	}

	keys, ok = r.URL.Query()["status"]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request"))
		return
	}

	if keys[0] != "success" {
		_, _ = w.Write([]byte("Request declined"))
		return
	}

	keys, ok = r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request"))
		return
	}
	subscriptionId := keys[0]

	err = model.InsertChatSubscription(chat, subscriptionId)
	if err != nil {
		log.Err(err).Msg("Could not insert chat subscription")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Subscribed!"))
}
