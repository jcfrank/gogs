// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-martini/martini"

	"github.com/gogits/gogs/modules/base"
	"github.com/gogits/gogs/modules/log"
	"github.com/gogits/gogs/modules/middleware/binding"
)

// Web form interface.
type Form interface {
	Name(field string) string
}

type RegisterForm struct {
	UserName     string `form:"username" binding:"Required;AlphaDashDot;MaxSize(30)"`
	Email        string `form:"email" binding:"Required;Email;MaxSize(50)"`
	Password     string `form:"passwd" binding:"Required;MinSize(6);MaxSize(30)"`
	RetypePasswd string `form:"retypepasswd"`
	LoginType    string `form:"logintype"`
	LoginName    string `form:"loginname"`
}

func (f *RegisterForm) Name(field string) string {
	names := map[string]string{
		"UserName":     "Username",
		"Email":        "E-mail address",
		"Password":     "Password",
		"RetypePasswd": "Re-type password",
	}
	return names[field]
}

func (f *RegisterForm) Validate(errs *binding.Errors, req *http.Request, ctx martini.Context) {
	data := ctx.Get(reflect.TypeOf(base.TmplData{})).Interface().(base.TmplData)
	validate(errs, data, f)
}

type LogInForm struct {
	UserName string `form:"username" binding:"Required;MaxSize(35)"`
	Password string `form:"passwd" binding:"Required;MinSize(6);MaxSize(30)"`
	Remember bool   `form:"remember"`
}

func (f *LogInForm) Name(field string) string {
	names := map[string]string{
		"UserName": "Username",
		"Password": "Password",
	}
	return names[field]
}

func (f *LogInForm) Validate(errs *binding.Errors, req *http.Request, ctx martini.Context) {
	data := ctx.Get(reflect.TypeOf(base.TmplData{})).Interface().(base.TmplData)
	validate(errs, data, f)
}

func GetMinMaxSize(field reflect.StructField) string {
	for _, rule := range strings.Split(field.Tag.Get("binding"), ";") {
		if strings.HasPrefix(rule, "MinSize(") || strings.HasPrefix(rule, "MaxSize(") {
			return rule[8 : len(rule)-1]
		}
	}
	return ""
}

func validate(errs *binding.Errors, data base.TmplData, f Form) {
	if errs.Count() == 0 {
		return
	} else if len(errs.Overall) > 0 {
		for _, err := range errs.Overall {
			log.Error("%s: %v", reflect.TypeOf(f), err)
		}
		return
	}

	data["HasError"] = true
	AssignForm(f, data)

	typ := reflect.TypeOf(f)
	val := reflect.ValueOf(f)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		}

		if err, ok := errs.Fields[field.Name]; ok {
			data["Err_"+field.Name] = true
			switch err {
			case binding.BindingRequireError:
				data["ErrorMsg"] = f.Name(field.Name) + " cannot be empty"
			case binding.BindingAlphaDashError:
				data["ErrorMsg"] = f.Name(field.Name) + " must be valid alpha or numeric or dash(-_) characters"
			case binding.BindingAlphaDashDotError:
				data["ErrorMsg"] = f.Name(field.Name) + " must be valid alpha or numeric or dash(-_) or dot characters"
			case binding.BindingMinSizeError:
				data["ErrorMsg"] = f.Name(field.Name) + " must contain at least " + GetMinMaxSize(field) + " characters"
			case binding.BindingMaxSizeError:
				data["ErrorMsg"] = f.Name(field.Name) + " must contain at most " + GetMinMaxSize(field) + " characters"
			case binding.BindingEmailError:
				data["ErrorMsg"] = f.Name(field.Name) + " is not a valid e-mail address"
			case binding.BindingUrlError:
				data["ErrorMsg"] = f.Name(field.Name) + " is not a valid URL"
			default:
				data["ErrorMsg"] = "Unknown error: " + err
			}
			return
		}
	}
}

// AssignForm assign form values back to the template data.
func AssignForm(form interface{}, data base.TmplData) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		}

		data[fieldName] = val.Field(i).Interface()
	}
}

type InstallForm struct {
	Database        string `form:"database" binding:"Required"`
	Host            string `form:"host"`
	User            string `form:"user"`
	Passwd          string `form:"passwd"`
	DatabaseName    string `form:"database_name"`
	SslMode         string `form:"ssl_mode"`
	DatabasePath    string `form:"database_path"`
	RepoRootPath    string `form:"repo_path"`
	RunUser         string `form:"run_user"`
	Domain          string `form:"domain"`
	AppUrl          string `form:"app_url"`
	AdminName       string `form:"admin_name" binding:"Required;AlphaDashDot;MaxSize(30)"`
	AdminPasswd     string `form:"admin_pwd" binding:"Required;MinSize(6);MaxSize(30)"`
	AdminEmail      string `form:"admin_email" binding:"Required;Email;MaxSize(50)"`
	SmtpHost        string `form:"smtp_host"`
	SmtpEmail       string `form:"mailer_user"`
	SmtpPasswd      string `form:"mailer_pwd"`
	RegisterConfirm string `form:"register_confirm"`
	MailNotify      string `form:"mail_notify"`
}

func (f *InstallForm) Name(field string) string {
	names := map[string]string{
		"Database":    "Database name",
		"AdminName":   "Admin user name",
		"AdminPasswd": "Admin password",
		"AdminEmail":  "Admin e-maill address",
	}
	return names[field]
}

func (f *InstallForm) Validate(errors *binding.Errors, req *http.Request, context martini.Context) {
	data := context.Get(reflect.TypeOf(base.TmplData{})).Interface().(base.TmplData)
	validate(errors, data, f)
}
