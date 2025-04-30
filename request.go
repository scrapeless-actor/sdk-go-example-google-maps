package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var c *http.Client

func InitProxyClient(proxy string) {
	// proxy addr
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}

	// custom Transport
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	// create Client
	client := &http.Client{
		Transport: transport,
	}
	c = client
}

func doMapsManage(ctx context.Context, params *RequestParam) (*Response, error) {
	err := params.FieldValidation()
	if err != nil {
		return nil, err
	}

	var urlStr string
	if params.PlaceId == "" {
		if params.Type == "search" {
			if params.Ll != "" {
				urlStr = fmt.Sprintf("https://www.%s/maps/search/%s/%s?hl=%s&gl=%s", params.GoogleDomain, params.Q, params.Ll, params.Hl, params.Gl)
				if params.Start != "" {
					latitude, longitude, err := extractLatLong(params.Ll)
					if err != nil {
						return nil, fmt.Errorf(err.Error())
					}
					pb := fmt.Sprintf("!4m9!1m3!1d295%s!2d%s!3d%s!2m0!3m2!1i784!2i644!4f13.1!7i20!8i%s!10b1!12m8!1m1!18b1!2m3!5m1!6e2!20e3!10b1!16b1!19m4!2m3!1i360!2i120!4i8!20m57!2m2!1i203!2i100!3m2!2i4!5b1!6m6!1m2!1i86!2i86!1m2!1i408!2i240!7m42!1m3!1e1!2b0!3e3!1m3!1e2!2b1!3e2!1m3!1e2!2b0!3e3!1m3!1e3!2b0!3e3!1m3!1e8!2b0!3e3!1m3!1e3!2b1!3e2!1m3!1e9!2b1!3e2!1m3!1e10!2b0!3e3!1m3!1e10!2b1!3e2!1m3!1e10!2b0!3e4!2b1!4b1!9b0!22m6!1sg3qzXsG-JpeGoATHyYKQBw%%3A1!2zMWk6MSx0OjExODg3LGU6MCxwOmczcXpYc0ctSnBlR29BVEh5WUtRQnc6MQ!7e81!12e3!17sg3qzXsG-JpeGoATHyYKQBw%%3A110!18e15!24m46!1m12!13m6!2b1!3b1!4b1!6i1!8b1!9b1!18m4!3b1!4b1!5b1!6b1!2b1!5m5!2b1!3b1!5b1!6b1!7b1!10m1!8e3!14m1!3b1!17b1!20m2!1e3!1e6!24b1!25b1!26b1!30m1!2b1!36b1!43b1!52b1!55b1!56m2!1b1!3b1!65m5!3m4!1m3!1m2!1i224!2i298!26m4!2m3!1i80!2i92!4i8!30m28!1m6!1m2!1i0!2i0!2m2!1i458!2i644!1m6!1m2!1i734!2i0!2m2!1i784!2i644!1m6!1m2!1i0!2i0!2m2!1i784!2i20!1m6!1m2!1i0!2i624!2m2!1i784!2i644!34m13!2b1!3b1!4b1!6b1!8m3!1b1!3b1!4b1!9b1!12b1!14b1!20b1!23b1!37m1!1e81!42b1!47m0!49m1!3b1!50m4!2e2!3m2!1b1!3b0!65m0",
						latitude, longitude, latitude, params.Start)
					urlStr = fmt.Sprintf("https://www.google.com/search?q=%s&hl=%s&tbm=map&tch=1&pb=%s", url.PathEscape(params.Q), params.Hl, pb)
				}
			} else {
				urlStr = fmt.Sprintf("https://www.%s/maps/search/%s?hl=%s&gl=%s", params.GoogleDomain, params.Q, params.Hl, params.Gl)
			}
		} else if params.Type == "place" {
			urlStr = fmt.Sprintf("https://www.%s/maps/place/%s/data=%s?hl=%s&gl=%s", params.GoogleDomain, params.Q, params.Data, params.Hl, params.Gl)
		}
	} else {
		urlStr = fmt.Sprintf("https://www.%s/maps/place/?q=place_id:%s&hl=%s&gl=%s", params.GoogleDomain, params.PlaceId, params.Hl, params.Gl)
	}
	req, reqError := http.NewRequest("GET", urlStr, nil)
	if reqError != nil {
		return nil, fmt.Errorf("reqError: %v", reqError)
	}
	req.Header = http.Header{
		"accept":                            {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-language":                   {"en-US,en;q=0.9"},
		"cache-control":                     {"no-cache"},
		"pragma":                            {"no-cache"},
		"priority":                          {"u=0, i"},
		"sec-ch-ua":                         {`"Not A(Brand";v="8", "Chromium";v="132", "Google Chrome";v="132"`},
		"sec-ch-ua-arch":                    {`"x86"`},
		"sec-ch-ua-bitness":                 {`"64"`},
		"sec-ch-ua-form-factors":            {`"Desktop"`},
		"sec-ch-ua-full-version":            {`"132.0.6834.83"`},
		"sec-ch-ua-full-version-list":       {`"Not A(Brand";v="8.0.0.0", "Chromium";v="132.0.6834.83", "Google Chrome";v="132.0.6834.83"`},
		"sec-ch-ua-mobile":                  {"?0"},
		"sec-ch-ua-platform":                {`"Windows"`},
		"sec-ch-ua-platform-version":        {`"10.0.0"`},
		"sec-ch-ua-wow64":                   {"?0"},
		"sec-fetch-dest":                    {"document"},
		"sec-fetch-mode":                    {"navigate"},
		"sec-fetch-site":                    {"none"},
		"sec-fetch-user":                    {"?1"},
		"service-worker-navigation-preload": {"true"},
		"upgrade-insecure-requests":         {"1"},
		"user-agent":                        {"Mozilla/5.0 (Windows NT 10.0; Win64; x64}, AppleWebKit/537.36 (KHTML, like Gecko}, Chrome/132.0.0.0 Safari/537.36"},
	}
	resp, respError := c.Do(req)
	if respError != nil {
		return nil, fmt.Errorf("respError: %v", respError)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("resp.StatusCode: %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	bodyStr := string(bodyText)
	if err != nil {
		return nil, fmt.Errorf("io resp.Body.ReadAll: %v", err)
	}
	var dataStr string
	isNotStart := params.Start == ""
	if isNotStart {
		dataStr = extractBodyByHtml(bodyText)
	} else {
		// 去除 /*""*/
		cleanJSON := strings.TrimSuffix(bodyStr, "/*\"\"*/")
		// 解析 JSON 数据到结构体
		var data jsonData
		err := json.Unmarshal([]byte(cleanJSON), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %v", err)
		}
		dataStr = strings.ReplaceAll(data.D, ")]}'\n", "")
	}

	searchResultFunc := func(mapsData string, isNotStart bool) (Response, error) {
		var response Response // 最终返回的结构体

		var localResultStr string
		if isNotStart {
			localResultStr = jsoniter.Get([]byte(mapsData), 64).ToString()
		} else {
			localResultStr = jsoniter.Get([]byte(mapsData), 0, 1).ToString()
		}

		var localResultArr []any
		if localResultStr != "" {
			if err = json.Unmarshal([]byte(localResultStr), &localResultArr); err != nil {
				return Response{}, fmt.Errorf("failed to parse JSON: %v", err)
			}
		} else {
			return Response{}, nil
		}
		var count int
		var localResultsList []LocalResults
		for i := range len(localResultArr) {
			var dataResultStr string
			if isNotStart {
				dataResultStr = jsoniter.Get([]byte(localResultStr), i, 1).ToString()
			} else {
				dataResultStr = jsoniter.Get([]byte(localResultStr), i, 14).ToString()
			}

			title := jsoniter.Get([]byte(dataResultStr), 11).ToString()
			if title == "" {
				continue
			}
			reviews, _ := strconv.Atoi(extractNumbersUsingMap(jsoniter.Get([]byte(dataResultStr), 4, 8).ToString()))
			gpsCoordinates := &GpsCoordinates{
				Latitude:  jsoniter.Get([]byte(dataResultStr), 9, 2).ToFloat64(),
				Longitude: jsoniter.Get([]byte(dataResultStr), 9, 3).ToFloat64(),
			}
			var types []string
			_ = json.Unmarshal([]byte(jsoniter.Get([]byte(dataResultStr), 13).ToString()), &types)
			var typeIdsArr []any
			_ = json.Unmarshal([]byte(jsoniter.Get([]byte(dataResultStr), 76).ToString()), &typeIdsArr)
			var typeIds []string
			if len(typeIdsArr) > 0 {
				for i2 := range len(typeIdsArr) {
					str, _ := json.Marshal(typeIdsArr[i2])
					var strArr []any
					_ = json.Unmarshal(str, &strArr)
					for i3 := range strArr {
						str2, _ := json.Marshal(strArr[i3])
						typeIds = append(typeIds, string(str2))
					}
				}
			}

			operatingHours := make(map[string]interface{})
			openTimeStr := jsoniter.Get([]byte(dataResultStr), 203, 0).ToString()
			var openTimeArr []any
			_ = json.Unmarshal([]byte(openTimeStr), &openTimeArr)
			if len(openTimeArr) > 0 {
				for i2 := range len(openTimeArr) {
					key := jsoniter.Get([]byte(dataResultStr), 203, 0, i2, 0).ToString()
					value := jsoniter.Get([]byte(dataResultStr), 203, 0, i2, 3, 0, 0).ToString()
					operatingHours[key] = value
				}
			}

			count++
			localResults := LocalResults{
				Position:       count,
				Title:          jsoniter.Get([]byte(dataResultStr), 11).ToString(),
				PlaceID:        jsoniter.Get([]byte(dataResultStr), 78).ToString(),
				DataID:         jsoniter.Get([]byte(dataResultStr), 10).ToString(),
				DataCid:        jsoniter.Get([]byte(dataResultStr), 57, 8).ToString(),
				ReviewsLink:    jsoniter.Get([]byte(dataResultStr), 4, 3, 0).ToString(),
				PhotosLink:     jsoniter.Get([]byte(dataResultStr), 37, 0, 0, 6, 0).ToString(),
				ProviderID:     jsoniter.Get([]byte(dataResultStr), 89).ToString(),
				Rating:         jsoniter.Get([]byte(dataResultStr), 4, 7).ToFloat64(),
				Reviews:        reviews,
				GpsCoordinates: gpsCoordinates,
				PlaceIDSearch:  jsoniter.Get([]byte(dataResultStr), 174, 0).ToString(),
				Price:          jsoniter.Get([]byte(dataResultStr), 4, 2).ToString(),
				Type:           jsoniter.Get([]byte(dataResultStr), 13, 0).ToString(),
				Types:          types,
				TypeID:         jsoniter.Get([]byte(dataResultStr), 76, 0, 0).ToString(),
				TypeIds:        typeIds,
				Address:        jsoniter.Get([]byte(dataResultStr), 18).ToString(),
				OpenState:      jsoniter.Get([]byte(dataResultStr), 203, 1, 4, 0).ToString(),
				Hours:          jsoniter.Get([]byte(dataResultStr), 203, 1, 4, 0).ToString(),
				Phone:          jsoniter.Get([]byte(dataResultStr), 178, 0, 0).ToString(),
				Website:        jsoniter.Get([]byte(dataResultStr), 7, 0).ToString(),
				Description:    jsoniter.Get([]byte(dataResultStr), 32, 1, 1).ToString(),
				OperatingHours: operatingHours,
				OrderOnline:    jsoniter.Get([]byte(dataResultStr), 75, 0, 0, 5, 1, 2, 0).ToString(),
				Thumbnail:      jsoniter.Get([]byte(dataResultStr), 157).ToString(),
			}

			var serviceOptionsArr []any
			_ = json.Unmarshal([]byte(jsoniter.Get([]byte(dataResultStr), 100, 1, 0, 2).ToString()), &serviceOptionsArr)
			var serviceOptions []string
			for i2 := range len(serviceOptionsArr) {
				serviceOptions = append(serviceOptions, jsoniter.Get([]byte(dataResultStr), 100, 1, 0, 2, i2, 1).ToString())
			}
			if len(serviceOptions) > 0 {
				localResults.ServiceOptions = serviceOptions
			}

			// if !localResults.IsEmpty() {
			localResultsList = append(localResultsList, localResults)
			// }
		}

		response.LocalResults = localResultsList

		return response, nil
	}

	// type = place or placeId != ""
	placeResultFunc := func(mapsData string) (Response, error) {
		var response Response
		placeResults := PlaceResults{
			Title:          jsoniter.Get([]byte(mapsData), 6, 11).ToString(),
			PlaceId:        jsoniter.Get([]byte(mapsData), 6, 78).ToString(),
			DataId:         jsoniter.Get([]byte(mapsData), 6, 10).ToString(),
			DataCid:        jsoniter.Get([]byte(mapsData), 6, 37, 0, 0, 29, 1).ToString(),
			ProviderId:     jsoniter.Get([]byte(mapsData), 6, 89).ToString(),
			Thumbnail:      jsoniter.Get([]byte(mapsData), 6, 72, 0, 0, 6, 0).ToString(),
			Rating:         jsoniter.Get([]byte(mapsData), 6, 4, 7).ToFloat64(),
			Reviews:        jsoniter.Get([]byte(mapsData), 6, 4, 8).ToInt(),
			Price:          jsoniter.Get([]byte(mapsData), 6, 4, 2).ToString(),
			Description:    jsoniter.Get([]byte(mapsData), 6, 32, 1, 1).ToString(),
			Phone:          jsoniter.Get([]byte(mapsData), 6, 178, 0, 0).ToString(),
			OpenState:      jsoniter.Get([]byte(mapsData), 6, 34, 4, 4).ToString(),
			PlusCode:       jsoniter.Get([]byte(mapsData), 6, 183, 2, 2, 0).ToString(),
			WebResultsLink: jsoniter.Get([]byte(mapsData), 6, 174, 0).ToString(),
		}
		gpsCoordinates := GpsCoordinates{
			Latitude:  jsoniter.Get([]byte(mapsData), 6, 9, 2).ToFloat64(),
			Longitude: jsoniter.Get([]byte(mapsData), 6, 9, 3).ToFloat64(),
		}
		if !gpsCoordinates.IsEmpty() {
			placeResults.GpsCoordinates = &gpsCoordinates
		}

		var ratingSummaryList []RatingSummary
		var ratingSummaryArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 175, 3).ToString()), &ratingSummaryArr)
		for i := range len(ratingSummaryArr) {
			ratingSummary := RatingSummary{
				Starts: i + 1,
				Amount: jsoniter.Get([]byte(mapsData), 6, 175, 3, i).ToInt(),
			}
			if !ratingSummary.IsEmpty() {
				ratingSummaryList = append(ratingSummaryList, ratingSummary)
			}
		}
		if len(ratingSummaryList) > 0 {
			placeResults.RatingSummary = ratingSummaryList
		}

		var types []string
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 13).ToString()), &types)
		placeResults.Type = types

		var typeIds []string
		var typeIdsArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 76).ToString()), &typeIdsArr)
		for i := range len(typeIdsArr) {
			typeIds = append(typeIds, jsoniter.Get([]byte(mapsData), 6, 76, i, 0).ToString())
		}
		placeResults.TypeIds = typeIds

		menu := Menu{
			Link:   jsoniter.Get([]byte(mapsData), 6, 38, 0).ToString(),
			Source: jsoniter.Get([]byte(mapsData), 6, 38, 1).ToString(),
		}
		if menu.Link != "" {
			menu.Link = fmt.Sprintf("https://www.google.com%s", menu.Link)
		}
		if !menu.IsEmpty() {
			placeResults.Menu = &menu
		}

		if jsoniter.Get([]byte(mapsData), 6, 49, 1).ToString() != "" {
			placeResults.UnclaimedListing = true
		}

		serviceOptions := ServiceOptions{}
		if jsoniter.Get([]byte(mapsData), 6, 100, 3, 0, 2, 1, 0, 0).ToInt() == 1 {
			serviceOptions.DineIn = true
		} else {
			serviceOptions.DineIn = false
		}
		if jsoniter.Get([]byte(mapsData), 6, 100, 3, 1, 2, 1, 0, 0).ToInt() == 1 {
			serviceOptions.Takeout = true
		} else {
			serviceOptions.Takeout = false
		}
		if jsoniter.Get([]byte(mapsData), 6, 100, 3, 2, 2, 1, 0, 0).ToInt() == 1 {
			serviceOptions.Delivery = true
		} else {
			serviceOptions.Delivery = false
		}
		if !serviceOptions.IsEmpty() {
			placeResults.ServiceOptions = &serviceOptions
		}

		var extensions []map[string]interface{}
		var unsupportedExtensions []map[string]interface{}
		var extensionsArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 100, 1).ToString()), &extensionsArr)
		for i := range len(extensionsArr) {
			extensionMap := make(map[string]interface{})
			unsupportedExtensionMap := make(map[string]interface{})
			var extensionMapValue []string
			var unsupportedExtensionMapValue []string
			var extensionMapValueArr []any
			_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 100, 1, i, 2).ToString()), &extensionMapValueArr)
			if len(extensionMapValueArr) == 0 {
				continue
			}
			for i2 := range len(extensionMapValueArr) {
				if jsoniter.Get([]byte(mapsData), 6, 100, 1, i, 2, i2, 2, 2, 0).ToString() != "" && jsoniter.Get([]byte(mapsData), 6, 100, 1, i, 2, i2, 2, 2, 0).ToInt() == 0 {
					unsupportedExtensionMapValue = append(unsupportedExtensionMapValue, jsoniter.Get([]byte(mapsData), 6, 100, 1, i, 2, i2, 1).ToString())
				} else if jsoniter.Get([]byte(mapsData), 6, 100, 1, i, 2, i2, 2, 2, 0).ToInt() == 1 {
					extensionMapValue = append(extensionMapValue, jsoniter.Get([]byte(mapsData), 6, 100, 1, i, 2, i2, 1).ToString())
				}
			}
			if len(extensionMapValue) > 0 {
				extensionMap[jsoniter.Get([]byte(mapsData), 6, 100, 1, i, 0).ToString()] = extensionMapValue
				extensions = append(extensions, extensionMap)
			}
			if len(unsupportedExtensionMapValue) > 0 {
				unsupportedExtensionMap[jsoniter.Get([]byte(mapsData), 6, 100, 1, i, 0).ToString()] = unsupportedExtensionMapValue
				unsupportedExtensions = append(unsupportedExtensions, unsupportedExtensionMap)
			}
		}
		if len(extensions) > 0 {
			placeResults.Extensions = extensions
		}
		if len(unsupportedExtensions) > 0 {
			placeResults.UnsupportedExtensions = unsupportedExtensions
		}

		var addressArr []string
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 2).ToString()), &addressArr)
		placeResults.Address = strings.Join(addressArr, ",")

		if jsoniter.Get([]byte(mapsData), 6, 7, 1).ToString() != "" {
			placeResults.Website = fmt.Sprintf("https://www.%s", jsoniter.Get([]byte(mapsData), 6, 7, 1).ToString())
		}

		var hours []map[string]interface{}
		var hoursArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 203, 0).ToString()), &hoursArr)
		for i := range len(hoursArr) {
			hourMap := make(map[string]interface{})
			hourMap[jsoniter.Get([]byte(mapsData), 6, 203, 0, i, 0).ToString()] = jsoniter.Get([]byte(mapsData), 6, 203, 0, i, 3, 0, 0).ToString()
			hours = append(hours, hourMap)
		}
		if len(hours) > 0 {
			placeResults.Hours = hours
		}

		var imageList []Image
		var imageArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 171, 0).ToString()), &imageArr)
		for i := range len(imageArr) {
			image := Image{
				Title:     jsoniter.Get([]byte(mapsData), 6, 171, 0, i, 2).ToString(),
				Thumbnail: jsoniter.Get([]byte(mapsData), 6, 171, 0, i, 3, 0, 6, 0).ToString(),
			}
			if !image.IsEmpty() {
				placeResults.Images = append(placeResults.Images, image)
			}
		}
		if len(imageList) > 0 {
			placeResults.Images = imageList
		}

		var questionsAndAnswerList []QuestionsAndAnswers
		var questionsAndAnswerArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 126, 0).ToString()), &questionsAndAnswerArr)
		for i := range len(questionsAndAnswerArr) {
			questionsAndAnswer := QuestionsAndAnswers{TotalAnswers: jsoniter.Get([]byte(mapsData), 6, 126, 0, i, 2).ToInt()}
			question := Question{
				Text:     jsoniter.Get([]byte(mapsData), 6, 126, 0, i, 0, 2).ToString(),
				Data:     jsoniter.Get([]byte(mapsData), 6, 126, 0, i, 0, 7).ToString(),
				Language: jsoniter.Get([]byte(mapsData), 6, 126, 0, i, 0, 12).ToString(),
			}
			if !question.IsEmpty() {
				questionsAndAnswer.Question = &question
			}
			var answerList []Answer
			var answerArr []any
			_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 126, 0, i, 1).ToString()), &answerArr)
			for i2 := range len(answerArr) {
				answer := Answer{
					Text:     jsoniter.Get([]byte(mapsData), 6, 126, 0, i, 1, i2, 2).ToString(),
					Data:     jsoniter.Get([]byte(mapsData), 6, 126, 0, i, 1, i2, 7).ToString(),
					Language: jsoniter.Get([]byte(mapsData), 6, 126, 0, i, 1, i2, 12).ToString(),
				}
				if !answer.IsEmpty() {
					answerList = append(answerList, answer)
				}
			}
			if len(answerList) > 0 {
				questionsAndAnswer.Answer = answerList
			}
			if !questionsAndAnswer.IsEmpty() {
				questionsAndAnswerList = append(questionsAndAnswerList, questionsAndAnswer)
			}
		}
		if len(questionsAndAnswerList) > 0 {
			placeResults.QuestionsAndAnswers = questionsAndAnswerList
		}

		var userReview UserReview
		var userReviewArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 175, 9, 0, 0).ToString()), &userReviewArr)
		var mostRelevantList []MostRelevant
		for i := range len(userReviewArr) {
			mostRelevant := MostRelevant{
				Username:      jsoniter.Get([]byte(mapsData), 6, 175, 9, 0, 0, i, 0, 1, 4, 5, 0).ToString(),
				Rating:        jsoniter.Get([]byte(mapsData), 6, 175, 9, 0, 0, i, 0, 2, 0, 0).ToInt(),
				ContributorId: jsoniter.Get([]byte(mapsData), 6, 175, 9, 0, 0, i, 0, 1, 4, 5, 3).ToString(),
				Description:   jsoniter.Get([]byte(mapsData), 6, 175, 9, 0, 0, i, 0, 2, 15, 0, 0).ToString(),
				Link:          jsoniter.Get([]byte(mapsData), 6, 175, 9, 0, 0, i, 0, 4, 3, 0).ToString(),
				Date:          jsoniter.Get([]byte(mapsData), 6, 175, 9, 0, 0, i, 0, 1, 6).ToString(),
			}
			var imageList2 []Image
			var imageArr2 []any
			_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 175, 9, 0, 0, i, 0, 2, 2).ToString()), &imageArr2)
			for i2 := range len(imageArr2) {
				image := Image{Thumbnail: jsoniter.Get([]byte(mapsData), 6, 175, 9, 0, 0, i, 0, 2, 2, i2, 1, 6, 0).ToString()}
				if !image.IsEmpty() {
					imageList2 = append(imageList2, image)
				}
			}
			if len(imageList2) > 0 {
				mostRelevant.Images = imageList2
			}
			if !mostRelevant.IsEmpty() {
				mostRelevantList = append(mostRelevantList, mostRelevant)
			}
		}
		if len(mostRelevantList) > 0 {
			userReview.MostRelevant = mostRelevantList
		}
		var summaryList []Summary
		var summaryArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 31, 1).ToString()), &summaryArr)
		for i := range len(summaryArr) {
			summary := Summary{Snippet: jsoniter.Get([]byte(mapsData), 6, 31, 1, i, 1).ToString()}
			if !summary.IsEmpty() {
				summaryList = append(summaryList, summary)
			}
		}
		if len(summaryList) > 0 {
			userReview.Summary = summaryList
		}
		if !userReview.IsEmpty() {
			placeResults.UserReviews = &userReview
		}

		var peopleAlsoSearchList []PeopleAlsoSearchFor
		var peopleAlsoSearchArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 99, 0).ToString()), &peopleAlsoSearchArr)
		for i := range len(peopleAlsoSearchArr) {
			peopleAlsoSearch := PeopleAlsoSearchFor{SearchTerm: jsoniter.Get([]byte(mapsData), 6, 99, 0, i, 0).ToString()}
			var localResultList []LocalResult
			var localResultArr []any
			_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 99, 0, i, 1).ToString()), &localResultArr)
			var count int
			for i2 := range len(localResultArr) {
				if jsoniter.Get([]byte(mapsData), 6, 99, 0, i, 1, i2, 1).ToString() == "" {
					continue
				}
				count++
				localResult := LocalResult{
					Position:  count,
					Title:     jsoniter.Get([]byte(mapsData), 6, 99, 0, i, 1, i2, 1, 11).ToString(),
					DataID:    jsoniter.Get([]byte(mapsData), 6, 99, 0, i, 1, i2, 1, 10).ToString(),
					Thumbnail: jsoniter.Get([]byte(mapsData), 6, 99, 0, i, 1, i2, 1, 37, 0, 0, 6, 0).ToString(),
					Rating:    jsoniter.Get([]byte(mapsData), 6, 99, 0, i, 1, i2, 1, 4, 7).ToFloat64(),
					Reviews:   jsoniter.Get([]byte(mapsData), 6, 99, 0, i, 1, i2, 1, 4, 8).ToInt(),
				}
				gpsCoordinates2 := GpsCoordinates{
					Latitude:  jsoniter.Get([]byte(mapsData), 6, 99, 0, i, 1, i2, 1, 9, 2).ToFloat64(),
					Longitude: jsoniter.Get([]byte(mapsData), 6, 99, 0, i, 1, i2, 1, 9, 3).ToFloat64(),
				}
				if !gpsCoordinates2.IsEmpty() {
					localResult.ReviewsLink = fmt.Sprintf("https://www.google.com/maps/place/data=!4m7!3m6!1s%s!5m2!4m1!1i2!9m1!1b1", localResult.DataID)
					localResult.GPSCoordinates = &gpsCoordinates2
				}
				_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 99, 0, i, 1, i2, 1, 13).ToString()), &localResult.Type)
				if !localResult.IsEmpty() {
					localResultList = append(localResultList, localResult)
				}
			}
			if len(localResultList) > 0 {
				peopleAlsoSearch.LocalResults = localResultList
			}
			if !peopleAlsoSearch.IsEmpty() {
				peopleAlsoSearchList = append(peopleAlsoSearchList, peopleAlsoSearch)
			}
		}
		if len(peopleAlsoSearchList) > 0 {
			placeResults.PeopleAlsoSearchFor = peopleAlsoSearchList
		}

		popularTimes := PopularTimes{}
		GraphResults := make(map[string]interface{})
		var graphResultsInfoArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 84, 0).ToString()), &graphResultsInfoArr)
		for i := range len(graphResultsInfoArr) {
			weekInt := jsoniter.Get([]byte(mapsData), 6, 84, 0, i, 0).ToInt()
			key := weekMapping(weekInt)
			var valueArr []any
			_ = json.Unmarshal([]byte(jsoniter.Get([]byte(mapsData), 6, 84, 0, i, 1).ToString()), &valueArr)
			var graphResultsInfoList []GraphResultsInfo
			for i2 := range len(valueArr) {
				graphResultsInfo := GraphResultsInfo{
					Time:          jsoniter.Get([]byte(mapsData), 6, 84, 0, i, 1, i2, 4).ToString(),
					Info:          jsoniter.Get([]byte(mapsData), 6, 84, 0, i, 1, i2, 2).ToString(),
					BusynessScore: jsoniter.Get([]byte(mapsData), 6, 84, 0, i, 1, i2, 1).ToInt(),
				}
				if !graphResultsInfo.IsEmpty() {
					graphResultsInfoList = append(graphResultsInfoList, graphResultsInfo)
				}
			}
			GraphResults[key] = graphResultsInfoList
		}
		if len(GraphResults) > 0 {
			popularTimes.GraphResults = GraphResults
		}
		liveHash := LiveHash{
			Info:      jsoniter.Get([]byte(mapsData), 6, 117, 1).ToString(),
			TimeSpent: jsoniter.Get([]byte(mapsData), 6, 117, 0).ToString(),
		}
		if !liveHash.IsEmpty() {
			popularTimes.LiveHash = &liveHash
		}
		if !popularTimes.IsEmpty() {
			placeResults.PopularTimes = &popularTimes
		}

		if !placeResults.IsEmpty() {
			placeResults.ReviewsLink = fmt.Sprintf("https://www.google.com/maps/place/data=!4m7!3m6!1s%s!5m2!4m1!1i2!9m1!1b1", placeResults.DataId)
			response.PlaceResults = &placeResults
		}
		return response, nil
	}

	var response Response

	var mapsData string
	if params.PlaceId != "" || params.Type == "place" {
		mapsData = strings.ReplaceAll(jsoniter.Get([]byte(dataStr), 3, 6).ToString(), ")]}'\n", "")
		response, err = placeResultFunc(mapsData)
		if err != nil {
			return nil, err
		}
	} else {
		if isNotStart {
			mapsData = strings.ReplaceAll(jsoniter.Get([]byte(dataStr), 3, 2).ToString(), ")]}'\n", "")
		} else {
			mapsData = dataStr
		}
		response, err = searchResultFunc(mapsData, isNotStart)
		if err != nil {
			return nil, fmt.Errorf("searchResultFunc err: %v", err)
		}
	}
	return &response, nil
}

