package api_helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const content_type_json = "application/json"
const content_type_form = "application/x-www-form-urlencoded"
const content_type_multipart = "multipart/form-data"

func new_request(method, a_url string, body_reader io.Reader) (*http.Request, error) {
	return http.NewRequest(method, a_url, body_reader)
}
func Fresh_request(method, a_url string) *http.Request {
	return new(http.Request)
}

func Add_method(request *http.Request, method string) {
	request.Method = method
}
func Add_url(request *http.Request, a_url string) (err error) {
	request.URL, err = url.Parse(a_url)
	return err
}
func Add_cookie(request *http.Request, name, value string) {
	cookie := &http.Cookie{Name: name, Value: value}
	request.AddCookie(cookie)
}
func Add_data(request *http.Request) {

}

func Post_json_data_to_url[data_type any](a_url string, data data_type) (response *http.Response, err error) {
	data_json, err := json.Marshal(data)
	if err != nil {
		return response, err
	}
	data_reader := bytes.NewBuffer(data_json)
	response, err = http.Post(a_url, content_type_json, data_reader)
	return response, err
}
func Post_form_data_to_url(a_url string, data map[string]any) (response *http.Response, err error) {
	post_form := url.Values{}
	for key, value := range data {
		switch value.(type) {
		case string:
			post_form.Set(key, value.(string))
		case []string:
			list, _ := value.([]string)
			post_form[key] = list
			// post_form[key] = value.([]string)
		default:
			err = fmt.Errorf(fmt.Sprintf("need string or []string but got type %T", value))
		}
	}
	fmt.Println("just before posting", post_form)
	if err != nil {
		return response, err
	}
	response, err = http.PostForm(a_url, post_form)
	return response, err
}

func Post_form_data_with_cookie_to_url(a_url string, data map[string]any, cookie_data map[string]string) (response *http.Response, err error) {
	// var b bytes.Buffer
	post_form := url.Values{}
	for key, value := range data {
		switch value.(type) {
		case string:
			post_form.Set(key, value.(string))
		case []string:
			post_form[key] = value.([]string)
		default:
			err = fmt.Errorf(fmt.Sprintf("need string or []string but got type %T", value))
		}
	}
	if err != nil {
		return response, err
	}
	fmt.Println("form", post_form)
	data_reader := strings.NewReader(post_form.Encode())
	request, err := http.NewRequest("POST", a_url, data_reader)
	request.PostForm = post_form
	request.Header.Set("Content-Type", content_type_form)
	Add_cookie(request, cookie_data["Name"], cookie_data["Value"])
	return http.DefaultClient.Do(request)
}
func Post_multipart_data_to_url(a_url string, data map[string]any) (response *http.Response, err error) {
	body_buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(body_buffer)
	fmt.Println("data", data)
	for key, value := range data {
		switch value.(type) {
		case string:
			err = writer.WriteField(key, value.(string))
		case []string:
			err = writer.WriteField(key, fmt.Sprintf("%v", value))
		default:
			panic("cant make it string")
		}
		if err != nil {
			panic("failed")
		}
	}
	writer.Close()
	return http.Post(a_url, writer.FormDataContentType(), body_buffer)
}

func Get_data_from_request_json[data_type any](request *http.Request) (data_type, error) {
	decoder := json.NewDecoder(request.Body)
	var data data_type
	err := decoder.Decode(&data)
	return data, err
}
func Get_data_from_response_json[data_type any](request *http.Response) (data_type, error) {
	decoder := json.NewDecoder(request.Body)
	var data data_type
	err := decoder.Decode(&data)
	return data, err
}

func Set_cookie_into_response(writer http.ResponseWriter, name, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(writer, &cookie)
}
func Get_cookie_data_from_request(request *http.Request, cookie_name string) (map[string]string, error) {
	cookie_struct, err := request.Cookie(cookie_name)
	cookie_map := map[string]string{}
	cookie_map["Name"] = cookie_struct.Name
	cookie_map["Value"] = cookie_struct.Value
	return cookie_map, err
}

func Get_cookie_data_from_response(response *http.Response, cookie_name string) (cookie_map map[string]string, err error) {
	cookies := response.Cookies()
	if len(cookies) == 0 {
		return
	}
	cookie_map = map[string]string{}
	for _, cookie := range cookies {
		if cookie.Name == cookie_name {
			cookie_map["Name"] = cookie.Name
			cookie_map["Value"] = cookie.Value
		}
	}
	if len(cookie_map) == 0 {
		return cookie_map, fmt.Errorf("did not find named cookie")
	}
	return cookie_map, err
}

