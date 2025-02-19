// Code generated by go-swagger; DO NOT EDIT.

package measurements

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"network.golang/measurements/models"
)

// GetMeasurementsOKCode is the HTTP code returned for type GetMeasurementsOK
const GetMeasurementsOKCode int = 200

/*
GetMeasurementsOK List of measurements

swagger:response getMeasurementsOK
*/
type GetMeasurementsOK struct {

	/*
	  In: Body
	*/
	Payload []*models.Measurement `json:"body,omitempty"`
}

// NewGetMeasurementsOK creates GetMeasurementsOK with default headers values
func NewGetMeasurementsOK() *GetMeasurementsOK {

	return &GetMeasurementsOK{}
}

// WithPayload adds the payload to the get measurements o k response
func (o *GetMeasurementsOK) WithPayload(payload []*models.Measurement) *GetMeasurementsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get measurements o k response
func (o *GetMeasurementsOK) SetPayload(payload []*models.Measurement) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetMeasurementsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.Measurement, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
