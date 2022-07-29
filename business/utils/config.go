package utils

import (
	"io/ioutil"

	model "github.com/ec2ainun/poc-nats/foundation/messaging"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func Read(filename string) (model.Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return model.Config{}, errors.Wrapf(err, "err read config file: %s", filename)
	}

	var config model.Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		return model.Config{}, errors.Wrapf(err, "err unmarshal config file: %s", filename)
	}

	return config, nil
}
