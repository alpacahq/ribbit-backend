package mobile

// Service is the interface to our mobile service
type Service interface {
	GenerateSMSToken(countryCode, mobile string) error
	CheckCode(countryCode, mobile, code string) error
}
