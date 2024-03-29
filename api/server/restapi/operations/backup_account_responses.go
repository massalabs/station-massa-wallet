// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/massalabs/station-massa-wallet/api/server/models"
)

// BackupAccountNoContentCode is the HTTP code returned for type BackupAccountNoContent
const BackupAccountNoContentCode int = 204

/*
BackupAccountNoContent Account backup

swagger:response backupAccountNoContent
*/
type BackupAccountNoContent struct {
}

// NewBackupAccountNoContent creates BackupAccountNoContent with default headers values
func NewBackupAccountNoContent() *BackupAccountNoContent {

	return &BackupAccountNoContent{}
}

// WriteResponse to the client
func (o *BackupAccountNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// BackupAccountBadRequestCode is the HTTP code returned for type BackupAccountBadRequest
const BackupAccountBadRequestCode int = 400

/*
BackupAccountBadRequest Bad request.

swagger:response backupAccountBadRequest
*/
type BackupAccountBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewBackupAccountBadRequest creates BackupAccountBadRequest with default headers values
func NewBackupAccountBadRequest() *BackupAccountBadRequest {

	return &BackupAccountBadRequest{}
}

// WithPayload adds the payload to the backup account bad request response
func (o *BackupAccountBadRequest) WithPayload(payload *models.Error) *BackupAccountBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the backup account bad request response
func (o *BackupAccountBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BackupAccountBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// BackupAccountUnauthorizedCode is the HTTP code returned for type BackupAccountUnauthorized
const BackupAccountUnauthorizedCode int = 401

/*
BackupAccountUnauthorized Unauthorized - The request requires user authentication.

swagger:response backupAccountUnauthorized
*/
type BackupAccountUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewBackupAccountUnauthorized creates BackupAccountUnauthorized with default headers values
func NewBackupAccountUnauthorized() *BackupAccountUnauthorized {

	return &BackupAccountUnauthorized{}
}

// WithPayload adds the payload to the backup account unauthorized response
func (o *BackupAccountUnauthorized) WithPayload(payload *models.Error) *BackupAccountUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the backup account unauthorized response
func (o *BackupAccountUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BackupAccountUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// BackupAccountNotFoundCode is the HTTP code returned for type BackupAccountNotFound
const BackupAccountNotFoundCode int = 404

/*
BackupAccountNotFound Account Not found.

swagger:response backupAccountNotFound
*/
type BackupAccountNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewBackupAccountNotFound creates BackupAccountNotFound with default headers values
func NewBackupAccountNotFound() *BackupAccountNotFound {

	return &BackupAccountNotFound{}
}

// WithPayload adds the payload to the backup account not found response
func (o *BackupAccountNotFound) WithPayload(payload *models.Error) *BackupAccountNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the backup account not found response
func (o *BackupAccountNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BackupAccountNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// BackupAccountUnprocessableEntityCode is the HTTP code returned for type BackupAccountUnprocessableEntity
const BackupAccountUnprocessableEntityCode int = 422

/*
BackupAccountUnprocessableEntity Unprocessable Entity - syntax is correct, but the server was unable to process the contained instructions.

swagger:response backupAccountUnprocessableEntity
*/
type BackupAccountUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewBackupAccountUnprocessableEntity creates BackupAccountUnprocessableEntity with default headers values
func NewBackupAccountUnprocessableEntity() *BackupAccountUnprocessableEntity {

	return &BackupAccountUnprocessableEntity{}
}

// WithPayload adds the payload to the backup account unprocessable entity response
func (o *BackupAccountUnprocessableEntity) WithPayload(payload *models.Error) *BackupAccountUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the backup account unprocessable entity response
func (o *BackupAccountUnprocessableEntity) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BackupAccountUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// BackupAccountInternalServerErrorCode is the HTTP code returned for type BackupAccountInternalServerError
const BackupAccountInternalServerErrorCode int = 500

/*
BackupAccountInternalServerError Internal Server Error - The server has encountered a situation it does not know how to handle.

swagger:response backupAccountInternalServerError
*/
type BackupAccountInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewBackupAccountInternalServerError creates BackupAccountInternalServerError with default headers values
func NewBackupAccountInternalServerError() *BackupAccountInternalServerError {

	return &BackupAccountInternalServerError{}
}

// WithPayload adds the payload to the backup account internal server error response
func (o *BackupAccountInternalServerError) WithPayload(payload *models.Error) *BackupAccountInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the backup account internal server error response
func (o *BackupAccountInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BackupAccountInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
