package service

import (
	"cache/internal/model"
	"cache/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Prepare() *service.CacheService {
	return service.New()
}

func TestGetNotExistentKeyValue(t *testing.T) {
	s := Prepare()
	getResultChan := s.ReceiveCommand(model.NewGetCommand(model.Key("test1")))
	getResult := <-getResultChan

	assert.Equal(t, false, getResult.Ok)
	assert.Equal(t, nil, getResult.Value.EncodedRawValue)
}

func TestSetAndGetKeyValue(t *testing.T) {
	s := Prepare()
	setResultChan := s.ReceiveCommand(model.NewSetCommand(model.Key("test1"), 2523, 0))
	setResult := <-setResultChan
	assert.NoError(t, setResult.Error)

	getResultChan := s.ReceiveCommand(model.NewGetCommand(model.Key("test1")))
	getResult := <-getResultChan

	assert.Equal(t, true, getResult.Ok)
	assert.Equal(t, 2523, getResult.Value.EncodedRawValue)
}

func TestDeleteKeyValue(t *testing.T) {
	s := Prepare()
	setResultChan := s.ReceiveCommand(model.NewSetCommand(model.Key("test1"), 2523, 0))
	setResult := <-setResultChan
	assert.NoError(t, setResult.Error)

	getResultChan := s.ReceiveCommand(model.NewGetCommand(model.Key("test1")))
	getResult := <-getResultChan

	assert.Equal(t, true, getResult.Ok)
	assert.Equal(t, 2523, getResult.Value.EncodedRawValue)

	delResultChan := s.ReceiveCommand(model.NewDeleteCommand(model.Key("test1")))
	delResult := <-delResultChan
	assert.Equal(t, true, delResult.Ok)

	getResultChan = s.ReceiveCommand(model.NewGetCommand(model.Key("test1")))
	getResult = <-getResultChan

	assert.Equal(t, false, getResult.Ok)
	assert.Equal(t, nil, getResult.Value.EncodedRawValue)
}

func TestAutomaticDeleteKeyValue(t *testing.T) {
	s := Prepare()
	// two second delay
	setResultChan := s.ReceiveCommand(model.NewSetCommand(model.Key("test1"), 2523, 2_000))
	setResult := <-setResultChan
	assert.NoError(t, setResult.Error)

	getResultChan := s.ReceiveCommand(model.NewGetCommand(model.Key("test1")))
	getResult := <-getResultChan

	assert.Equal(t, true, getResult.Ok)
	assert.Equal(t, 2523, getResult.Value.EncodedRawValue)

	time.Sleep(time.Second * 5)

	getResultChan = s.ReceiveCommand(model.NewGetCommand(model.Key("test1")))
	getResult = <-getResultChan

	assert.Equal(t, false, getResult.Ok)
	assert.Equal(t, nil, getResult.Value.EncodedRawValue)
}
