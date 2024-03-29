// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/massalabs/station-massa-wallet/api/server/models"
)

// ExportAccountFileOKCode is the HTTP code returned for type ExportAccountFileOK
const ExportAccountFileOKCode int = 200

/*
ExportAccountFileOK Download the account file

swagger:response exportAccountFileOK
*/
type ExportAccountFileOK struct {

	/*
	  In: Body
	*/
	Payload io.ReadCloser `json:"body,omitempty"`
}

// NewExportAccountFileOK creates ExportAccountFileOK with default headers values
func NewExportAccountFileOK() *ExportAccountFileOK {

	return &ExportAccountFileOK{}
}

// WithPayload adds the payload to the export account file o k response
func (o *ExportAccountFileOK) WithPayload(payload io.ReadCloser) *ExportAccountFileOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the export account file o k response
func (o *ExportAccountFileOK) SetPayload(payload io.ReadCloser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExportAccountFileOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ExportAccountFileBadRequestCode is the HTTP code returned for type ExportAccountFileBadRequest
const ExportAccountFileBadRequestCode int = 400

/*
ExportAccountFileBadRequest Bad request

swagger:response exportAccountFileBadRequest
*/
type ExportAccountFileBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewExportAccountFileBadRequest creates ExportAccountFileBadRequest with default headers values
func NewExportAccountFileBadRequest() *ExportAccountFileBadRequest {

	return &ExportAccountFileBadRequest{}
}

// WithPayload adds the payload to the export account file bad request response
func (o *ExportAccountFileBadRequest) WithPayload(payload *models.Error) *ExportAccountFileBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the export account file bad request response
func (o *ExportAccountFileBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExportAccountFileBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ExportAccountFileNotFoundCode is the HTTP code returned for type ExportAccountFileNotFound
const ExportAccountFileNotFoundCode int = 404

/*
ExportAccountFileNotFound Not found.

swagger:response exportAccountFileNotFound
*/
type ExportAccountFileNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewExportAccountFileNotFound creates ExportAccountFileNotFound with default headers values
func NewExportAccountFileNotFound() *ExportAccountFileNotFound {

	return &ExportAccountFileNotFound{}
}

// WithPayload adds the payload to the export account file not found response
func (o *ExportAccountFileNotFound) WithPayload(payload *models.Error) *ExportAccountFileNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the export account file not found response
func (o *ExportAccountFileNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExportAccountFileNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ExportAccountFileInternalServerErrorCode is the HTTP code returned for type ExportAccountFileInternalServerError
const ExportAccountFileInternalServerErrorCode int = 500

/*
ExportAccountFileInternalServerError Internal Server Error - The server has encountered a situation it does not know how to handle.

swagger:response exportAccountFileInternalServerError
*/
type ExportAccountFileInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewExportAccountFileInternalServerError creates ExportAccountFileInternalServerError with default headers values
func NewExportAccountFileInternalServerError() *ExportAccountFileInternalServerError {

	return &ExportAccountFileInternalServerError{}
}

// WithPayload adds the payload to the export account file internal server error response
func (o *ExportAccountFileInternalServerError) WithPayload(payload *models.Error) *ExportAccountFileInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the export account file internal server error response
func (o *ExportAccountFileInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExportAccountFileInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
