package comparator

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type config struct {
	Spec struct {
		Template struct {
			Spec struct {
				Containers []struct {
					Name string `yaml:"name"`
					Env  []struct {
						Name  string `yaml:"name"`
						Value string `yaml:"value"`
					} `yaml:"env"`
				} `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

type EnvProblem struct {
	Name string
	Val1 string
	Val2 string
}

type CompareResult struct {
	Name string
	Envs []*EnvProblem
}

func parseFile(ctx context.Context, file string) ([]*config, error) {
	logger := log.Ctx(ctx).With().Str("file", file).Logger()

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Error().Err(err).Msg("Failed read file")
		return nil, err
	}

	return parseYaml(ctx, yamlFile)
}

func parseYaml(ctx context.Context, data []byte) ([]*config, error) {
	logger := log.Ctx(ctx).With().Logger()

	var conf []*config
	parts := strings.Split(string(data), "---\n")
	for _, p := range parts {
		if p == "" {
			continue
		}
		var c config
		if err := yaml.Unmarshal([]byte(p), &c); err != nil {
			logger.Error().Err(err).Msg("Failed unmarshal data")
			return nil, err
		}
		conf = append(conf, &c)
	}
	return conf, nil
}

func CompareYaml(ctx context.Context, file1, file2 string) ([]*CompareResult, error) {
	logger := log.Ctx(ctx)
	p1, err := parseFile(ctx, file1)
	if err != nil {
		logger.Error().Err(err).Str("file", file1).Msg("Failed parse file")
		return nil, err
	}
	p2, err := parseFile(ctx, file2)
	if err != nil {
		logger.Error().Err(err).Str("file", file2).Msg("Failed parse file")
		return nil, err
	}
	var result []*CompareResult
	for _, conf1 := range p1 {
		for _, cont1 := range conf1.Spec.Template.Spec.Containers {
			if cont1.Name == "" {
				continue
			}
			r := CompareResult{
				Name: cont1.Name,
			}
			for _, conf2 := range p2 {
				for _, cont2 := range conf2.Spec.Template.Spec.Containers {
					if cont2.Name == cont1.Name {
						for _, env1 := range cont1.Env {
							found := false
							for _, env2 := range cont2.Env {
								if env1.Name == env2.Name {
									found = true
									if env1.Value != env2.Value {
										r.Envs = append(r.Envs, &EnvProblem{
											Name: env1.Name,
											Val1: env1.Value,
											Val2: env2.Value,
										})
									}
									break
								}
							}
							if !found {
								r.Envs = append(r.Envs, &EnvProblem{
									Name: env1.Name,
									Val1: env1.Value,
								})
							}
						}

						for _, env2 := range cont2.Env {
							found := false
							for _, env1 := range cont1.Env {
								if env1.Name == env2.Name {
									found = true
									break
								}
							}
							if !found {
								r.Envs = append(r.Envs, &EnvProblem{
									Name: env2.Name,
									Val2: env2.Value,
								})
							}
						}
					}
				}
			}
			if len(r.Envs) != 0 {
				result = append(result, &r)
			}
		}
	}
	return result, nil
}