func DoMaps(ctx context.Context, params *RequestParam) (*Response, error) {
	response, err := doMapsManage(ctx, params)
	if err != nil {
		return &Response{}, err
	}
	return response, nil
}

func DoMapsAutocomplete(ctx context.Context, params *RequestParam) (*Response, error) {
	response := &Response{}
	llSplit := strings.Split(params.Ll, ",")
	if len(llSplit) < 2 {
		return nil, fmt.Errorf("ll format error")
	}

	pb := fmt.Sprintf("!2i6!4m9!1m3!1d362730.1311737605!2d%s!3d%s!2m0!3m2!1i2160!2i1440!4f13.1!7i20!10b1!12m16!1m1!18b1!2m3!5m1!6e2!20e3!10b1!12b1!13b1!16b1!17m1!3e1!20m3!5e2!6b1!14b1!19m4!2m3!1i360!2i120!4i8!20m57!2m2!1i203!2i100!3m2!2i4!5b1!6m6!1m2!1i86!2i86!1m2!1i408!2i240!7m42!1m3!1e1!2b0!3e3!1m3!1e2!2b1!3e2!1m3!1e2!2b0!3e3!1m3!1e8!2b0!3e3!1m3!1e10!2b0!3e3!1m3!1e10!2b1!3e2!1m3!1e9!2b1!3e2!1m3!1e10!2b0!3e3!1m3!1e10!2b1!3e2!1m3!1e10!2b0!3e4!2b1!4b1!9b0!22m2!1saSjN8nTmk850sCPMcoo3o-8!7e81!23m2!4b1!10b1!24m82!1m29!13m9!2b1!3b1!4b1!6i1!8b1!9b1!14b1!20b1!25b1!18m18!3b1!4b1!5b1!6b1!9b1!12b1!13b1!14b1!15b1!17b1!20b1!21b1!22b0!25b1!27m1!1b0!28b0!30b0!2b1!5m6!2b1!3b1!5b1!6b1!7b1!10b1!10m1!8e3!11m1!3e1!14m1!3b1!17b1!20m2!1e3!1e6!24b1!25b1!26b1!29b1!30m1!2b1!36b1!39m3!2m2!2i1!3i1!43b1!52b1!54m1!1b1!55b1!56m2!1b1!3b1!65m5!3m4!1m3!1m2!1i224!2i298!71b1!72m4!1m2!3b1!5b1!4b1!89b1!103b1!113b1!26m4!2m3!1i80!2i92!4i8!34m18!2b1!3b1!4b1!6b1!8m6!1b1!3b1!4b1!5b1!6b1!7b1!9b1!12b1!14b1!20b1!23b1!25b1!26b1!37m1!1e81!47m0!49m6!3b1!6m2!1b1!2b1!7m1!1e3!67m2!7b1!10b1!69i648", llSplit[1], strings.ReplaceAll(llSplit[0], "@", ""))

	var urlQuery = url.Values{
		"q":     {params.Q},
		"hl":    {params.Hl},
		"gl":    {params.Gl},
		"gs_ri": {"maps"},
		"pb":    {pb},
	}
	req, err := http.NewRequest("GET", "https://www.google.com/s?"+urlQuery.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	req.Header = http.Header{
		"accept":                      {"*/*"},
		"accept-language":             {"en-US,en;q=0.9"},
		"cache-control":               {"no-cache"},
		"pragma":                      {"no-cache"},
		"priority":                    {"u=1, i"},
		"referer":                     {"https://www.google.com/"},
		"sec-ch-ua":                   {`"Not A(Brand";v="8", "Chromium";v="132", "Google Chrome";v="132"`},
		"sec-ch-ua-arch":              {`"x86"`},
		"sec-ch-ua-form-factors":      {`"Desktop"`},
		"sec-ch-ua-full-version":      {`"132.0.6834.83"`},
		"sec-ch-ua-full-version-list": {`"Not A(Brand";v="8.0.0.0", "Chromium";v="132.0.6834.83", "Google Chrome";v="132.0.6834.83"`},
		"sec-ch-ua-platform":          {`"Windows"`},
		"sec-ch-ua-platform-version":  {`"10.0.0"`},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64}, AppleWebKit/537.36 (KHTML, like Gecko}, Chrome/132.0.0.0 Safari/537.36"},
	}
	resp, respError := c.Do(req)
	if respError != nil {
		return nil, fmt.Errorf("request %s fail error:%v", req.URL.String(), respError.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request %s fail error:%v", req.URL.String(), respError.Error())
	}
	bodyText, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read body fail error:%v", err)
	}
	dataBody := strings.ReplaceAll(string(bodyText), ")]}'\n", "")
	var mapsAutocompleteArr []any
	err = json.Unmarshal([]byte(jsoniter.Get([]byte(dataBody), 0, 1).ToString()), &mapsAutocompleteArr)
	if err != nil {
		return nil, fmt.Errorf("unmarshal body fail error:%v", err)
	}
	var suggestionsList []Suggestions
	for i := range len(mapsAutocompleteArr) {
		one := jsoniter.Get([]byte(dataBody), 0, 1, i).ToString()
		suggestions := Suggestions{
			Value:      jsoniter.Get([]byte(one), 22, 1, 0).ToString(),
			Subtext:    jsoniter.Get([]byte(one), 22, 0, 0).ToString(),
			Type:       params.Type,
			Latitude:   jsoniter.Get([]byte(one), 22, 11, 2).ToFloat64(),
			Longitude:  jsoniter.Get([]byte(one), 22, 11, 3).ToFloat64(),
			DataID:     jsoniter.Get([]byte(one), 22, 13, 0, 0).ToString(),
			ProviderId: jsoniter.Get([]byte(one), 22, 13, 0, 10).ToString(),
			PhotosLink: jsoniter.Get([]byte(one), 22, 24, 6, 0).ToString(),
		}
		if !suggestions.IsEmpty() {
			suggestionsList = append(suggestionsList, suggestions)
		}
	}

	if len(suggestionsList) > 0 {
		response.Suggestions = suggestionsList
	}

	return response, nil
}

