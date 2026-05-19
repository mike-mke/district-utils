package utils

import (
	"fmt"
	"strings"
)

const DEFAULT_DISTRICT_NO = "38" // AA district # for Milwaukee, WI

// NOTE: All of the following constants will be used in the JSON as attribute names after converting to lowercase.
// E.g., "name" is the json attribute
const COL_NAME = "NAME"      // Group Name
const COL_ID = "ID"          // NY ID
const COL_DISTRICT = "DIST"  // AA District number
const COL_PLACE = "LOCATION" // Name of the place where the Group meets
const COL_STREET = "STREET"
const COL_CITY = "CITY"
const COL_STATE = "STATE"
const COL_ZIP = "ZIP"
const COL_CODES = "CODES"
const COL_NOTES = "NOTES"
const COL_TIME = "TIME"
const COL_DAY = "DAY"
const COL_MOTM = "MOTM" // Meeting Of The Month (TRUE/FALSE)
const COL_SEC_FIRST_NAME = "SEC_FIRST_NAME"
const COL_EMAIL = "EMAIL" // email of Group Secretary
const COL_PHONE = "PHONE" // phone of Group Secretary
const COL_ZOOM_ID = "ZOOM_ID"
const COL_ZOOM_PASSCODE = "ZOOM_PASSCODE"
const COL_ZOOM_URL = "ZOOM_URL"

var columnHeaders = []string{
	COL_NAME,
	COL_ID,
	COL_DISTRICT,
	COL_PLACE,
	COL_STREET,
	COL_CITY,
	COL_STATE,
	COL_ZIP,
	COL_CODES,
	COL_NOTES,
	COL_TIME,
	COL_DAY,
	COL_MOTM,
	COL_SEC_FIRST_NAME,
	COL_EMAIL,
	COL_PHONE,
	COL_ZOOM_ID,
	COL_ZOOM_PASSCODE,
	COL_ZOOM_URL,
}

// MeetingCode defines the structure for meeting attributes
type MeetingCode struct {
	InternalID string `json:"lid"`
	Code       string `json:"code"`
	Label      string `json:"label"`
}

var MEETING_CODES = []MeetingCode{
	{InternalID: "code14", Code: "WB", Label: "Online Available"},
	{InternalID: "code01", Code: "A", Label: "Agnostics"},
	{InternalID: "code02", Code: "*", Label: "Al-Anon also meets"},
	{InternalID: "code04", Code: "B", Label: "Beginner's Class"},
	{InternalID: "code05", Code: "CC", Label: "Child Care Available"},
	{InternalID: "code06", Code: "DD", Label: "Dual Diagnosis"},
	{InternalID: "code07", Code: "G", Label: "Gay/Lesbian"},
	{InternalID: "code08", Code: "W", Label: "Handicap Access"},
	{InternalID: "code09", Code: "L", Label: "Ladies/Women"},
	{InternalID: "code10", Code: "M", Label: "Men's"},
	{InternalID: "code11", Code: "I", Label: "Native American"},
	{InternalID: "code15", Code: "PO", Label: "Polish Interpreter Available"},
	{InternalID: "code16", Code: "P", Label: "Professionals"},
	{InternalID: "code17", Code: "ASL", Label: "ASL Interpreter"},
	{InternalID: "code18", Code: "S", Label: "Spanish Speaking"},
	{InternalID: "code19", Code: "O", Label: "Weekly/Monthly Open Meeting"},
	{InternalID: "code20", Code: "Y", Label: "Young People Welcome"},
}

type AaGroup struct {
	data         map[string]string
	key          string
	meetingCodes []MeetingCode
}

func NewAaGroup() AaGroup {
	dataMap := make(map[string]string)
	codes := make([]MeetingCode, 0)
	group := AaGroup{data: dataMap, key: "", meetingCodes: codes}
	return group
}

func findKey(aKey string) *string {
	for _, key := range columnHeaders {
		if strings.EqualFold(key, aKey) {
			return &key
		}
	}
	return nil
}

func (group *AaGroup) GetDistrictNumber() string {
	value, ok := group.data[COL_DISTRICT]
	if ok && "" != strings.TrimSpace(value) {
		return value
	}
	return DEFAULT_DISTRICT_NO
}

func (group *AaGroup) SetValue(aKey string, aValue string) {
	trimmedValue := strings.TrimSpace(aValue)
	if len(trimmedValue) == 0 {
		return
	}
	key := findKey(aKey)
	if key == nil {
		group.data[aKey] = trimmedValue
	} else {
		group.data[*key] = aValue
	}
}

func CreateAaGroup(headers []string, values []string) AaGroup {
	group := NewAaGroup()
	for idx, colCell := range values {
		group.SetValue(headers[idx], colCell)
	}
	group.SetKey()

	value, ok := group.data[COL_CODES]
	if ok && "" != strings.TrimSpace(value) {
		mtgCodes := group.createMeetingCodes(value)
		group.meetingCodes = mtgCodes
	}
	return group
}

