package factual

import(
	"log"
	"fmt"
	"errors"
	"strconv"
	// "net/url"
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

type Hours map[string][][]string

// Factual gives its hours as a seperatly encoded entity
func (hours *Hours) UnmarshalJSON(data []byte) (err error){
	// intermediary to avoid an infinite loop
	i := map[string][][]string{}
	s, err := strconv.Unquote(string(data))
	err = json.Unmarshal([]byte(s), &i)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid hours in JSON: %s (%s)", string(data), err))
	}
	*hours = Hours(i)
	return
}

type Place struct{
	Tel string `json:"tel"`
	Name string `json:"name"`
	Email string `json:"email"`
	Website string `json:"website"`
	Hours Hours `json:"hours"`
	HoursDisplay string `json:"hours_display"`
	FactualId string `json:"factual_id"`
	Address string `json:"address"`
	Neighborhood []string `json:"neighborhood"`
	Region string `json:"region"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Postcode string `json:"postcode"`
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

type Query struct{
	Factual Factual
	Table *Table
	Id string
}

func (this *Table) Filters(f F) *Query{
	return &Query{this.Factual, this, ""}
}

func (this *Table) Id(id string) *Query{
	return &Query{this.Factual, this, id}
}

// func (this *Table) Search(s string) *Query{

// }

func (this *Query) url() (u string){
	u = this.Table.url()
	if this.Id != ""{
		u += "/" + this.Id
	}
	return
}

func (this *Query) params() map[string]string{
	return map[string]string{}
}

func (this *Query) Iter() *Iter{
	return &Iter{this.Factual, this}
}

func (this *Query) All() error{
	return this.Iter().All()
}

func (this *Query) One(r *interface{}) error{
	return this.Iter().One(r)
}

type Iter struct{
	Factual Factual
	Query *Query
	// Results interface{}
}

func (this *Iter) fetch() (r Response, err error){
	t := this.Factual.Token
	r = Response{}
	log.Println(this.Query.url())
	resp, err := this.Factual.Consumer.Get(this.Query.url(), this.Query.params(), t)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return
	}

	err = json.Unmarshal(body, &r)

	if r.ErrorType != ""{
		err = FactualError{r.ErrorType, r.Message}
	}

	return
}

// func (this *Iter) Next(r interface{}) error{

// }

func (this *Iter) All() error{
	_, err := this.fetch()
	return err
	// check the user passed in an array pointer
	// resultv := reflect.ValueOf(result)
	// if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
	// 	panic("result argument must be a slice address")
	// }

	// for{
	// 	v = &
	// 	if err = this.Next(v); err == nil{
	// 		result = append(result, v)
	// 	}
	// }

	// return 
}

func (this *Iter) One(r interface{}) (err error){
	resp, err := this.fetch()
	err = json.Unmarshal(resp.Response.Data[0], r)
	return
}