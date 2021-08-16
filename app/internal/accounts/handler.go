package accounts

import (
	"encoding/json"
	"fmt"

	"github.com/charopevez/eob-accountant-worker/internal/apperror"
	"github.com/charopevez/eob-accountant-worker/pkg/logging"
	"github.com/julienschmidt/httprouter"

	"net/http"
)

const (
	registerURL = "/api/register"
	accountURL  = "/api/account/:uuid"
	loginURL    = "/api/login"
)

type Handler struct {
	Logger            logging.Logger
	AccountantService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, loginURL, apperror.Middleware(h.Authenticate))
	router.HandlerFunc(http.MethodPost, registerURL, apperror.Middleware(h.CreateAccount))
	router.HandlerFunc(http.MethodGet, accountURL, apperror.Middleware(h.GetAccount))
	router.HandlerFunc(http.MethodPatch, accountURL, apperror.Middleware(h.UpdateAccount))
	router.HandlerFunc(http.MethodPut, accountURL, apperror.Middleware(h.UpdateCredentials))
	router.HandlerFunc(http.MethodDelete, accountURL, apperror.Middleware(h.DeleteAccount))
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET USER ACCOUNT BY EMAIL AND PASSWORD")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("decode credentials dto")
	var cred CredentialsDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}

	account, err := h.AccountantService.AuthenticateAccount(r.Context(), cred)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshal user account")
	accountBytes, err := json.Marshal(account)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(accountBytes)

	return nil
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("CREATE USER ACCOUNT")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("decode create account dto")
	var crAcc CreateAccountDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&crAcc); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}

	accountUUID, err := h.AccountantService.Create(r.Context(), crAcc)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%s", registerURL, accountUUID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *Handler) UpdateCredentials(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("UPDATE USER CREDENTIALS")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	accountUUID := params.ByName("uuid")

	h.Logger.Debug("decode update credentials dto")
	var updAccount UpdateCredentialsDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&updAccount); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}
	updAccount.UUID = accountUUID

	err := h.AccountantService.UpdateCredentials(r.Context(), updAccount)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("UPDATE USER ACCOUNT")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	accountUUID := params.ByName("uuid")

	h.Logger.Debug("decode update account dto")
	var updAccount UpdateAccountDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&updAccount); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}
	updAccount.UUID = accountUUID

	err := h.AccountantService.UpdateAccount(r.Context(), updAccount)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET ACCOUNT")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	accountUUID := params.ByName("uuid")

	account, err := h.AccountantService.GetAccount(r.Context(), accountUUID)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshal account")
	accountBytes, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("failed to marshall account. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(accountBytes)
	return nil
}

func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE ACCOUNT")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	accountUUID := params.ByName("uuid")

	err := h.AccountantService.Delete(r.Context(), accountUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
