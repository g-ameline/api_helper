package api_helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

const https_protocol = "http://"
const domain = "localhost"
const port = ":3000"
const root = "/"
const main_endpoint = ""
const http_url = https_protocol + domain + port // http://localhost:3000
const home_url = http_url + root                // http://localhost:3000/

func Test(t *testing.T) {
	var server *http.ServeMux = http.NewServeMux()
	// start a server to receive data request
	go func() {
		err := http.ListenAndServe(port, server)
		if err != nil {
			panic(2)
		}
	}()
	server.HandleFunc(root+"data_in_body", print_data_from_request)
	server.HandleFunc(root+"data_as_post_form", print_data_from_post_form_from_request)
	server.HandleFunc(root+"whatever", print_data_from_whatever_from_request)
	server.HandleFunc(root+"a_file", check_if_file)
	server.HandleFunc(root+"multipart", print_data_from_request_multipart)
	server.HandleFunc(root+"data_files", print_data_and_files_from_request_multipart)
	server.HandleFunc(root+"respond_json", respond_json)

	// next we send stuff to the server

	some_data := map[string]any{}
	some_data["patate"] = []int{1, 2, 3, 4, 6}
	some_data["courgette"] = "rododindron"
	println("\nsend data through request body")
	response, err := Post_json_data_to_url(home_url+"data_in_body", some_data)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	just_a_word := "voila"
	response, err = Post_json_data_to_url(home_url+"data_in_body", just_a_word)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("\nsend data as post form")
	response, err = Post_form_data_to_url(home_url+"data_as_post_form", some_data)
	// fmt.Println("response from sending data", response.StatusCode)
	fmt.Println("error", err)
	println("\nsend data through request body as multipart")
	new_data := map[string]any{}
	new_data["cacahouette"] = "banane"
	new_data["rutabaga"] = []string{"2", "1", "3", "carotte"}
	response, err = Post_multipart_data_to_url(home_url+"multipart", new_data)
	fmt.Println("error", err)
	response, err = Post_form_data_to_url(home_url+"data_as_post_form", new_data)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("\nsend data in whatever way")
	fmt.Println("data before anything", new_data)
	println("post form")
	response, err = Post_form_data_to_url(home_url+"whatever", new_data)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("")
	println("post body")
	response, err = Post_json_data_to_url(home_url+"whatever", new_data)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("\nsend a file in post form")
	response, err = Post_saved_file_to_url(home_url+"a_file", "./doc.jpeg", "image")
	fmt.Println("response from sending file", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	file_to_send, err := Open_saved_file("./doc.jpeg")
	response, err = Post_file_to_url(home_url+"a_file", file_to_send, "image")
	// file_to_send.Close()
	fmt.Println("response from sending file", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("\nsend a file in post form get in with whatever")
	response, err = Post_saved_file_to_url(home_url+"whatever", "./doc.jpeg", "image")
	fmt.Println("response from sending file", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("\nsend post form and cookie")
	cookie_data := map[string]string{}
	cookie_data["Value"] = "TAPOUERE!"
	cookie_data["Name"] = "PUREE!"
	response, err = Post_form_data_with_cookie_to_url(home_url+"data_as_post_form", new_data, cookie_data)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("\nsend stuff however and get it anyever")
	response, err = Post_json_data_to_url(home_url+"whatever", new_data)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	response, err = Post_form_data_to_url(home_url+"whatever", new_data)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	response, err = Post_multipart_data_to_url(home_url+"whatever", new_data)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("\nsend a file in post with form data")
	file_to_send, err = Open_saved_file("./doc.jpeg")
	files_to_send := map[string]multipart.File{}
	files_to_send["doc"] = file_to_send
	response, err = Post_data_and_files(home_url+"data_files", new_data, files_to_send)
	// file_to_send.Close()
	fmt.Println("response from sending file", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	files_to_send = map[string]multipart.File{}
	response, err = Post_data_and_files(home_url+"data_files", new_data, files_to_send)
	// file_to_send.Close()
	fmt.Println("response from sending file", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("\nget json response from server")
	response, err = Post_json_data_to_url(home_url+"respond_json", new_data)
	fmt.Println("response from sending file", response.StatusCode)
	print_response_body(response)

}
func print_response_body(response *http.Response) {
	body, err := io.ReadAll(response.Body)
	text := string(body[:])
	fmt.Println("body of response", err, text)
	response.Body.Close()

}
func print_data_from_request(responder http.ResponseWriter, request *http.Request) {
	data, err := Get_data_from_request_json[any](request)
	fmt.Println("data received from request", data)
	cookies := request.Cookies()
	if len(cookies) > 0 {
		fmt.Println("cookies", cookies)
	}
	fmt.Printf("of type %T\n", data)
	fmt.Println("the error", err)
	responder.Write([]byte("putain prout"))
}
func print_data_from_request_multipart(responder http.ResponseWriter, request *http.Request) {
	data, err := Get_data_from_request_multipart(request)
	cookies := request.Cookies()
	if len(cookies) > 0 {
		fmt.Println("cookies", cookies)
	}
	fmt.Printf("of type %T\n", data)
	fmt.Println("the error", err)
	responder.Write([]byte("putain prout"))
}
func print_data_from_post_form_from_request(responder http.ResponseWriter, request *http.Request) {
	data, err := Get_data_from_request_form(request)
	fmt.Println("data received from postform from request", data)
	cookies := request.Cookies()
	if len(cookies) > 0 {
		fmt.Println("cookies", cookies)
	}
	fmt.Printf("of type %T\n", data)
	fmt.Println("the error", err)
	responder.Write([]byte("putain dieu"))
}
func print_data_and_files_from_request_multipart(responder http.ResponseWriter, request *http.Request) {
	data, files, err := Get_data_and_files_from_request_multipart(request)
	fmt.Println("data received from multipart from request", data)
	fmt.Println("file received from multipart from request", files)
	cookies := request.Cookies()
	if len(cookies) > 0 {
		fmt.Println("cookies", cookies)
	}
	fmt.Printf("of type %T\n", data)
	fmt.Printf("of type %T\n", data)
	fmt.Println("the error", err)
	responder.Write([]byte("putain prout"))
}
func print_data_from_whatever_from_request(responder http.ResponseWriter, request *http.Request) {
	data := map[string]any{}
	err := *new(error)
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
		// data, err := Get_files_from_request_multipart(request)
		// fmt.Println("file", data, err)
	}
	fmt.Println("data received from request", data)
	fmt.Printf("of type %T\n", data)
	fmt.Println("the error", err)
	responder.Write([]byte("putain con"))
}
func check_if_file(responder http.ResponseWriter, request *http.Request) {
	data, err := Get_files_from_request_multipart(request)
	fmt.Println("data received from postform from request", data)
	fmt.Printf("of type %T\n", data)
	fmt.Println("the error", err)
	responder.Write([]byte("putain merde"))
}
func respond_json(responder http.ResponseWriter, request *http.Request) {
	data := map[string]string{}
	data["truc"] = "chose"
	data["bidule"] = "machin"
	Respond_with_json_data(responder, data)
}