func DoMapsContributorReviews(ctx context.Context, params *RequestParam) (*Response, error) {
	urlStr := fmt.Sprintf("https://www.%s/maps/contrib/%s/reviews?hl=%s&gl=%s", params.GoogleDomain, params.ContributorId, params.Hl, params.Gl)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	response := &Response{}
	req.Header = http.Header{
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-language":           {"en-US,en;q=0.9"},
		"cache-control":             {"no-cache"},
		"pragma":                    {"no-cache"},
		"priority":                  {"u=0, i"},
		"sec-ch-ua":                 {`"Not A(Brand";v="8", "Chromium";v="132", "Google Chrome";v="132"`},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"sec-fetch-dest":            {"document"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-user":            {"?1"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64}, AppleWebKit/537.36 (KHTML, like Gecko}, Chrome/132.0.0.0 Safari/537.36"},
	}
	resp, respError := c.Do(req)
	if respError != nil {
		return nil, fmt.Errorf("request %s fail error:%v", req.URL.String(), respError.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request %s fail error:%v", req.URL.String(), respError.Error())
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body fail error:%v", err)
	}

	pattern := regexp.MustCompile(`window\.APP_INITIALIZATION_STATE\s*=\s*(.*?);\s*window\.APP_FLAGS`)
	match := pattern.FindStringSubmatch(string(bodyText))
	var dataStr string
	if len(match) > 1 {
		dataStr = strings.TrimSpace(match[1])
	}

	var arrData string
	if params.NextPageToken != "" {
		arrData2 := strings.ReplaceAll(jsoniter.Get([]byte(dataStr), 3, 8).ToString(), ")]}'\n", "")
		token := jsoniter.Get([]byte(arrData2), 1, 0).ToString()

		pb := fmt.Sprintf("!1s%s!2m3!1s%s!7e81!15i14416!6m2!4b1!7b1!9m0!16m6!1i20!4b1!5b1!6B%s!9b1!14b1!17m0!18m15!1m3!1d4078050.2656404977!2d-80.83924470000001!3d36.4484085!2m3!1f0!2f0!3f0!3m2!1i668!2i953!4f13.1!6m2!1f0!2f0", params.ContributorId, token, params.NextPageToken)
		nextPageUrlQuery := url.Values{
			"authuser": {"0"},
			"hl":       {params.Hl},
			"gl":       {params.Gl},
			"pb":       {pb},
		}

		nextPageUrl := fmt.Sprintf("https://www.%s/locationhistory/preview/mas?", params.GoogleDomain) + nextPageUrlQuery.Encode()
		nextPageReq, _ := http.NewRequest("GET", nextPageUrl, nil)
		nextPageReq.Header = http.Header{
			"accept":                       {"*/*"},
			"accept-language":              {"en-US,en;q=0.9"},
			"cache-control":                {"no-cache"},
			"downlink":                     {"10"},
			"pragma":                       {"no-cache"},
			"priority":                     {"u=1, i"},
			"referer":                      {"https://www.google.com/"},
			"rtt":                          {"200"},
			"sec-ch-ua":                    {`"Not A(Brand";v="8", "Chromium";v="132", "Google Chrome";v="132"`},
			"sec-ch-ua-mobile":             {"?0"},
			"sec-ch-ua-platform":           {`"Windows"`},
			"sec-fetch-dest":               {"empty"},
			"sec-fetch-mode":               {"cors"},
			"sec-fetch-site":               {"same-origin"},
			"user-agent":                   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64}, AppleWebKit/537.36 (KHTML, like Gecko}, Chrome/132.0.0.0 Safari/537.36"},
			"x-maps-diversion-context-bin": {"CAE="},
		}

		nextPageResp, respError := c.Do(nextPageReq)
		if respError != nil {
			return nil, fmt.Errorf("request %s fail error:%v", req.URL.String(), respError.Error())
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("request %s fail error:%v", req.URL.String(), respError.Error())
		}
		defer nextPageResp.Body.Close()
		bodyText2, err2 := io.ReadAll(nextPageResp.Body)
		if err2 != nil {
			return nil, fmt.Errorf("read body fail error:%v", err2.Error())
		}
		arrData = strings.ReplaceAll(string(bodyText2), ")]}'\n", "")
	} else {
		arrData = strings.ReplaceAll(jsoniter.Get([]byte(dataStr), 3, 9).ToString(), ")]}'\n", "")
	}

	contributor := Contributor{
		Name:      jsoniter.Get([]byte(arrData), 16, 0).ToString(),
		Thumbnail: jsoniter.Get([]byte(arrData), 16, 1, 6, 0).ToString(),
		Points:    jsoniter.Get([]byte(arrData), 16, 8, 1, 0).ToInt(),
		Level:     jsoniter.Get([]byte(arrData), 16, 8, 1, 1).ToInt(),
	}

	var contributionsArr []any
	_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 16, 8, 0).ToString()), &contributionsArr)
	contributions := make(map[string]interface{})
	for i := range len(contributionsArr) {
		key := jsoniter.Get([]byte(arrData), 16, 8, 0, i, 6).ToString()
		value := jsoniter.Get([]byte(arrData), 16, 8, 0, i, 9).ToString()
		contributions[key] = value
	}
	if len(contributionsArr) > 0 {
		contributor.Contributions = contributions
	}
	if !contributor.IsEmpty() {
		response.Contributor = &contributor
	}

	var reviewList []Reviews
	var reviewArr []any
	_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 24, 0).ToString()), &reviewArr)
	for i := range len(reviewArr) {
		review := Reviews{
			PlaceInfo: &PlaceInfo{
				Title:   jsoniter.Get([]byte(arrData), 24, 0, i, 1, 2).ToString(),
				Address: jsoniter.Get([]byte(arrData), 24, 0, i, 1, 3).ToString(),
				GpsCoords: &GpsCoordinates{
					Latitude:  jsoniter.Get([]byte(arrData), 24, 0, i, 1, 0, 2).ToFloat64(),
					Longitude: jsoniter.Get([]byte(arrData), 24, 0, i, 1, 0, 3).ToFloat64(),
				},
				Type:   jsoniter.Get([]byte(arrData), 24, 0, i, 1, 19).ToString(),
				DataID: jsoniter.Get([]byte(arrData), 24, 0, i, 1, 30).ToString(),
			},
			Date:     jsoniter.Get([]byte(arrData), 24, 0, i, 6, 1, 6).ToString(),
			Snippet:  jsoniter.Get([]byte(arrData), 24, 0, i, 6, 2, 15, 0, 0).ToString(),
			ReviewID: jsoniter.Get([]byte(arrData), 24, 0, i, 6, 0).ToString(),
			Rating:   jsoniter.Get([]byte(arrData), 24, 0, i, 6, 2, 0, 0).ToFloat64(),
			Likes:    jsoniter.Get([]byte(arrData), 24, 0, i, 6, 4, 1).ToInt(),
			Link:     jsoniter.Get([]byte(arrData), 24, 0, i, 6, 1, 4, 2, 0).ToString(),
		}

		var detailsArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 24, 0, i, 6, 2, 6).ToString()), &detailsArr)
		details := make(map[string]interface{})
		for i2 := range len(detailsArr) {
			key := jsoniter.Get([]byte(arrData), 24, 0, i, 6, 2, 6, i2, 5).ToString()
			value := jsoniter.Get([]byte(arrData), 24, 0, i, 6, 2, 6, i2, 2, 0, 0, 1).ToString()
			details[key] = value
		}
		if len(detailsArr) > 0 {
			review.Details = details
		}

		var imagesArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 24, 0, i, 2).ToString()), &imagesArr)
		var imagesList []ReviewImage
		for i2 := range len(imagesArr) {
			image := ReviewImage{
				Title:     review.PlaceInfo.Title,
				Thumbnail: jsoniter.Get([]byte(arrData), 24, 0, i, 2, i2, 6, 0).ToString(),
				Date: jsoniter.Get([]byte(arrData), 24, 0, i, 2, i2, 21, 6, 8, 0).ToString() +
					"-" + jsoniter.Get([]byte(arrData), 24, 0, i, 2, i2, 21, 6, 8, 1).ToString() +
					"-" + jsoniter.Get([]byte(arrData), 24, 0, i, 2, i2, 21, 6, 8, 2).ToString(),
				Video: jsoniter.Get([]byte(arrData), 24, 0, i, 2, i2, 21, 2, 10, 1, 0, 3).ToString(),
			}
			if !image.IsEmpty() {
				imagesList = append(imagesList, image)
			}
		}
		if len(imagesArr) > 0 {
			review.Images = imagesList
		}

		reviewResponse := ReviewResponse{
			Date:    jsoniter.Get([]byte(arrData), 24, 0, i, 6, 3, 3).ToString(),
			Snippet: jsoniter.Get([]byte(arrData), 24, 0, i, 6, 3, 14, 0, 0).ToString(),
		}
		if !reviewResponse.IsEmpty() {
			review.Response = &reviewResponse
		}

		if !review.IsEmpty() {
			reviewList = append(reviewList, review)
		}
	}
	if len(reviewList) > 0 {
		response.Reviews = reviewList
	}
	response.NextPageToken = jsoniter.Get([]byte(arrData), 24, 3).ToString()
	return response, nil
}
func DoMapsReviews(ctx context.Context, params *RequestParam) (*Response, error) {
	urlStr := fmt.Sprintf("https://www.google.com/maps/place/data=!4m7!3m6!1s%s!5m2!4m1!1i2!9m1!1b1?hl=%s&gl=%s", params.DataId, params.Hl, params.Gl)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("NewRequest err:%v", err)
	}
	req.Header = http.Header{
		"accept":                            {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-language":                   {"en-US,en;q=0.9"},
		"cache-control":                     {"no-cache"},
		"pragma":                            {"no-cache"},
		"priority":                          {"u=0, i"},
		"sec-ch-ua":                         {`"Not A(Brand";v="8", "Chromium";v="132", "Google Chrome";v="132"`},
		"sec-ch-ua-arch":                    {`"x86"`},
		"sec-ch-ua-bitness":                 {`"64"`},
		"sec-ch-ua-form-factors":            {`"Desktop"`},
		"sec-ch-ua-full-version":            {`"132.0.6834.83"`},
		"sec-ch-ua-full-version-list":       {`"Not A(Brand";v="8.0.0.0", "Chromium";v="132.0.6834.83", "Google Chrome";v="132.0.6834.83"`},
		"sec-ch-ua-mobile":                  {"?0"},
		"sec-ch-ua-platform":                {`"Windows"`},
		"sec-ch-ua-platform-version":        {`"10.0.0"`},
		"sec-ch-ua-wow64":                   {"?0"},
		"sec-fetch-dest":                    {"document"},
		"sec-fetch-mode":                    {"navigate"},
		"sec-fetch-site":                    {"none"},
		"sec-fetch-user":                    {"?1"},
		"service-worker-navigation-preload": {"true"},
		"upgrade-insecure-requests":         {"1"},
		"user-agent":                        {"Mozilla/5.0 (Windows NT 10.0; Win64; x64}, AppleWebKit/537.36 (KHTML, like Gecko}, Chrome/132.0.0.0 Safari/537.36"},
	}

	resp, respError := c.Do(req)
	if respError != nil {
		return nil, fmt.Errorf("do request err:%v", respError)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("do request err:%v", resp.Status)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body err:%v", err)
	}
	response := &Response{}
	pattern := regexp.MustCompile(`window\.APP_INITIALIZATION_STATE\s*=\s*(.*?);\s*window\.APP_FLAGS`)
	match := pattern.FindStringSubmatch(string(bodyText))
	var dataStr string
	if len(match) > 1 {
		dataStr = strings.TrimSpace(match[1])
	}
	arrData := strings.ReplaceAll(jsoniter.Get([]byte(dataStr), 3, 6).ToString(), ")]}'\n", "")

	paramFunc := func(dataStr string) (string, error) {
		var numStr string
		switch params.SortBy {
		case "newestFirst":
			numStr = "2"
		case "ratingHigh":
			numStr = "3"
		case "ratingLow":
			numStr = "4"
		default:
			numStr = "1"
		}
		arrData2 := strings.ReplaceAll(jsoniter.Get([]byte(dataStr), 3, 5).ToString(), ")]}'\n", "")
		tokenParam := jsoniter.Get([]byte(arrData2), 13, 0).ToString()

		var paramTopic string
		if params.TopicId != "" {
			paramTopic = "!5m2!1m1!1z" + base64.StdEncoding.EncodeToString([]byte(params.TopicId))
		}
		if params.Num == "" {
			params.Num = "10"
		}
		pb := fmt.Sprintf("!1m9!1s%s%s!6m4!4m1!1e1!4m1!1e3!2m2!1i%s!2s%s!5m2!1s%s!7e81!8m9!2b1!3b1!5b1!7b1!12m4!1b1!2b1!4m1!1e1!11m0!13m1!1e%s", params.DataId, paramTopic, params.Num, params.NextPageToken, tokenParam, numStr)

		sortByUrlQuery := url.Values{
			"authuser": {"0"},
			"hl":       {params.Hl},
			"gl":       {params.Gl},
			"pb":       {pb},
		}
		paramUrl := fmt.Sprintf("https://www.%s/maps/rpc/listugcposts?", params.GoogleDomain) + sortByUrlQuery.Encode()

		paramReq, _ := http.NewRequest("GET", paramUrl, nil)
		paramReq.Header = http.Header{
			"accept":                       {"*/*"},
			"accept-language":              {"en-US,en;q=0.9"},
			"cache-control":                {"no-cache"},
			"downlink":                     {"10"},
			"pragma":                       {"no-cache"},
			"priority":                     {"u=1, i"},
			"referer":                      {"https://www.google.com/"},
			"rtt":                          {"250"},
			"sec-ch-ua":                    {`"Not A(Brand";v="8", "Chromium";v="132", "Google Chrome";v="132"`},
			"sec-ch-ua-mobile":             {"?0"},
			"sec-ch-ua-platform":           {`"Windows"`},
			"sec-fetch-dest":               {"empty"},
			"sec-fetch-mode":               {"cors"},
			"sec-fetch-site":               {"same-origin"},
			"user-agent":                   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64}, AppleWebKit/537.36 (KHTML, like Gecko}, Chrome/132.0.0.0 Safari/537.36"},
			"x-maps-diversion-context-bin": {"CAE="},
		}

		paramResp, err2 := c.Do(paramReq)
		if err2 != nil {
			return "", fmt.Errorf("request err:%v", err2)
		}
		defer paramResp.Body.Close()
		bodyText2, err2 := io.ReadAll(paramResp.Body)
		if err2 != nil {
			return "", fmt.Errorf("read response body err:%v", err2)
		}
		return strings.ReplaceAll(string(bodyText2), ")]}'\n", ""), nil
	}

	placeInfo := &PlaceInfo{
		Title:   jsoniter.Get([]byte(arrData), 6, 11).ToString(),
		Address: jsoniter.Get([]byte(arrData), 6, 39).ToString(),
		Ratings: jsoniter.Get([]byte(arrData), 6, 4, 7).ToFloat64(),
		Reviews: jsoniter.Get([]byte(arrData), 6, 4, 8).ToInt(),
		Type:    jsoniter.Get([]byte(arrData), 6, 13, 0).ToString(),
	}
	if !placeInfo.IsEmpty() {
		response.PlaceInfo = placeInfo
	}

	var topicsList []Topics
	var topicsArr []any
	_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 6, 153, 0).ToString()), &topicsArr)
	for i := range len(topicsArr) {
		topics := Topics{
			Keyword:  jsoniter.Get([]byte(arrData), 6, 153, 0, i, 1).ToString(),
			Mentions: jsoniter.Get([]byte(arrData), 6, 153, 0, i, 3, 4).ToInt(),
			Id:       jsoniter.Get([]byte(arrData), 6, 153, 0, i, 0, 0).ToString(),
		}
		if !topics.IsEmpty() {
			topicsList = append(topicsList, topics)
		}
	}
	if len(topicsList) > 0 {
		response.Topics = topicsList
	}

	var reviewsArr []any
	var reviewsStr string

	if params.NextPageToken == "" && params.SortBy == "" && params.TopicId == "" {
		reviewsStr = jsoniter.Get([]byte(arrData), 6, 175, 9, 0, 0).ToString()
		_ = json.Unmarshal([]byte(reviewsStr), &reviewsArr)
	} else {
		resultDataStr, err2 := paramFunc(dataStr)
		if err2 != nil {
			return response, fmt.Errorf("request err:%v", err2)
		}
		reviewsStr = jsoniter.Get([]byte(resultDataStr), 2).ToString()
		_ = json.Unmarshal([]byte(reviewsStr), &reviewsArr)
	}

	var reviewsList []Reviews
	for i := range len(reviewsArr) {
		reviews := Reviews{
			Link:   jsoniter.Get([]byte(reviewsStr), i, 0, 4, 3, 0).ToString(),
			Rating: jsoniter.Get([]byte(reviewsStr), i, 0, 2, 0, 0).ToFloat64(),
			Date:   jsoniter.Get([]byte(reviewsStr), i, 0, 1, 6).ToString(),
			IsoDate: jsoniter.Get([]byte(reviewsStr), i, 0, 2, 2, 0, 1, 21, 6, 8, 0).ToString() +
				"-" + jsoniter.Get([]byte(reviewsStr), i, 0, 2, 2, 0, 1, 21, 6, 8, 1).ToString() +
				"-" + jsoniter.Get([]byte(reviewsStr), i, 0, 2, 2, 0, 1, 21, 6, 8, 2).ToString(),
			Source:   jsoniter.Get([]byte(reviewsStr), i, 0, 1, 13, 0).ToString(),
			ReviewID: jsoniter.Get([]byte(reviewsStr), i, 0, 0).ToString(),
			User: &User{
				Name:          jsoniter.Get([]byte(reviewsStr), i, 0, 1, 4, 5, 0).ToString(),
				Link:          jsoniter.Get([]byte(reviewsStr), i, 0, 1, 4, 5, 2, 0).ToString(),
				ContributorID: jsoniter.Get([]byte(reviewsStr), i, 0, 1, 4, 5, 3).ToString(),
				Thumbnail:     jsoniter.Get([]byte(reviewsStr), i, 0, 1, 4, 5, 1).ToString(),
				Reviews:       jsoniter.Get([]byte(reviewsStr), i, 0, 1, 4, 5, 5).ToInt(),
				Photos:        jsoniter.Get([]byte(reviewsStr), i, 0, 1, 4, 5, 6).ToInt(),
			},
			Snippet: jsoniter.Get([]byte(reviewsStr), i, 0, 2, 15, 0, 0).ToString(),
			Likes:   jsoniter.Get([]byte(reviewsStr), i, 0, 1, 4, 5, 9).ToInt(),
		}
		if reviews.IsoDate == "--" {
			reviews.IsoDate = ""
		}
		reviews.IsoDateOfLastEdit = reviews.IsoDate
		var imagesList []ReviewImage
		var imagesArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(reviewsStr), i, 0, 2, 2).ToString()), &imagesArr)
		for i2 := range len(imagesArr) {
			image := ReviewImage{
				Thumbnail: jsoniter.Get([]byte(reviewsStr), i, 0, 2, 2, i2, 1, 6, 0).ToString(),
			}
			if !image.IsEmpty() {
				imagesList = append(imagesList, image)
			}
		}
		if len(imagesList) > 0 {
			reviews.Images = imagesList
		}

		details := make(map[string]interface{})
		var detailsArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(reviewsStr), i, 0, 2, 6).ToString()), detailsArr)
		for i2 := range len(detailsArr) {
			key := jsoniter.Get([]byte(reviewsStr), i, 0, 2, 6, i2, 5).ToString()
			value := jsoniter.Get([]byte(reviewsStr), i, 0, 2, 6, i2, 2, 0, 0, 1).ToString()
			details[key] = value
		}
		if details != nil && len(details) > 0 {
			reviews.Details = details
		}

		if !reviews.IsEmpty() {
			reviewsList = append(reviewsList, reviews)
		}

		if i == len(reviewsArr)-1 {
			response.NextPageToken = jsoniter.Get([]byte(reviewsStr), i, 2).ToString()
		}
	}

	if len(reviewsList) > 0 {
		response.Reviews = reviewsList
	}

	return response, nil
}

