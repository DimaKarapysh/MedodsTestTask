package delivery

import (
	"MedodsTestTask/app/core"
	"MedodsTestTask/app/rest"
	"MedodsTestTask/domain"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

type DeliveryService struct {
	log  core.Logger
	iter domain.Iter
	v    validator.Validate
}

func NewDeliveryService(log core.Logger, i domain.Iter, v validator.Validate) *DeliveryService {
	return &DeliveryService{
		log:  log,
		iter: i,
		v:    v,
	}
}

func (t *DeliveryService) GetById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := uuid.Parse(id)
	if err != nil {
		rest.ValidationError(w, "Invalid GUID format")
		return
	}
	Id := Convert(id)
	if Id == -1 {
		rest.ValidationError(w, "Can not parse Id")
		return
	}

	ip := r.RemoteAddr

	rToken, aToken, err := t.iter.GetById(Id, ip)
	if err != nil {
		rest.ServerError(w, err)
		return
	}

	rest.ServerSuccessStruct(w, DTOTokens(DTOR(rToken), DTOA(aToken)))
}

func (t *DeliveryService) AddUser(w http.ResponseWriter, r *http.Request) {
	var User User
	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		rest.ValidationError(w, "Cannot parse json")
		return
	}

	err = t.v.Struct(User)
	if err != nil {
		rest.ValidationError(w, err.Error())
		return
	}

	err = t.iter.InsertUser(User.DTOUser(r))
	if err != nil {
		rest.ServerError(w, errors.Wrap(err, "AddTask"))
		return
	}

	rest.ServerSuccessStruct(w, User)
}

func (t *DeliveryService) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	tokens := &TokenPairDTO{AccessToken: AccessTokenDTO{}, RefreshToken: RefreshTokenDTO{}}

	err := json.NewDecoder(r.Body).Decode(&tokens)
	if err != nil {
		rest.ValidationError(w, "Cannot parse json")
		return
	}

	err = t.v.Struct(tokens)
	if err != nil {
		rest.ValidationError(w, err.Error())
		return
	}

	tok := tokens.Dto()
	ip := r.RemoteAddr

	rToken, aToken, err := t.iter.Refresh(tok.RefreshToken.Token, tok.AccessToken.Token, ip)
	if err != nil {
		rest.ServerError(w, err)
		return
	}

	rest.ServerSuccessStruct(w, DTOTokens(DTOR(rToken), DTOA(aToken)))
}