func (group *AaGroup) GetCityStateZip() string {
	result := ""
	value, ok := group.data[COL_CITY]
	if ok {
		result = (value + ", ")
	}
	value, ok = group.data[COL_STATE]
	if ok {
		result += (value + " ")
	}
	value, ok = group.data[COL_ZIP]
	if ok {
		result += (value + " ")
	}
	return result
}

func (group *AaGroup) HasZoomInfo() bool {
	_, ok := group.data[COL_ZOOM_ID]
	if ok {
		return ok
	}

	_, ok = group.data[COL_ZOOM_URL]
	return ok
}

func (group *AaGroup) IsMotm() bool {
	value, ok := group.data[COL_ZOOM_ID]
	if ok {
		if strings.EqualFold(value, "TRUE") {
			return true
		}
	}
	return false
}

func (group *AaGroup) Print() {
	for idx, value := range columnHeaders {
		fmt.Printf("Idx: %d  Value: %s \n", idx, value)
	}

	if group.HasZoomInfo() {
		fmt.Printf("Group has Zoom info! \n")
	}
	if group.IsMotm() {
		fmt.Printf("Group is the 'Meeting of the Month'! \n")
	}
}

func (group *AaGroup) PrintEmail() {
	secFirstName := group.data[COL_SEC_FIRST_NAME]
	secEmail := group.data[COL_EMAIL]
	groupName := group.data[COL_NAME]
	groupDay := group.data[COL_DAY]
	groupTime := group.data[COL_TIME]

	if len(secFirstName) > 0 && len(secEmail) > 0 {
		fmt.Printf("To: %s \n", secEmail)
		fmt.Printf("Subject: %s \n", groupName)
		fmt.Printf("Hi %s - \n\n", secFirstName)
		fmt.Printf("I'm writing to you about group '%s' that meets each %s at %s. \n", groupName, groupDay, groupTime)
		fmt.Printf("Milwaukee Central Office has you listed as the group's secretary. Is this still true? \n\n")
		fmt.Printf("Additionally, does your group have a GSR (General Service Representative)? If so, could you please let me know the person's name and email address? \n\n")
		fmt.Printf("Thank you for your service.\n\nWarm regards,\n\nMike\n")
	}
}

func (group *AaGroup) PrintHtml() {
	if group.IsMotm() {
		fmt.Printf("<div id='motm' class='aaGroup'>\n")
	} else {
		fmt.Printf("<div class='aaGroup'>\n")
	}
	fmt.Printf("<b>%s (%s) </b><br/> \n", group.data[COL_NAME], group.data[COL_ID])
	place, ok := group.data[COL_PLACE]
	if ok {
		fmt.Printf("%s <br/>\n", place)
	}
	street, ok := group.data[COL_STREET]
	if ok {
		fmt.Printf("%s <br/>\n", street)
	}
	csz := group.GetCityStateZip()
	if len(csz) > 0 {
		fmt.Printf("%s <br/>\n", csz)
	}
	instructions, ok := group.data[COL_NOTES]
	if ok {
		fmt.Printf("%s <br/>\n", instructions)
	}

	fmt.Printf("<b>%s \t%s </b><br/>\n", group.data[COL_DAY], group.data[COL_TIME])
	if group.HasZoomInfo() {
		fmt.Printf("\n<p><b>Zoom Info</b><br/>\n")
		zoomUrl, ok := group.data[COL_ZOOM_URL]
		if ok {
			fmt.Printf("%s<br/>\n", zoomUrl)
		}
		zoomId, ok := group.data[COL_ZOOM_ID]
		if ok {
			fmt.Printf("Meeting ID: %s &nbsp; &nbsp; Passcode: %s \n", zoomId, group.data[COL_ZOOM_PASSCODE])
		}
		fmt.Printf("</p>\n")
		fmt.Printf("<hr style='border:0.75px solid #bebebe; width:500px; margin-left:0; margin-top:-5px'>\n")
	} else {
		fmt.Printf("<hr style='border:0.75px solid #bebebe; width:500px; margin-left:0; margin-top:10px'>\n")
	}

	fmt.Printf("</div>\n")
}

func (group *AaGroup) PrintJavaScript() {
	group.printJsonImpl(true)
}

func (group *AaGroup) PrintJson() {
	group.printJsonImpl(false)
}

