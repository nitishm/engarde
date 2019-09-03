package parser

import (
	"github.com/mitchellh/mapstructure"
	"github.com/vjeantet/grok"
)

type AccessLog struct {
	Authority           string `mapstructure:"authority" json:"authority,omitempty"`
	BytesReceived       string `mapstructure:"bytes_received" json:"bytes_received,omitempty"`
	BytesSent           string `mapstructure:"bytes_sent" json:"bytes_sent,omitempty"`
	Duration            string `mapstructure:"duration" json:"duration,omitempty"`
	ForwardedFor        string `mapstructure:"forwarded-for" json:"forwarded_for,omitempty"`
	Method              string `mapstructure:"method" json:"method,omitempty"`
	Protocol            string `mapstructure:"protocol" json:"protocol,omitempty"`
	RequestId           string `mapstructure:"request_id" json:"request_id,omitempty"`
	ResponseFlags       string `mapstructure:"response_flags" json:"response_flags,omitempty"`
	StatusCode          string `mapstructure:"status_code" json:"status_code,omitempty"`
	TcpServiceTime      string `mapstructure:"tcp_service_time" json:"tcp_service_time,omitempty"`
	Timestamp           string `mapstructure:"timestamp" json:"timestamp,omitempty"`
	UpstreamService     string `mapstructure:"upstream_service" json:"upstream_service,omitempty"`
	UpstreamServiceTime string `mapstructure:"upstream_service_time" json:"upstream_service_time,omitempty"`
	UriParam            string `mapstructure:"uri_param" json:"uri_param,omitempty"`
	UriPath             string `mapstructure:"uri_path" json:"uri_path,omitempty"`
	UserAgent           string `mapstructure:"user_agent" json:"user_agent,omitempty"`
	OriginalMessage     string `json:"original_message,omitempty"`
	ParseError          string `json:"parse_error,omitempty"`
}

type Parser struct {
	g *grok.Grok
}

var (
	envoyPattern = `\[%{TIMESTAMP_ISO8601:timestamp}\] \"%{DATA:method} (?:%{URIPATH:uri_path}(?:%{URIPARAM:uri_param})?|%{DATA:}) %{DATA:protocol}\" %{NUMBER:status_code} %{DATA:response_flags} %{NUMBER:bytes_sent} %{NUMBER:bytes_received} %{NUMBER:duration} (?:%{NUMBER:upstream_service_time}|%{DATA:tcp_service_time}) \"%{DATA:forwarded_for}\" \"%{DATA:user_agent}\" \"%{DATA:request_id}\" \"%{DATA:authority}\" \"%{DATA:upstream_service}\"`
)

func New() *Parser {
	g, _ := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	return &Parser{g}
}

func (p *Parser) Parse(text string) (*AccessLog, error) {
	accessLog := &AccessLog{
		OriginalMessage: text,
	}
	m, err := p.g.Parse(envoyPattern, text)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(m, accessLog)
	if err != nil {
		return nil, err
	}

	return accessLog, nil
}