func DoMapsDirections(ctx context.Context, params *RequestParam) (*Response, error) {
	startParam := &RequestParam{
		Q:      params.StartAddr,
		Engine: GoogleMapsDirections,
		Type:   "search",
	}
	startResponse, err2 := doMapsManage(ctx, startParam)
	if startResponse == nil || err2 != nil {
		return nil, fmt.Errorf("request err:%v", err2)
	}
	var startDataId string
	var startLatitude string
	var startLongitude string
	// var startTitle string
	if len(startResponse.LocalResults) > 0 {
		startDataId = startResponse.LocalResults[0].DataID
		if !startResponse.LocalResults[0].GpsCoordinates.IsEmpty() {
			startLatitude = fmt.Sprintf("%0.15f", startResponse.LocalResults[0].GpsCoordinates.Latitude)
			startLongitude = fmt.Sprintf("%0.7f", startResponse.LocalResults[0].GpsCoordinates.Longitude)
		}
		// startTitle = startResponse.LocalResults[0].Title
	}

	endParam := &RequestParam{
		Q:      params.EndAddr,
		Engine: GoogleMapsDirections,
		Type:   "search",
	}
	endResponse, err2 := doMapsManage(ctx, endParam)
	if err2 != nil {
		return nil, fmt.Errorf("request err:%v", err2)
	}
	var endDataId string
	var endLatitude string
	var endLongitude string
	// var endTitle string
	if len(endResponse.LocalResults) > 0 {
		endDataId = endResponse.LocalResults[0].DataID
		if !endResponse.LocalResults[0].GpsCoordinates.IsEmpty() {
			endLatitude = fmt.Sprintf("%0.15f", endResponse.LocalResults[0].GpsCoordinates.Latitude)
			endLongitude = fmt.Sprintf("%0.7f", endResponse.LocalResults[0].GpsCoordinates.Longitude)
		}
	}

	var travelMode string
	if params.TravelMode != "" {
		travelMode = fmt.Sprintf("!3e%s", params.TravelMode)
	}
	var route string
	if params.Route != "" {
		route = fmt.Sprintf("!4e%s", params.Route)
	}
	var prefer string
	var preferLen int
	if params.Prefer != "" {
		var preferArr []string
		for _, value := range strings.Split(params.Prefer, ",") {
			switch value {
			case "bus":
				preferArr = append(preferArr, "!5e0")
			case "subway":
				preferArr = append(preferArr, "!5e1")
			case "train":
				preferArr = append(preferArr, "!5e2")
			case "tram_light_rail":
				preferArr = append(preferArr, "!5e3")
			}
		}
		prefer = strings.Join(preferArr, "")
		preferLen = len(preferArr)
	}
	var avoid string
	var avoidLen int
	if params.Avoid != "" {
		var avoidArr []string
		for _, value := range strings.Split(params.Avoid, ",") {
			switch value {
			case "highways":
				avoidArr = append(avoidArr, "!1b1")
			case "tolls":
				avoidArr = append(avoidArr, "!2b1")
			case "ferries":
				avoidArr = append(avoidArr, "!3b1")
			default:
				avoidArr = append(avoidArr, "")
			}
		}
		avoid = strings.Join(avoidArr, "")
		avoidLen = len(avoidArr)
	}
	var paramTime string
	if params.Time != "" {
		timeArr := strings.Split(params.Time, ":")
		if len(timeArr) > 1 {
			switch timeArr[0] {
			case "depart_at":
				paramTime = fmt.Sprintf("!6e0!7e2!8j%s", timeArr[1])
			case "arrive_by":
				paramTime = fmt.Sprintf("!6e1!7e2!8j%s", timeArr[1])
			}
		}
		if params.Time == "last_available" {
			paramTime = "!6e2"
		}
	}

	var m4Param string
	var m2Param string
	if params.Route == "" && params.Prefer == "" && params.Avoid == "" && params.Time == "" {
		m4Param = "!4m14!4m13"
		m2Param = ""
	}

	if preferLen > 0 || avoidLen > 0 || len(params.Route) > 0 || len(params.Time) > 0 {
		paramNum := preferLen + avoidLen + len(params.Route)
		if len(params.Time) > 0 {
			if params.Time == "last_available" {
				paramNum = paramNum + 1
			} else {
				paramNum = paramNum + 3
			}
		}
		var m2ParamNum = paramNum
		if len(params.Time) > 0 {
			if params.Time != "last_available" {
				m2ParamNum = paramNum - 2
			}
		}
		m2Param = fmt.Sprintf("!2m%s", strconv.Itoa(m2ParamNum))
		m4Param = fmt.Sprintf("!4m%s!4m%s", strconv.Itoa(14+paramNum+1), strconv.Itoa(14+paramNum))
	}

	urlStr := fmt.Sprintf("https://www.%s/maps/dir/%s/%s/data=!3m1!4b1%s!1m5!1m1!1s%s!2m2!1d%s!2d%s!1m5!1m1!1s%s!2m2!1d%s!2d%s%s%s%s%s%s%s?hl=%s&gl=%s",
		params.GoogleDomain, params.StartAddr, params.EndAddr, m4Param, startDataId, startLongitude, startLatitude, endDataId, endLongitude, endLatitude, m2Param, avoid, route, prefer, paramTime, travelMode, params.Hl, params.Gl)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest err: %v", err)
	}
	req.Header = http.Header{
		"accept":                            {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-language":                   {"en-US,en;q=0.9"},
		"cache-control":                     {"no-cache"},
		"pragma":                            {"no-cache"},
		"priority":                          {"u=0, i"},
		"rtt":                               {"50"},
		"sec-ch-ua":                         {`"Not A(Brand";v="8", "Chromium";v="132", "Google Chrome";v="132"`},
		"sec-ch-ua-mobile":                  {"?0"},
		"sec-ch-ua-platform":                {`"Windows"`},
		"sec-fetch-dest":                    {"document"},
		"sec-fetch-mode":                    {"navigate"},
		"sec-fetch-site":                    {"none"},
		"sec-fetch-user":                    {"?1"},
		"service-worker-navigation-preload": {"true"},
		"upgrade-insecure-requests":         {"1"},
		"user-agent":                        {"Mozilla/5.0 (Windows NT 10.0; Win64; x64}, AppleWebKit/537.36 (KHTML, like Gecko}, Chrome/132.0.0.0 Safari/537.36"},
		"x-browser-channel":                 {"stable"},
		"x-browser-copyright":               {"Copyright 2025 Google LLC. All rights reserved."},
	}
	resp, respError := c.Do(req)
	if respError != nil {
		return nil, fmt.Errorf("c.Do err: %v", respError)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("resp.StatusCode err: %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("resp.Body.ReadAll err: %v", err)
	}
	var response = &Response{}
	pattern := regexp.MustCompile(`window\.APP_INITIALIZATION_STATE\s*=\s*(.*?);\s*window\.APP_FLAGS`)
	match := pattern.FindStringSubmatch(string(bodyText))
	var dataStr string
	if len(match) > 1 {
		dataStr = strings.TrimSpace(match[1])
	}

	arrData := strings.ReplaceAll(jsoniter.Get([]byte(dataStr), 3, 4).ToString(), ")]}'\n", "")

	var placeInfosList []PlaceInfo
	var placeInfosArr []any
	_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 0, 0).ToString()), &placeInfosArr)
	for i := range len(placeInfosArr) {
		placeInfo := PlaceInfo{
			Address: jsoniter.Get([]byte(arrData), 0, 0, i, 0, 0, 0).ToString(),
			DataID:  jsoniter.Get([]byte(arrData), 0, 0, i, 0, 0, 1).ToString(),
			GpsCoords: &GpsCoordinates{
				Latitude:  jsoniter.Get([]byte(arrData), 0, 0, i, 0, 0, 2, 2).ToFloat64(),
				Longitude: jsoniter.Get([]byte(arrData), 0, 0, i, 0, 0, 2, 3).ToFloat64(),
			},
		}
		if !placeInfo.IsEmpty() {
			placeInfosList = append(placeInfosList, placeInfo)
		}
	}
	if len(placeInfosList) > 0 {
		response.PlaceInfos = placeInfosList
	}

	var directionsList []Directions
	var directionsArr []any
	_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 0, 1).ToString()), &directionsArr)

	for i := range len(directionsArr) {
		direction := Directions{
			TravelMode:           travelModeMapping(jsoniter.Get([]byte(arrData), 0, 1, i, 0, 0).ToInt()),
			Via:                  jsoniter.Get([]byte(arrData), 0, 1, i, 0, 1).ToString(),
			StartTime:            jsoniter.Get([]byte(arrData), 0, 1, i, 0, 5, 0, 2).ToString(),
			EndTime:              jsoniter.Get([]byte(arrData), 0, 1, i, 0, 5, 1, 2).ToString(),
			Distance:             jsoniter.Get([]byte(arrData), 0, 1, i, 0, 2, 0).ToInt(),
			Duration:             jsoniter.Get([]byte(arrData), 0, 1, i, 0, 3, 0).ToInt(),
			TypicalDurationRange: jsoniter.Get([]byte(arrData), 0, 1, i, 0, 10, 4, 2).ToString(),
			FormattedDistance:    jsoniter.Get([]byte(arrData), 0, 1, i, 0, 2, 1).ToString(),
			FormattedDuration:    jsoniter.Get([]byte(arrData), 0, 1, i, 0, 10, 0, 1).ToString(),
			Cost:                 jsoniter.Get([]byte(arrData), 0, 1, i, 0, 11, 0).ToInt(),
			Currency:             jsoniter.Get([]byte(arrData), 0, 1, i, 0, 11, 2).ToString(),
			ArriveAround:         jsoniter.Get([]byte(arrData), 0, 1, i, 0, 10, 4, 2).ToInt(),
		}

		if params.DistanceUnit == "0" {
			kmNumStr := strings.TrimSpace(strings.ReplaceAll(direction.FormattedDistance, "miles", ""))
			kmF, _ := strconv.ParseFloat(kmNumStr, 64)
			direction.FormattedDistance = fmt.Sprintf("%.1f", math.Round(kmF*1.60934*10)/10) + "KM"
		}

		var extensions []string
		var extensionsArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 0, 1, i, 0, 9).ToString()), &extensionsArr)
		for i2 := range len(extensionsArr) {
			extensions = append(extensions, jsoniter.Get([]byte(arrData), 0, 1, i, 0, 9, i2, 1).ToString())
		}
		if len(extensions) > 0 {
			direction.Extensions = extensions
		}

		elevationProfile := ElevationProfile{
			Ascent:               jsoniter.Get([]byte(arrData), 0, 1, i, 17, 4, 0).ToInt(),
			Descent:              jsoniter.Get([]byte(arrData), 0, 1, i, 17, 5, 0).ToInt(),
			MaxAltitude:          jsoniter.Get([]byte(arrData), 0, 1, i, 17, 1, 0).ToInt(),
			MinAltitude:          jsoniter.Get([]byte(arrData), 0, 1, i, 17, 2, 0).ToInt(),
			FormattedAscent:      jsoniter.Get([]byte(arrData), 0, 1, i, 17, 4, 1).ToString(),
			FormattedDescent:     jsoniter.Get([]byte(arrData), 0, 1, i, 17, 5, 1).ToString(),
			FormattedMaxAltitude: jsoniter.Get([]byte(arrData), 0, 1, i, 17, 1, 1).ToString(),
			FormattedMinAltitude: jsoniter.Get([]byte(arrData), 0, 1, i, 17, 2, 1).ToString(),
		}
		if !elevationProfile.IsEmpty() {
			direction.ElevationProfile = &elevationProfile
		}

		var flightLink string
		if jsoniter.Get([]byte(arrData), 0, 1, i, 0, 13, 15, 0).ToString() != "" {
			flightLink = fmt.Sprintf("https://www.%s%s", params.GoogleDomain, jsoniter.Get([]byte(arrData), 0, 1, i, 0, 13, 15, 0).ToString())
		}
		flight := FlightDetails{
			Departure:                   jsoniter.Get([]byte(arrData), 0, 1, i, 0, 13, 0).ToString(),
			Arrival:                     jsoniter.Get([]byte(arrData), 0, 1, i, 0, 13, 1).ToString(),
			Date:                        jsoniter.Get([]byte(arrData), 0, 1, i, 0, 13, 2, 0).ToString(),
			RoundTripPrice:              0, // 目前这俩字段，搜索场景没有 无法判断在数组的哪个位置
			Currency:                    "",
			NonstopDuration:             jsoniter.Get([]byte(arrData), 0, 1, i, 0, 13, 4, 0).ToString(),
			FormattedNonstopDuration:    jsoniter.Get([]byte(arrData), 0, 1, i, 0, 13, 4, 1).ToString(),
			ConnectingDuration:          jsoniter.Get([]byte(arrData), 0, 1, i, 0, 13, 6, 0).ToString(),
			FormattedConnectingDuration: jsoniter.Get([]byte(arrData), 0, 1, i, 0, 13, 6, 1).ToString(),
			GoogleFlightsLink:           flightLink,
		}

		var airlines []string
		var airlinesArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 0, 1, i, 0, 13, 10).ToString()), &airlinesArr)
		for i2 := range len(airlinesArr) {
			airlines = append(airlines, jsoniter.Get([]byte(arrData), 0, 1, i, 0, 13, 10, i2).ToString())
		}
		if len(airlines) > 0 {
			flight.Airlines = airlines
		}
		if !flight.IsEmpty() {
			direction.Flight = &flight
		}

		var tripList []Trip
		var tripArr []any
		_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 0, 1, i, 1, 0, 1).ToString()), &tripArr)
		for i2 := range len(tripArr) {
			one := jsoniter.Get([]byte(arrData), 0, 1, i, 1, 0, 1).ToString()
			trip := Trip{
				TravelMode:        travelModeMapping(jsoniter.Get([]byte(one), i2, 0, 0).ToInt()),
				Distance:          jsoniter.Get([]byte(one), i2, 0, 2, 0).ToInt(),
				Duration:          jsoniter.Get([]byte(one), i2, 0, 3, 0).ToInt(),
				FormattedDistance: jsoniter.Get([]byte(one), i2, 0, 2, 1).ToString(),
				FormattedDuration: jsoniter.Get([]byte(one), i2, 0, 3, 1).ToString(),
			}
			var titleArr []any
			var titles []string
			_ = json.Unmarshal([]byte(jsoniter.Get([]byte(one), i2, 0, 14).ToString()), &titleArr)
			for i3 := range len(titleArr) {
				titles = append(titles, jsoniter.Get([]byte(one), i2, 0, 14, i3, 1, 0).ToString())
			}
			if len(titles) > 0 {
				trip.Title = strings.Join(titles, " ")
			}

			var detailsList []DirectionDetail
			var detailsArr []any
			_ = json.Unmarshal([]byte(jsoniter.Get([]byte(one), i2, 1).ToString()), &detailsArr)
			for i3 := range len(detailsArr) {
				detail := DirectionDetail{
					Action:            jsoniter.Get([]byte(one), i2, 1, i3, 2, 1).ToString(),
					Distance:          jsoniter.Get([]byte(one), i2, 1, i3, 0, 2, 0).ToInt(),
					Duration:          jsoniter.Get([]byte(one), i2, 1, i3, 0, 3, 0).ToInt(),
					FormattedDistance: jsoniter.Get([]byte(one), i2, 1, i3, 0, 2, 1).ToString(),
					FormattedDuration: jsoniter.Get([]byte(one), i2, 1, i3, 0, 3, 1).ToString(),
					GeoPhoto:          jsoniter.Get([]byte(one), i2, 1, i3, 0, 7, 5, 6, 0).ToString(),
					GPSCoordinates: &GpsCoordinates{
						Latitude:  jsoniter.Get([]byte(one), i2, 1, i3, 0, 7, 1, 0, 2).ToFloat64(),
						Longitude: jsoniter.Get([]byte(one), i2, 1, i3, 0, 7, 1, 0, 3).ToFloat64(),
					},
				}
				var detailTitleArr []any
				_ = json.Unmarshal([]byte(jsoniter.Get([]byte(one), i2, 1, i3, 0, 14).ToString()), &detailTitleArr)
				var detailDetailTitles []string
				for i4 := range len(detailTitleArr) {
					detailDetailTitles = append(detailDetailTitles, jsoniter.Get([]byte(one), i2, 1, i3, 0, 14, i4, 1, 0).ToString())
				}
				if len(detailTitleArr) > 0 {
					detail.Title = strings.Join(detailDetailTitles, " ")
				}

				var detailExtensions []string
				var detailExtensionsArr []any
				_ = json.Unmarshal([]byte(jsoniter.Get([]byte(one), i2, 1, i3, 0, 9).ToString()), &detailExtensionsArr)
				for i4 := range len(detailExtensionsArr) {
					detailExtensions = append(detailExtensions, jsoniter.Get([]byte(one), i2, 1, i3, 0, 9, i4, 1).ToString())
				}
				if len(detailExtensions) > 0 {
					detail.Extensions = detailExtensions
				}

				if !detail.IsEmpty() {
					detailsList = append(detailsList, detail)
				}

			}
			if len(detailsList) > 0 {
				trip.Details = detailsList
			}

			startStop := Stops{
				Name:   jsoniter.Get([]byte(one), i2, 5, 0, 0).ToString(),
				StopId: jsoniter.Get([]byte(one), i2, 5, 0, 1).ToString(),
				Time:   jsoniter.Get([]byte(one), i2, 5, 0, 3, 2).ToString(),
				DataId: jsoniter.Get([]byte(one), i2, 5, 0, 6).ToString(),
			}
			if !startStop.IsEmpty() {
				trip.StartStop = &startStop
			}

			endStop := Stops{
				Name:   jsoniter.Get([]byte(one), i2, 5, 1, 0).ToString(),
				StopId: jsoniter.Get([]byte(one), i2, 5, 1, 1).ToString(),
				Time:   jsoniter.Get([]byte(one), i2, 5, 1, 2, 2).ToString(),
				DataId: jsoniter.Get([]byte(one), i2, 5, 1, 6).ToString(),
			}
			if !endStop.IsEmpty() {
				trip.EndStop = &endStop
			}

			var stopsList []Stops
			var stopsArr []any
			_ = json.Unmarshal([]byte(jsoniter.Get([]byte(one), i2, 5, 7).ToString()), &stopsArr)
			for i3 := range len(stopsArr) {
				stop := Stops{
					Name:   jsoniter.Get([]byte(one), i2, 5, 7, i3, 0).ToString(),
					StopId: jsoniter.Get([]byte(one), i2, 5, 7, i3, 1).ToString(),
					Time:   jsoniter.Get([]byte(one), i2, 5, 7, i3, 2, 2).ToString(),
					DataId: jsoniter.Get([]byte(one), i2, 5, 7, i3, 6).ToString(),
				}
				if !stop.IsEmpty() {
					stopsList = append(stopsList, stop)
				}
			}
			if len(stopsList) > 0 {
				trip.Stops = stopsList
			}

			if !trip.IsEmpty() {
				tripList = append(tripList, trip)
			}
		}
		if len(tripList) > 0 {
			direction.Trips = tripList
		}

		if !direction.IsEmpty() {
			directionsList = append(directionsList, direction)
		}
	}

	if len(directionsList) > 0 {
		response.Directions = directionsList
	}

	return response, nil
}

