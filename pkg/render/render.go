package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/macices/bookings/pkg/config"
	"github.com/macices/bookings/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(writer http.ResponseWriter, tmplName string, templateData *models.TemplateData) {
	var templateCache map[string]*template.Template

	if app.UseCache {
		templateCache = app.TemplateCache

	} else {
		templateCache, _ = CreateTemplateCache()
	}

	//get requested template from cache
	template, ok := templateCache[tmplName]
	if !ok {
		log.Fatal("could not get template from templateCache")
	}

	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData)

	_ = template.Execute(buf, templateData)

	_, err := buf.WriteTo(writer)
	if err != nil {
		log.Println("Error writing template to bbrowser", err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	//get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	//range throug all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet
	}

	return myCache, nil
}

//var templateCache = make(map[string]*template.Template)

//func RenderTemplate(writer http.ResponseWriter, tmpl string) {
//	var templt *template.Template
//	var err error
//
/*	_, inMap := templateCache[tmpl]
	if !inMap {
		log.Println("creating template and adding to cache")
		err = createTemplateCache(tmpl)
		if err != nil {
			log.Println(err)
		}

	} else {
		log.Println("using cache template")
	}

	templt = templateCache[tmpl]
	err = templt.Execute(writer, nil)
	if err != nil {
		log.Println(err)
	}

}

func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t), "./templates/base.layout.tmpl",
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
//	}
//
//	templateCache[t] = tmpl
//
//	return nil
//}*/
