package addition

import (
	"encoding/json"
	"future_today/internal/config"
	"io"
	"net/http"
)

type Addition struct {
	http.Client
	agifyURL       string
	nationalizeURL string
	genderizeURL   string
}

// no async
type AdditionResponse struct {
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Country []struct {
		CountryID string `json:"country_id"`
	} `json:"country"`
}

// async
type AgeResponse struct {
	Age int `json:"age"`
}

type GenderResponse struct {
	Gender string `json:"gender"`
}

type NationResponse struct {
	Country []struct {
		CountryID string `json:"country_id"`
	} `json:"country"`
}

func NewAddition(cfg *config.Config) *Addition {
	return &Addition{agifyURL: cfg.AgifyURL, nationalizeURL: cfg.NationalizeURL, genderizeURL: cfg.GenderizeURL}
}

// добавить асинхронности
// func (add *Addition) GetAddition(name string) (*AdditionResponse, error) {
// 	var addResp AdditionResponse
// 	resp_age, err := add.Get(add.agifyURL + "/?name=" + name)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting age by url:%v", err)
// 	}
// 	bytes, err := io.ReadAll(resp_age.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error in unmarshalling:%v", err)
// 	}
// 	err = json.Unmarshal(bytes, &addResp)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading body:%v", err)
// 	}
// 	respGender, err := add.Get(add.genderizeURL + "/?name=" + name)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting gender by url:%v", err)
// 	}
// 	bytes, err = io.ReadAll(respGender.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error in unmarshalling:%v", err)
// 	}
// 	err = json.Unmarshal(bytes, &addResp)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading body:%v", err)
// 	}
// 	respNation, err := add.Get(add.nationalizeURL + "/?name=" + name)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting nation by url:%v", err)
// 	}
// 	bytes, err = io.ReadAll(respNation.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error in unmarshalling:%v", err)
// 	}
// 	err = json.Unmarshal(bytes, &addResp)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading body:%v", err)
// 	}

// 	return &addResp, nil

// }

func (add *Addition) GetAdditionAsync(name string) (int, string, string, error) {
	var ageResp AgeResponse
	var genderResp GenderResponse
	var natResp NationResponse

	ageChan := make(chan int)
	nationChan := make(chan string)
	genderChan := make(chan string)
	errChan := make(chan error)

	go func() {
		resp_age, err := add.Get(add.agifyURL + "/?name=" + name)
		if err != nil {
			errChan <- err
			return
		}
		bytes, err := io.ReadAll(resp_age.Body)
		if err != nil {
			errChan <- err
			return
		}
		err = json.Unmarshal(bytes, &ageResp)
		if err != nil {
			errChan <- err
			return
		}
		ageChan <- ageResp.Age
	}()

	go func() {
		respGender, err := add.Get(add.genderizeURL + "/?name=" + name)
		if err != nil {
			errChan <- err
			return
		}
		bytes, err := io.ReadAll(respGender.Body)
		if err != nil {
			errChan <- err
			return
		}
		err = json.Unmarshal(bytes, &genderResp)
		if err != nil {
			errChan <- err
			return
		}
		genderChan <- genderResp.Gender
	}()

	go func() {
		respNation, err := add.Get(add.nationalizeURL + "/?name=" + name)
		if err != nil {
			errChan <- err
			return
		}
		bytes, err := io.ReadAll(respNation.Body)
		if err != nil {
			errChan <- err
			return
		}
		err = json.Unmarshal(bytes, &natResp)
		if err != nil {
			errChan <- err
			return
		}
		nationChan <- natResp.Country[0].CountryID
	}()
	var age int
	var gender, nationality string
	var err error
	for i := 0; i < 3; i++ {
		select {
		case age = <-ageChan:
		case gender = <-genderChan:
		case nationality = <-nationChan:
		case err = <-errChan:
			return 0, "", "", err
		}
	}
	return age, gender, nationality, nil
}
