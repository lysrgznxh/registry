package model

import (
	"time"
)

// Status represents the lifecycle status of a server
type Status string

const (
	StatusActive     Status = "active"
	StatusDeprecated Status = "deprecated"
	StatusDeleted    Status = "deleted"
)

// Transport represents transport configuration with optional URL templating
type Transport struct {
	Type    string          `json:"type"`
	URL     string          `json:"url,omitempty"`
	Headers []KeyValueInput `json:"headers,omitempty"`
	Jsons   []Argument      `json:"jsons,omitempty"`
	Posts   []Argument      `json:"posts,omitempty"`
}

// Package represents a package configuration
type Package struct {
	// RegistryType indicates how to download packages (e.g., "npm", "pypi", "oci", "mcpb")
	RegistryType string `json:"registryType" minLength:"1"`
	// RegistryBaseURL is the base URL of the package registry
	RegistryBaseURL string `json:"registryBaseUrl,omitempty"`
	// Identifier is the package identifier - either a package name (for registries) or URL (for direct downloads)
	Identifier           string          `json:"identifier" minLength:"1"`
	Version              string          `json:"version" minLength:"1"`
	FileSHA256           string          `json:"fileSha256,omitempty"`
	RunTimeHint          string          `json:"runtimeHint,omitempty"`
	Transport            Transport       `json:"transport,omitempty"`
	RuntimeArguments     []Argument      `json:"runtimeArguments,omitempty"`
	PackageArguments     []Argument      `json:"packageArguments,omitempty"`
	EnvironmentVariables []KeyValueInput `json:"environmentVariables,omitempty"`
}

// Repository represents a source code repository as defined in the spec
type Repository struct {
	URL       string `json:"url"`
	Source    string `json:"source"`
	ID        string `json:"id,omitempty"`
	Subfolder string `json:"subfolder,omitempty"`
}

// Format represents the input format type
type Format string
type FormatRequire string

const (
	FormatString   Format = "string"
	FormatNumber   Format = "number"
	FormatBoolean  Format = "boolean"
	FormatFilePath Format = "file_path"
	FormatJson     Format = "json"
)

const (
	FormatRequireDocument FormatRequire = "document"
	FormatRequireImage    FormatRequire = "image"
	FormatRequireAudio    FormatRequire = "audio"
	FormatRequireVideo    FormatRequire = "video"
)

// Input represents a configuration input
type Input struct {
	Description        string           `json:"description,omitempty"`
	IsRequired         bool             `json:"isRequired,omitempty"`
	Format             Format           `json:"format,omitempty"`
	FormatRequirements []FormatRequire  `json:"formatRequire,omitempty"`
	Value              string           `json:"value,omitempty"`
	IsSecret           bool             `json:"isSecret,omitempty"`
	IsArray            bool             `json:"isArray,omitempty"`
	Default            string           `json:"default,omitempty"`
	Choices            []string         `json:"choices,omitempty"`
	Variables          map[string]Input `json:"variables,omitempty"`
}

// InputWithVariables represents an input that can contain variables
type InputWithVariables struct {
	Input     `json:",inline"`
	Variables map[string]Input `json:"variables,omitempty"`
}

// KeyValueInput represents a named input with variables
type KeyValueInput struct {
	InputWithVariables `json:",inline"`
	Name               string `json:"name"`
}

// ArgumentType represents the type of argument
type ArgumentType string

const (
	ArgumentTypePositional ArgumentType = "positional"
	ArgumentTypeNamed      ArgumentType = "named"
)

// Argument defines a type that can be either a PositionalArgument or a NamedArgument
type Argument struct {
	InputWithVariables `json:",inline"`
	Type               ArgumentType `json:"type,omitempty"`
	Name               string       `json:"name,omitempty"`
	IsRepeated         bool         `json:"isRepeated,omitempty"`
	ValueHint          string       `json:"valueHint,omitempty"`
}

// RegistryExtensions represents registry-generated metadata
type RegistryExtensions struct {
	ServerID    string    `json:"serverId"`  // Consistent ID across all versions of a server
	VersionID   string    `json:"versionId"` // Unique ID for this specific version
	PublishedAt time.Time `json:"publishedAt"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	IsLatest    bool      `json:"isLatest"`
}

// ServerListResponse represents the paginated server list response
type ServerListResponse struct {
	Servers  []ServerJSON `json:"servers"`
	Metadata Metadata     `json:"metadata"`
}

// ServerMeta represents the structured metadata with known extension fields
type ServerMeta struct {
	Official          *RegistryExtensions    `json:"io.modelcontextprotocol.registry/official,omitempty"`
	PublisherProvided map[string]interface{} `json:"io.modelcontextprotocol.registry/publisher-provided,omitempty"`
}

// ServerJSON represents complete server information as defined in the MCP spec, with extension support
type ServerJSON struct {
	Schema      string      `json:"$schema,omitempty"`
	Name        string      `json:"name" minLength:"1" maxLength:"200"`
	Description string      `json:"description" minLength:"1" maxLength:"100"`
	Status      Status      `json:"status,omitempty" minLength:"1"`
	Repository  Repository  `json:"repository,omitempty"`
	Version     string      `json:"version"`
	WebsiteURL  string      `json:"websiteUrl,omitempty"`
	Packages    []Package   `json:"packages,omitempty"`
	Remotes     []Transport `json:"remotes,omitempty"`
	Meta        *ServerMeta `json:"_meta,omitempty"`
}

// Metadata represents pagination metadata
type Metadata struct {
	NextCursor string `json:"next_cursor,omitempty"`
	Count      int    `json:"count"`
}

func (s *ServerJSON) GetServerID() string {
	if s.Meta != nil && s.Meta.Official != nil {
		return s.Meta.Official.ServerID
	}
	return ""
}

func (s *ServerJSON) GetVersionID() string {
	if s.Meta != nil && s.Meta.Official != nil {
		return s.Meta.Official.VersionID
	}
	return ""
}

type AgentType string

const (
	AgentTypeAgent AgentType = "agent"
	AgentTypeMcp   AgentType = "mcp"
)
