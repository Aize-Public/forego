package openapi

func NewService(title string) *Service {
	out := &Service{
		OpenAPI: "3.0.0",
		Paths:   map[string]*Path{},
	}
	out.Info.Title = title
	out.Info.License.Name = "private"
	out.Info.Version = "0.0"
	out.Components.SecurityScheme = map[string]*SecurityScheme{}
	out.Components.SecurityScheme["jwt"] = &SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	}
	return out
}

type Service struct {
	OpenAPI string `json:"openapi"` // should be 3.0.0
	Info    struct {
		Version string `json:"version,omitempty"`
		Title   string `json:"title,omitempty"`
		License struct {
			Name string `json:"name,omitempty"`
		} `json:"license,omitempty"`
	} `json:"info"`

	Paths      map[string]*Path `json:"paths"`
	Components Component        `json:"components"`
}

type Component struct {
	SecurityScheme map[string]*SecurityScheme `json:"securitySchemes"`
	Parameters     map[string]*Parameter      `json:"parameters,omitempty"`
	Schemas        map[string]*Schema         `json:"schemas,omitempty"`
}
type SecurityScheme struct {
	Type        string `json:"type"` // apiKey, http, mutualTLS, oauth2, openIdConnect
	Description string `json:"description"`

	// http
	Scheme       string `json:"scheme,omitempty"`       // Bearer (https://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml)
	BearerFormat string `json:"bearerFormat,omitempty"` // likely JWT

	// apiKey
	In   string `json:"in,omitempty"` // query, header, cookie
	Name string `json:"name,omitempty"`
}

type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"` // header, query, path, cookie
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Deprecated  bool    `json:"deprecated,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}
