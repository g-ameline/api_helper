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
)

const content_type_json = "application/json"

func Post_data_to_url[data_type any](a_url string, data data_type) (response *http.Response, err error) {
	data_json, err := json.Marshal(data)
	if err != nil {
		return response, err
	}
	data_reader := bytes.NewBuffer(data_json)
	response, err = http.Post(a_url, content_type_json, data_reader)
	return response, err
}
func PostForm_data_to_url(a_url string, data map[string]any) (response *http.Response, err error) {
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
	response, err = http.PostForm(a_url, post_form)
	return response, err
}
func Get_data_from_request[data_type any](request *http.Request) (data_type, error) {
	decoder := json.NewDecoder(request.Body)
	var data data_type
	err := decoder.Decode(&data)
	return data, err
}
func Get_cookie_data_from_request(request *http.Request, cookie_name string) (map[string]string, error) {
	cookie_struct, err := request.Cookie(cookie_name)
	cookie_map := map[string]string{}
	cookie_map["Name"] = cookie_struct.Name
	cookie_map["Value"] = cookie_struct.Value
	return cookie_map, err
}

func Get_data_from_post_form_from_request(request *http.Request) (map[string]any, error) {
	err := request.ParseForm()
	data := map[string]any{}
	if err == nil {
		for k, v := range request.PostForm {
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
	client := &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		return response, err
	}
	defer response.Body.Close()
	return response, err
}

func Get_files_from_post_form(r *http.Request) (files map[string]multipart.File, err error) {
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
