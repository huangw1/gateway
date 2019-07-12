/**
 * @Author: huangw1
 * @Date: 2019/7/10 18:14
 */

package config

import (
	"errors"
	"fmt"
	"github.com/huangw1/gateway/encoding"
	"log"
	"regexp"
	"strings"
	"time"
)

type ServerConfig struct {
	Endpoints []*EndpointConfig `mapstructure:"endpoints"`
	Timeout   time.Duration     `mapstructure:"timeout"`
	CacheTTL  time.Duration     `mapstructure:"cache_ttl"`
	Host      []string          `mapstructure:"host"`
	Port      int               `mapstructure:"port"`
	Version   int               `mapstructure:"version"`
	Debug     bool
}

type EndpointConfig struct {
	Endpoint        string        `mapstructure:"endpoint"`
	Method          string        `mapstructure:"method"`
	Backend         []*Backend    `mapstructure:"backend"`
	ConcurrentCalls int           `mapstructure:"concurrent_calls"`
	Timeout         time.Duration `mapstructure:"timeout"`
	CacheTTL        time.Duration `mapstructure:"cache_ttl"`
	QueryString     []string      `mapstructure:"querystring_params"`
}

type Backend struct {
	Group           string            `mapstructure:"group"`
	Method          string            `mapstructure:"method"`
	Host            []string          `mapstructure:"host"`
	URLPattern      string            `mapstructure:"url_pattern"`
	Blacklist       []string          `mapstructure:"blacklist"`
	Whitelist       []string          `mapstructure:"whitelist"`
	Mapping         map[string]string `mapstructure:"mapping"`
	Encoding        string            `mapstructure:"encoding"`
	Target          string            `mapstructure:"target"`
	URLKeys         []string
	ConcurrentCalls int
	Timeout         time.Duration
	IsCollection    bool `mapstructure:"is_collection"`
	Decoder         encoding.Decoder
}

var (
	simpleURLKeysPattern   = regexp.MustCompile(`\{([a-zA-Z\-_0-9]+)\}`)
	endpointURLKeysPattern = regexp.MustCompile(`/\{([a-zA-Z\-_0-9]+)\}`)
	hostPattern            = regexp.MustCompile(`(https?://)?([a-zA-Z0-9\._\-]+)(:[0-9]{2,6})?/?`)
	errInvalidHost         = errors.New("invalid host")
	defaultPort            = 8080
)

