package TestDeliveryGo

type TrackingStatus int

const (
	UnknownStatus    TrackingStatus = -1
	Pending          TrackingStatus = 1
	Ready            TrackingStatus = 2
	PickupCompelete  TrackingStatus = 3
	Loading          TrackingStatus = 4
	Unloading        TrackingStatus = 5
	DeleveryStart    TrackingStatus = 51
	DeleveryComplete TrackingStatus = 91
	DoNOtDelevery    TrackingStatus = 99
)

const (
	NoCode           string = "NO_CODE_AVAILABLE"
	NoTrackingInfo   string = "NO_TRACKING_INFO"
	ParseError       string = "PARSE_ERROR"
	RequestPageError string = "REQUEST_PAGE_ERROR"
)
