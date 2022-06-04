package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/MetroReviews/metro-integrase/types"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

/*
 * Package ``lib`` provides the main code for the integrase system.
 */

func authReq(r *http.Request, cfg types.ListConfig) bool {
	if r == nil {
		if cfg.RequestLogs {
			log.Error("Request is nil")
		}
		return false
	}

	if r.Header.Get("Authorization") == "" || r.Header.Get("Authorization") != cfg.SecretKey {
		if cfg.RequestLogs {
			log.Error("Authorization header is missing or invalid")
		}
		return false
	}

	return true
}

type ListFunction func(bot *types.Bot) error

func coreHandler(fn ListFunction, cfg types.ListConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}

		if !authReq(r, cfg) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			log.Warning("Unauthorized request: ", r.URL.Path)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			if cfg.RequestLogs {
				log.Error(err, " at URL ", r.URL.Path)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request. See logs for more information if you have them enabled:"))
			return
		}

		var bot types.Bot
		err = json.Unmarshal(body, &bot)

		if err != nil {
			if cfg.RequestLogs {
				log.Error(err, " at URL ", r.URL.Path)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Serialization error occured. See logs for more information if you have them enabled"))
			return
		}

		adpErr := fn(&bot)

		if adpErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Request handle error:" + adpErr.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK :)"))
	}
}

// Patches a list
func PatchList(cfg types.ListConfig, data types.ListPatch) (*types.ListPatchResp, error) {
	payload, err := json.Marshal(data)

	if err != nil && cfg.RequestLogs {
		log.Error("PatchList error:", err)
		return nil, err
	}

	client := &http.Client{Timeout: 20 * time.Second}

	req, err := http.NewRequest("PATCH", types.APIUrl+"/lists/"+cfg.ListID, bytes.NewBuffer(payload))

	if err != nil && cfg.RequestLogs {
		log.Error("PatchList error (in making new request):", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", cfg.SecretKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Error("PatchList error (in performing new request):", err)
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var respData types.ListPatchResp

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Error("PatchList error (in reading response body):", err)
			return nil, err
		}

		err = json.Unmarshal(body, &respData)

		if err != nil {
			log.Error("PatchList error (in unmarshalling response body):", err)
			return nil, err
		}

		return &respData, nil
	} else {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Error("PatchList error (in reading response body):", err)
			return nil, err
		}

		return nil, errors.New(string(body))
	}
}

// Starts a web server handling all core integrase functions
func StartServer(adp types.ListAdapter, r *mux.Router) {
	cfg := adp.GetConfig()

	if cfg.StartupLogs {
		log.Info("Starting integrase server")
	}

	if cfg.ListID == "" {
		panic("List ID not set")
	} else if cfg.SecretKey == "" {
		panic("Secret Key not set")
	}

	r.HandleFunc("/claim", coreHandler(adp.ClaimBot, cfg))
	r.HandleFunc("/unclaim", coreHandler(adp.UnclaimBot, cfg))
	r.HandleFunc("/approve", coreHandler(adp.ApproveBot, cfg))
	r.HandleFunc("/deny", coreHandler(adp.DenyBot, cfg))
	r.HandleFunc("/data-request", func(w http.ResponseWriter, r *http.Request) {
		if !authReq(r, cfg) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		botId := r.URL.Query().Get("bot_id")

		if botId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bot ID is missing"))
			return
		}

		bot, err := adp.DataRequest(botId)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Request handle error:" + err.Error()))
			return
		}

		botStr, err := json.Marshal(bot)

		if err != nil {
			if cfg.RequestLogs {
				log.Error(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Serialization error occured. See logs for more information if you have them enabled"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(botStr))
	})

	r.HandleFunc("/data-delete", func(w http.ResponseWriter, r *http.Request) {
		if !authReq(r, cfg) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		botId := r.URL.Query().Get("bot_id")

		if botId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bot ID is missing"))
			return
		}

		err := adp.DataDelete(botId)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Request handle error:" + err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("All associated data has been deleted from this list according to the lists adapter"))
	})

	if cfg.DomainName != "" {
		if cfg.StartupLogs {
			log.Info("Updating Metro Reviews with new routes")
		}

		patched, err := PatchList(cfg, types.ListPatch{
			ClaimBotAPI:     cfg.DomainName + "/claim",
			UnclaimBotAPI:   cfg.DomainName + "/unclaim",
			ApproveBotAPI:   cfg.DomainName + "/approve",
			DenyBotAPI:      cfg.DomainName + "/deny",
			DataRequestAPI:  cfg.DomainName + "/data-request",
			DataDeletionAPI: cfg.DomainName + "/data-delete",
		})

		if err != nil {
			log.Error("Metro Reviews update failed: ", err)

		} else {
			log.Info("Metro Reviews update successful with ", patched.HasUpdated, " updated")
		}
	}

	if cfg.BindAddr == "" {
		cfg.BindAddr = ":8080"
	}

	if cfg.StartupLogs {
		log.Info("Integrase server now going to start listening on address ", cfg.BindAddr)
	}

	handledRouter := handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(r))

	err := http.ListenAndServe(cfg.BindAddr, handledRouter)

	if err != nil {
		log.Error("Integrase server error: ", err)
	}
}
