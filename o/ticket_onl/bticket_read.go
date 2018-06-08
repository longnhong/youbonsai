package video

import (
	"LongTM/basic/x/math"
	"LongTM/basic/x/rest"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

func CheckCustomerCode(customerCode string, branchID string) (tk *Video, err error) {
	var timeBeginDay = math.BeginningOfDay().Unix()
	var tiemEnOfday = math.EndOfDay().Unix()
	var queryMatch = bson.M{
		"customer_code": customerCode,
		"branch_id":     branchID,
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		//"status": BookingStateCreated,
	}
	err = VideoTable.FindOne(queryMatch, &tk)
	if err != nil {
		if err.Error() == "not found" {
			err = rest.BadRequestValid(errors.New("Code không tồn tại!"))
		}
		return
	}
	if tk.Status == BookingStateSancelled || tk.Status == BookingStateFinished {
		err = errors.New("Vé đã kết thúc!")
	} else if tk.CheckInAt != 0 {
		err = rest.BadRequestValid(errors.New("Code đã sử dụng!"))
	}
	return
}

func GetCustomerMySchedule(customerId string) (btks []*RateVideo, err error) {
	var status = []string{string(BookingStateSancelled), (string(BookingStateFinished))}
	var queryMatch = bson.M{
		"customer_id": customerId,
		"updated_at":  bson.M{"$ne": 0},
		"status":      bson.M{"$in": status},
	}
	var query = []bson.M{}
	var joinRate = bson.M{
		"from":         "rate",
		"localField":   "_id",
		"foreignField": "bVideo_id",
		"as":           "rate",
	}
	var unWindRate = bson.M{"path": "$rate", "preserveNullAndEmptyArrays": true}
	query = []bson.M{
		{"$match": queryMatch},
		{"$lookup": joinRate},
		{"$unwind": unWindRate},
	}
	err = VideoTable.Pipe(query).All(&btks)
	rest.IsErrorRecord(err)
	return btks, err
}

func GetAllVideoCus(customerId string) (btks []*RateVideo, err error) {
	var queryMatch = bson.M{
		"customer_id": customerId,
		"updated_at":  bson.M{"$ne": 0},
	}
	var query = []bson.M{}
	var joinRate = bson.M{
		"from":         "rate",
		"localField":   "_id",
		"foreignField": "bVideo_id",
		"as":           "rate",
	}
	var unWindRate = bson.M{"path": "$rate", "preserveNullAndEmptyArrays": true}
	query = []bson.M{
		{"$match": queryMatch},
		{"$lookup": joinRate},
		{"$unwind": unWindRate},
	}
	return btks, VideoTable.Pipe(query).All(&btks)
}

func CheckVideoByDay(customerId string) (btks []*Video, err error) {
	var timeBeginDay = math.BeginningOfDay().Unix()
	var tiemEnOfday = math.EndOfDay().Unix()
	var queryMatch = bson.M{
		"customer_id": customerId,
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		"status": BookingStateCreated,
	}
	err = VideoTable.FindWhere(queryMatch, &btks)
	return btks, rest.IsErrorRecord(err)
}

func GetVideoNear(customerId string) (btk *RateVideo, err error) {
	var queryMatch = bson.M{
		"customer_id": customerId,
		"status":      BookingStateFinished,
	}
	var query = []bson.M{}
	var sort = bson.M{
		"created_at": -1,
	}
	var joinRate = bson.M{
		"from":         "rate",
		"localField":   "_id",
		"foreignField": "bVideo_id",
		"as":           "rate",
	}
	var unWindRate = bson.M{"path": "$rate", "preserveNullAndEmptyArrays": true}
	query = []bson.M{
		{"$match": queryMatch},
		{"$lookup": joinRate},
		{"$unwind": unWindRate},
		{"$sort": sort},
	}

	var btks []*RateVideo
	err = VideoTable.Pipe(query).All(&btks)
	if err == nil && len(btks) > 0 {
		btk = btks[0]
	}
	rest.IsErrorRecord(err)
	return btk, err
}

func (tk *Video) UpdateTimeCheckIn() error {
	var timeNow = math.GetTimeNowVietNam().Unix()
	var tracks = tk.updateTrack(tk.ServiceID, tk.BranchID, BookingStateConfirmed, timeNow)
	var up = bson.M{
		"updated_at":  timeNow,
		"check_in_at": timeNow,
		"status":      BookingStateConfirmed,
		"tracks":      tracks,
	}
	var err = VideoTable.UnsafeUpdateByID(tk.ID, up)
	if err == nil {
		tk.CheckInAt = timeNow
		tk.Status = BookingStateConfirmed
		tk.Tracks = tracks
	}
	return err
}

func (tk *Video) UpdateByCnumCetm(cnum string, idCetm string) error {
	var up = bson.M{
		"cnum_cetm":      cnum,
		"id_Video_cetm": idCetm,
		"status":         BookingStateConfirmed,
	}
	var err = VideoTable.UnsafeUpdateByID(tk.ID, up)
	if err == nil {
		tk.CnumCetm = cnum
		tk.IdVideoCetm = idCetm
		tk.Status = BookingStateConfirmed
	}
	return err
}

func GetVideoInBranch(branchID string, timeStart int64, timeEnd int64) (btks []*VideoUser, err error) {
	var queryMatch = bson.M{
		"branch_id": branchID,
		"time_go_bank": bson.M{
			"$gte": timeStart,
			"$lte": timeEnd,
		},
		"status": BookingStateCreated,
	}
	var query = []bson.M{}
	var joinUser = bson.M{
		"from":         "user",
		"localField":   "customer_id",
		"foreignField": "_id",
		"as":           "customer",
	}
	var unWindCus = "$customer"
	query = []bson.M{
		{"$match": queryMatch},
		{"$lookup": joinUser},
		{"$unwind": unWindCus},
	}
	err = VideoTable.Pipe(query).All(&btks)
	rest.IsErrorRecord(err)
	return btks, err
}
func GetVideoTimeInBranch(branchID string, timeStart int64, timeEnd int64) (btks []*Video, err error) {
	var queryMatch = bson.M{
		"branch_id": branchID,
		"time_go_bank": bson.M{
			"$gte": timeStart,
			"$lte": timeEnd,
		},
		"status": BookingStateCreated,
	}

	err = VideoTable.FindWhere(queryMatch, &btks)
	return btks, err
}

func GetAllVideo() (btks []*Video, err error) {
	var timeBeginDay = math.BeginningOfDay().Unix()
	var tiemEnOfday = math.EndOfDay().Unix()
	var queryMatch = bson.M{
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		"status": BookingStateCreated,
	}
	return btks, VideoTable.FindWhere(queryMatch, &btks)
}

func GetVideoByUser(cusId string) (btks []*Video, err error) {
	var timeBeginDay = math.BeginningOfDay().Unix()
	var tiemEnOfday = math.EndOfDay().Unix()
	var queryMatch = bson.M{
		"customer_id": cusId,
		"time_go_bank": bson.M{
			"$gte": timeBeginDay,
			"$lte": tiemEnOfday,
		},
		"status": BookingStateCreated,
	}
	return btks, VideoTable.FindWhere(queryMatch, &btks)
}

func GetAllVideoByTimeSearch(timeSearch int64, typeVideo TypeVideo) (btks []*Video, err error) {
	var start, end = math.BeginAndEndDay(timeSearch)
	var queryMatch = bson.M{
		"time_go_bank": bson.M{
			"$gte": start,
			"$lte": end,
		},
		"status": BookingStateCreated,
	}
	if typeVideo == TypeSchedule {
		queryMatch["type_Video"] = TypeSchedule
	}
	return btks, VideoTable.FindWhere(queryMatch, &btks)
}

func SearchVideo(idBranchs []string, timeStart int64, timeEnd int64) (btks []*VideoSchedule, err error) {
	//var start, end = math.BeginAndEndDay(timeSearch)
	var queryMatch = bson.M{
		"branch_id": bson.M{"$in": idBranchs},
		"time_go_bank": bson.M{
			"$gte": timeStart,
			"$lte": timeEnd,
		},
		"status": BookingStateCreated,
	}
	var group = bson.M{"_id": "$branch_id", "count": bson.M{"$sum": 1}}
	var query = []bson.M{
		{"$match": queryMatch},
		{"$group": group},
	}
	return btks, VideoTable.Pipe(query).All(&btks)
}

func GetByID(id string) (tk *Video, err error) {
	err = VideoTable.FindByID(id, &tk)
	rest.IsErrorRecord(err)
	return
}

func GetVideoByUserNeedFeedback(userId string) (tks []*Video, err error) {
	var queryMatch = bson.M{
		"customer_id": userId,
		"status":      BookingStateFinished,
	}
	err = VideoTable.FindWhere(queryMatch, &tks)
	rest.IsErrorRecord(err)
	return
}

func UpdateRate(id string, numRate TypeRate) (err error) {
	var up = bson.M{
		"is_rate": numRate,
	}
	err = VideoTable.UnsafeUpdateByID(id, up)
	rest.IsErrorRecord(err)
	return
}

func GetVideoReport(branchIds []string, timeStart int64, timeEnd int64, skip int64, limit int64) (btks []*VideoUser, err error) {
	btks = make([]*VideoUser, 0)
	var queryMatch = bson.M{
		"branch_id": bson.M{"$in": branchIds},
		"created_at": bson.M{
			"$gte": timeStart,
			"$lte": timeEnd,
		},
	}
	var query = []bson.M{}
	var joinUser = bson.M{
		"from":         "user",
		"localField":   "customer_id",
		"foreignField": "_id",
		"as":           "customer",
	}
	var unWindCus = bson.M{"path": "$customer", "preserveNullAndEmptyArrays": true}
	query = []bson.M{
		{"$match": queryMatch},
		{"$lookup": joinUser},
		{"$unwind": unWindCus},
		{"$skip": skip},
		{"$limit": limit},
	}
	err = VideoTable.Pipe(query).All(&btks)
	rest.IsErrorRecord(err)
	return btks, err
}

type VideoDetailReport struct {
	ID      int           `json:"id" bson:"_id"`
	Videos []*VideoUser `json:"Videos" bson:"Videos"`
}

func GetDetailReport(branchIds []string, timeStart int64, timeEnd int64) (btks []*VideoDetailReport, err error) {
	btks = make([]*VideoDetailReport, 0)
	var queryMatch = bson.M{
		"branch_id": bson.M{"$in": branchIds},
		"created_at": bson.M{
			"$gte": timeStart,
			"$lte": timeEnd,
		},
	}
	var query = []bson.M{}
	var joinUser = bson.M{
		"from":         "user",
		"localField":   "customer_id",
		"foreignField": "_id",
		"as":           "customer",
	}
	var unWindCus = bson.M{"path": "$customer", "preserveNullAndEmptyArrays": true}
	var project1 = bson.M{
		"date":        bson.M{"$add": []interface{}{math.NewTimeVN(), bson.M{"$multiply": []interface{}{"$time_go_bank", 1000}}}}, //7*60*1000
		"customer_id": 1, "service_id": 1, "service_name": 1, "branch_id": 1, "branch_address": 1,
		"type_Video": 1, "lang": 1, "customer_code": 1, "check_in_at": 1, "avatar_teller": 1, "id_Video_cetm": 1,
		"branch_name": 1, "tracks": 1, "cnum_cetm": 1, "teller_id": 1, "teller": 1,
		"serving_time": 1, "waiting_time": 1, "is_rate": 1, "status": 1, "customer": 1, "time_go_bank": 1, "updated_at": 1, "created_at": 1,
	}
	var group = bson.M{
		"_id":     bson.M{"$hour": "$date"},
		"Videos": bson.M{"$push": "$$ROOT"},
	}
	var pro2 = bson.M{
		"_id":     bson.M{"$add": []interface{}{"$_id", 7}},
		"Videos": "$Videos",
	}
	query = []bson.M{
		{"$match": queryMatch},
		{"$lookup": joinUser},
		{"$unwind": unWindCus},
		{"$project": project1},
		{"$group": group},
		{"$project": pro2},
	}
	err = VideoTable.Pipe(query).All(&btks)
	rest.IsErrorRecord(err)

	return btks, err
}

func GetVideoReportByTime(branchIds []string, timeStart int64, timeEnd int64) (int, error) {
	var queryMatch = bson.M{
		"branch_id": bson.M{"$in": branchIds},
		"created_at": bson.M{
			"$gte": timeStart,
			"$lte": timeEnd,
		},
	}
	var count, err = VideoTable.CountWhere(queryMatch)
	rest.IsErrorRecord(err)
	return count, err
}
