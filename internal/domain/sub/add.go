package sub

type AddSubscriptionParams struct {
	userID int64
	URL    string
}

type AddSubscriptionResult struct {
	subscriptionID int64
}

func (s *Subscriptioner) AddSubscription(params AddSubscriptionParams) (AddSubscriptionResult, error) {
	return AddSubscriptionResult{}, nil
}
