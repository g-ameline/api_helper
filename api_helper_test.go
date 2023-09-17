package api_helper

import (
	"fmt"
	"io"
	"net/http"
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

	// next we send stuff to the server

	some_data := map[string]any{}
	some_data["patate"] = []int{1, 2, 3, 4, 6}
	some_data["courgette"] = "rododindron"
	println("\nsend data through request body")
	response, err := Post_data_to_url(home_url+"data_in_body", some_data)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	just_a_word := "voila"
	response, err = Post_data_to_url(home_url+"data_in_body", just_a_word)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("\nsend data as post form")
	response, err = PostForm_data_to_url(home_url+"data_as_post_form", some_data)
	// fmt.Println("response from sending data", response.StatusCode)
	fmt.Println("error", err)
	new_data := map[string]any{}
	new_data["cacahouette"] = "banane"
	new_data["rutabaga"] = []string{"2,1,3,carotte"}
	response, err = PostForm_data_to_url(home_url+"data_as_post_form", new_data)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("\nsend data in whatever way")
	fmt.Println("data before anything", new_data)
	println("post form")
	response, err = PostForm_data_to_url(home_url+"whatever", new_data)
	fmt.Println("response from sending data", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("")
	println("post body")
	response, err = Post_data_to_url(home_url+"whatever", new_data)
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
	fmt.Println("response from sending file", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)
	println("\nsend a file in post form get in with whatever")
	response, err = Post_saved_file_to_url(home_url+"whatever", "./doc.jpeg", "image")
	fmt.Println("response from sending file", response.StatusCode)
	print_response_body(response)
	fmt.Println("error", err)

}
func print_response_body(response *http.Response) {
	body, err := io.ReadAll(response.Body)
	text := string(body[:])
	fmt.Println("body of response", err, text)
	response.Body.Close()

}
func print_data_from_request(responder http.ResponseWriter, request *http.Request) {
	data, err := Get_data_from_request[any](request)
	fmt.Println("data received from request", data)
	fmt.Printf("of type %T\n", data)
	fmt.Println("the error", err)
	responder.Write([]byte("putain prout"))
}
func print_data_from_post_form_from_request(responder http.ResponseWriter, request *http.Request) {
	data, err := Get_data_from_post_form_from_request(request)
	fmt.Println("data received from postform from request", data)
	fmt.Printf("of type %T\n", data)
	fmt.Println("the error", err)
	responder.Write([]byte("putain dieu"))
}
func print_data_from_whatever_from_request(responder http.ResponseWriter, request *http.Request) {
	data_two, err := Get_data_from_post_form_from_request(request)
	fmt.Println("post form", data_two, err)
	data_three, err := Get_files_from_post_form(request)
	fmt.Println("file", data_three, err)
	data_one, err := Get_data_from_request[map[string]any](request)
	fmt.Println("body", data_one, err)
	data := map[string]any{}
	for k, v := range data_one {
		data[k] = v
	}
	for k, v := range data_two {
		data[k] = v
	}
	for k, v := range data_three {
		data[k] = v
	}
	fmt.Println("data received from postform from request", data)
	fmt.Printf("of type %T\n", data)
	fmt.Println("the error", err)
	responder.Write([]byte("putain con"))
}
func check_if_file(responder http.ResponseWriter, request *http.Request) {
	data, err := Get_files_from_post_form(request)
	fmt.Println("data received from postform from request", data)
	fmt.Printf("of type %T\n", data)
	fmt.Println("the error", err)
	responder.Write([]byte("putain merde"))
}
