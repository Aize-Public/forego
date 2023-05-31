package openapi

type Path struct {
	Post    *PathItem `json:"post,omitempty"`
	Get     *PathItem `json:"get,omitempty"`
	Upgrade *PathItem `json:"upgrade,omitempty"`
}

type PathItem struct {
	Summary     string              `json:"summary"`     // this is more like a Name
	Description string              `json:"description"` // displayed on the right
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Servers     []Server            `json:"servers,omitempty"`
	Responses   map[string]Response `json:"responses"` // "default", "200", ...
	Security    Securities          `json:"security,omitempty"`
	Deprecated  bool                `json:"deprecated,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
}

// add a jwt security, and optionally an empty security unless required is true
func (this *PathItem) SetJWT(required bool) {
	this.Security.Append(Security{
		"jwt": []string{},
	})
	if !required {
		this.Security.Append(Security{})
	}
}

type RequestBody struct {
	Content  map[string]MediaType `json:"content"`
	Required bool                 `json:"required"`
}

type MediaType struct {
	Schema *Schema `json:"schema"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type Response struct {
	Description string             `json:"description"`
	Content     map[string]Content `json:"content,omitempty"`
}

type Content struct {
	Schema *Schema `json:"schema"`
}

type Securities []Security

func (this *Securities) Append(s Security) {
	*this = append(*this, s)
}

type Security map[string]any

// allow for quick setting of expected req/res type as JSON
// set req to nil for GET
// set res to nil for 204
/*func (this *Path) MustSetJSON(c ctx.C, req, res any) *PathItem {
	pi := this.Post
	if req != nil {
		s, err := this.SchemaFromType(c, reflect.TypeOf(req), nil)
		if err != nil {
			panic(err)
		}
		pi.RequestBody.Content = map[string]MediaType{
			"application/json": {Schema: s},
		}
	} else {
		this.Get = &PathItem{}
		pi = this.Get
		this.Post = nil
	}
	if res != nil {
		s, err := this.SchemaFromType(c, reflect.TypeOf(res), nil)
		if err != nil {
			panic(err)
		}
		pi.Responses = map[string]Response{
			"200": {
				Content: map[string]Content{
					"application/json": {Schema: s},
				},
			},
		}
	} else {
		pi.Responses = map[string]Response{
			"204": {},
		}
	}
	return pi
}*/
