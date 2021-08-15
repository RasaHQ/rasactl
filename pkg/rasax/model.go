/*
Copyright © 2021 Rasa Technologies GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package rasax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	rtypes "github.com/RasaHQ/rasactl/pkg/types/rasax"
	"github.com/schollz/progressbar/v3"
)

func (r *RasaX) ModelUpload() error {
	file, err := os.Open(r.Flags.Model.Upload.File)
	stat, _ := file.Stat()
	if err != nil {
		return err
	}
	defer file.Close()

	bar := r.progressBarBytes(
		stat.Size(),
		fmt.Sprintf("Sending %s", filepath.Base(file.Name())),
	)

	buffer := new(bytes.Buffer)
	body := progressbar.NewReader(buffer, bar)
	writer := multipart.NewWriter(buffer)
	part, err := writer.CreateFormFile("model", filepath.Base(file.Name()))
	if err != nil {
		return err
	}

	io.Copy(part, file)
	writer.Close()
	//buffer
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/projects/default/models", r.URL), &body)
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.BearerToken))
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 201:
		fmt.Println("Successfully uploaded.")
	case 409:
		fmt.Println("A model with that name already exists.")
	default:
		content, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("%s", content)
	}

	return nil
}

func (r *RasaX) ModelDownload() error {
	request, err := http.NewRequest("GET",
		fmt.Sprintf("%s/api/projects/default/models/%s", r.URL, r.Flags.Model.Download.Name), nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.BearerToken))
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var file string = r.Flags.Model.Download.FilePath
	if r.Flags.Model.Download.FilePath == "" {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		file = fmt.Sprintf("%s/%s.tar.gz", dir, r.Flags.Model.Download.Name)
	}
	f, _ := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	switch resp.StatusCode {
	case 200:
		bar := r.progressBarBytes(
			resp.ContentLength,
			fmt.Sprintf("Downloading %s", r.Flags.Model.Download.Name),
		)
		io.Copy(io.MultiWriter(f, bar), resp.Body)
		fmt.Println("Model has been downloaded successfully.")
	case 404:
		return fmt.Errorf("model '%s' not found", r.Flags.Model.Download.Name)
	default:
		content, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%s", content)
	}
	return nil
}

func (r *RasaX) ModelList() (*rtypes.ModelsListEndpointResponse, error) {
	bodyData := &rtypes.ModelsListEndpointResponse{}
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/default/models", r.URL), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.BearerToken))
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(body, &bodyData.Models)
		return bodyData, nil
	case 404:
		return nil, fmt.Errorf("unauthorized, use the 'rasactl auth login' command to authorized")
	default:
		content, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s", content)
	}
}

func (r *RasaX) ModelTag() error {
	r.Log.Info("Sending a request to Rasa X")
	request, err := http.NewRequest("PUT",
		fmt.Sprintf("%s/api/projects/default/models/%s/tags/%s", r.URL, r.Flags.Model.Tag.Model, r.Flags.Model.Tag.Name), nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.BearerToken))
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 204:
		fmt.Println("Model has been tagged successfully.")
		return nil
	default:
		content, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%s", content)
	}
}

func (r *RasaX) ModelDelete() error {
	r.Log.Info("Sending a request to Rasa X")
	request, err := http.NewRequest("DELETE",
		fmt.Sprintf("%s/api/projects/default/models/%s", r.URL, r.Flags.Model.Delete.Name), nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.BearerToken))
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 204:
		fmt.Println("Model has been deleted successfully.")
		return nil
	case 404:
		return fmt.Errorf("model '%s' not found", r.Flags.Model.Delete.Name)
	default:
		content, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%s", content)
	}
}
