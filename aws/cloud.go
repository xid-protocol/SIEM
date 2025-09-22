package aws

type AWSCloud struct {
	GuardDutyChan chan GdEvent
}

func NewAWSCloud() *AWSCloud {
	return &AWSCloud{
		GuardDutyChan: make(chan GdEvent, 100),
	}
}
