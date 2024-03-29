// +build !ignore_autogenerated

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"./pkg/apis/jakubbujny/v1alpha1.JenkinsPipeline":       schema_pkg_apis_jakubbujny_v1alpha1_JenkinsPipeline(ref),
		"./pkg/apis/jakubbujny/v1alpha1.JenkinsPipelineSpec":   schema_pkg_apis_jakubbujny_v1alpha1_JenkinsPipelineSpec(ref),
		"./pkg/apis/jakubbujny/v1alpha1.JenkinsPipelineStatus": schema_pkg_apis_jakubbujny_v1alpha1_JenkinsPipelineStatus(ref),
	}
}

func schema_pkg_apis_jakubbujny_v1alpha1_JenkinsPipeline(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "JenkinsPipeline is the Schema for the jenkinspipelines API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/jakubbujny/v1alpha1.JenkinsPipelineSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/jakubbujny/v1alpha1.JenkinsPipelineStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/jakubbujny/v1alpha1.JenkinsPipelineSpec", "./pkg/apis/jakubbujny/v1alpha1.JenkinsPipelineStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_jakubbujny_v1alpha1_JenkinsPipelineSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "JenkinsPipelineSpec defines the desired state of JenkinsPipeline",
				Properties: map[string]spec.Schema{
					"microservice": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
				Required: []string{"microservice"},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_jakubbujny_v1alpha1_JenkinsPipelineStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "JenkinsPipelineStatus defines the observed state of JenkinsPipeline",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
