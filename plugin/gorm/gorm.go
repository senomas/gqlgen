package gorm

import (
	"go/types"
	"regexp"

	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/vektah/gqlparser/v2/ast"
)

var tref = regexp.MustCompile("^(.*){{ref:([^\\s]+)\\s+([^}]+)}}(.*)$")
var trefTag = regexp.MustCompile("^(.*){{refTag:([^}]+)}}(.*)$")

type RefType struct {
	Name string
}

func (r *RefType) Underlying() types.Type {
	return nil
}

func (r *RefType) String() string {
	return r.Name
}

func MutateHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	for _, model := range b.Models {
		var fields []*modelgen.Field
		for _, field := range model.Fields {
			if ms := tref.FindAllStringSubmatch(field.Tag, -1); len(ms) != 0 {
				ms0 := ms[0]
				field.Tag = ms0[1] + ms0[4]
				if ns := trefTag.FindAllStringSubmatch(field.Tag, -1); len(ns) != 0 {
					ns0 := ns[0]
					field.Tag = ns0[1] + ns0[3]
					fields = append(fields, field)
					fields = append(fields, &modelgen.Field{
						Name:   ms0[2],
						GoName: ms0[2],
						Type:   &RefType{Name: ms0[3]},
						Tag:    "json:\"-\" gorm:\"" + ns0[2] + "\"",
					})
				} else {
					fields = append(fields, field)
					fields = append(fields, &modelgen.Field{
						Name:   ms0[2],
						GoName: ms0[2],
						Type:   &RefType{Name: ms0[3]},
						Tag:    "json:\"-\"",
					})
				}
			} else {
				fields = append(fields, field)
			}
		}
		model.Fields = fields
	}

	return b
}

func FieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	c := fd.Directives.ForName("gorm")
	if c != nil {
		tag := c.Arguments.ForName("tag")
		if tag != nil {
			f.Tag += " gorm:" + tag.Value.String()
		}
		ref := c.Arguments.ForName("ref")
		if ref != nil {
			f.Tag += "{{ref:" + ref.Value.Raw + "}}"
		}
		refTag := c.Arguments.ForName("refTag")
		if refTag != nil {
			f.Tag += "{{refTag:" + refTag.Value.Raw + "}}"
		}
	}
	return f, nil
}
