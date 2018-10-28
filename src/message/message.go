package message

import (
	"Utils"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"log"
	"net/http"
	"time"
)

const email_address_key = "emailValue"

var keyspaceAndTable = Utils.GetConfig(Utils.KEYSPACE_KEY) + "." + Utils.TABLE_NAME

type DBMessageModel struct {
	Id              gocql.UUID
	EmailAddress	string
	Title			string
	Content			string
	MagicNumber		int
}

type POSTMessageModel struct {
	EmailAddress	string `json:"email"`
	Title			string `json:"title"`
	Content			string `json:"content"`
	MagicNumber		int `json:"magic_number"`
}

type HttpResp struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

type SendMessageReqModel struct {
	MagicNumber		int `json:"magic_number"`
}

func GetEmailsByAddress(response http.ResponseWriter, request *http.Request) {
	var messages []DBMessageModel
	var err error

	vars := mux.Vars(request)
	emailAddress := vars[email_address_key]

	if err := validateMailAddress(emailAddress); err != nil {
		json.NewEncoder(response).Encode(HttpResp{Status: 404, Description: Utils.INVALID_EMAIL_ADDRESS, Body: err.Error()})
		return
	}

	messages, err = getByAddress(&emailAddress)

	if len(messages) == 0 {
		json.NewEncoder(response).Encode(HttpResp{Status: 404, Description: Utils.NOT_FOUND + "; " + Utils.NO_MESSAGE_FOR_MAIL_ADDRESS , Body: ""})
		return
	}

	if err != nil {
		fmt.Fprintln(response, err.Error())
		return
	}

	json.NewEncoder(response).Encode(messages)
}

func SendMessage(response http.ResponseWriter, request *http.Request) {
	var reqModel SendMessageReqModel
	var messages []DBMessageModel
	var err error
	decoder := json.NewDecoder(request.Body)

	if err := decoder.Decode(&reqModel); err != nil{
		json.NewEncoder(response).Encode(HttpResp{Status: 400, Description: Utils.BAD_REQUEST, Body: err.Error()})
		return
	}

	messages, err  = getByMagicNum(&reqModel.MagicNumber)

	if err != nil {
		fmt.Fprintln(response, err.Error())
		return
	}

	if len(messages) == 0 {
		json.NewEncoder(response).Encode(HttpResp{Status: 404, Description: Utils.NOT_FOUND + "; " + Utils.NO_MESSAGE_FOR_MAGIC_NUM , Body: ""})
		return
	}

	if err := sendMails(messages); err != nil {
		fmt.Fprintln(response, err.Error())
		return
	}

	json.NewEncoder(response).Encode(HttpResp{Status: 200, Description: Utils.SUCCESS_SEND_MESSAGE, Body: ""})

}

func CreateMessage(response http.ResponseWriter, request *http.Request) {

	var reqMessage POSTMessageModel

	decoder := json.NewDecoder(request.Body)

	if err := decoder.Decode(&reqMessage); err != nil {
		json.NewEncoder(response).Encode(HttpResp{Status: 400, Description: Utils.BAD_REQUEST, Body: err.Error()})
		return
	}

	if err := validateMessageModel(&reqMessage); err != nil{
		json.NewEncoder(response).Encode(HttpResp{Status: 404, Description: Utils.NOT_FOUND, Body: err.Error()})
		return
	}

	dbMessage := reqMessage.toDBModel()

	err := save(dbMessage)

	if err != nil{
		json.NewEncoder(response).Encode(HttpResp{Status: 500, Description: Utils.FAILED_INSERT_MESSAGE, Body: ""})
		return
	}

	json.NewEncoder(response).Encode(HttpResp{Status: 200, Description: Utils.SUCCESS_INSER_MESSAGE, Body: ""})
}

func getByAddress(emailAddress *string) ([]DBMessageModel, error) {
	var messages []DBMessageModel
	stmt, names := qb.Select(keyspaceAndTable).
		Where(qb.Eq("email_address")).AllowFiltering().
		ToCql()


	q := gocqlx.Query(Session.Query(stmt), names).BindMap(qb.M{
		"email_address": emailAddress,
	})

	if err := q.SelectRelease(&messages); err != nil {
		log.Println(err)
		return nil, err
	}

	return messages, nil
}

func getByMagicNum(magicNum *int) ([]DBMessageModel, error) {
	stmt, names := qb.Select(keyspaceAndTable).
		Where(qb.Eq("magic_number")).AllowFiltering().ToCql()

	var message []DBMessageModel

	q := gocqlx.Query(Session.Query(stmt), names).BindMap(qb.M{
		"magic_number": magicNum,
	})

	if err := q.SelectRelease(&message); err != nil {
		log.Println(err)
		return nil, err
	}

	return message, nil
}

func save(reqMessage *DBMessageModel) error {

	stmt, names := qb.Insert(keyspaceAndTable).
		Columns("id","email_address", "title", "content", "magic_number").TTL(5 * time.Minute).
		ToCql()

	if err := gocqlx.Query(Session.Query(stmt), names).BindStruct(&reqMessage).ExecRelease(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func delete(id gocql.UUID) error {

	stmt, names := qb.Delete(keyspaceAndTable).Where(qb.Eq("id")).ToCql()

	q := gocqlx.Query(Session.Query(stmt), names).BindMap(qb.M{
		"id": id,
	})

	if err := q.ExecRelease(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (model *POSTMessageModel) toDBModel() *DBMessageModel  {
	dbModel := DBMessageModel{EmailAddress: model.EmailAddress, Title: model.Title, Content: model.Content, MagicNumber: model.MagicNumber}
	dbModel.Id = gocql.TimeUUID()
	return &dbModel
}