func (s *ServerConfig) Init() error {
	if s.Version != 1 {
		return fmt.Errorf("Unsupported version: %d\n", s.Version)
	}
	if s.Port == 0 {
		s.Port = defaultPort
	}
	s.Host = s.cleanHosts(s.Host)
	for i, e := range s.Endpoints {
		e.Endpoint = s.cleanPath(e.Endpoint)
		if err := e.validate(); err != nil {
			return err
		}

		inputParams := s.extractPlaceHoldersFromURLTemplate(e.Endpoint, endpointURLKeysPattern)
		inputSet := map[string]interface{}{}
		for ip := range inputParams {
			inputSet[inputParams[ip]] = nil
		}
		e.Endpoint = s.getEndpointPath(e.Endpoint, inputParams)

		s.initEndpointDefaults(i)

		for j, b := range e.Backend {

			s.initBackendDefaults(i, j)

			b.Method = strings.ToTitle(b.Method)

			if err := s.initBackendURLMappings(i, j, inputSet); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ServerConfig) initBackendURLMappings(e, b int, inputParams map[string]interface{}) error {
	backend := s.Endpoints[e].Backend[b]

	backend.URLPattern = s.cleanPath(backend.URLPattern)

	outputParams := s.extractPlaceHoldersFromURLTemplate(backend.URLPattern, simpleURLKeysPattern)

	outputSet := map[string]interface{}{}
	for op := range outputParams {
		outputSet[outputParams[op]] = nil
	}

	if len(outputSet) > len(inputParams) {
		return fmt.Errorf("Too many output params! input: %v, output: %v\n", outputSet, outputParams)
	}

	tmp := backend.URLPattern
	backend.URLKeys = make([]string, len(outputParams))
	for o := range outputParams {
		if _, ok := inputParams[outputParams[o]]; !ok {
			return fmt.Errorf("Undefined output param [%s]! input: %v, output: %v\n", outputParams[o], inputParams, outputParams)
		}
		tmp = strings.Replace(tmp, "{"+outputParams[o]+"}", "{{."+strings.Title(outputParams[o])+"}}", -1)
		backend.URLKeys = append(backend.URLKeys, strings.Title(outputParams[o]))
	}
	backend.URLPattern = tmp
	return nil
}

func (s *ServerConfig) initBackendDefaults(e, b int) {
	endpoint := s.Endpoints[e]
	backend := endpoint.Backend[b]
	if len(backend.Host) == 0 {
		backend.Host = s.Host
	} else {
		backend.Host = s.cleanHosts(backend.Host)
	}
	if backend.Method == "" {
		backend.Method = endpoint.Method
	}
	backend.Timeout = endpoint.Timeout
	backend.ConcurrentCalls = endpoint.ConcurrentCalls
	switch strings.ToLower(backend.Encoding) {
	case encoding.XML:
		backend.Decoder = encoding.NewXMLDecoder(backend.IsCollection)
	default:
		backend.Decoder = encoding.NewJSONDecoder(backend.IsCollection)
	}
}

func (s *ServerConfig) initEndpointDefaults(e int) {
	endpoint := s.Endpoints[e]
	if endpoint.Method == "" {
		endpoint.Method = "GET"
	} else {
		endpoint.Method = strings.ToTitle(endpoint.Method)
	}
	if s.CacheTTL != 0 && endpoint.CacheTTL == 0 {
		endpoint.CacheTTL = s.CacheTTL
	}
	if s.Timeout != 0 && endpoint.Timeout == 0 {
		endpoint.Timeout = s.Timeout
	}
	if endpoint.ConcurrentCalls == 0 {
		endpoint.ConcurrentCalls = 1
	}
}

func (s *ServerConfig) getEndpointPath(path string, params []string) string {
	result := path
	for p := range params {
		result = strings.Replace(result, "/{"+params[p]+"}", "/:"+params[p], -1)
	}
	return result
}

func (s *ServerConfig) extractPlaceHoldersFromURLTemplate(subject string, pattern *regexp.Regexp) []string {
	matches := pattern.FindAllStringSubmatch(subject, -1)
	keys := make([]string, len(matches))
	for k, v := range matches {
		keys[k] = v[1]
	}
	return keys
}

func (s *ServerConfig) cleanPath(path string) string {
	return "/" + strings.TrimPrefix(path, "/")
}

func (s *ServerConfig) cleanHosts(hosts []string) []string {
	cleans := make([]string, 0)
	for _, host := range hosts {
		cleans = append(cleans, s.cleanHost(host))
	}
	return cleans
}

func (s *ServerConfig) cleanHost(host string) string {
	matches := hostPattern.FindAllStringSubmatch(host, -1)
	if len(matches) != 1 {
		panic(errInvalidHost)
	}
	keys := matches[0][1:]
	if keys[0] == "" {
		keys[0] = "http://"
	}
	return strings.Join(keys, "")
}

func (e *EndpointConfig) validate() error {
	matched, err := regexp.MatchString("^[^/]|/__debug(/.*)?$", e.Endpoint)
	if err != nil {
		log.Printf("ERROR: parsing the endpoint url [%s]: %s. Ignoring\n", e.Endpoint, err.Error())
		return err
	}
	if matched {
		return fmt.Errorf("ERROR: the endpoint url path [%s] is not a valid one!!! Ignoring\n", e.Endpoint)
	}

	if len(e.Backend) == 0 {
		return fmt.Errorf("WARNING: the [%s] endpoint has 0 backends defined! Ignoring\n", e.Endpoint)
	}
	return nil
}
