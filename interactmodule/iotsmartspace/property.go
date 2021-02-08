package iotsmartspace

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/globalsign/mgo/bson"
	"github.com/h3c/iotzigbeeserver-go/config"
	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globalmsgtype"
	"github.com/h3c/iotzigbeeserver-go/models"
	"github.com/h3c/iotzigbeeserver-go/publicstruct"
	"github.com/h3c/iotzigbeeserver-go/zcl/common"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/h3c/iotzigbeeserver-go/zcl/zclmsgdown"
	"github.com/lib/pq"
)

// genOnOffSocketCommandDown 处理clusterID 0x0006
func genOnOffSocketCommandDown(devEUI string, value string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][genOnOffSocketCommandDown] value: %s", devEUI, value)
	endpointIndex := string([]byte(value)[:1])
	index, _ := strconv.Atoi(endpointIndex)
	onOffCommand := string([]byte(value)[1:2])
	var commandID uint8
	if onOffCommand == "0" {
		commandID = 0x00
	} else {
		commandID = 0x01
	}
	cmd := common.Command{
		Cmd:              commandID,
		DstEndpointIndex: index,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SocketCommand,
		ClusterID:   uint16(cluster.OnOff),
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// genOnOffSwitchCommandDown 处理clusterID 0x0006
func genOnOffSwitchCommandDown(devEUI string, value string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][genOnOffSwitchCommandDown] value: %s", devEUI, value)
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][genOnOffSwitchCommandDown] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
		return
	}
	if terminalInfo != nil {
		endpointTemp := pq.StringArray{}
		if constant.Constant.UsePostgres {
			endpointTemp = terminalInfo.EndpointPG
		} else {
			endpointTemp = terminalInfo.Endpoint
		}
		endpoint := string([]byte(value)[:1])
		endpoint = "0" + endpoint
		var dstEndpointIndex int
		for i, v := range endpointTemp {
			if v == endpoint {
				dstEndpointIndex = i
			}
		}
		onOffCommand := string([]byte(value)[1:2])
		var commandID uint8
		if onOffCommand == "0" {
			commandID = 0x00
		} else {
			commandID = 0x01
		}
		cmd := common.Command{
			Cmd:              commandID,
			DstEndpointIndex: dstEndpointIndex,
		}
		zclDownMsg := common.ZclDownMsg{
			MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
			DevEUI:      devEUI,
			CommandType: common.SwitchCommand,
			ClusterID:   uint16(cluster.OnOff),
			Command:     cmd,
			MsgID:       msgID,
		}
		zclmsgdown.ProcZclDownMsg(zclDownMsg)
	}
}

// honyarSocketCommandDown 处理鸿雁sokcet command down
func honyarSocketCommandDown(devEUI string, value string, msgID interface{}, commandType string) {
	globallogger.Log.Infof("[devEUI: %v][honyarSocketCommandDown] value: %s", devEUI, value)
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][honyarSocketCommandDown] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
		return
	}
	if terminalInfo != nil {
		endpointTemp := pq.StringArray{}
		if constant.Constant.UsePostgres {
			endpointTemp = terminalInfo.EndpointPG
		} else {
			endpointTemp = terminalInfo.Endpoint
		}
		endpoint := string([]byte(value)[:1])
		endpoint = "0" + endpoint
		var dstEndpointIndex int
		for i, v := range endpointTemp {
			if v == endpoint {
				dstEndpointIndex = i
			}
		}
		cmd := common.Command{
			DstEndpointIndex: dstEndpointIndex,
		}
		var clusterID uint16
		switch commandType {
		case common.ChildLock:
			clusterID = uint16(cluster.OnOff)
			cmd.Value = string([]byte(value)[1:2]) == "1"
			cmd.AttributeID = 0x8000
			cmd.AttributeName = "ChildLock"
		case common.OnOffUSB:
			clusterID = uint16(cluster.OnOff)
			if string([]byte(value)[1:2]) == "0" {
				cmd.Cmd = 0x00
			} else {
				cmd.Cmd = 0x01
			}
			cmd.DstEndpointIndex = 1
		case common.BackGroundLight:
			clusterID = uint16(cluster.OnOff)
			cmd.Value = string([]byte(value)[1:2]) == "1"
			cmd.AttributeID = 0x8005
			cmd.AttributeName = "BackGroundLight"
		case common.PowerOffMemory:
			clusterID = uint16(cluster.OnOff)
			cmd.Value = string([]byte(value)[1:2]) == "0"
			cmd.AttributeID = 0x8006
			cmd.AttributeName = "PowerOffMemory"
		case common.HistoryElectricClear:
			clusterID = uint16(cluster.SmartEnergyMetering)
			cmd.Value = string([]byte(value)[1:2]) == "1"
			cmd.AttributeID = 0x8000
			cmd.AttributeName = "HistoryElectricClear"
		}
		if commandType == common.OnOffUSB {
			commandType = common.SocketCommand
		}
		zclDownMsg := common.ZclDownMsg{
			MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
			DevEUI:      devEUI,
			CommandType: commandType,
			ClusterID:   clusterID,
			Command:     cmd,
			MsgID:       msgID,
		}
		zclmsgdown.ProcZclDownMsg(zclDownMsg)
	}
}

