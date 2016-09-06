package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/barnettZQG/alertCenter/core"
	"github.com/barnettZQG/alertCenter/models"
	"github.com/barnettZQG/alertCenter/util"
)

type APIController struct {
	beego.Controller
	session      *core.MongoSession
	alertService *core.AlertService
	teamServcie  *core.TeamService
}

func (e *APIController) ReceiveAlert() {

	data := e.Ctx.Input.RequestBody
	if data != nil && len(data) > 0 {
		var AlertMessage *models.AlertReceive = &models.AlertReceive{}
		err := json.Unmarshal(data, AlertMessage)
		if err == nil {
			util.Info("get a alert message,receiver:" + AlertMessage.Receiver)
			go core.HandleMessage(AlertMessage)
			e.Data["json"] = util.GetSuccessJson("receive alert success")
		} else {
			util.Error("Parse the received message to AlertMessage faild." + err.Error())
		}
	} else {
		util.Error("receive a unknow data")
		//util.Info(string(data))
	}
	e.ServeJSON()
}
func (e *APIController) Receive() {

	data := e.Ctx.Input.RequestBody
	if data != nil && len(data) > 0 {
		var Alerts []*models.Alert = make([]*models.Alert, 0)
		err := json.Unmarshal(data, &Alerts)
		if err == nil {
			core.HandleAlerts(Alerts)
			e.Data["json"] = util.GetSuccessJson("receive alert success")
		} else {
			util.Error("Parse the received message to Alerts faild." + err.Error())
		}
	} else {
		util.Error("receive a unknow data")
		//util.Info(string(data))
	}
	e.ServeJSON()
}
func (e *APIController) AddTag() {
	we := &core.WeAlertSend{}
	if ok := we.GetAllTags(); ok {
		e.Data["json"] = util.GetSuccessJson("get weiTag success")
	} else {
		e.Data["json"] = util.GetFailJson("get weiTag faild")
	}
	e.ServeJSON()
}

func (e *APIController) HandleAlert() {
	ID := e.GetString(":ID")
	Type := e.GetString(":type")
	message := e.GetString("message")
	if len(ID) == 0 || len(Type) == 0 {
		e.Data["json"] = util.GetErrorJson("参数格式错误")
	} else {
		session := core.GetMongoSession()
		defer session.Close()
		alertService := core.GetAlertService(session)
		alert := alertService.FindByID(ID)
		if alert == nil {
			e.Data["json"] = util.GetFailJson("报警信息不存在，id信息错误")
		} else {
			if Type == "handle" {
				alert.IsHandle = 1
			} else if Type == "miss" {
				alert.IsHandle = -1
			}
			alert.HandleDate = time.Now()
			alert.HandleMessage = message
			if ok := alertService.Update(alert); ok {
				e.Data["json"] = util.GetSuccessJson("登记成功")
			} else {
				e.Data["json"] = util.GetFailJson("登记失败")
			}
		}
	}
	e.ServeJSON()
}

func (e *APIController) GetAlerts() {
	receiver := e.GetString(":receiver")
	e.session = core.GetMongoSession()
	if e.session == nil {
		e.Data["json"] = util.GetFailJson("get database session faild.")
	} else {
		e.alertService = core.GetAlertService(e.session)
		defer e.session.Close()
		if len(receiver) != 0 && receiver != "all" {
			alerts := e.alertService.FindByUser(receiver)
			util.Info("Get" + strconv.Itoa(len(alerts)) + " alerts")
			if alerts == nil {
				e.Data["json"] = util.GetFailJson("get database collection faild or receiver is error ")
			} else {
				e.Data["json"] = util.GetSuccessReJson(alerts)
			}
		} else if receiver == "all" {
			alerts := e.alertService.FindAll()
			if alerts == nil {
				e.Data["json"] = util.GetFailJson("get database collection faild")
			} else {
				e.Data["json"] = util.GetSuccessReJson(alerts)
			}
		} else {
			e.Data["json"] = util.GetErrorJson("api use error,please provide receiver")
		}
	}
	e.ServeJSON()
}
func (e *APIController) GetTeams() {

	e.session = core.GetMongoSession()
	if e.session == nil {
		e.Data["json"] = util.GetFailJson("get database session faild.")
	} else {
		relation := &core.Relation{}
		teams := relation.GetAllTeam()
		if teams == nil {
			e.Data["json"] = util.GetFailJson("There is no info of team")
		} else {
			e.Data["json"] = util.GetSuccessReJson(teams)
		}
	}
	e.ServeJSON()
}

func (e *APIController) AddTeam() {
	data := e.Ctx.Input.RequestBody
	if data != nil && len(data) > 0 {
		var team *models.Team = &models.Team{}
		err := json.Unmarshal(data, team)
		if err == nil {
			relation := &core.Relation{}
			relation.SetTeam(team)
			e.Data["json"] = util.GetSuccessJson("receive team info success")
		} else {
			util.Error("Parse the received message to teams faild." + err.Error())
			e.Data["json"] = util.GetFailJson("Parse the received message to teams faild.")
		}
	} else {
		util.Error("receive a unknow data")
		e.Data["jaon"] = util.GetErrorJson("receive a unknow data")
	}
	e.ServeJSON()
}
