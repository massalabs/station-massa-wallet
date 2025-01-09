// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/massalabs/station-massa-wallet/api/server/models"
)

// AddSignRuleOKCode is the HTTP code returned for type AddSignRuleOK
const AddSignRuleOKCode int = 200

/*
AddSignRuleOK signRule Id.

swagger:response addSignRuleOK
*/
type AddSignRuleOK struct {

	/*
	  In: Body
	*/
	Payload *models.AddSignRuleResponse `json:"body,omitempty"`
}

// NewAddSignRuleOK creates AddSignRuleOK with default headers values
func NewAddSignRuleOK() *AddSignRuleOK {

	return &AddSignRuleOK{}
}

// WithPayload adds the payload to the add sign rule o k response
func (o *AddSignRuleOK) WithPayload(payload *models.AddSignRuleResponse) *AddSignRuleOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add sign rule o k response
func (o *AddSignRuleOK) SetPayload(payload *models.AddSignRuleResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddSignRuleOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddSignRuleBadRequestCode is the HTTP code returned for type AddSignRuleBadRequest
const AddSignRuleBadRequestCode int = 400

/*
AddSignRuleBadRequest Bad request.

swagger:response addSignRuleBadRequest
*/
type AddSignRuleBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddSignRuleBadRequest creates AddSignRuleBadRequest with default headers values
func NewAddSignRuleBadRequest() *AddSignRuleBadRequest {

	return &AddSignRuleBadRequest{}
}

// WithPayload adds the payload to the add sign rule bad request response
func (o *AddSignRuleBadRequest) WithPayload(payload *models.Error) *AddSignRuleBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add sign rule bad request response
func (o *AddSignRuleBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddSignRuleBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddSignRuleUnauthorizedCode is the HTTP code returned for type AddSignRuleUnauthorized
const AddSignRuleUnauthorizedCode int = 401

/*
AddSignRuleUnauthorized Unauthorized - The request requires user authentication.

swagger:response addSignRuleUnauthorized
*/
type AddSignRuleUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddSignRuleUnauthorized creates AddSignRuleUnauthorized with default headers values
func NewAddSignRuleUnauthorized() *AddSignRuleUnauthorized {

	return &AddSignRuleUnauthorized{}
}

// WithPayload adds the payload to the add sign rule unauthorized response
func (o *AddSignRuleUnauthorized) WithPayload(payload *models.Error) *AddSignRuleUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add sign rule unauthorized response
func (o *AddSignRuleUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddSignRuleUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddSignRuleNotFoundCode is the HTTP code returned for type AddSignRuleNotFound
const AddSignRuleNotFoundCode int = 404

/*
AddSignRuleNotFound Account Not found.

swagger:response addSignRuleNotFound
*/
type AddSignRuleNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddSignRuleNotFound creates AddSignRuleNotFound with default headers values
func NewAddSignRuleNotFound() *AddSignRuleNotFound {

	return &AddSignRuleNotFound{}
}

// WithPayload adds the payload to the add sign rule not found response
func (o *AddSignRuleNotFound) WithPayload(payload *models.Error) *AddSignRuleNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add sign rule not found response
func (o *AddSignRuleNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddSignRuleNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddSignRuleUnprocessableEntityCode is the HTTP code returned for type AddSignRuleUnprocessableEntity
const AddSignRuleUnprocessableEntityCode int = 422

/*
AddSignRuleUnprocessableEntity Unprocessable Entity - syntax is correct, but the server was unable to process the contained instructions.

swagger:response addSignRuleUnprocessableEntity
*/
type AddSignRuleUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddSignRuleUnprocessableEntity creates AddSignRuleUnprocessableEntity with default headers values
func NewAddSignRuleUnprocessableEntity() *AddSignRuleUnprocessableEntity {

	return &AddSignRuleUnprocessableEntity{}
}

// WithPayload adds the payload to the add sign rule unprocessable entity response
func (o *AddSignRuleUnprocessableEntity) WithPayload(payload *models.Error) *AddSignRuleUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add sign rule unprocessable entity response
func (o *AddSignRuleUnprocessableEntity) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddSignRuleUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddSignRuleInternalServerErrorCode is the HTTP code returned for type AddSignRuleInternalServerError
const AddSignRuleInternalServerErrorCode int = 500

/*
AddSignRuleInternalServerError Internal Server Error - The server has encountered a situation it does not know how to handle.

swagger:response addSignRuleInternalServerError
*/
type AddSignRuleInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddSignRuleInternalServerError creates AddSignRuleInternalServerError with default headers values
func NewAddSignRuleInternalServerError() *AddSignRuleInternalServerError {

	return &AddSignRuleInternalServerError{}
}

// WithPayload adds the payload to the add sign rule internal server error response
func (o *AddSignRuleInternalServerError) WithPayload(payload *models.Error) *AddSignRuleInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add sign rule internal server error response
func (o *AddSignRuleInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddSignRuleInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
