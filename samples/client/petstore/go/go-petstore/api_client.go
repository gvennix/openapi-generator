package swagger

import (
    "strings"
    "github.com/go-resty/resty"
    "errors"
    "reflect"
    "bytes"
    "path/filepath"
)

type ApiClient struct {

}

func (c *ApiClient) SelectHeaderContentType(contentTypes []string) string {
    if (len(contentTypes) == 0){
        return ""
    }
    if contains(contentTypes,"application/json") {
        return "application/json"
    }

    return contentTypes[0] // use the first content type specified in 'consumes'
}

func (c *ApiClient) SelectHeaderAccept(accepts []string) string {
    if (len(accepts) == 0){        
        return ""
    }

    if contains(accepts,"application/json"){        
        return "application/json"
    }

    return strings.Join(accepts,",")
}

func contains(source []string, containvalue string) bool {
    for _, a := range source {
        if strings.ToLower(a) == strings.ToLower(containvalue) {
            return true
        }
    }
    return false
}


func (c *ApiClient) CallApi(path string, method string,
    postBody interface{},
    headerParams map[string]string,
    queryParams map[string]string,
    formParams map[string]string,
    fileName string,
    fileBytes []byte) (*resty.Response, error) {

    //set debug flag
    configuration := NewConfiguration()
    resty.SetDebug(configuration.GetDebug())

    request := prepareRequest(postBody, headerParams, queryParams, formParams,fileName,fileBytes)

    switch strings.ToUpper(method) {
    case "GET":
        response, err := request.Get(path)
        return response, err
    case "POST":
        response, err := request.Post(path)
        return response, err
    case "PUT":
        response, err := request.Put(path)
        return response, err
    case "PATCH":
        response, err := request.Patch(path)
        return response, err
    case "DELETE":
        response, err := request.Delete(path)
        return response, err
    }

    return nil, errors.New("Invalid method " + method)
}

func (c *ApiClient) ParameterToString(obj interface{}) string {
    if reflect.TypeOf(obj).String() == "[]string" {
        return strings.Join(obj.([]string), ",")
    } else {
        return obj.(string)
    }
}

func (c *ApiClient) GetApiResponse(httpResp interface{}) *ApiResponse{
  httpResponse := httpResp.(*resty.Response)
  apiResponse := new(ApiResponse) 
  apiResponse.Code = int32(httpResponse.StatusCode()) 
  apiResponse.Message = httpResponse.Status()

  return apiResponse
}

func (c *ApiClient) SetErrorApiResponse(errorMessage string) *ApiResponse{

  apiResponse := new(ApiResponse) 
  apiResponse.Code = int32(400) 
  apiResponse.Message = errorMessage

  return apiResponse
}

func prepareRequest(postBody interface{},
    headerParams map[string]string,
    queryParams map[string]string,
    formParams map[string]string,
    fileName string,
    fileBytes []byte) *resty.Request {

    request := resty.R()

    request.SetBody(postBody)

    // add header parameter, if any
    if len(headerParams) > 0 {
        request.SetHeaders(headerParams)
    }
    
    // add query parameter, if any
    if len(queryParams) > 0 {
        request.SetQueryParams(queryParams)
    }

    // add form parameter, if any
    if len(formParams) > 0 {
        request.SetFormData(formParams)
    }
    
    if len(fileBytes) > 0 && fileName != "" {
        _, fileNm := filepath.Split(fileName)
        request.SetFileReader("file", fileNm, bytes.NewReader(fileBytes))
    }
    return request
}
