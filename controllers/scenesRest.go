package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"ngrinder-sampling/models"
	"strconv"
	"time"
)

type ScenesController struct {
	beego.Controller
}

// @router	/getScenesById	[get]
func (s *ScenesController) GetScenesById() {
	result := make(map[string]interface{})
	result["code"] = 0
	id, err := s.GetInt64("id")
	if err != nil {
		result["errMsg"] = err
		result["code"] = 1
		s.responseRst(result)
	}
	//获取testPms
	testPms, err := models.TestPmsGetById(id)
	if err != nil {
		result["code"] = 1
		result["errMsg"] = "the scenes of id " + strconv.FormatInt(id, 10) + " not exist"
		s.responseRst(result)
	}
	//获取requestPms
	requestPmsList, err := models.RequestPmsGetByTestPmsId(id)
	if err != nil {
		result["code"] = 1
		result["errMsg"] = err
		s.responseRst(result)
	}
	//组装数据
	sencesRequestBean, err := models.BuildScenesBean(testPms, &requestPmsList)
	if err != nil {
		result["code"] = 1
		result["errMsg"] = err
		s.responseRst(result)
	}

	result["data"] = sencesRequestBean
	s.responseRst(result)
}

// @router /create	[post]
func (s *ScenesController) Create() {
	result := make(map[string]interface{})
	result["code"] = 0

	var sencesRequestBean models.SencesRequestBean
	data := s.Ctx.Input.RequestBody
	//json数据封装到user对象中
	err := json.Unmarshal(data, &sencesRequestBean)
	if err != nil {
		result["code"] = 1
		result["errMsg"] = err.Error()
		s.responseRst(result)
	}

	//对参数完整性进行校验
	if err := models.ValidSencesParams(&sencesRequestBean); err != nil {
		result["code"] = 1
		result["errMsg"] = err.Error()
		s.responseRst(result)
	}
	//开启事务
	o := orm.NewOrm()
	//获取testPms
	testPms, err := models.GetTestPmsBean(sencesRequestBean)
	if err != nil {
		result["code"] = 1
		result["errMsg"] = err.Error()
		s.responseRst(result)
	}
	testPms.CreateTime = time.Now()
	testPms.UpdateTime = testPms.CreateTime

	o.Begin()
	id, dbErr := models.TestPmsSave(testPms, o)
	if dbErr != nil {
		result["code"] = 1
		result["errMsg"] = dbErr.Error()
		s.responseRst(result)
	}
	testPms.Id = id

	//获取requestPms
	requestPmsArray, err := models.GetRequestPmsBean(sencesRequestBean)
	if err != nil {
		result["code"] = 1
		result["errMsg"] = dbErr.Error()
		s.responseRst(result)
	}
	for i, requestPms := range *requestPmsArray {
		requestPms.TestPmsId = id
		//设置生成脚本时的函数名
		requestPms.FunName = "test" + strconv.Itoa(i);
		pid, dbError := models.RequestPmsSave(&requestPms, o)
		if dbError != nil {
			//事务回滚
			o.Rollback()
			result["code"] = 1
			result["errMsg"] = dbError.Error()
			s.responseRst(result)
		}
		requestPms.Id = pid
	}
	o.Commit()

	result["data"] = id
	s.responseRst(result)
}

// @router	/update	[post]
func (s *ScenesController) Update() {
	result := make(map[string]interface{})
	result["code"] = 0

	var sencesRequestBean models.SencesRequestBean
	data := s.Ctx.Input.RequestBody
	//json数据封装到user对象中
	err := json.Unmarshal(data, &sencesRequestBean)
	if err != nil {
		result["code"] = 1
		result["errMsg"] = err.Error()
		s.responseRst(result)
	}

	//对参数完整性进行校验
	if err := models.ValidUpdateSencesParams(&sencesRequestBean); err != nil {
		result["code"] = 1
		result["errMsg"] = err.Error()
		s.responseRst(result)
	}
	//开启事务
	o := orm.NewOrm()
	//获取testPms
	testPms, err := models.GetTestPmsBean(sencesRequestBean)
	if err != nil {
		result["code"] = 1
		result["errMsg"] = err.Error()
		s.responseRst(result)
	}
	testPms.UpdateTime = time.Now()

	o.Begin()
	dbErr := models.TestPmsUpdate(testPms, o)
	if dbErr != nil {
		result["code"] = 1
		result["errMsg"] = dbErr.Error()
		s.responseRst(result)
	}

	//获取requestPms
	requestPmsArray, err := models.GetRequestPmsBean(sencesRequestBean)
	if err != nil {
		result["code"] = 1
		result["errMsg"] = dbErr.Error()
		s.responseRst(result)
	}
	for _, requestPms := range *requestPmsArray {
		dbError := models.RequestPmsUpdate(&requestPms, o)
		if dbError != nil {
			//事务回滚
			o.Rollback()
			result["code"] = 1
			result["errMsg"] = dbError.Error()
			s.responseRst(result)
		}
	}
	o.Commit()

	s.responseRst(result)
}

// @router	/delete	[get]
func (s *ScenesController) Delete() {
	result := make(map[string]interface{})
	result["code"] = 0
	id, err := s.GetInt64("id")
	if err != nil {
		result["errMsg"] = err
		result["code"] = 1
		s.responseRst(result)
	}
	//获取testPms
	testPms, err := models.TestPmsGetById(id)
	if err != nil {
		result["code"] = 1
		result["errMsg"] = "the scenes of id " + strconv.FormatInt(id, 10) + " not exist"
		s.responseRst(result)
	}
	//开启事务，删除requestPms、testPms
	o := orm.NewOrm()
	o.Begin()
	_, err = models.RequestPmsDeleteByTestPmsId(id, o)
	if err != nil {
		o.Rollback()
		result["errMsg"] = err
		result["code"] = 1
		s.responseRst(result)
	}
	_, err = models.TestPmsDelete(testPms, o)
	if err != nil {
		o.Rollback()
		result["errMsg"] = err
		result["code"] = 1
		s.responseRst(result)
	}
	o.Commit()

	//增加请求，删除svn上的脚本和压测数据源文件
	requestBean, err := models.BuildScenesBean(testPms, nil)
	if err != nil {
		result["errMsg"] = err
		result["code"] = 1
		s.responseRst(result)
	}
	var path string
	if len(requestBean.FileDataList) > 0 {
		for index, fileData := range requestBean.FileDataList {
			if index != 0 {
				path += ","
			}
			path += fileData.Path
		}
	}
	ngrinderUrl := beego.AppConfig.String("ngrinder.serverurl")
	apiUrl := beego.AppConfig.String("ngrinder.api.create")
	ngrinderUrl += apiUrl + "?id=" + strconv.FormatInt(id, 10)
	if len(path) > 0 {
		ngrinderUrl += "&path=" + path
	}
	req := httplib.Get(ngrinderUrl)
	var js struct {
		Code   int
		ErrMsg string
	}
	rst, err := req.String()
	if err != nil {
		result["code"] = 1
		result["errMsg"] = err
	} else {
		err = json.Unmarshal([]byte(rst), &js)
		if err != nil {
			result["code"] = 1
			result["errMsg"] = rst
		} else {
			result["code"] = js.Code
			result["errMsg"] = js.ErrMsg
		}
	}

	s.responseRst(result)
}

//handle the result
func (s *ScenesController) responseRst(result map[string]interface{}) {
	s.Data["json"] = result
	s.ServeJSON()
	s.StopRun()
}
