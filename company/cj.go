package company

type CJ struct{}

func (t CJ) Code() string {
	return "CJ"
}

func (t CJ) Name() string {
	return "CJ대한통운"
}

func (t CJ) TrackingUrl() string {
	return "http://nexs.cjgls.com/web/info.jsp?slipno=%s"
}
