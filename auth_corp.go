package monobank

// TODO: add NewCorporateAuthorizer https://api.monobank.ua/docs/corporate.html
// type CorporateAuthorizer struct {
// 	keyID     string // X-Key-Id
// 	requestID string // X-Request-Id
// 	sign      string // X-Sign (X-time | X-Request-Id | URL)
// }
//
// Sign (X-time | X-Request-Id | URL)
//  - X-time=1561560962 (Wed, 26 Jun 2019 14:56:02 GMT)
//  - X-Request-Id=acW5k2ERnupgnWFyBkCY0nA
//  - URL=/personal/client-info
//  Sign (1561560962acW5k2ERnupgnWFyBkCY0nA/personal/client-info)=eaBHn+T18kr7w6uAzLJ037o1w/JMpAfV81yNaHaXJUxv9bbi/cORuA0gGsazwG+VxCq2Y+TmIb81zbGbuiaRQA==