func DoMapsPhotos(ctx context.Context, params *RequestParam) (*Response, error) {
	var urlStr string
	if params.NextPageToken != "" {
		urlStr = fmt.Sprintf("https://www.%s/maps/rpc/photo/listentityphotos?hl=%s&pb=!1e2!3m3!1s%s!9e0!11s!5m59!2m2!1i203!2i100!3m3!2i20!3s%s!5b1!7m50!1m3!1e1!2b0!3e3!1m3!1e2!2b1!3e2!1m3!1e2!2b0!3e3!1m3!1e3!2b0!3e3!1m3!1e8!2b0!3e3!1m3!1e3!2b1!3e2!1m3!1e10!2b0!3e3!1m3!1e10!2b1!3e2!1m3!1e9!2b1!3e2!1m3!1e10!2b0!3e3!1m3!1e10!2b1!3e2!1m3!1e10!2b0!3e4!2b1!4b1!9b0!6m3!1s!2z!7e81!16m4!1m1!1BCgIgAQ!2b1!4e1", params.GoogleDomain, params.Hl, params.DataId, params.NextPageToken)
	} else {
		urlStr = fmt.Sprintf("https://www.%s/maps/rpc/photo/listentityphotos?hl=%s&pb=!1e2!3m3!1s%s!9e0!11s!5m58!2m2!1i203!2i100!3m2!2i20!5b1!7m50!1m3!1e1!2b0!3e3!1m3!1e2!2b1!3e2!1m3!1e2!2b0!3e3!1m3!1e3!2b0!3e3!1m3!1e8!2b0!3e3!1m3!1e3!2b1!3e2!1m3!1e10!2b0!3e3!1m3!1e10!2b1!3e2!1m3!1e9!2b1!3e2!1m3!1e10!2b0!3e3!1m3!1e10!2b1!3e2!1m3!1e10!2b0!3e4!2b1!4b1!9b0!6m3!1s!2z!7e81!16m4!1m1!1BCgIgAQ!2b1!4e1", params.GoogleDomain, params.Hl, params.DataId)
	}
	if params.CategoryId != "" {
		urlStr = regexp.MustCompile(`!1B(.*?)!2b1`).ReplaceAllString(urlStr, "!1B"+params.CategoryId+"!2b1")
	}
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header = http.Header{
		"accept":                            {"text/html,application/xhtml+xml,application/xml;q0.9,image/avif,image/webp,image/apng,*/*;q0.8,application/signed-exchange;vb3;q0.7"},
		"accept-language":                   {"en-US,en;q=0.9"},
		"cache-control":                     {"no-cache"},
		"downlink":                          {"10"},
		"pragma":                            {"no-cache"},
		"priority":                          {"u0, i"},
		"rtt":                               {"250"},
		"sec-ch-ua":                         {`"Not A(Brand";v"8", "Chromium";v"132", "Google Chrome";v"132"`},
		"sec-ch-ua-mobile":                  {"?0"},
		"sec-ch-ua-platform":                {`"Windows"`},
		"sec-fetch-dest":                    {"document"},
		"sec-fetch-mode":                    {"navigate"},
		"sec-fetch-site":                    {"none"},
		"sec-fetch-user":                    {"?1"},
		"service-worker-navigation-preload": {"true"},
		"upgrade-insecure-requests":         {"1"},
		"user-agent":                        {"Mozilla/5.0 (Windows NT 10.0; Win64; x64}, AppleWebKit/537.36 (KHTML, like Gecko}, Chrome/132.0.0.0 Safari/537.36"},
		"x-browser-channel":                 {"stable"},
	}

	resp, respError := c.Do(req)
	if respError != nil {
		return nil, fmt.Errorf("request error: %s", respError.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request status code: %s", resp.Status)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %s", err.Error())
	}
	response := &Response{}
	arrData := strings.ReplaceAll(string(bodyText), ")]}'\n", "")

	var categoriesList []Categories
	var categoryDataArr []any
	_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 12, 0).ToString()), &categoryDataArr)

	for i := range len(categoryDataArr) {
		category := Categories{
			Title: jsoniter.Get([]byte(arrData), 12, 0, i, 2).ToString(),
			Id:    strings.ReplaceAll(jsoniter.Get([]byte(arrData), 12, 0, i, 0).ToString(), "=", ""),
		}
		if !category.IsEmpty() {
			categoriesList = append(categoriesList, category)
		}
	}
	if len(categoriesList) > 0 {
		response.Categories = categoriesList
	}

	var photosList []Photos
	var photosArr []any
	_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 0).ToString()), &photosArr)
	for i := range len(photosArr) {
		photo := Photos{
			Thumbnail:   jsoniter.Get([]byte(arrData), 0, i, 6, 0).ToString(),
			Video:       jsoniter.Get([]byte(arrData), 0, i, 26, 1, 0, 3).ToString(),
			PhotoDataId: jsoniter.Get([]byte(arrData), 0, i, 0).ToString(),
			DataId:      jsoniter.Get([]byte(arrData), 0, i, 15, 0, 0, 0).ToString(),
		}
		if !photo.IsEmpty() {
			photosList = append(photosList, photo)
		}
	}
	if len(photosList) > 0 {
		response.Photos = photosList
	}

	response.NextPageToken = jsoniter.Get([]byte(arrData), 5).ToString()

	return response, nil
}

