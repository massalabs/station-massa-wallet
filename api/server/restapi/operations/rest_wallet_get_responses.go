// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
)

// RestWalletGetOKCode is the HTTP code returned for type RestWalletGetOK
const RestWalletGetOKCode int = 200

/*
RestWalletGetOK Wallet retrieved.

swagger:response restWalletGetOK
*/
type RestWalletGetOK struct {

	/*
	  In: Body
	*/
	Payload *models.Wallet `json:"body,omitempty"`
}

// NewRestWalletGetOK creates RestWalletGetOK with default headers values
func NewRestWalletGetOK() *RestWalletGetOK {

	return &RestWalletGetOK{}
}

// WithPayload adds the payload to the rest wallet get o k response
func (o *RestWalletGetOK) WithPayload(payload *models.Wallet) *RestWalletGetOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the rest wallet get o k response
func (o *RestWalletGetOK) SetPayload(payload *models.Wallet) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RestWalletGetOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RestWalletGetBadRequestCode is the HTTP code returned for type RestWalletGetBadRequest
const RestWalletGetBadRequestCode int = 400

/*
RestWalletGetBadRequest Bad request.

swagger:response restWalletGetBadRequest
*/
type RestWalletGetBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewRestWalletGetBadRequest creates RestWalletGetBadRequest with default headers values
func NewRestWalletGetBadRequest() *RestWalletGetBadRequest {

	return &RestWalletGetBadRequest{}
}

// WithPayload adds the payload to the rest wallet get bad request response
func (o *RestWalletGetBadRequest) WithPayload(payload *models.Error) *RestWalletGetBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the rest wallet get bad request response
func (o *RestWalletGetBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RestWalletGetBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RestWalletGetNotFoundCode is the HTTP code returned for type RestWalletGetNotFound
const RestWalletGetNotFoundCode int = 404

/*
RestWalletGetNotFound Bad request.

swagger:response restWalletGetNotFound
*/
type RestWalletGetNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewRestWalletGetNotFound creates RestWalletGetNotFound with default headers values
func NewRestWalletGetNotFound() *RestWalletGetNotFound {

	return &RestWalletGetNotFound{}
}

// WithPayload adds the payload to the rest wallet get not found response
func (o *RestWalletGetNotFound) WithPayload(payload *models.Error) *RestWalletGetNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the rest wallet get not found response
func (o *RestWalletGetNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RestWalletGetNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RestWalletGetUnprocessableEntityCode is the HTTP code returned for type RestWalletGetUnprocessableEntity
const RestWalletGetUnprocessableEntityCode int = 422

/*
RestWalletGetUnprocessableEntity Unprocessable Entity - syntax is correct, but the server was unable to process the contained instructions.

swagger:response restWalletGetUnprocessableEntity
*/
type RestWalletGetUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewRestWalletGetUnprocessableEntity creates RestWalletGetUnprocessableEntity with default headers values
func NewRestWalletGetUnprocessableEntity() *RestWalletGetUnprocessableEntity {

	return &RestWalletGetUnprocessableEntity{}
}

// WithPayload adds the payload to the rest wallet get unprocessable entity response
func (o *RestWalletGetUnprocessableEntity) WithPayload(payload *models.Error) *RestWalletGetUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the rest wallet get unprocessable entity response
func (o *RestWalletGetUnprocessableEntity) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RestWalletGetUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RestWalletGetInternalServerErrorCode is the HTTP code returned for type RestWalletGetInternalServerError
const RestWalletGetInternalServerErrorCode int = 500

/*
RestWalletGetInternalServerError Internal Server Error - The server has encountered a situation it does not know how to handle.

swagger:response restWalletGetInternalServerError
*/
type RestWalletGetInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewRestWalletGetInternalServerError creates RestWalletGetInternalServerError with default headers values
func NewRestWalletGetInternalServerError() *RestWalletGetInternalServerError {

	return &RestWalletGetInternalServerError{}
}

// WithPayload adds the payload to the rest wallet get internal server error response
func (o *RestWalletGetInternalServerError) WithPayload(payload *models.Error) *RestWalletGetInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the rest wallet get internal server error response
func (o *RestWalletGetInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RestWalletGetInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
