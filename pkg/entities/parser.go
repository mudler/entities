package entities

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	UserKind = "user"
)

type EntitiesParser interface {
	ReadEntity(entity string) (Entity, error)
}

type Signature struct {
	Kind string `yaml:"kind"`
}

type Parser struct{}

func (p Parser) ReadEntity(entity string) (Entity, error) {
	var signature Signature
	var e DefaultEntity
	yamlFile, err := ioutil.ReadFile(entity)
	if err != nil {
		return &e, errors.Wrap(err, "Failed while reading entity file")
	}

	err = yaml.Unmarshal(yamlFile, &signature)
	if err != nil {
		return &e, errors.Wrap(err, "Failed while parsing entity file")
	}

	switch signature.Kind {
	case UserKind:
		var user UserPasswd

		err = yaml.Unmarshal(yamlFile, &user)
		if err != nil {
			return &e, errors.Wrap(err, "Failed while parsing entity file")
		}
		e.User = &user
	}

	return &e, nil
}