func DoMapsPhotoMeta(ctx context.Context, params *RequestParam) (*Response, error) {
	urlStr := fmt.Sprintf("https://www.%s/maps/photometa/v1?pb=!1m4!1smaps_sv.tactile!11m2!2m1!1b1!2m2!1sen!2smy!3m5!1m2!1e10!2s%s!2m1!5s%s!4m57!1e1!1e2!1e3!1e4!1e5!1e6!1e8!1e12!2m1!1e1!4m1!1i48!5m1!1e1!5m1!1e2!6m1!1e1!6m1!1e2!9m36!1m3!1e2!2b1!3e2!1m3!1e2!2b0!3e3!1m3!1e3!2b1!3e2!1m3!1e3!2b0!3e3!1m3!1e8!2b0!3e3!1m3!1e1!2b0!3e3!1m3!1e4!2b0!3e3!1m3!1e10!2b1!3e2!1m3!1e10!2b0!3e3", params.GoogleDomain, params.PhotoDataId, "0x31cc376b910e0ea9:0xf47bc10a17c8b2ab")
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header = http.Header{
		"accept":                            {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-language":                   {"en-US,en;q=0.9"},
		"cache-control":                     {"no-cache"},
		"downlink":                          {"10"},
		"pragma":                            {"no-cache"},
		"priority":                          {"u=0, i"},
		"rtt":                               {"250"},
		"sec-ch-prefers-color-scheme":       {"dark"},
		"sec-ch-ua":                         {`"Not A(Brand";v="8", "Chromium";v="132", "Google Chrome";v="132"`},
		"sec-ch-ua-arch":                    {`"x86"`},
		"sec-ch-ua-bitness":                 {`"64"`},
		"sec-ch-ua-form-factors":            {`"Desktop"`},
		"sec-ch-ua-full-version":            {`"132.0.6834.83"`},
		"sec-ch-ua-full-version-list":       {`"Not A(Brand";v="8.0.0.0", "Chromium";v="132.0.6834.83", "Google Chrome";v="132.0.6834.83"`},
		"sec-ch-ua-mobile":                  {"?0"},
		"sec-ch-ua-platform":                {`"Windows"`},
		"sec-ch-ua-platform-version":        {`"10.0.0"`},
		"sec-ch-ua-wow64":                   {"?0"},
		"sec-fetch-dest":                    {"document"},
		"sec-fetch-mode":                    {"navigate"},
		"sec-fetch-site":                    {"none"},
		"sec-fetch-user":                    {"?1"},
		"service-worker-navigation-preload": {"true"},
		"upgrade-insecure-requests":         {"1"},
		"user-agent":                        {"Mozilla/5.0 (Windows NT 10.0; Win64; x64}, AppleWebKit/537.36 (KHTML, like Gecko}, Chrome/132.0.0.0 Safari/537.36"},
	}

	resp, respError := c.Do(req)
	if respError != nil {
		return nil, fmt.Errorf("request failed with %s", respError.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with %s", resp.Status)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed with %s", err.Error())
	}
	response := &Response{}
	arrData := strings.ReplaceAll(string(bodyText), ")]}'\n", "")
	user := User{
		Name:   jsoniter.Get([]byte(arrData), 1, 0, 4, 1, 0, 0, 0).ToString(),
		UserId: jsoniter.Get([]byte(arrData), 1, 0, 4, 1, 0, 5).ToString(),
	}
	if jsoniter.Get([]byte(arrData), 1, 0, 4, 1, 0, 1).ToString() != "" {
		user.Link = "https:" + jsoniter.Get([]byte(arrData), 1, 0, 4, 1, 0, 1).ToString()
	}
	if jsoniter.Get([]byte(arrData), 1, 0, 4, 1, 0, 2).ToString() != "" {
		user.Image = "https:" + jsoniter.Get([]byte(arrData), 1, 0, 4, 1, 0, 2).ToString()
	}
	if !user.IsEmpty() {
		response.User = &user
	}

	location := Location{
		Latitude:  jsoniter.Get([]byte(arrData), 1, 0, 5, 0, 1, 0, 2).ToFloat64(),
		Longitude: jsoniter.Get([]byte(arrData), 1, 0, 5, 0, 1, 0, 3).ToFloat64(),
		Name:      jsoniter.Get([]byte(arrData), 1, 0, 5, 0, 9, 0, 2, 0).ToString(),
		Type:      jsoniter.Get([]byte(arrData), 1, 0, 5, 0, 9, 0, 3, 0).ToString(),
	}
	if !location.IsEmpty() {
		response.Location = &location
	}

	var dateArr []int
	_ = json.Unmarshal([]byte(jsoniter.Get([]byte(arrData), 1, 0, 6, 7).ToString()), &dateArr)
	var dateArrStr []string
	if len(dateArr) > 0 {
		for i := range len(dateArr) {
			dateArrStr = append(dateArrStr, strconv.Itoa(dateArr[i]))
		}
	}
	response.Date = strings.Join(dateArrStr, "-")
	return response, nil
}

func weekMapping(flag int) string {
	switch flag {
	case 7:
		return "Sunday"
	case 6:
		return "saturday"
	case 5:
		return "friday"
	case 4:
		return "thursday"
	case 3:
		return "wednesday"
	case 2:
		return "tuesday"
	case 1:
		return "monday"
	default:
		return ""
	}
}

func travelModeMapping(flag int) string {
	switch flag {
	case 6:
		return "Best"
	case 0:
		return "Driving"
	case 9:
		return "Two-wheeler"
	case 3:
		return "Transit"
	case 2:
		return "Walking"
	case 1:
		return "Cycling"
	case 4:
		return "Flight"
	default:
		return ""
	}
}

func extractNumbersUsingMap(s string) string {
	filter := func(r rune) rune {
		if unicode.IsDigit(r) || r == '.' {
			return r
		}
		return -1
	}
	return strings.Map(filter, s)
}

func extractLatLong(input string) (latitude, longitude string, err error) {
	trimmed := strings.TrimPrefix(input, "@")
	trimmed = strings.TrimSuffix(trimmed, "z")

	parts := strings.Split(trimmed, ",")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid input format")
	}

	latitude = parts[0]
	longitude = parts[1]

	return latitude, longitude, nil
}

type jsonData struct {
	C int    `json:"c"`
	D string `json:"d"`
	E string `json:"e"`
	P bool   `json:"p"`
	U string `json:"u"`
}

func extractBodyByHtml(bodyText []byte) string {
	pattern := regexp.MustCompile(`window\.APP_INITIALIZATION_STATE\s*=\s*(.*?);\s*window\.APP_FLAGS`)
	match := pattern.FindStringSubmatch(string(bodyText))
	var dataStr string
	if len(match) > 1 {
		dataStr = strings.TrimSpace(match[1])
	}
	return dataStr
}
