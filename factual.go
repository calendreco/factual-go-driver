package factual

import(
	"log"
	"fmt"
	"reflect"
	"io/ioutil"
	"encoding/json"
	"github.com/mrjones/oauth"
)

const API_URL = "http://api.v3.factual.com"
const DRIVER_VERSION_TAG = "factual-golang-driver-0.1"

type FactualError struct{
	Type string
	Message string
}

func (e FactualError) Error() string {
	return fmt.Sprintf("%v: %v", e.Type, e.Message)
}

type Response struct{
	Version int `json:"version"`
	Status string `json:"status"`
	ErrorType string `json:"error_type"`
	Message string `json:"message"`
	Response struct{
		Data []json.RawMessage `json:"data"`
		IncludedRows int `json:"included_rows"`
	}
}

type F map[string]interface{}

type Factual struct{
	Consumer *oauth.Consumer
	Token *oauth.AccessToken
}

func New(key string, secret string) Factual{
	consumer := oauth.NewConsumer(key, secret, oauth.ServiceProvider{})
	return Factual{
		Consumer: consumer,
		Token: &oauth.AccessToken{},
	}
}

// type Table struct {Factual string, Name string, }
type Table struct{
	Factual Factual
	Name string
}

func (this Factual) Table(name string) *Table{
	return &Table{this, name}
}

func (this *Table) url() string{
	return API_URL + "/t/" + this.Name
}

type Query interface{
	Params() map[string]string
	Url() string
	Iter() *Iter
}

type Iter struct{
	Factual Factual
	Query Query
}

func (this *Iter) fetch() (r Response, err error){
	t := this.Factual.Token
	url := this.Query.Url()
	params := this.Query.Params()
	params = map[string]string{}
	log.Println("request", url, params)
	resp, err := this.Factual.Consumer.Get(url, params, t)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return
	}

	r = Response{}
	err = json.Unmarshal(body, &r)

	if r.ErrorType != ""{
		err = FactualError{r.ErrorType, r.Message}
	}

	return
}

func (this *Iter) All(rs interface{}) (err error){
	resp, err := this.fetch()

	ref := reflect.ValueOf(rs).Elem()
	elemt := ref.Type().Elem()

	for _, r := range resp.Response.Data{
		// Generate a new pointer, resolve it
		elemp := reflect.New(elemt)
		// Generate an interface for unmarshalling
		i := elemp.Interface()
		err = json.Unmarshal(r, &i)
		// Add it to the results
		ref = reflect.Append(ref, elemp.Elem())
	}

	reflect.ValueOf(rs).Elem().Set(ref)

	return
}

func (this *Iter) One(r interface{}) (err error){
	resp, err := this.fetch()
	err = json.Unmarshal(resp.Response.Data[0], r)
	return
}