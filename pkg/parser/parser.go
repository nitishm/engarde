package parser

import (
	"github.com/mitchellh/mapstructure"
	"github.com/vjeantet/grok"
)

// AccessLog defines all possible fields than can be parsed through the configured log pattern
type AccessLog struct {
	// Authority is the request authority header %REQ(:AUTHORITY)%
	Authority string `mapstructure:"authority" json:"authority,omitempty"`
	// BytesReceived in response to the request %BYTES_RECEIVED%
	BytesReceived string `mapstructure:"bytes_received" json:"bytes_received,omitempty"`
	// BytesSent as part of the request body %BYTES_SENT%
	BytesSent string `mapstructure:"bytes_sent" json:"bytes_sent,omitempty"`
	// Duration of the request %DURATION%
	Duration string `mapstructure:"duration" json:"duration,omitempty"`
	// ForwardedFor is the X-Forwarded-For header value %REQ(FORWARDED-FOR)%
	ForwardedFor string `mapstructure:"forwarded_for" json:"forwarded_for,omitempty"`
	// Method is the HTTP method %REQ(:METHOD)%
	Method string `mapstructure:"method" json:"method,omitempty"`
	// Protocol can either be HTTP or TCP %PROTOCOL%
	Protocol string `mapstructure:"protocol" json:"protocol,omitempty"`
	// RequestId is the envoy generated X-REQUEST-ID header "%REQ(X-REQUEST-ID)%"
	RequestId string `mapstructure:"request_id" json:"request_id,omitempty"`
	// ResponseFlags provide any additional details about the response or connection, if any. %RESPONSE_FLAGS%
	ResponseFlags string `mapstructure:"response_flags" json:"response_flags,omitempty"`
	// StatusCode is the response status code %RESPONSE_CODE%
	StatusCode string `mapstructure:"status_code" json:"status_code,omitempty"`
	// TcpServiceTime is the time the tcp request took
	TcpServiceTime string `mapstructure:"tcp_service_time" json:"tcp_service_time,omitempty"`
	// Timestamp is the Start Time %START_TIME%
	Timestamp string `mapstructure:"timestamp" json:"timestamp,omitempty"`
	// UpstreamService is the upstream host the request is intended for %UPSTREAM_HOST%
	UpstreamService string `mapstructure:"upstream_service" json:"upstream_service,omitempty"`
	// UpstreamServiceTime is the time taken to reach target host %RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%
	UpstreamServiceTime string `mapstructure:"upstream_service_time" json:"upstream_service_time,omitempty"`
	// UpstreamCluster is the upstream envoy cluster being reached %UPSTREAM_CLUSTER%
	UpstreamCluster string `mapstructure:"upstream_cluster" json:"upstream_cluster,omitempty"`
	// UpstreamLocal is the local address of the upstream connection %UPSTREAM_LOCAL_ADDRESS%
	UpstreamLocal string `mapstructure:"upstream_local" json:"upstream_local,omitempty"`
	// DownstreamLocal is the local address of the downstream connection %DOWNSTREAM_LOCAL_ADDRESS%
	DownstreamLocal string `mapstructure:"downstream_local" json:"downstream_local,omitempty"`
	// DownstreamRemote is the remote address of the downstream connection %DOWNSTREAM_REMOTE_ADDRESS%
	DownstreamRemote string `mapstructure:"downstream_remote" json:"downstream_remote,omitempty"`
	// RequestedServer is the String value set on ssl connection socket for Server Name Indication (SNI) %REQUESTED_SERVER_NAME%
	RequestedServer string `mapstructure:"requested_server" json:"requested_server,omitempty"`
	// RouteName is the name of the VirtualService route which matched this request %ROUTE_NAME%
	RouteName string `mapstructure:"route_name" json:"route_name,omitempty"`
	// UpstreamFailureReason is the upstream transport failure reason %UPSTREAM_TRANSPORT_FAILURE_REASON%
	UpstreamFailureReason string `mapstructure:"upstream_failure_reason" json:"upstream_failure_reason,omitempty"`
	// UriParam is the params field of the request path
	UriParam string `mapstructure:"uri_param" json:"uri_param,omitempty"`
	// UriPath is the base request path
	UriPath string `mapstructure:"uri_path" json:"uri_path,omitempty"`
	// UserAgent is the request User Agent field %REQ(USER-AGENT)%"
	UserAgent string `mapstructure:"user_agent" json:"user_agent,omitempty"`
	// MixerStatus is the dynamic metadata information for the mixer status %DYNAMIC_METADATA(mixer:status)%
	MixerStatus string `mapstructure:"mixer_status" json:"mixer_status,omitempty"`
	// OriginalMessage is the original raw log line.
	OriginalMessage string `json:"original_message,omitempty"`
	// ParseError provides a string value if a parse error occured.
	ParseError string `json:"parse_error,omitempty"`
}

// Pattern captures various supported grok patterns for different flavors of default envoy access log messages
type Pattern string

const (
	// EnvoyAccessLogsPattern is the default envoy access log format
	EnvoyAccessLogsPattern Pattern = `\[%{TIMESTAMP_ISO8601:timestamp}\] \"%{DATA:method} (?:%{URIPATH:uri_path}(?:%{URIPARAM:uri_param})?|%{DATA}) %{DATA:protocol}\" %{NUMBER:status_code} %{DATA:response_flags} %{NUMBER:bytes_received} %{NUMBER:bytes_sent} %{NUMBER:duration} (?:%{NUMBER:upstream_service_time}|%{DATA:tcp_service_time}) \"%{DATA:forwarded_for}\" \"%{DATA:user_agent}\" \"%{DATA:request_id}\" \"%{DATA:authority}\" \"%{DATA:upstream_service}\"`
	// IstioProxyAccessLogsPattern is the default istio-proxy (envoy) access log format in Istio Service Mesh (matching Istio 1.1, 1.2, and 1.3+ formats)
	IstioProxyAccessLogsPattern Pattern = `\[%{TIMESTAMP_ISO8601:timestamp}\] \"%{DATA:method} (?:(?:%{URIPATH:uri_path}(?:%{URIPARAM:uri_param})?)|%{DATA}) %{DATA:protocol}\" %{NUMBER:status_code} %{DATA:response_flags} \"%{DATA:mixer_status}\"(?: \"%{DATA:upstream_failure_reason}\")? %{NUMBER:bytes_received} %{NUMBER:bytes_sent} %{NUMBER:duration} (?:%{NUMBER:upstream_service_time}|%{DATA:tcp_service_time}) \"%{DATA:forwarded_for}\" \"%{DATA:user_agent}\" \"%{DATA:request_id}\" \"%{DATA:authority}\" \"%{DATA:upstream_service}\" %{DATA:upstream_cluster} %{DATA:upstream_local} %{DATA:downstream_local} %{DATA:downstream_remote} %{DATA:requested_server}(?: %{DATA:route_name})?$`
)

// Parser implements the parsing logic
type Parser struct {
	g       *grok.Grok
	pattern Pattern
}

// New returns a new instance of the Parser object
func New(pattern Pattern) *Parser {
	g, _ := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	return &Parser{g, pattern}
}

// Parse is used to parse a single newline terminated log message with the configured parser pattern
func (p *Parser) Parse(text string) (*AccessLog, error) {
	accessLog := &AccessLog{
		OriginalMessage: text,
	}
	m, err := p.g.Parse(string(p.pattern), text)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(m, accessLog)
	if err != nil {
		return nil, err
	}

	return accessLog, nil
}
