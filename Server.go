package main

import ("fmt"
		"net"
		"net/rpc"
		"net/rpc/jsonrpc"
		"net/http"
		"io/ioutil"
		"encoding/json"
		"strings"
		"strconv")

type Server struct{}
var IncrementId int 
var storageMap map[int]string
type Request struct{
	StockSymbolAndPercentage []InnerRequest `json:"stockSymbolAndPercentage"`
	Budget float32 `json:"budget"`
}

type SecondRequest struct{
	Tradeid int `json:"tradeid"`
}
type InnerRequest struct{
	Fields ActualFields `json:"fields"`
}

type ActualFields struct{
	Name string `json:"name"`
	Percentage int `json:"perecentage"`
}

type Response struct{
	Stocks []InnerResponse `json:"stocks"`
	TradeId int `json:"tradeid"`
	UnvestedAmount float32 `json:"unvestedAmount"`
}

type InnerResponse struct{
	ResponseFields ActualResponseFields `json:"fields"`
}

type ActualResponseFields struct{
	Name string `json:"name"`
	Number int `json:"number"`
	Price string `json:"price"`
}

type SecondResponse struct{
	Stocks []InnerResponse `json:"stocks"`
	CurrentMarketValue float32 `json:"currentMarketValue"`
	UnvestedAmount float32 `json:"unvestedAmount"`
}


func (this *Server) PrintMessage(msg string,reply *string) error{
		var jsonInt interface{}
		var structResponse Response
		var jsonMsg Request
		var companyName string
		var remainder float32=0.0
		json.Unmarshal([]byte(msg),&jsonMsg)
		for _, i:= range jsonMsg.StockSymbolAndPercentage{
			companyName += i.Fields.Name +","
		}
		companyName=strings.Trim(companyName,",")
		response,err:= http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+companyName+"/quote?format=json")
		if(err!=nil){
			fmt.Println(err)
		}else{
			defer response.Body.Close()
			contents,err:= ioutil.ReadAll(response.Body)
			json.Unmarshal(contents,&jsonInt)
			for i,index := range (jsonInt.(map[string]interface{})["list"]).(map[string]interface{})["resources"].([]interface{}){ 
				price := index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["price"]
				price1,_ := strconv.ParseFloat(price.(string),64)
				Remainder1:=(float64(jsonMsg.StockSymbolAndPercentage[i].Fields.Percentage) * float64(jsonMsg.Budget))/100
				name := index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["symbol"]
				number := int( Remainder1/price1)
				remainder = remainder + (float32(price1)*float32(number))
				structActualResponseFields:=ActualResponseFields{Name:name.(string),Number:number,Price:strconv.FormatFloat(price1,'f',-1,64)}
				structInnerResponse := InnerResponse{ResponseFields:structActualResponseFields}
				structResponse.Stocks = append(structResponse.Stocks,structInnerResponse)
			}
			remainder=jsonMsg.Budget-remainder
			result1 := &Response{
    		TradeId:IncrementId,
        	Stocks: structResponse.Stocks,
        	UnvestedAmount:remainder} //Map the values to Request structure
    		result2, _ := json.Marshal(result1) //Convert the Request to JSON
    		*reply = string(result2)
			storageMap[IncrementId]=string(result2)
			IncrementId++
			if(err!=nil){
				fmt.Println(err)
			}
				
		}
		
		return nil
}

func (this *Server) LossOrGain(msg string,reply *string) error{
	var jsonReq SecondRequest
	var jsonMsg Response
	var jsonInt interface{}
	var companyName string
	var price []float64
	var structSecondResponse SecondResponse
	json.Unmarshal([]byte(msg),&jsonReq)
	tradeid:= jsonReq.Tradeid
	data:= storageMap[tradeid]
	json.Unmarshal([]byte(data),&jsonMsg)
	for _,index:= range jsonMsg.Stocks{
		companyName += index.ResponseFields.Name +","
	}
	companyName=strings.Trim(companyName,",")
	response,err:= http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+companyName+"/quote?format=json")
	if(err!=nil){
		fmt.Println(err)
	}else{
		defer response.Body.Close()
		contents,_:= ioutil.ReadAll(response.Body)
		json.Unmarshal(contents,&jsonInt)
		for _,index := range (jsonInt.(map[string]interface{})["list"]).(map[string]interface{})["resources"].([]interface{}){ 
				price1,_ := strconv.ParseFloat((index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["price"]).(string),64)
				price = append(price,price1)
			}
		var value float32=0.0
		var strprice string
		for i,index := range jsonMsg.Stocks{
				temp,_:= strconv.ParseFloat(index.ResponseFields.Price,64)
				fmt.Println(price[i],temp)
				if price[i] > temp{
					strprice = "$+"+strconv.FormatFloat(price[i],'f',-1,64)

				}
				if price[i] < temp {
					strprice = "$-"+strconv.FormatFloat(price[i],'f',-1,64)
				}else {
					strprice = "$"+strconv.FormatFloat(price[i],'f',-1,64)
				}
				structActualResponseFields:=ActualResponseFields{Name:index.ResponseFields.Name,Number:index.ResponseFields.Number,Price:strprice}
				structInnerResponse := InnerResponse{ResponseFields:structActualResponseFields}
				structSecondResponse.Stocks = append(structSecondResponse.Stocks,structInnerResponse)
				value = value + (float32(index.ResponseFields.Number) * float32(price[i]))
		}
		result1 := &SecondResponse{
    	CurrentMarketValue:value,
        Stocks: structSecondResponse.Stocks,
        UnvestedAmount:jsonMsg.UnvestedAmount} //Map the values to Request structure
    	result2, _ := json.Marshal(result1) //Convert the Request to JSON
    	*reply = string(result2)
	
	}		
	return nil
}

func main(){
	IncrementId++
	storageMap=make(map[int]string)
	rpc.Register(new(Server))
	hear,err:= net.Listen("tcp",":8080")
	if(err!=nil){
		fmt.Println(err)
		return
	}
	for{
		c,error:= hear.Accept()
		if(error!=nil){
			continue
		}
		go jsonrpc.ServeConn(c)
	}

}
