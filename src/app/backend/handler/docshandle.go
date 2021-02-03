package handler

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
)

func CreateApiDocsHTTPHandler(wsContainer *restful.Container, specURL string, next http.Handler) http.Handler {

	config := restfulspec.Config{
		WebServices:                   wsContainer.RegisteredWebServices(),
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}

	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	opts := middleware.RedocOpts{SpecURL: specURL}
	sh := middleware.Redoc(opts, next)

	return sh
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Fission OpenAPI 2.0",
			Description: "TEST",
			Version:     "v1",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{Name: "dhkang"},
			},
		},
	}
	swo.Tags = []spec.Tag{
		{
			TagProps: spec.TagProps{
				Name:        "Users",
				Description: "Managing users",
			},
		}}
}