// heimanWarningDeviceProcAlarm 处理clusterID 0x0502
func heimanWarningDeviceProcAlarm(devEUI string, alarmTime string, alarmStart string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanWarningDeviceProcAlarm] alarmTime: %s, alarmStart: %s", devEUI, alarmTime, alarmStart)
	time, err := strconv.ParseUint(alarmTime, 10, 16)
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanWarningDeviceProcAlarm] ParseUint time err: %+v", devEUI, err)
	}
	var warningControl uint8 = 0x0000
	if alarmStart == "0" {
		warningControl = 0x0014
	}
	endpointIndex := 0
	cmd := common.Command{
		DstEndpointIndex: endpointIndex,
		WarningControl:   warningControl,
		WarningTime:      uint16(time),
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.StartWarning,
		ClusterID:   0x0502,
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// heimanTemperatureAlarmThresholdSet 处理clusterID 0xfc81
func heimanTemperatureAlarmThresholdSet(devEUI string, value string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanTemperatureAlarmThresholdSet] value: %s", devEUI, value)
	valueArr := strings.Split(value, ",")
	if len(valueArr) == 1 {
		valueArr = strings.Split(value, "，")
	}
	minValue, err := strconv.ParseInt(valueArr[0], 10, 16)
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanTemperatureAlarmThresholdSet] ParseInt minValue err: %+v", devEUI, err)
	}
	maxValue, err := strconv.ParseInt(valueArr[1], 10, 16)
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanTemperatureAlarmThresholdSet] ParseInt maxValue err: %+v", devEUI, err)
	}
	minValue = minValue * 100
	maxValue = maxValue * 100
	if minValue < 0 {
		minValue = minValue + 65536
	}
	if maxValue < 0 {
		maxValue = maxValue + 65536
	}
	var terminalInfo *config.TerminalInfo
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanTemperatureAlarmThresholdSet] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
	}
	var disturb uint16 = 0x00f0
	if terminalInfo.Disturb == "80c0" {
		disturb = 0x80c0
	}
	endpointTemp := pq.StringArray{}
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	endpointIndex := 0
	if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
		endpointIndex = 1
	}
	cmd := common.Command{
		DstEndpointIndex: endpointIndex,
		MaxTemperature:   int16(maxValue),
		MinTemperature:   int16(minValue),
		Disturb:          disturb,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SetThreshold,
		ClusterID:   0xfc81,
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// heimanHumidityAlarmThresholdSet 处理clusterID 0xfc81
func heimanHumidityAlarmThresholdSet(devEUI string, value string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanHumidityAlarmThresholdSet] value: %s", devEUI, value)
	valueArr := strings.Split(value, ",")
	if len(valueArr) == 1 {
		valueArr = strings.Split(value, "，")
	}
	minValue, err := strconv.ParseUint(valueArr[0], 10, 16)
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanHumidityAlarmThresholdSet] ParseUint minValue err: %+v", devEUI, err)
	}
	maxValue, err := strconv.ParseUint(valueArr[1], 10, 16)
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanHumidityAlarmThresholdSet] ParseUint maxValue err: %+v", devEUI, err)
	}
	minValue = minValue * 100
	maxValue = maxValue * 100
	var terminalInfo *config.TerminalInfo
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanHumidityAlarmThresholdSet] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
	}
	var disturb uint16 = 0x00f0
	if terminalInfo.Disturb == "80c0" {
		disturb = 0x80c0
	}
	endpointTemp := pq.StringArray{}
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	endpointIndex := 0
	if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
		endpointIndex = 1
	}
	cmd := common.Command{
		DstEndpointIndex: endpointIndex,
		MaxHumidity:      uint16(maxValue),
		MinHumidity:      uint16(minValue),
		Disturb:          disturb,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SetThreshold,
		ClusterID:   0xfc81,
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// heimanDisturbAlarmThresholdSet 处理clusterID 0xfc81
func heimanDisturbAlarmThresholdSet(devEUI string, value string, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanDisturbAlarmThresholdSet] value: %s", devEUI, value)
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][heimanDisturbAlarmThresholdSet] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
	}
	endpointTemp := pq.StringArray{}
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	endpointIndex := 0
	if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
		endpointIndex = 1
	}
	var disturb uint16 = 0x00f0
	if value == "1" { // 开启勿扰模式
		disturb = 0x80c0
	}
	var maxTemperature int64 = 6000
	var minTemperature int64 = -1000
	var maxHumidity uint64 = 10000
	var minHumidity uint64 = 0
	if terminalInfo.MaxTemperature != "" {
		maxTemperature, _ = strconv.ParseInt(terminalInfo.MaxTemperature, 10, 16)
	}
	if terminalInfo.MinTemperature != "" {
		minTemperature, _ = strconv.ParseInt(terminalInfo.MinTemperature, 10, 16)
	}
	if terminalInfo.MaxHumidity != "" {
		maxHumidity, _ = strconv.ParseUint(terminalInfo.MaxHumidity, 10, 16)
	}
	if terminalInfo.MinHumidity != "" {
		minHumidity, _ = strconv.ParseUint(terminalInfo.MinHumidity, 10, 16)
	}
	cmd := common.Command{
		DstEndpointIndex: endpointIndex,
		MaxTemperature:   int16(maxTemperature),
		MinTemperature:   int16(minTemperature),
		MaxHumidity:      uint16(maxHumidity),
		MinHumidity:      uint16(minHumidity),
		Disturb:          disturb,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SetThreshold,
		ClusterID:   0xfc81,
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// heimanIRControlEMSendKeyCommand 处理clusterID 0xfc82
func heimanIRControlEMSendKeyCommand(devEUI string, value cluster.HEIMANInfraredRemoteSendKeyCommand, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMSendKeyCommand] value: %+v", devEUI, value)
	cmd := common.Command{
		InfraredRemoteID:      value.ID,
		InfraredRemoteKeyCode: value.KeyCode,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SendKeyCommand,
		ClusterID:   0xfc82,
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// heimanIRControlEMStudyKey 处理clusterID 0xfc82
func heimanIRControlEMStudyKey(devEUI string, value cluster.HEIMANInfraredRemoteStudyKey, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMStudyKey] value: %+v", devEUI, value)
	cmd := common.Command{
		InfraredRemoteID:      value.ID,
		InfraredRemoteKeyCode: value.KeyCode,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.StudyKey,
		ClusterID:   0xfc82,
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// heimanIRControlEMDeleteKey 处理clusterID 0xfc82
func heimanIRControlEMDeleteKey(devEUI string, value cluster.HEIMANInfraredRemoteDeleteKey, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMDeleteKey] value: %+v", devEUI, value)
	cmd := common.Command{
		InfraredRemoteID:      value.ID,
		InfraredRemoteKeyCode: value.KeyCode,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.DeleteKey,
		ClusterID:   0xfc82,
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// heimanIRControlEMCreateID 处理clusterID 0xfc82
func heimanIRControlEMCreateID(devEUI string, value cluster.HEIMANInfraredRemoteCreateID, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMCreateID] value: %+v", devEUI, value)
	cmd := common.Command{
		InfraredRemoteModelType: value.ModelType,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.CreateID,
		ClusterID:   0xfc82,
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// heimanIRControlEMGetIDAndKeyCodeList 处理clusterID 0xfc82
func heimanIRControlEMGetIDAndKeyCodeList(devEUI string, value cluster.HEIMANInfraredRemoteGetIDAndKeyCodeList, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanIRControlEMGetIDAndKeyCodeList] value: %+v", devEUI, value)
	cmd := common.Command{}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.GetIDAndKeyCodeList,
		ClusterID:   0xfc82,
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// heimanHS2AQSetLanguage 处理clusterID 0xfc81
func heimanHS2AQSetLanguage(devEUI string, language uint8, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanHS2AQSetLanguage] language: %+v", devEUI, language)
	cmd := common.Command{
		Cmd: language,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SetLanguage,
		ClusterID:   0xfc81,
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// heimanHS2AQSetUnitOfTemperature 处理clusterID 0xfc81
func heimanHS2AQSetUnitOfTemperature(devEUI string, unitOfTemperature uint8, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][heimanHS2AQSetUnitOfTemperature] unitOfTemperature: %+v", devEUI, unitOfTemperature)
	unitoftemperature := "C"
	if unitOfTemperature == 0 {
		unitoftemperature = "F"
	}
	if constant.Constant.UsePostgres {
		models.FindTerminalAndUpdatePG(map[string]interface{}{"deveui": devEUI}, map[string]interface{}{"unitoftemperature": unitoftemperature})
	} else {
		models.FindTerminalAndUpdate(bson.M{"devEUI": devEUI}, bson.M{"unitOfTemperature": unitoftemperature})
	}
	cmd := common.Command{
		Cmd: unitOfTemperature,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: common.SetUnitOfTemperature,
		ClusterID:   0xfc81,
		Command:     cmd,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// procPropertyIotware 处理属性消息
func procPropertyIotware(mqttMsg publicstruct.RPCIotware) {
	globallogger.Log.Infof("[procPropertyIotware]: mqttMsg.Data.Params: %+v", mqttMsg.Data.Params)
	method := mqttMsg.Data.Method
	mapParams := mqttMsg.Data.Params.(map[string]interface{})
	devEUI := mqttMsg.Device
	switch method {
	case IotwarePropertyOnOff:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyOnOff:
				if value.(string) == "00" {
					genOnOffSwitchCommandDown(devEUI, "10", mqttMsg.Data.ID)
				} else if value.(string) == "01" {
					genOnOffSwitchCommandDown(devEUI, "11", mqttMsg.Data.ID)
				}
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertyOnOffOne:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyOnOffOne:
				if value.(string) == "00" {
					genOnOffSwitchCommandDown(devEUI, "10", mqttMsg.Data.ID)
				} else if value.(string) == "01" {
					genOnOffSwitchCommandDown(devEUI, "11", mqttMsg.Data.ID)
				}
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertyOnOffTwo:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyOnOffTwo:
				if value.(string) == "00" {
					genOnOffSwitchCommandDown(devEUI, "20", mqttMsg.Data.ID)
				} else if value.(string) == "01" {
					genOnOffSwitchCommandDown(devEUI, "21", mqttMsg.Data.ID)
				}
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertyOnOffThree:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyOnOffThree:
				if value.(string) == "00" {
					genOnOffSwitchCommandDown(devEUI, "30", mqttMsg.Data.ID)
				} else if value.(string) == "01" {
					genOnOffSwitchCommandDown(devEUI, "31", mqttMsg.Data.ID)
				}
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertyDownWarning:
		for key, value := range mapParams {
			switch key {
			case "time":
				heimanWarningDeviceProcAlarm(devEUI, strconv.FormatUint(uint64(value.(uint16)), 10), "0", mqttMsg.Data.ID)
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertyClearWarning:
		heimanWarningDeviceProcAlarm(devEUI, "1", "1", mqttMsg.Data.ID)
	case IotwarePropertyTThreshold:
		var lower uint16
		var upper uint16
		for key, value := range mapParams {
			switch key {
			case "lower":
				lower = value.(uint16)
			case "upper":
				upper = value.(uint16)
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
		heimanTemperatureAlarmThresholdSet(devEUI, strconv.FormatUint(uint64(lower), 10)+","+strconv.FormatUint(uint64(upper), 10), mqttMsg.Data.ID)
	case IotwarePropertyHThreshold:
		var lower uint16
		var upper uint16
		for key, value := range mapParams {
			switch key {
			case "lower":
				lower = value.(uint16)
			case "upper":
				upper = value.(uint16)
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
		heimanHumidityAlarmThresholdSet(devEUI, strconv.FormatUint(uint64(lower), 10)+","+strconv.FormatUint(uint64(upper), 10), mqttMsg.Data.ID)
	case IotwarePropertyNotDisturb:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyNotDisturb:
				heimanDisturbAlarmThresholdSet(devEUI, value.(string), mqttMsg.Data.ID)
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertySendKeyCommand:
		var IRId uint8 = 0xff
		var keyCode uint8 = 0xff
		for key, value := range mapParams {
			switch key {
			case "IRId":
				IRId = value.(uint8)
			case "keyCode":
				keyCode = value.(uint8)
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
		heimanIRControlEMSendKeyCommand(devEUI, cluster.HEIMANInfraredRemoteSendKeyCommand{ID: IRId, KeyCode: keyCode}, mqttMsg.Data.ID)
	case IotwarePropertyStudyKey:
		var IRId uint8 = 0xff
		var keyCode uint8 = 0xff
		for key, value := range mapParams {
			switch key {
			case "IRId":
				IRId = value.(uint8)
			case "keyCode":
				keyCode = value.(uint8)
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
		heimanIRControlEMStudyKey(devEUI, cluster.HEIMANInfraredRemoteStudyKey{ID: IRId, KeyCode: keyCode}, mqttMsg.Data.ID)
	case IotwarePropertyDeleteKey:
		var IRId uint8 = 0xff
		var keyCode uint8 = 0xff
		for key, value := range mapParams {
			switch key {
			case "IRId":
				IRId = value.(uint8)
			case "keyCode":
				keyCode = value.(uint8)
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
		heimanIRControlEMDeleteKey(devEUI, cluster.HEIMANInfraredRemoteDeleteKey{ID: IRId, KeyCode: keyCode}, mqttMsg.Data.ID)
	case IotwarePropertyCreateID:
		var modelType uint8
		for key, value := range mapParams {
			switch key {
			case "modelType":
				modelType = value.(uint8)
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
		heimanIRControlEMCreateID(devEUI, cluster.HEIMANInfraredRemoteCreateID{ModelType: modelType}, mqttMsg.Data.ID)
	case IotwarePropertyGetIDAndKeyCodeList:
		heimanIRControlEMGetIDAndKeyCodeList(devEUI, cluster.HEIMANInfraredRemoteGetIDAndKeyCodeList{}, mqttMsg.Data.ID)
	case IotwarePropertyChildLock:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyChildLock:
				if value.(string) == "00" {
					honyarSocketCommandDown(devEUI, "10", mqttMsg.Data.ID, common.ChildLock)
				} else if value.(string) == "01" {
					honyarSocketCommandDown(devEUI, "11", mqttMsg.Data.ID, common.ChildLock)
				}
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertyUSBSwitch:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyUSBSwitch:
				if value.(string) == "00" {
					honyarSocketCommandDown(devEUI, "10", mqttMsg.Data.ID, common.OnOffUSB)
				} else if value.(string) == "01" {
					honyarSocketCommandDown(devEUI, "11", mqttMsg.Data.ID, common.OnOffUSB)
				}
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertyBackLightSwitch:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyBackLightSwitch:
				if value.(string) == "00" {
					honyarSocketCommandDown(devEUI, "10", mqttMsg.Data.ID, common.BackGroundLight)
				} else if value.(string) == "01" {
					honyarSocketCommandDown(devEUI, "11", mqttMsg.Data.ID, common.BackGroundLight)
				}
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertyMemorySwitch:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyMemorySwitch:
				if value.(string) == "00" {
					honyarSocketCommandDown(devEUI, "10", mqttMsg.Data.ID, common.PowerOffMemory)
				} else if value.(string) == "01" {
					honyarSocketCommandDown(devEUI, "11", mqttMsg.Data.ID, common.PowerOffMemory)
				}
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertyHistoryClear:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyHistoryClear:
				if value.(string) == "0" {
					honyarSocketCommandDown(devEUI, "10", mqttMsg.Data.ID, common.HistoryElectricClear)
				} else if value.(string) == "1" {
					honyarSocketCommandDown(devEUI, "11", mqttMsg.Data.ID, common.HistoryElectricClear)
				}
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertySwitchLeftUpLogo,
		IotwarePropertySwitchLeftMidLogo,
		IotwarePropertySwitchLeftDownLogo,
		IotwarePropertySwitchRightUpLogo,
		IotwarePropertySwitchRightMidLogo,
		IotwarePropertySwitchRightDownLogo:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertySwitchLeftUpLogo,
				IotwarePropertySwitchLeftMidLogo,
				IotwarePropertySwitchLeftDownLogo,
				IotwarePropertySwitchRightUpLogo,
				IotwarePropertySwitchRightMidLogo,
				IotwarePropertySwitchRightDownLogo:
				honyarAddSceneCommandDown(devEUI, value.(string), mqttMsg.Data.ID, key, common.AddScene)
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertyLanguage:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyLanguage:
				if value.(string) == "0" {
					heimanHS2AQSetLanguage(devEUI, 0x00, mqttMsg.Data.ID)
				} else if value.(string) == "1" {
					heimanHS2AQSetLanguage(devEUI, 0x01, mqttMsg.Data.ID)
				}
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	case IotwarePropertyUnitOfTemperature:
		for key, value := range mapParams {
			switch key {
			case IotwarePropertyUnitOfTemperature:
				if value.(string) == "0" {
					heimanHS2AQSetUnitOfTemperature(devEUI, 0x01, mqttMsg.Data.ID)
				} else if value.(string) == "1" {
					heimanHS2AQSetUnitOfTemperature(devEUI, 0x00, mqttMsg.Data.ID)
				}
			default:
				globallogger.Log.Warnf("[procPropertyIotware]: invalid key: %s", key)
			}
		}
	default:
		globallogger.Log.Warnf("[procPropertyIotware]: invalid method: %s", method)
	}
}

// procProperty 处理属性消息
func procProperty(mqttMsg mqttMsgSt) {
	globallogger.Log.Infof("[procProperty]: mqttMsg.Params: %+v", mqttMsg.Params)
	mapParams := mqttMsg.Params.(map[string]interface{})
	devEUI := mapParams["terminalId"]
	var alarmTime string
	var alarmStart string
	for key, value := range mapParams {
		if key == "terminalId" {
			continue
		}
		switch key {
		case HeimanSmartPlugPropertyOnOff:
			fallthrough
		case HeimanESocketPropertyOnOff, HonyarSocket000a0c3cPropertyOnOff,
			HonyarSocketHY0105PropertyOnOff, HonyarSocketHY0106PropertyOnOff:
			globallogger.Log.Errorf("[devEUI: %v][procProperty] subscribe time: %+v", devEUI, time.Now())
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(devEUI.(string), "10", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(devEUI.(string), "11", mqttMsg.ID)
			}
		case HeimanHS2SW1LEFR30PropertyOnOff1, HonyarSingleSwitch00500c32PropertyOnOff1:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(devEUI.(string), "10", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(devEUI.(string), "11", mqttMsg.ID)
			}
		case HeimanHS2SW2LEFR30PropertyOnOff1, HonyarDoubleSwitch00500c33PropertyOnOff1:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(devEUI.(string), "10", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(devEUI.(string), "11", mqttMsg.ID)
			}
		case HeimanHS2SW2LEFR30PropertyOnOff2, HonyarDoubleSwitch00500c33PropertyOnOff2:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(devEUI.(string), "20", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(devEUI.(string), "21", mqttMsg.ID)
			}
		case HeimanHS2SW3LEFR30PropertyOnOff1, HonyarTripleSwitch00500c35PropertyOnOff1:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(devEUI.(string), "10", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(devEUI.(string), "11", mqttMsg.ID)
			}
		case HeimanHS2SW3LEFR30PropertyOnOff2, HonyarTripleSwitch00500c35PropertyOnOff2:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(devEUI.(string), "20", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(devEUI.(string), "21", mqttMsg.ID)
			}
		case HeimanHS2SW3LEFR30PropertyOnOff3, HonyarTripleSwitch00500c35PropertyOnOff3:
			if value.(string) == "00" {
				genOnOffSwitchCommandDown(devEUI.(string), "30", mqttMsg.ID)
			} else if value.(string) == "01" {
				genOnOffSwitchCommandDown(devEUI.(string), "31", mqttMsg.ID)
			}
		case HeimanWarningDevicePropertyAlarmTime:
			alarmTime = value.(string)
			if alarmStart != "" {
				heimanWarningDeviceProcAlarm(devEUI.(string), alarmTime, alarmStart, mqttMsg.ID)
			}
		case HeimanWarningDevicePropertyAlarmStart:
			alarmStart = value.(string)
			if alarmTime != "" {
				heimanWarningDeviceProcAlarm(devEUI.(string), alarmTime, alarmStart, mqttMsg.ID)
			}
			if alarmStart == "1" {
				alarmTime = "1"
				heimanWarningDeviceProcAlarm(devEUI.(string), alarmTime, alarmStart, mqttMsg.ID)
			}
		case HeimanHS2AQPropertyTemperatureAlarmThreshold:
			heimanTemperatureAlarmThresholdSet(devEUI.(string), value.(string), mqttMsg.ID)
		case HeimanHS2AQPropertyHumidityAlarmThreshold:
			heimanHumidityAlarmThresholdSet(devEUI.(string), value.(string), mqttMsg.ID)
		case HeimanHS2AQPropertyUnDisturb:
			heimanDisturbAlarmThresholdSet(devEUI.(string), value.(string), mqttMsg.ID)
		case HeimanIRControlEMPropertySendKeyCommand:
			heimanIRControlEMSendKeyCommand(devEUI.(string), value.(cluster.HEIMANInfraredRemoteSendKeyCommand), mqttMsg.ID)
		case HeimanIRControlEMPropertyStudyKey:
			heimanIRControlEMStudyKey(devEUI.(string), value.(cluster.HEIMANInfraredRemoteStudyKey), mqttMsg.ID)
		case HeimanIRControlEMPropertyDeleteKey:
			heimanIRControlEMDeleteKey(devEUI.(string), value.(cluster.HEIMANInfraredRemoteDeleteKey), mqttMsg.ID)
		case HeimanIRControlEMPropertyCreateID:
			heimanIRControlEMCreateID(devEUI.(string), value.(cluster.HEIMANInfraredRemoteCreateID), mqttMsg.ID)
		case HeimanIRControlEMPropertyGetIDAndKeyCodeList:
			heimanIRControlEMGetIDAndKeyCodeList(devEUI.(string), value.(cluster.HEIMANInfraredRemoteGetIDAndKeyCodeList), mqttMsg.ID)
		case HonyarSocket000a0c3cPropertyChildLock, HonyarSocketHY0105PropertyChildLock, HonyarSocketHY0106PropertyChildLock:
			if value.(string) == "00" {
				honyarSocketCommandDown(devEUI.(string), "10", mqttMsg.ID, common.ChildLock)
			} else if value.(string) == "01" {
				honyarSocketCommandDown(devEUI.(string), "11", mqttMsg.ID, common.ChildLock)
			}
		case HonyarSocket000a0c3cPropertyUSB, HonyarSocketHY0105PropertyUSB:
			if value.(string) == "00" {
				honyarSocketCommandDown(devEUI.(string), "10", mqttMsg.ID, common.OnOffUSB)
			} else if value.(string) == "01" {
				honyarSocketCommandDown(devEUI.(string), "11", mqttMsg.ID, common.OnOffUSB)
			}
		case HonyarSocketHY0105PropertyBackGroundLight, HonyarSocketHY0106PropertyBackGroundLight:
			if value.(string) == "00" {
				honyarSocketCommandDown(devEUI.(string), "10", mqttMsg.ID, common.BackGroundLight)
			} else if value.(string) == "01" {
				honyarSocketCommandDown(devEUI.(string), "11", mqttMsg.ID, common.BackGroundLight)
			}
		case HonyarSocketHY0105PropertyPowerOffMemory, HonyarSocketHY0106PropertyPowerOffMemory:
			if value.(string) == "00" {
				honyarSocketCommandDown(devEUI.(string), "10", mqttMsg.ID, common.PowerOffMemory)
			} else if value.(string) == "01" {
				honyarSocketCommandDown(devEUI.(string), "11", mqttMsg.ID, common.PowerOffMemory)
			}
		case HonyarSocketHY0105PropertyHistoryElectricClear, HonyarSocketHY0106PropertyHistoryElectricClear:
			if value.(string) == "0" {
				honyarSocketCommandDown(devEUI.(string), "10", mqttMsg.ID, common.HistoryElectricClear)
			} else if value.(string) == "1" {
				honyarSocketCommandDown(devEUI.(string), "11", mqttMsg.ID, common.HistoryElectricClear)
			}
		case Honyar6SceneSwitch005f0c3bPropertyLeftUpAddScene,
			Honyar6SceneSwitch005f0c3bPropertyLeftMiddleAddScene,
			Honyar6SceneSwitch005f0c3bPropertyLeftDownAddScene,
			Honyar6SceneSwitch005f0c3bPropertyRightUpAddScene,
			Honyar6SceneSwitch005f0c3bPropertyRightMiddleAddScene,
			Honyar6SceneSwitch005f0c3bPropertyRightDownAddScene:
			honyarAddSceneCommandDown(devEUI.(string), value.(string), mqttMsg.ID, key, common.AddScene)
		default:
			globallogger.Log.Warnf("[procProperty]: invalid key: %s", key)
		}
	}
}

// procScenesCommand 处理scenes command
func procScenesCommand(devEUI string, commandType string, command common.Command, msgID interface{}) {
	globallogger.Log.Infof("[devEUI: %v][procScenesCommand] command: %+v", devEUI, command)
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI)
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI)
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][procScenesCommand] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
	}
	endpointTemp := pq.StringArray{}
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	endpointIndex := 0
	if len(endpointTemp) > 1 && endpointTemp[1] == "01" {
		endpointIndex = 1
	}
	command.DstEndpointIndex = endpointIndex
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI,
		CommandType: commandType,
		ClusterID:   0x0005,
		Command:     command,
		MsgID:       msgID,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}

// honyarAddSceneCommandDown 处理场景消息
func honyarAddSceneCommandDown(devEUI string, value string, msgID interface{}, key string, commandType string) {
	value = strings.ToUpper(hex.EncodeToString([]byte(mahonia.NewEncoder("gbk").ConvertString(value))))
	command := common.Command{}
	command.GroupID = 1
	command.TransitionTime = 0
	command.SceneName = "0" + strconv.FormatInt(int64(len(value)/2), 16) + value
	command.KeyID = 1
	switch key {
	case IotwarePropertySwitchLeftUpLogo, Honyar6SceneSwitch005f0c3bPropertyLeftUpAddScene:
		command.SceneID = 1
	case IotwarePropertySwitchLeftMidLogo, Honyar6SceneSwitch005f0c3bPropertyLeftMiddleAddScene:
		command.SceneID = 2
	case IotwarePropertySwitchLeftDownLogo, Honyar6SceneSwitch005f0c3bPropertyLeftDownAddScene:
		command.SceneID = 3
	case IotwarePropertySwitchRightUpLogo, Honyar6SceneSwitch005f0c3bPropertyRightUpAddScene:
		command.SceneID = 4
	case IotwarePropertySwitchRightMidLogo, Honyar6SceneSwitch005f0c3bPropertyRightMiddleAddScene:
		command.SceneID = 5
	case IotwarePropertySwitchRightDownLogo, Honyar6SceneSwitch005f0c3bPropertyRightDownAddScene:
		command.SceneID = 6
	}
	procScenesCommand(devEUI, commandType, command, msgID)
}

// procScenes 处理场景消息
func procScenes(mqttMsg mqttMsgSt) {
	globallogger.Log.Infof("[procScenes]: mqttMsg.Params: %+v", mqttMsg.Params)
	mapParams := mqttMsg.Params.(map[string]interface{})
	devEUI := mapParams["terminalId"]
	command := common.Command{}
	for k, v := range mapParams {
		switch k {
		case "GroupID":
			tempGroupID, _ := strconv.ParseUint(v.(string), 10, 16)
			command.GroupID = uint16(tempGroupID)
		case "SceneID":
			tempSceneID, _ := strconv.ParseUint(v.(string), 10, 8)
			command.SceneID = uint8(tempSceneID)
		case "TransitionTime":
			tempTransitionTime, _ := strconv.ParseUint(v.(string), 10, 16)
			command.TransitionTime = uint16(tempTransitionTime)
		case "SceneName":
			command.SceneName = v.(string)
		case "KeyID":
			tempKeyID, _ := strconv.ParseUint(v.(string), 10, 8)
			command.KeyID = uint8(tempKeyID)
		case "Mode":
			tempMode, _ := strconv.ParseUint(v.(string), 10, 8)
			command.Mode = uint8(tempMode)
		case "FromGroupID":
			tempFromGroupID, _ := strconv.ParseUint(v.(string), 10, 16)
			command.FromGroupID = uint16(tempFromGroupID)
		case "FromSceneID":
			tempFromSceneID, _ := strconv.ParseUint(v.(string), 10, 16)
			command.FromSceneID = uint16(tempFromSceneID)
		case "ToGroupID":
			tempToGroupID, _ := strconv.ParseUint(v.(string), 10, 16)
			command.ToGroupID = uint16(tempToGroupID)
		case "ToSceneID":
			tempToSceneID, _ := strconv.ParseUint(v.(string), 10, 16)
			command.ToSceneID = uint16(tempToSceneID)
		}
	}
	for key := range mapParams {
		if key == "terminalId" || key == "GroupID" || key == "SceneID" || key == "TransitionTime" || key == "SceneName" ||
			key == "KeyID" || key == "Mode" || key == "FromGroupID" || key == "FromSceneID" || key == "ToGroupID" || key == "ToSceneID" {
			continue
		}
		switch key {
		case common.AddScene, common.ViewScene, common.RemoveScene, common.RemoveAllScenes, common.StoreScene,
			common.RecallScene, common.GetSceneMembership, common.EnhancedAddScene, common.EnhancedViewScene, common.CopyScene:
			procScenesCommand(devEUI.(string), key, command, mqttMsg.ID)
		default:
			globallogger.Log.Warnf("[procScenes]: invalid key: %s", key)
		}
	}
}

// procReadAttribute 处理read req
func procReadAttribute(mqttMsg mqttMsgSt) {
	globallogger.Log.Infof("[procReadAttribute]: mqttMsg.Params: %+v", mqttMsg.Params)
	mapParams := mqttMsg.Params.(map[string]interface{})
	devEUI := mapParams["terminalId"]
	var clusterID uint16
	var endpoint string
	for k, v := range mapParams {
		switch k {
		case "clusterID":
			tempClusterID, _ := strconv.ParseUint(v.(string), 10, 16)
			clusterID = uint16(tempClusterID)
		case "endpoint":
			endpoint = v.(string)
		}
	}
	var terminalInfo *config.TerminalInfo
	var err error
	if constant.Constant.UsePostgres {
		terminalInfo, err = models.GetTerminalInfoByDevEUIPG(devEUI.(string))
	} else {
		terminalInfo, err = models.GetTerminalInfoByDevEUI(devEUI.(string))
	}
	if err != nil {
		globallogger.Log.Errorf("[devEUI: %v][procReadAttribute] GetTerminalInfoByDevEUI err: %+v", devEUI, err)
	}
	endpointTemp := pq.StringArray{}
	if constant.Constant.UsePostgres {
		endpointTemp = terminalInfo.EndpointPG
	} else {
		endpointTemp = terminalInfo.Endpoint
	}
	endpointIndex := 0
	if len(endpointTemp) > 1 {
		for index, EP := range endpointTemp {
			if endpoint == EP {
				endpointIndex = index
			}
		}
	}
	cmd := common.Command{
		DstEndpointIndex: endpointIndex,
	}
	zclDownMsg := common.ZclDownMsg{
		MsgType:     globalmsgtype.MsgType.DOWNMsg.ZigbeeCmdRequestEvent,
		DevEUI:      devEUI.(string),
		CommandType: common.ReadAttribute,
		ClusterID:   clusterID,
		Command:     cmd,
	}
	zclmsgdown.ProcZclDownMsg(zclDownMsg)
}