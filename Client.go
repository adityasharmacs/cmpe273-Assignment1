package main

import ("fmt"
        "net/rpc/jsonrpc"
        "strings"
        "bufio"
        "os"
        "strconv"
        "encoding/json")

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

func PurchaseStocks(){

	caller,err:= jsonrpc.Dial("tcp","127.0.0.1:8080")
	if err!=nil{
		fmt.Println(err)
		return
	}
	var replyGiven string
	var structRequestGiven Request
	var msg,data,newData []string
	fmt.Println("Enter the request")
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	msg = strings.SplitN(line," ",-1)
	data = strings.SplitN(msg[0],":",2)
	newData = strings.SplitN(msg[1],":",2)
	bValue,err:=strconv.ParseFloat(strings.TrimSpace(newData[1]),64)
	data[1]= strings.Replace(data[1],"\"","",-1)
	data[1]= strings.Replace(data[1],"%","",-1)
	fields := strings.SplitN(data[1],",",-1)
	for _,index:=range fields{
			c:= strings.SplitN(index,":",-1)
			a,_:=strconv.Atoi(c[1])
			structFields := ActualFields{Name:c[0],Percentage:a} 
			structInnerRequest := InnerRequest {Fields:structFields}
			structRequestGiven.StockSymbolAndPercentage =append(structRequestGiven.StockSymbolAndPercentage,structInnerRequest)
	}
	result1 := &Request{
    	Budget:float32(bValue),
        StockSymbolAndPercentage: structRequestGiven.StockSymbolAndPercentage} //Map the values to Request structure
    result2, _ := json.Marshal(result1) //Convert the Request to JSON
	err = caller.Call("Server.PrintMessage",string(result2),&replyGiven)
	var jsonMsg Response
	var outputValue string
	outputValue = "\"tradeid\":"
	json.Unmarshal([]byte(replyGiven),&jsonMsg)
	outputValue+=strconv.Itoa(jsonMsg.TradeId)+"\n"+"\"stocks\":\""
	for _, i:= range jsonMsg.Stocks{
		outputValue += i.ResponseFields.Name +":"+strconv.Itoa(i.ResponseFields.Number)+":"+"$"+i.ResponseFields.Price+","
	}
	outputValue=strings.Trim(outputValue,",")
	outputValue+="\"\n\"unvestedAmount\":$"+strconv.FormatFloat(float64(jsonMsg.UnvestedAmount),'f',-1,32)		
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println("\nResponse:\n")
		fmt.Println(outputValue)
	}
}

func CheckPortfolio(){

	callerValue,err:= jsonrpc.Dial("tcp","127.0.0.1:8080")
	if err!=nil{
		fmt.Println(err)
		return
	}
	structSecondRequest:=new(SecondRequest)
	fmt.Println("Enter the request")
	var secondRequest string
	fmt.Scanf("%s",&secondRequest)
	secondRequest= strings.Replace(secondRequest,"\"","",-1)
	newsRequest:=strings.SplitN(secondRequest,":",-1)
	structSecondRequest.Tradeid,_= strconv.Atoi(newsRequest[1])
	result3 := &SecondRequest{
		Tradeid: structSecondRequest.Tradeid}
	result4,_:= json.Marshal(result3)
	var jsonMsg2 SecondResponse
	var reply string
	err = callerValue.Call("Server.LossOrGain",string(result4),&reply)
	var outputValue string
	outputValue = "\"stocks\":"
	json.Unmarshal([]byte(reply),&jsonMsg2)
	for _, i:= range jsonMsg2.Stocks{
		outputValue += i.ResponseFields.Name +":"+strconv.Itoa(i.ResponseFields.Number)+":"+i.ResponseFields.Price+","
	}
	outputValue=strings.Trim(outputValue,",")
	outputValue+="\"\n\"currentMarketValue\":$"+strconv.FormatFloat(float64(jsonMsg2.CurrentMarketValue),'f',-1,32)
	outputValue+="\n\"TotalUnInvestedAmount\":$"+strconv.FormatFloat(float64(jsonMsg2.UnvestedAmount),'f',-1,32)
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println("\nResponse:\n")
		fmt.Println(outputValue)
	}
}

func main(){
	fmt.Println("Enter your choice\n1.Buying Stocks\n2.Checking your portfolio")
	var choice int64 
	fmt.Scanf("%d\n",&choice)
	switch choice{
		case 1:
			PurchaseStocks()
			break
		case 2:
			CheckPortfolio()
			break
		default:
			fmt.Println("Please enter a valid choice")
			break
		}
}