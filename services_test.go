package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPopulateTLSConfig(t *testing.T) {

}

func TestPopulateTLSConfigSetsInsecureSkipVerify(t *testing.T) {
	assert := assert.New(t)

	tls_config := TLSConfig{
		"InsecureSkipVerify": "true",
	}
	config, err := populateTLSConfig(tls_config)
	assert.Nil(err)
	assert.Equal(true, config.InsecureSkipVerify)
}