func Get_data_from_request_form(request *http.Request) (map[string]any, error) {
	err := request.ParseForm()
	fmt.Println("post form fromt requst", request.PostForm)
	data := map[string]any{}
	if err == nil {
		for k, v := range request.PostForm {
			fmt.Println("into form", k, v, len(v))
			switch len(v) {
			case 1:
				data[k] = v[0]
			case 0:
				continue
			default:
				data[k] = v
			}
		}
	}
	return data, err
}

func Open_saved_file(file_path string) (file multipart.File, err error) {
	file, err = os.Open(file_path)
	return file, err
}

func Post_file_to_url(a_url string, file multipart.File, file_name string) (response *http.Response, err error) {
	defer file.Close()
	body_buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(body_buffer)
	part, err := writer.CreateFormFile("File", file_name)
	if err != nil {
		return response, err
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return response, err
	}
	request, err := http.NewRequest("POST", a_url, body_buffer)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		return response, err
	}
	defer response.Body.Close()
	return response, err
}

func Post_saved_file_to_url(a_url string, file_path, file_type string) (response *http.Response, err error) {
	body_buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(body_buffer)
	file, err := os.Open(file_path)
	if err != nil {
		return response, err
	}
	defer file.Close()
	part, err := writer.CreateFormFile("File", filepath.Base(file_path))
	if err != nil {
		return response, err
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return response, err
	}
	request, err := http.NewRequest("POST", a_url, body_buffer)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	// defer body_buffer.Close()
	return http.DefaultClient.Do(request)
}

func Get_files_from_request_multipart(r *http.Request) (files map[string]multipart.File, err error) {
	files = map[string]multipart.File{}
	err = r.ParseMultipartForm(32 << 10)
	if err != nil {
		return files, err
	}
	for k, v := range r.MultipartForm.File {
		file, err := v[0].Open()
		if err == nil {
			files[k] = file
		}
	}
	return files, err
}
func Get_data_and_files_from_request_multipart(r *http.Request) (data map[string]any, files map[string]multipart.File, err error) {
	data = map[string]any{}
	files = map[string]multipart.File{}
	err = r.ParseMultipartForm(32 << 10)
	if err != nil {
		return data, files, err
	}
	for k, v := range r.MultipartForm.File {
		file, err := v[0].Open()
		if err == nil {
			files[k] = file
		}
	}
	for k, v := range r.MultipartForm.Value {
		switch len(v) {
		case 1:
			data[k] = v[0]
		case 0:
			continue
		default:
			data[k] = v
		}
	}
	return data, files, err
}
func Get_data_from_request_multipart(r *http.Request) (data map[string]any, err error) {
	data = map[string]any{}
	err = r.ParseMultipartForm(32 << 10)
	if err != nil {
		return data, err
	}
	fmt.Println("form of multipart", r.MultipartForm)
	for k, v := range r.MultipartForm.Value {
		switch len(v) {
		case 1:
			data[k] = v[0]
		case 0:
			continue
		default:
			data[k] = v
		}
	}
	return data, err
}
func decide_on_response(response *http.Response) (err error) {
	switch response.StatusCode {
	case 200:
	case 404:
		println("\033[31m", "not found!")
	default:
		err = fmt.Errorf(response.Status)
	}
	return err
}
func Post_data_and_files(a_url string, data map[string]any, files map[string]multipart.File) (response *http.Response, err error) {
	// Prepare a form that you will submit to that URL.
	body_buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(body_buffer)
	for key, value := range data {
		switch value.(type) {
		case string:
			err = writer.WriteField(key, value.(string))
		case []string:
			err = writer.WriteField(key, fmt.Sprintf("%v", value))
		default:
			panic("cant make it string")
		}
		if err != nil {
			panic("failed")
		}
	}
	// prepare writing files
	for file_name, file := range files {
		part, err := writer.CreateFormFile("File", file_name)
		if err != nil {
			return response, err
		}
		_, err = io.Copy(part, file)
		defer file.Close()
		if err != nil {
			return response, err
		}

	}
	// If you don't close it, your request will be missing the terminating boundary.
	err = writer.Close()
	// Don't forget to set the content type, this will contain the boundary.
	request, err := http.NewRequest("POST", a_url, body_buffer)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	// defer response.Body.Close()
	return http.DefaultClient.Do(request)
}
func Get_data_from_request(request *http.Request) (data map[string]any, err error) {
	data = map[string]any{}
	fmt.Println("content type ", request.Header.Get("content-type"))
	only_content_type := strings.Split(request.Header.Get("Content-Type"), ";")[0]
	switch only_content_type {
	case content_type_json:
		data, err = Get_data_from_request_json[map[string]any](request)
		fmt.Println("body", data, err)
	case content_type_form:
		data, err = Get_data_from_request_form(request)
		fmt.Println("post form", data, err)
	case content_type_multipart:
		data, data["files"], err = Get_data_and_files_from_request_multipart(request)
		fmt.Println("multipart", data, err)
	}
	return data, err
}