//	{
//	 "name": "Sample Group Name",
//	 "id": "Group ID",
//	 "day": "Monday",
//	 "time": "7:00 PM",
//	 "place": "Location name",
//	 "street": "street address",
//	 "city": "city",
//	 "state": "WI",
//	 "zip": "53092",
//	 "notes": "formerly called special instructions",
//	 "codes": ["w","s"]
//	 "zoom_id": "a zoom id",
//	 "zoom_passcode": "a zoom passcode",
//	 "zoom_url": "a zoom URL",
//		}
func (group *AaGroup) printJsonImpl(dontQuote bool) {
	fmt.Printf("{")

	fmt.Printf("%s,", getJson(COL_NAME, group.data[COL_NAME], dontQuote))
	fmt.Printf("%s,", getJson(COL_ID, group.data[COL_ID], dontQuote))
	fmt.Printf("%s,", getJson(COL_DAY, group.data[COL_DAY], dontQuote))
	fmt.Printf("%s,", getJson(COL_TIME, group.data[COL_TIME], dontQuote))
	fmt.Printf("%s,", getJson(COL_PLACE, group.data[COL_PLACE], dontQuote))
	fmt.Printf("%s,", getJson(COL_STREET, group.data[COL_STREET], dontQuote))
	fmt.Printf("%s,", getJson(COL_CITY, group.data[COL_CITY], dontQuote))
	fmt.Printf("%s,", getJson(COL_STATE, group.data[COL_STATE], dontQuote))
	fmt.Printf("%s,", getJson(COL_ZIP, group.data[COL_ZIP], dontQuote))
	fmt.Printf("%s,", getJson(COL_NOTES, group.data[COL_NOTES], dontQuote))

	codes := transformMeetingCodes(group.meetingCodes)
	if dontQuote {
		fmt.Printf("%s: [ %s ],", strings.ToLower(COL_CODES), codes)
	} else {
		fmt.Printf("\"%s\": [ %s ],", strings.ToLower(COL_CODES), codes)
	}

	fmt.Printf("%s,", getJson(COL_ZOOM_ID, group.data[COL_ZOOM_ID], dontQuote))
	fmt.Printf("%s,", getJson(COL_ZOOM_PASSCODE, group.data[COL_ZOOM_PASSCODE], dontQuote))
	fmt.Printf("%s,", getJson(COL_ZOOM_URL, group.data[COL_ZOOM_URL], dontQuote))
	fmt.Printf("%s", getJson(COL_MOTM, group.data[COL_MOTM], dontQuote))

	fmt.Printf("}")
}

// Example input: "NAME", "Some Group Name"
// This method will return: "name": "Some Group Name" (and that will be the single string returned (ie, the return value is wrapped in double quotes))
func getJson(columnName string, value string, dontQuote bool) string {
	name := strings.ToLower(columnName)
	val := strings.TrimSpace(value)
	if "" == val {
		if dontQuote {
			return name + ": null"
		} else {
			return "\"" + name + "\": null"
		}
	}
	if dontQuote {
		return name + ": \"" + val + "\""
	} else {
		return "\"" + name + "\": \"" + val + "\""
	}
}

func findMeetingCode(code string) *MeetingCode {
	key := strings.ToUpper(strings.TrimSpace(code))
	if "" == key {
		return nil
	} else {
		for _, c := range MEETING_CODES {
			if key == c.Code {
				return &c
			}
		}
	}
	return nil
}

func (group *AaGroup) createMeetingCodes(codes string) []MeetingCode {
	results := []MeetingCode{}
	for _, s := range strings.Split(codes, ",") {
		meetingCode := findMeetingCode(s)
		if meetingCode != nil {
			results = append(results, *meetingCode)
		}
	}
	return results
}

func transformMeetingCodes(meetingCodes []MeetingCode) string {
	if len(meetingCodes) == 0 {
		return ""
	}
	results := ""
	for idx, meetingCode := range meetingCodes {
		if idx > 0 {
			results += ", "
		}
		results += "\"" + meetingCode.Code + "\""
	}
	return results
}

func (group *AaGroup) PrintItem() {
	fmt.Printf("%s (%s) \n", group.data[COL_NAME], group.data[COL_ID])
	place, ok := group.data[COL_PLACE]
	if ok {
		fmt.Printf("%s \n", place)
	}
	street, ok := group.data[COL_STREET]
	if ok {
		fmt.Printf("%s \n", street)
	}
	csz := group.GetCityStateZip()
	if len(csz) > 0 {
		fmt.Printf("%s \n", csz)
	}
	instructions, ok := group.data[COL_NOTES]
	if ok {
		fmt.Printf("%s \n", instructions)
	}

	fmt.Printf("%s \t%s \n", group.data[COL_DAY], group.data[COL_TIME])
	if group.HasZoomInfo() {
		fmt.Printf("Zoom Info \n")
		zoomId, ok := group.data[COL_ZOOM_ID]
		if ok {
			fmt.Printf("Meeting ID: %s pwd: %s", zoomId, group.data[COL_ZOOM_PASSCODE])
		}
		zoomUrl, ok := group.data[COL_ZOOM_URL]
		if ok {
			fmt.Printf(" URL: %s \n", zoomUrl)
		}
		fmt.Printf("\n")
	}
}

func (group *AaGroup) GetKey() string {
	return group.key
}

func (group *AaGroup) SetKey() {
	grpTime, ok := group.data[COL_TIME]
	if ok {
		time, err := NewAaTime(grpTime)
		d := 0
		switch strings.ToLower(group.data[COL_DAY]) {
		case "monday":
			d = 1
		case "tuesday":
			d = 2
		case "wednesday":
			d = 3
		case "thursday":
			d = 4
		case "friday":
			d = 5
		case "saturday":
			d = 6
		case "sunday":
			d = 7
		}
		if err == nil {
			key := fmt.Sprintf("%d:%04s|%s", d, time.GetKey(), group.data[COL_ID])
			group.key = key
		} else {
			key := fmt.Sprintf("%d:????|%s", d, group.data[COL_ID])
			group.key = key
		}
	} else {
		group.key = ""
	}
}
